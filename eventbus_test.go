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
	`sync`
	`testing`
)

func TestEventBus_Publish(t *testing.T) {
	type fields struct {
		topicNodes map[string]*node
	}
	type args struct {
		topic string
		data  any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test eventbus publish",
			fields: fields{
				topicNodes: map[string]*node{
					"eventbusgo": newNode(),
				},
			},
			args: args{
				topic: "eventbusgo",
				data:  "say hello",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := &EventBus{
				topicNodes: tt.fields.topicNodes,
				lock:       sync.RWMutex{},
			}

			ch := make(chan Event)
			bus.Subscribe(tt.args.topic, ch)

			if err := bus.Publish(tt.args.topic, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			select {
			case task := <-ch:
				if tt.args.topic != task.topic {
					t.Errorf("onEvent() got = %v, want = %v", task.topic, tt.args.topic)
				}
				if tt.args.data != task.data {
					t.Errorf("onEvent() got = %v, want = %v", task.data, tt.args.data)
				}
			}
		})
	}

	funcTests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test eventbus publish by topic-publisher",
			fields: fields{
				topicNodes: map[string]*node{
					"topic-publisher": newNode(),
				},
			},
			args: args{
				topic: "topic-publisher",
				data:  "say hello -> topic-publisher",
			},
			wantErr: false,
		},
	}

	for _, tt := range funcTests {
		t.Run(tt.name, func(t *testing.T) {
			bus := &EventBus{
				topicNodes: tt.fields.topicNodes,
				lock:       sync.RWMutex{},
			}

			ch := make(chan Event)
			bus.Subscribe(tt.args.topic, ch)

			// use PublishFunc -> publisher func
			publisher := bus.PublishFunc(tt.args.topic) // TopicPublisher
			if err := publisher(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			select {
			case task := <-ch:
				if tt.args.topic != task.topic {
					t.Errorf("onEvent() got = %v, want = %v", task.topic, tt.args.topic)
				}
				if tt.args.data != task.data {
					t.Errorf("onEvent() got = %v, want = %v", task.data, tt.args.data)
				}
			}
		})
	}
}
