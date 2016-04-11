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

// Package flat handles Flatbuffer based processing of a platform's uptime
// information: /proc/uptime.  Instead of returning a Go struct, it returns
// Flatbuffer serialized bytes.  A function to deserialize the Flatbuffer
// serialized bytes into a uptime.Uptime struct is provided.  After the first
// use, the flatbuffer builder is reused.
package flat

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/platform/uptime"
)

// Profiler is used to process the uptime information, /proc/uptime, using
// Flatbuffers.
type Profiler struct {
	Prof    *uptime.Profiler
	Builder *fb.Builder
}

// Initializes and returns an uptime information profiler that utilizes
// FlatBuffers.
func New() (prof *Profiler, err error) {
	p, err := uptime.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current uptime information as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	k, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current uptime information as Flatbuffer serialized bytes
// using the package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Serialize serializes uptime information using Flatbuffers.
func (prof *Profiler) Serialize(u uptime.Uptime) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	UptimeStart(prof.Builder)
	UptimeAddTotal(prof.Builder, u.Total)
	UptimeAddIdle(prof.Builder, u.Idle)
	prof.Builder.Finish(UptimeEnd(prof.Builder))
	return prof.Builder.Bytes[prof.Builder.Head():]
}

// Serialize serializes uptime information using Flatbuffers with the
// package's global Profiler.
func Serialize(u uptime.Uptime) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(u), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as uptime.Uptime.
func Deserialize(p []byte) uptime.Uptime {
	flatU := GetRootAsUptime(p, 0)
	var u uptime.Uptime
	u.Total = flatU.Total()
	u.Idle = flatU.Idle()
	return u
}
