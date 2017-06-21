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

// Package flat handles Flatbuffer based processing of a platform's loadavg
// information: /proc/loadavg. Instead of returning a Go struct, it returns
// Flatbuffer serialized bytes. A function to deserialize the Flatbuffer
// serialized bytes into a loadavg.Info struct is provided. After the first
// use, the flatbuffer builder is reused.
//
// Note: the package name is loadavg and not the final element of the import
// path (flat). 
package loadavg

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	l "github.com/mohae/joefriday/platform/loadavg"
	"github.com/mohae/joefriday/platform/loadavg/flat/structs"
)

// Profiler is used to process the loadavg information, /proc/loadavg, using
// Flatbuffers.
type Profiler struct {
	*l.Profiler
	*fb.Builder
}

// Initializes and returns an loadavg information profiler that utilizes
// FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := l.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current loadavg information as Flatbuffer serialized
// bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current loadavg information as Flatbuffer serialized bytes
// using the package's global Profiler.
func Get() (p []byte, err error) {
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

// Serialize serializes loadavg information using Flatbuffers.
func (prof *Profiler) Serialize(inf l.Info) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	structs.InfoStart(prof.Builder)
	structs.InfoAddTimestamp(prof.Builder, inf.Timestamp)
	structs.InfoAddMinute(prof.Builder, inf.Minute)
	structs.InfoAddFive(prof.Builder, inf.Five)
	structs.InfoAddFifteen(prof.Builder, inf.Fifteen)
	structs.InfoAddRunning(prof.Builder, inf.Running)
	structs.InfoAddTotal(prof.Builder, inf.Total)
	structs.InfoAddPID(prof.Builder, inf.PID)
	prof.Builder.Finish(structs.InfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes loadavg information using Flatbuffers with the
// package's global Profiler.
func Serialize(inf l.Info) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as loadavg.Info.
func Deserialize(p []byte) l.Info {
	flatInf := structs.GetRootAsInfo(p, 0)
	var inf l.Info
	inf.Timestamp = flatInf.Timestamp()
	inf.Minute = flatInf.Minute()
	inf.Five = flatInf.Five()
	inf.Fifteen = flatInf.Fifteen()
	inf.Running = flatInf.Running()
	inf.Total = flatInf.Total()
	inf.PID = flatInf.PID()
	return inf
}

// Ticker delivers the system's loadavg information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker continaing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			p, err := t.Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- p
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
