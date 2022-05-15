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
	"context"
	"testing"
	"time"
)

func TestEventBus_Publish(t *testing.T) {
	type args struct {
		topic string
		data  any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test eventbus publish",
			args: args{
				topic: "eventbusgo",
				data:  "say hello",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewEventBus()

			ch := make(chan Event)
			ch2 := make(chan Event)
			bus.Subscribe(tt.args.topic, ch)
			bus.Subscribe(tt.args.topic, ch2)

			if err := bus.Publish(tt.args.topic, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			go func() {
				select {
				case task := <-ch:
					t.Logf("tests.ch1: receive event: %+v", task)

					if tt.args.topic != task.topic {
						t.Errorf("onEvent() got = %v, want = %v", task.topic, tt.args.topic)
					}
					if tt.args.data != task.data {
						t.Errorf("onEvent() got = %v, want = %v", task.data, tt.args.data)
					}
				}
			}()

			go func() {
				select {
				case task := <-ch2:
					t.Logf("tests.ch2: receive event: %+v", task)
					if tt.args.topic != task.topic {
						t.Errorf("onEvent() got = %v, want = %v", task.topic, tt.args.topic)
					}
					if tt.args.data != task.data {
						t.Errorf("onEvent() got = %v, want = %v", task.data, tt.args.data)
					}
				}
			}()

			time.Sleep(2 * time.Second)
		})
	}

	funcTests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test eventbus publish by topic-publisher",
			args: args{
				topic: "topic-publisher",
				data:  "say hello -> topic-publisher",
			},
			wantErr: false,
		},
	}

	for _, tt := range funcTests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewEventBus()

			ch := make(chan Event)
			bus.Subscribe(tt.args.topic, ch)

			// use PublishFunc -> publisher func
			publisher := bus.PublishFunc(tt.args.topic) // TopicPublisher
			if err := publisher(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			ctx := context.Background()
			publisher2 := bus.PublishFunc(tt.args.topic, ctx) // TopicPublisher
			if err := publisher2(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			select {
			case task := <-ch:

				t.Logf("funcTests:ch1 receive task:%+v", task)

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
