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
package joefriday

import (
	"bytes"
	"strings"
	"testing"
)

var vals = []struct {
	val      []byte
	expected []byte
}{
	{[]byte{}, []byte{}},
	{[]byte(""), []byte("")},
	{[]byte("      "), []byte("")},
	{[]byte("hello"), []byte("hello")},
	{[]byte("salut   "), []byte("salut")},
	{[]byte("eamus catuli    "), []byte("eamus catuli")},
	{[]byte("hola                  "), []byte("hola")},
	{[]byte(" nihao   "), []byte(" nihao")},
	{[]byte("idographic space　  "), []byte("idographic space　")},
	{[]byte("punctuation space         "), []byte("punctuation space ")},
	{[]byte("EM Quad space  "), []byte("EM Quad space ")},
	{[]byte("OGHAM space    "), []byte("OGHAM space ")},
}

var stringVals = []struct {
	val      string
	expected string
}{
	{"      ", ""},
	{"hello", "hello"},
	{"salut   ", "salut"},
	{"eamus catuli    ", "eamus catuli"},
	{"hola                  ", "hola"},
	{" nihao   ", " nihao"},
	{"idographic space　  ", "idographic space　"},
	{"punctuation space         ", "punctuation space "},
	{"EM Quad space  ", "EM Quad space "},
	{"OGHAM space    ", "OGHAM space "},
}

func TestTrimTrailingSpaces(t *testing.T) {
	for i, test := range vals {
		tmp := TrimTrailingSpaces(test.val)
		if !bytes.Equal(tmp, test.expected) {
			t.Errorf("%d: got %q; want %q", i, tmp, test.expected)
		}
	}
}

func BenchmarkTrimTrailingSpaces(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for j := 2; j < len(vals); j++ {
			tmp = TrimTrailingSpaces(vals[j].val)
		}
	}
	_ = tmp
}

// benchmark with strings (no conversions)
func BenchmarkTrimSpaces(b *testing.B) {
	var tmp string
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(stringVals); j++ {
			tmp = strings.TrimSpace(stringVals[j].val)
		}
	}
	_ = tmp
}

// benchmark with bytes input and string returned val
func BenchmarkTrimSpacesByteInput(b *testing.B) {
	var tmp string
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(stringVals); j++ {
			tmp = strings.TrimSpace(string(vals[j].val))
		}
	}
	_ = tmp
}

// benchmark with bytes; everything converted
func BenchmarkTrimSpacesBytes(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(stringVals); j++ {
			tmp = []byte(strings.TrimSpace(string(vals[j].val)))
		}
	}
	_ = tmp
}
