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
func New() (prof *Profiler, err error) {
	p, err := loadavg.New()
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

// Ticker returns the current loadavg as Flatbuffer serialized bytes on a
// ticker.
func (prof *Profiler) Ticker(d time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan loadavg.LoadAvg)
	defer close(out)
	var (
		ok bool
		l  loadavg.LoadAvg
	)
	go prof.Profiler.Ticker(d, outCh, done, errs)
	for {
		select {
		case l, ok = <-outCh:
			if !ok {
				return
			}
			out <- prof.Serialize(l)
		}
	}
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current loadavg information as Flatbuffer serialized bytes
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

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// Serialize serializes loadavg information using Flatbuffers.
func (prof *Profiler) Serialize(l loadavg.LoadAvg) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	LoadAvgStart(prof.Builder)
	LoadAvgAddLastMinute(prof.Builder, l.LastMinute)
	LoadAvgAddLastFive(prof.Builder, l.LastFive)
	LoadAvgAddLastTen(prof.Builder, l.LastTen)
	LoadAvgAddRunningProcesses(prof.Builder, l.RunningProcesses)
	LoadAvgAddTotalProcesses(prof.Builder, l.TotalProcesses)
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
		std, err = New()
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
	l.LastMinute = flatL.LastMinute()
	l.LastFive = flatL.LastFive()
	l.LastTen = flatL.LastTen()
	l.RunningProcesses = flatL.RunningProcesses()
	l.TotalProcesses = flatL.TotalProcesses()
	l.PID = flatL.PID()
	return l
}
