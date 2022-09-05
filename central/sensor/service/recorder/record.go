package recorder

import (
	"sync"

	"github.com/stackrox/rox/generated/internalapi/central"
)

var (
	recorder EventRecorder
	once     sync.Once
)

type EventRecorder interface {
	Record(id string, event *central.MsgFromSensor)
	GetAllEvents() map[string][]*central.MsgFromSensor
	Clean() error
	SetEnabled(v bool)
}

func Singleton() EventRecorder {
	once.Do(func() {
		impl := &recorderImpl{
			injestionChannel: make(chan eventWrapper, 100),
			eventMutex:       &sync.Mutex{},
			isEnabled:        false,
		}
		go impl.startEventInjestion()
		recorder = impl
	})
	return recorder
}

type recorderImpl struct {
	sensorEvents     map[string][]*central.MsgFromSensor
	injestionChannel chan eventWrapper
	eventMutex       *sync.Mutex
	isEnabled        bool
}

type eventWrapper struct {
	event     *central.MsgFromSensor
	clusterId string
}

func (r *recorderImpl) Record(clusterId string, event *central.MsgFromSensor) {
	if r.isEnabled {
		r.injestionChannel <- eventWrapper{event, clusterId}
	}
}

func (r *recorderImpl) GetAllEvents() map[string][]*central.MsgFromSensor {
	return r.sensorEvents
}

func (r *recorderImpl) Clean() error {
	r.eventMutex.Lock()
	defer r.eventMutex.Unlock()
	r.sensorEvents = map[string][]*central.MsgFromSensor{}
	return nil
}

func (r *recorderImpl) startEventInjestion() {
	for {
		wrapper, more := <-r.injestionChannel
		if !more {
			return
		}
		r.eventMutex.Lock()
		r.sensorEvents[wrapper.clusterId] = append(r.sensorEvents[wrapper.clusterId], wrapper.event)
		r.eventMutex.Unlock()
	}
}

func (r *recorderImpl) SetEnabled(v bool) {
	r.isEnabled = v
}
