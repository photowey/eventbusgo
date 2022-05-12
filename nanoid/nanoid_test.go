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

package nanoid_test

import (
	`fmt`
	`regexp`
	`testing`
	`unicode`

	`github.com/photowey/eventbusgo/nanoid`
)

func TestNew(t *testing.T) {
	type args struct {
		sizes []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Test nanoid in golang",
			args: args{
				sizes: []int{32},
			},
			want:    32,
			wantErr: false,
		},
		{
			name: "Test nanoid in golang-default",
			args: args{
				sizes: []int{},
			},
			want:    21,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nanoid.New(tt.args.sizes...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type test func(src string, len int) bool

func TestGenerate(t *testing.T) {
	type args struct {
		alphabet string
		sizes    []int
		test     test
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Test nanoid in golang",
			args: args{
				alphabet: "0123456789",
				sizes:    []int{32},
				test:     isDigit,
			},
			want:    32,
			wantErr: false,
		},
		{
			name: "Test nanoid in golang-default",
			args: args{
				alphabet: "abcdefhgijklmnopqrstuvwxyz",
				sizes:    []int{},
				test:     alphabet,
			},
			want:    21,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nanoid.Generate(tt.args.alphabet, tt.args.sizes...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want || !tt.args.test(got, tt.want) {
				t.Errorf("Generate() got = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func isDigit(src string, len int) bool {
	for _, x := range []rune(src) {
		if !unicode.IsDigit(x) {
			return false
		}
	}

	return true
}

func alphabet(src string, len int) bool {
	var alpha = regexp.MustCompile(fmt.Sprintf("^[a-z]{%d}$", len))

	fmt.Println(alpha.MatchString(src))

	return true
}
