// Copyright 2016 The JoeFriday authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package joefriday gets facts.
package joefriday

import (
	"fmt"
	"strconv"
)

type Error struct {
	Type string
	Op   string
	Err  error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %q: %s", e.Type, e.Op, e.Err)
}

// Column returns a right justified string of width w.
func Column(w int, s string) string {
	pad := w - len(s)
	padding := make([]byte, pad)
	for i := 0; i < pad; i++ {
		padding[i] = 0x20
	}
	return fmt.Sprintf("%s%s", string(padding), s)
}

// Int64Column takes an int64 and returns a right justified string of width w.
func Int64Column(w int, v int64) string {
	s := strconv.FormatInt(v, 10)
	return Column(w, s)
}
