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
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestErrorCheck(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		isRead  bool
		isReset bool
		isParse bool
	}{
		{"standard error", errors.New("test"), false, false, false},
		{"read error", &ReadError{}, true, false, false},
		{"reset error", &ResetError{}, false, true, false},
		{"parse error", &ParseError{}, false, false, true},
	}
	for _, test := range tests {
		b := IsReadError(test.err)
		if b != test.isRead {
			t.Errorf("%s IsReadError: got %t want %t", test.name, b, test.isRead)
		}
		b = IsResetError(test.err)
		if b != test.isReset {
			t.Errorf("%s IsResetError: got %t want %t", test.name, b, test.isReset)
		}
		b = IsParseError(test.err)
		if b != test.isParse {
			t.Errorf("%s IsParseError: got %t want %t", test.name, b, test.isParse)
		}
	}
}

var trailingVals = []struct {
	val      []byte
	expected []byte
}{
	{[]byte{}, []byte{}},
	{[]byte(""), []byte("")},
	{[]byte("\n"), []byte("")},
	{[]byte("1\n"), []byte("1")},
	{[]byte("\t"), []byte("")},
	{[]byte("1\t"), []byte("1")},
	{[]byte("\t\n"), []byte("")},
	{[]byte("1\t\n"), []byte("1")},
	{[]byte("   \t\n"), []byte("")},
	{[]byte("1   \t\n"), []byte("1")},
	{[]byte("      "), []byte("")},
	{[]byte("hello"), []byte("hello")},
	{[]byte("salut   "), []byte("salut")},
	{[]byte("eamus catuli    \n"), []byte("eamus catuli")},
	{[]byte("hola                  "), []byte("hola")},
	{[]byte(" nihao   "), []byte(" nihao")},
	{[]byte("idographic space　  "), []byte("idographic space　")},
	{[]byte("punctuation space         "), []byte("punctuation space ")},
	{[]byte("EM Quad space  "), []byte("EM Quad space ")},
	{[]byte("OGHAM space    "), []byte("OGHAM space ")},
}

func TestTrimTrailingSpaces(t *testing.T) {
	for i, test := range trailingVals {
		tmp := TrimTrailingSpaces(test.val)
		if !bytes.Equal(tmp, test.expected) {
			t.Errorf("%d: got %q; want %q", i, tmp, test.expected)
		}
	}
}

var leadingVals = []struct {
	val      []byte
	expected []byte
}{
	{[]byte{}, []byte{}},
	{[]byte(""), []byte("")},
	{[]byte("\n"), []byte("\n")},
	{[]byte("  1"), []byte("1")},
	{[]byte("\t"), []byte("")},
	{[]byte("  \t1\t"), []byte("1\t")},
	{[]byte("\t "), []byte("")},
	{[]byte("1\t\n"), []byte("1\t\n")},
	{[]byte("   \t  "), []byte("")},
	{[]byte(" \t 1"), []byte("1")},
	{[]byte("      "), []byte("")},
	{[]byte("hello"), []byte("hello")},
	{[]byte("   salut"), []byte("salut")},
	{[]byte("     \teamus catuli"), []byte("eamus catuli")},
	{[]byte("      nihao "), []byte("nihao ")},
	{[]byte("    　ideographic space"), []byte("　ideographic space")},
	{[]byte("       punctuation space"), []byte(" punctuation space")},
	{[]byte(" EM Quad space"), []byte(" EM Quad space")},
	{[]byte(" OGHAM space"), []byte(" OGHAM space")},
}

func TestTrimLeadingSpaces(t *testing.T) {
	for i, test := range leadingVals {
		tmp := TrimLeadingSpaces(test.val)
		if !bytes.Equal(tmp, test.expected) {
			t.Errorf("%d: got %q; want %q", i, tmp, test.expected)
		}
	}
}

var trailingByteVals = [][]byte{
	[]byte("      "),
	[]byte("hello"),
	[]byte("salut   "),
	[]byte("eamus catuli    "),
	[]byte("hola                  "),
	[]byte(" nihao   "),
	[]byte("ideographic space　  "),
	[]byte("punctuation space         "),
	[]byte("EM Quad space  "),
	[]byte("OGHAM space    "),
}

func BenchmarkTrimTrailingSpaces(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for _, v := range trailingByteVals {
			tmp = TrimTrailingSpaces(v)
		}
	}
	_ = tmp
}

func BenchmarkTrimTrailingSpaceBytes(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for _, v := range trailingByteVals {
			tmp = bytes.TrimSpace(v)
		}
	}
	_ = tmp

}

var trailingStringVals = []string{
	"      ",
	"hello",
	"salut   ",
	"eamus catuli    ",
	"hola                  ",
	" nihao   ",
	"idographic space　  ",
	"punctuation space         ",
	"EM Quad space  ",
	"OGHAM space    ",
}

func BenchmarkTrimTrailingSpaceStrings(b *testing.B) {
	var tmp string
	for i := 0; i < b.N; i++ {
		for _, v := range trailingStringVals {
			tmp = strings.TrimSpace(v)
		}
	}
	_ = tmp

}

var leadingByteVals = [][]byte{
	[]byte("      "),
	[]byte("hello"),
	[]byte("   salut"),
	[]byte("    eamus catuli"),
	[]byte("                  hola"),
	[]byte("   nihao "),
	[]byte("  　ideographic space"),
	[]byte("         punctuation space"),
	[]byte("  EM Quad space"),
	[]byte("    OGHAM space"),
}

func BenchmarkTrimLeadingSpaces(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for _, v := range leadingByteVals {
			tmp = TrimLeadingSpaces(v)
		}
	}
	_ = tmp
}

func BenchmarkTrimLeadingSpaceBytes(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for _, v := range leadingByteVals {
			tmp = bytes.TrimSpace(v)
		}
	}
	_ = tmp
}

var leadingStringVals = []string{
	"      ",
	"hello",
	"   salut",
	"    eamus catuli",
	"                  hola",
	"   nihao ",
	"  　ideographic space",
	"punctuation space         ",
	"  EM Quad space",
	"    OGHAM space",
}

func BenchmarkTrimLeadingSpacesBytes(b *testing.B) {
	var tmp string
	for i := 0; i < b.N; i++ {
		for _, v := range leadingStringVals {
			tmp = strings.TrimSpace(v)
		}
	}
	_ = tmp

}

var byteVals = [][]byte{
	[]byte("      "),
	[]byte("hello"),
	[]byte("   salut   "),
	[]byte("    eamus catuli    "),
	[]byte("                  hola                  "),
	[]byte("   nihao   "),
	[]byte("  　ideographic space  "),
	[]byte("         punctuation space         "),
	[]byte("  EM Quad space  "),
	[]byte("    OGHAM space    "),
}

func BenchmarkTrimSpaces(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for _, v := range byteVals {
			tmp = TrimTrailingSpaces(TrimLeadingSpaces(v))
		}
	}
	_ = tmp
}

func BenchmarkTrimSpaceBytes(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		for _, v := range byteVals {
			tmp = bytes.TrimSpace(v)
		}
	}
	_ = tmp
}

var stringVals = []string{
	"      ",
	"hello",
	"   salut   ",
	"    eamus catuli    ",
	"                  hola                  ",
	"   nihao   ",
	"  　ideographic space  ",
	"         punctuation space         ",
	"  EM Quad space  ",
	"    OGHAM space    ",
}

// benchmark with strings (no conversions)
func BenchmarkTrimSpaceString(b *testing.B) {
	var tmp string
	for i := 0; i < b.N; i++ {
		for _, v := range stringVals {
			tmp = strings.TrimSpace(v)
		}
	}
	_ = tmp
}

func TestNewTempFileProc(t *testing.T) {
	data := "abcdefghijklmnopqrstuvwxyz"
	tests := []struct {
		prefix string
		name   string
	}{
		{"", ""},
		{"", "abc"},
		{"abc", ""},
		{"abc", "abc"},
	}
	for i, test := range tests {
		p, err := NewTempFileProc(test.prefix, test.name, []byte(data))
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if test.prefix == "" {
			if p.Dir != os.TempDir() {
				t.Errorf("%d: Dir: got %s; want %s", i, p.Dir, os.TempDir())
			}
		} else {
			dir := filepath.Base(p.Dir)
			if !strings.HasPrefix(dir, test.prefix) {
				t.Errorf("%d: Dir: expected %q to have a prefix of %q; it didn't", i, dir, test.prefix)
			}
		}
		if test.name == "" {
			if len(p.Name) != 12 {
				t.Errorf("%d: Name: expected random name to be 12 chars; got %d chars", i, len(p.Name))
			}
		} else {
			if test.name != p.Name {
				t.Errorf("%d: Name: got %q; want %q", i, p.Name, test.name)
			}
		}
		s, err := p.Buf.ReadString('z')
		if err != nil {
			t.Errorf("%d: unexpected error: %q", i, err)
			goto remove
		}
		if s != data {
			t.Errorf("%d: got %q; want %q", i, s, data)
		}

	remove:
		err = p.Remove()
		if err != nil {
			t.Errorf("%d: unexpected error: %q", i, err)
		}

		if p.Dir == os.TempDir() {
			_, err = os.Stat(p.FullPath())
			if !os.IsNotExist(err) {
				t.Errorf("%d: stat of %q: expected IsNotExist err; got %q", i, p.FullPath(), err)
			}
		} else {
			_, err = os.Stat(p.Dir)
			if !os.IsNotExist(err) {
				t.Errorf("%d: stat of %q: expected IsNotExist err; got %q", i, p.Dir, err)
			}
		}
	}
}
