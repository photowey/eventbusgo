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
	`github.com/photowey/eventbusgo/nanoid`
)

type Group []Subscriber

type Subscriber struct {
	id  string
	ech channel // event channel
}

func NewSubscribe(channel chan Event) Subscriber {
	id, _ := nanoid.New()
	return Subscriber{
		id:  id, // Subscriber's identifier
		ech: channel,
	}
}

func (sub *Subscriber) onEvent(event Event) {
	sub.ech <- event
}

func (sub *Subscriber) Await() (event Event) {
	return <-sub.ech
}
