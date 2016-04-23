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

// Package Release processes the OS Release information, /etc/os-release.
package release

import (
	"io"
	"sync"

	joe "github.com/mohae/joefriday"
)

const etcFile = "/etc/os-release"

// Release holds information about the OS release.
type Release struct {
	ID           string `json:"id"`
	IDLike       string `json:"id_like"`
	PrettyName   string `json:"pretty_name"`
	Version      string `json:"version"`
	VersionID    string `json:"version_id"`
	HomeURL      string `json:"home_url"`
	BugReportURL string `json:"bug_report_url"`
}

// Profiler processes the OS release information, /etc/os-release.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.New(etcFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get populates Release with /etc/os-release information.
func (prof *Profiler) Get() (r *Release, err error) {
	var (
		i, keyLen int
		v         byte
		release   Release
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		// The key is everything up to '='; 0x3D
		for i, v = range prof.Line {
			if v == 0x3D {
				prof.Val = prof.Line[:i]
				keyLen = len(prof.Val)
				// see if the value has quotes; if it does, elide them
				if prof.Line[i+1] == 0x22 {
					prof.Val = append(prof.Val, prof.Line[i+2:len(prof.Line)-2]...)
				} else {
					prof.Val = append(prof.Val, prof.Line[i+1:len(prof.Line)-1]...)
				}
				break
			}
		}
		v = prof.Val[0]
		if v == 'I' {
			if prof.Val[2] == '_' {
				release.IDLike = string(prof.Val[keyLen:])
				continue
			}
			release.ID = string(prof.Val[keyLen:])
			continue
		}
		if v == 'V' {
			if prof.Val[7] == '_' {
				release.VersionID = string(prof.Val[keyLen:])
				continue
			}
			release.Version = string(prof.Val[keyLen:])
			continue
		}
		if v == 'P' {
			release.PrettyName = string(prof.Val[keyLen:])
			continue
		}
		if v == 'H' {
			release.HomeURL = string(prof.Val[keyLen:])
			continue
		}
		if v == 'B' {
			release.BugReportURL = string(prof.Val[keyLen:])
			continue
		}
	}
	return &release, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get gets the OS release information using the package's global Profiler,
// which is lazily instantiated.
func Get() (r *Release, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}
