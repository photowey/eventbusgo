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
	"sync"
)

type node struct {
	subscribers Group
	lock        sync.RWMutex
}

func newNode() *node {
	return &node{
		subscribers: make(Group, 0),
		lock:        sync.RWMutex{},
	}
}

func (node *node) Length() int {
	return len(node.subscribers)
}

func (node *node) Usable() bool {
	return node.Length() > 0
}

func (node *node) remove(sub Subscriber) {
	length := len(node.subscribers)
	if length == 0 {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	idx := node.index(sub)
	if idx < 0 {
		return
	}

	copy(node.subscribers[:idx], node.subscribers[idx+1:])
	node.subscribers[length-1] = Subscriber{}
	node.subscribers = node.subscribers[:length-1]
}

func (node *node) index(src Subscriber) int {
	for idx, sub := range node.subscribers {
		if sub.id == src.id { // not: sub == src
			return idx
		}
	}

	return -1
}
