package processor

import (
	"github.com/stackrox/rox/central/notifiers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
)

// NotifierSet is a set that coordinates present policies and notifiers.
type NotifierSet interface {
	HasNotifiers() bool
	HasEnabledAuditNotifiers() bool

	ForEach(f func(notifiers.Notifier, AlertSet))

	UpsertNotifier(notifier notifiers.Notifier)
	RemoveNotifier(id string)
}

// NewNotifierSet returns a new instance of a NotifierSet
func NewNotifierSet() NotifierSet {
	return &notifierSetImpl{
		notifiers: make(map[string]notifiers.Notifier),
		failures:  make(map[string]AlertSet),
	}
}

// Implementation of the notifier set.
//////////////////////////////////////

type notifierSetImpl struct {
	lock sync.RWMutex

	notifiers map[string]notifiers.Notifier
	failures  map[string]AlertSet
}

// HasNotifiers returns if there are any notifiers in the set.
func (p *notifierSetImpl) HasNotifiers() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return len(p.notifiers) > 0
}

// HasEnabledAuditNotifiers returns if there are any enabled notifiers in the set.
func (p *notifierSetImpl) HasEnabledAuditNotifiers() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for _, n := range p.notifiers {
		auditN, ok := n.(notifiers.AuditNotifier)
		if ok && auditN.AuditLoggingEnabled() {
			return true
		}
	}
	return false
}

// ForEachesFailures performs a function on each notifier.
func (p *notifierSetImpl) ForEach(f func(notifiers.Notifier, AlertSet)) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for notifierID, notifier := range p.notifiers {
		f(notifier, p.failures[notifierID])
	}
}

// UpsertNotifier adds or updates a notifier in the set.
func (p *notifierSetImpl) UpsertNotifier(notifier notifiers.Notifier) {
	p.lock.Lock()
	defer p.lock.Unlock()

	notifierID := notifier.ProtoNotifier().GetId()
	if _, exists := p.failures[notifierID]; !exists {
		p.failures[notifierID] = NewAlertSet()
	}
	if knownNotifier := p.notifiers[notifierID]; knownNotifier != nil && knownNotifier != notifier {
		if err := knownNotifier.Close(); err != nil {
			log.Error("failed to close notifier instance", logging.Err(err))
		}
	}
	p.notifiers[notifierID] = notifier
}

// RemoveNotifier removes a notifier from the set.
func (p *notifierSetImpl) RemoveNotifier(id string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if notifier := p.notifiers[id]; notifier != nil {
		if err := notifier.Close(); err != nil {
			log.Error("failed to close notifier instance", logging.Err(err))
		}
	}

	delete(p.notifiers, id)
	delete(p.failures, id)
}
