// Copyright 2016 Joel Scoble and The JoeFriday authors.
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

// Package Uptime processes uptime information from the /proc/uptime file.
package uptime

import (
	"io"
	"strconv"
	"sync"

	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/uptime"

// Profiler processes the uptime information.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func New() (prof *Profiler, err error) {
	proc, err := joe.New(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get populates Uptime with /proc/uptime information.
func (prof *Profiler) Get() (u Uptime, err error) {
	err = prof.Reset()
	if err != nil {
		return u, err
	}
	var i int
	var v byte
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return u, joe.Error{Type: "platform", Op: "read /proc/version", Err: err}
		}
		// space delimits the two values
		for i, v = range prof.Line {
			if v == 0x20 {
				break
			}
		}
		u.Total, err = strconv.ParseFloat(string(prof.Line[:i]), 64)
		if err != nil {
			return u, err
		}
		u.Idle, err = strconv.ParseFloat(string(prof.Line[i+1:len(prof.Line)-1]), 64)
		if err != nil {
			return u, err
		}

	}
	return u, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get gets the uptime information using the package's global Profiler, which
// is lazily instantiated.
func Get() (u Uptime, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return u, err
		}
	}
	return std.Get()
}

// Uptime holds uptime information
type Uptime struct {
	Total float64
	Idle  float64
}
