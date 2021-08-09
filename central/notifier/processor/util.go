package processor

import (
	"context"

	"github.com/stackrox/rox/central/notifiers"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sac"
)

// Sending alerts.
//////////////////

func tryToAlert(ctx context.Context, notifier notifiers.Notifier, alert *storage.Alert) error {
	if alert.GetState() == storage.ViolationState_ACTIVE || alert.GetState() == storage.ViolationState_ATTEMPTED {
		alertNotifier, ok := notifier.(notifiers.AlertNotifier)
		if !ok {
			return nil
		}
		return sendNotification(ctx, alertNotifier, alert)
	}

	alertNotifier, ok := notifier.(notifiers.ResolvableAlertNotifier)
	if !ok {
		return nil
	}
	return sendResolvableNotification(alertNotifier, alert)
}

func sendNotification(ctx context.Context, notifier notifiers.AlertNotifier, alert *storage.Alert) error {
	err := notifier.AlertNotify(ctx, alert)
	if err != nil {
		logFailure(notifier, alert, err)
	}
	return err
}

func sendResolvableNotification(notifier notifiers.ResolvableAlertNotifier, alert *storage.Alert) error {
	// This is a background process so give it all access. If we're here the user already had access to resolve the alert.
	ctx := sac.WithAllAccess(context.Background())

	var err error
	switch alert.GetState() {
	case storage.ViolationState_SNOOZED:
		err = notifier.AckAlert(ctx, alert)
	case storage.ViolationState_RESOLVED:
		err = notifier.ResolveAlert(ctx, alert)
	}
	if err != nil {
		logFailure(notifier, alert, err)
	}
	return err
}

func logFailure(notifier notifiers.Notifier, alert *storage.Alert, err error) {
	protoNotifier := notifier.ProtoNotifier()
	log.Errorf("Unable to send %s notification to %s (%s) for alert %s: %v", alert.GetState().String(), protoNotifier.GetName(), protoNotifier.GetType(), alert.GetId(), err)
}
