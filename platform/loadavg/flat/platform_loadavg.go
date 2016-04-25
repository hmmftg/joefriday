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
// serialized bytes into a uptime.LoadAvg struct is provided.  After the first
// use, the flatbuffer builder is reused.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/platform/loadavg"
)

// Profiler is used to process the loadavg information, /proc/loadavg, using
// Flatbuffers.
type Profiler struct {
	*loadavg.Profiler
	*fb.Builder
}

// Initializes and returns an loadavg information profiler that utilizes
// FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := loadavg.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current loadavg information as Flatbuffer serialized
// bytes.
func (prof *Profiler) Get() ([]byte, error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k), nil
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
func (prof *Profiler) Serialize(l loadavg.LoadAvg) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	LoadAvgStart(prof.Builder)
	LoadAvgAddMinute(prof.Builder, l.Minute)
	LoadAvgAddFive(prof.Builder, l.Five)
	LoadAvgAddFifteen(prof.Builder, l.Fifteen)
	LoadAvgAddRunning(prof.Builder, l.Running)
	LoadAvgAddTotal(prof.Builder, l.Total)
	LoadAvgAddPID(prof.Builder, l.PID)
	prof.Builder.Finish(LoadAvgEnd(prof.Builder))
	return prof.Builder.Bytes[prof.Builder.Head():]
}

// Serialize serializes uptime information using Flatbuffers with the
// package's global Profiler.
func Serialize(u loadavg.LoadAvg) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(u), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as loadavg.LoadAvg.
func Deserialize(p []byte) loadavg.LoadAvg {
	flatL := GetRootAsLoadAvg(p, 0)
	var l loadavg.LoadAvg
	l.Minute = flatL.Minute()
	l.Five = flatL.Five()
	l.Fifteen = flatL.Fifteen()
	l.Running = flatL.Running()
	l.Total = flatL.Total()
	l.PID = flatL.PID()
	return l
}

// Ticker delivers the system's memory information at intervals.
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
