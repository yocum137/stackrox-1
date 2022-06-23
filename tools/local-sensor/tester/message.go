package tester

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
)

type filter func(sensor *central.MsgFromSensor) bool

func filterKind(kind string) filter {
	return func(msg *central.MsgFromSensor) bool {
		if msg.GetEvent() != nil {
			return msg.GetEvent().GetTiming().GetResource() == kind
		}
		return false
	}
}

func filterDeploymentName(name string) filter {
	return func(msg *central.MsgFromSensor) bool {
		if msg.GetEvent().GetDeployment() != nil {
			return msg.GetEvent().GetDeployment().Name == name
		}
		return false
	}
}

func getLastMessage(output []*central.MsgFromSensor, filters ...filter) *central.MsgFromSensor {
	var result []*central.MsgFromSensor
	for _, msg := range output {
		matched := true
		for _, f := range filters {
			if !f(msg) {
				matched = false
			}
		}
		if matched {
			result = append(result, msg)
		}
	}
	if len(result) > 0 {
		return result[len(result)-1]
	}
	return nil
}

type deploymentChecker struct {
	event *central.SensorEvent
}

var (
	checkerMap = map[string]func(event *central.SensorEvent) interface{} {
		"Deployment": func(event *central.SensorEvent) interface{} { return &deploymentChecker{event} },
	}
)

func (c *deploymentChecker) CheckPermissionLevel(value string) bool {
	v, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("incorrect non numeric value value: %s", value)
	}
	if c.event.GetDeployment() == nil {
		log.Fatal("event has no deployment property")
	}
	if c.event.GetDeployment().ServiceAccountPermissionLevel == storage.PermissionLevel(v) {
		return true
	} else {
		fmt.Printf(" expected: %s, actual: %s\n", storage.PermissionLevel(v), c.event.GetDeployment().ServiceAccountPermissionLevel)
		return false
	}
}

func CheckFields(event *central.SensorEvent, kind string, assertionFn string, expected string) (bool, error) {
	if checkerFactory, ok := checkerMap[kind]; !ok {
		return false, errors.Errorf("no checker configured for kind: %s", kind)
	} else {
		checker := checkerFactory(event)
		genericChecker := reflect.ValueOf(checker)
		method := genericChecker.MethodByName(assertionFn)
		result := method.Call([]reflect.Value{reflect.ValueOf(expected)})
		return result[0].Bool(), nil
	}
}
