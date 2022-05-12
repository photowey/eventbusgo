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

package nanoid

import (
	"github.com/matoous/go-nanoid/v2"
)

const (
	DefaultSize int = 21
	single          = 1
)

func New(sizes ...int) (string, error) {
	size := determineSize(sizes)

	return gonanoid.New(size)
}

func Generate(alphabet string, sizes ...int) (string, error) {
	size := determineSize(sizes)

	return gonanoid.Generate(alphabet, size)
}

func determineSize(sizes []int) int {
	size := DefaultSize
	switch len(sizes) {
	case single:
		if sizes[0] > 0 {
			size = sizes[0]
		}
	}

	return size
}
