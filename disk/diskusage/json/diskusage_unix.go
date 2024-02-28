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

// Package diskusage calculates IO usage the of block devices. Usage is
// calculated by taking the difference between tow snapshots of IO statistics
// of block devices, /proc/diskstats. The time elapsed between the two
// snapshots is stored in the TimeDelta field. Instead of returning a Go
// struct, it returns JSON serialized bytes. A function to deserialize the
// JSON serialized bytes into a struct.DiskUsage struct is provided.
//
// Note: the package name is diskusage and not the final element of the import
// path (json).
package diskusage

import (
	"encoding/json"
	"sync"
	"time"

	joe "github.com/hmmftg/joefriday"
	usage "github.com/hmmftg/joefriday/disk/diskusage"
	"github.com/hmmftg/joefriday/disk/structs"
)

// Profiler is used to process IO usage of the block devices using JSON.
type Profiler struct {
	*usage.Profiler
}

// Returns an initialized Profiler; ready to use. Upon creation, a
// /proc/diskstats snapshot is taken so that any Get() will return valid
// information.
func NewProfiler() (prof *Profiler, err error) {
	p, err := usage.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current IO usage of the block devices as JSON serialized
// bytes. Calculating usage requires two snapshots. This func gets the current
// snapshot of /proc/diskstats and calculates the difference between that and
// the prior snapshot. The current snapshot is stored for use as the prior
// snapshot on the next Get call. If ongoing usage information is desired, the
// Ticker should be used; it's better suited for ongoing usage information.
func (prof *Profiler) Get() (p []byte, err error) {
	u, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(u)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current IO usage of the block devices as JSON serialized
// bytes using the package's global Profiler.  The Profiler is instantiated
// lazily. If the profiler doesn't already exist, the first utilization
// information will not be useful due to minimal time elapsing between the
// initial and second snapshots used for utilization calculations; the results
// of the first call should be discarded.
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

// Serialize IO usage of the block devices using JSON.
func (prof *Profiler) Serialize(u *structs.DiskUsage) ([]byte, error) {
	return json.Marshal(u)
}

// Serialize IO usage of the block devices as JSON using the package's global
// Profiler.
func Serialize(u *structs.DiskUsage) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(u)
}

// Marshal is an alias for Serialize.
func (prof *Profiler) Marshal(u *structs.DiskUsage) ([]byte, error) {
	return prof.Serialize(u)
}

// Marshal is an alias for Serialize using the package global Profiler.
func Marshal(u *structs.DiskUsage) ([]byte, error) {
	return Serialize(u)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// structs.DiskUsage.
func Deserialize(p []byte) (*structs.DiskUsage, error) {
	u := &structs.DiskUsage{}
	err := json.Unmarshal(p, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Unmarshal is an alias for Deserialize.
func Unmarshal(p []byte) (*structs.DiskUsage, error) {
	return Deserialize(p)
}

// Ticker delivers the system's IO usage of the block devices at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
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
