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
	"time"

	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/uptime"

// Uptime holds uptime information
type Uptime struct {
	Timestamp int64
	Total     float64
	Idle      float64
}

// Profiler processes the uptime information.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
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
	u.Timestamp = time.Now().UTC().UnixNano()
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return u, &joe.ReadError{Err: err}
		}
		// space delimits the two values
		for i, v = range prof.Line {
			if v == 0x20 {
				break
			}
		}
		u.Total, err = strconv.ParseFloat(string(prof.Line[:i]), 64)
		if err != nil {
			return u, &joe.ParseError{Info: "total", Err: err}
		}
		u.Idle, err = strconv.ParseFloat(string(prof.Line[i+1:len(prof.Line)-1]), 64)
		if err != nil {
			return u, &joe.ParseError{Info: "idle", Err: err}
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
		std, err = NewProfiler()
		if err != nil {
			return u, err
		}
	}
	return std.Get()
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan Uptime
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan Uptime), Profiler: p}
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
