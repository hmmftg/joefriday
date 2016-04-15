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

// Package sysinfo returns system information using.  Instead of using proc
// files, syscalls are made.  There is some overlap between the data that
// this package provides and some of the other JoeFriday packages.
package sysinfo

import (
	"syscall"
	"time"

	joe "github.com/mohae/joefriday"
)

type Uptime struct {
	Timestamp int64
	Uptime    int64 // sorry for the stutter
}

func (u *Uptime) Get() error {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return err
	}
	u.Timestamp = time.Now().UTC().UnixNano()
	u.Uptime = sysinfo.Uptime
	return nil
}

func Get() (u Uptime, err error) {
	err = u.Get()
	return u, err
}

type Ticker struct {
	*joe.Ticker
	Data chan Uptime
}

// NewTicker returns a new Ticker containing a ticker channel, T,
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan Uptime)}
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
			u, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- u
		}
	}
}
