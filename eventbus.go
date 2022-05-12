/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package eventbusgo

import (
	`fmt`
	"sync"
)

type EventBus struct {
	topicNodes map[string]*node
	lock       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		topicNodes: make(map[string]*node),
		lock:       sync.RWMutex{},
	}
}

func (bus *EventBus) Length(topic string) (int, error) {
	bus.lock.Lock()
	if tn, ok := bus.topicNodes[topic]; ok {
		bus.lock.Unlock()
		tn.lock.RLock()
		defer tn.lock.RUnlock()
		return tn.Length(), nil
	} else {
		defer bus.lock.Unlock()
		return 0, fmt.Errorf("eventbus.Length: topic:%v not exist", topic)
	}
}

func (bus *EventBus) Publish(topic string, data any) error {
	bus.lock.RLock()

	if tn, ok := bus.topicNodes[topic]; ok {
		bus.lock.RUnlock()
		tn.lock.Lock()
		defer tn.lock.Unlock()

		makeGroup := make(Group, 0)
		newGroupx := append(makeGroup, tn.subscribers...)

		event := NewEvent(topic, data)
		go func(data Event, group Group) { // TODO Goroutine pool?
			for _, sub := range group {
				sub.onEvent(data)
			}
		}(event, newGroupx)

		return nil
	} else {
		defer bus.lock.RUnlock()
		return fmt.Errorf("eventbus.Publish: topic:%v not exist", topic)
	}
}

func (bus *EventBus) Subscribe(topic string, ch channel) {
	bus.lock.Lock()
	sub := NewSubscribe(ch)
	if tn, ok := bus.topicNodes[topic]; ok {
		bus.lock.Unlock()
		tn.lock.Lock()
		defer tn.lock.Unlock()
		tn.subscribers = append(tn.subscribers, sub)
	} else {
		defer bus.lock.Unlock()
		firstNode := newNode()
		bus.topicNodes[topic] = firstNode
		firstNode.subscribers = append(firstNode.subscribers, sub)
	}
}

func (bus *EventBus) UnSubscribe(topic string, sub Subscriber) {
	bus.lock.Lock()
	if tn, ok := bus.topicNodes[topic]; ok && tn.Usable() {
		bus.lock.Unlock()
		bus.topicNodes[topic].remove(sub) // bus.topicNodes[topic] -> instead of tn
	} else {
		defer bus.lock.Unlock()
		return
	}
}

type TopicPublisher func(data any) error

func (bus *EventBus) PublishFunc(topic string) TopicPublisher {
	return func(data any) error {
		return bus.Publish(topic, data)
	}
}
