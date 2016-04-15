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

// Package loadavg returns the system;s loadavg information, for 1, 5, and
// 15 minutes, using syscall.
package loadavg

import (
	"syscall"
	"time"

	joe "github.com/mohae/joefriday"
)

const LoadsScale = 65536

type LoadAvg struct {
	Timestamp int64
	One       float64
	Five      float64
	Fifteen   float64
}

// Get the load average for the last 1, 5, and 15 minutes.
func (l *LoadAvg) Get() error {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return err
	}
	l.Timestamp = time.Now().UTC().UnixNano()
	l.One = float64(sysinfo.Loads[0]) / LoadsScale
	l.Five = float64(sysinfo.Loads[1]) / LoadsScale
	l.Fifteen = float64(sysinfo.Loads[2]) / LoadsScale
	return nil
}

// Get returns a LoadAvg populated with the 1, 5, and 15 minute values.
func Get() (LoadAvg, error) {
	var l LoadAvg
	err := l.Get()
	return l, err
}

type Ticker struct {
	*joe.Ticker
	Data chan LoadAvg
}

// NewTicker returns a new Ticker containing a ticker channel, T,
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan LoadAvg)}
	go t.Run()
	return &t, nil
}

// Stop sends a signal to the done channel, stopping the goroutine.  This
// also closes all channels that the Ticker holds.
func (t *Ticker) Stop() {
	t.Done <- struct{}{}
}

func (t *Ticker) Run() {
	defer t.Close()
	defer close(t.Data)
	// read until done signal is received
	for {
		select {
		case <-t.Done:
			return
		case <-t.Ticker.C:
			s, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- s
		}
	}
}
