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

// Package uptime returns the system's uptime using syscall.
package uptime

import (
	"syscall"
	"time"

	joe "github.com/mohae/joefriday"
)

// Info holds the current uptime and timestamp.
type Info struct {
	Timestamp int64
	Uptime    int64 // sorry for the stutter
}

// Get gets the current uptime Info.
func (i *Info) Get() error {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return err
	}
	i.Timestamp = time.Now().UTC().UnixNano()
	i.Uptime = sysinfo.Uptime
	return nil
}

// Get gets the current uptime Info.
func Get() (i Info, err error) {
	err = i.Get()
	return i, err
}

// Ticker deliivers the uptime Info at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan Info
}

// NewTicker returns a new Ticker continaing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan Info)}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// read until done signal is received
	for {
		select {
		case <-t.Done:
			return
		case <-t.Ticker.C:
			u, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- u
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
