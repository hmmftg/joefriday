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

// Package mem returns memory information using syscalls, instead of proc
// files,  Only basic memory information is provided by this package.  for
// more detailed memory information, use the joefriday/mem packages.
package mem

import (
	"syscall"
	"time"

	joe "github.com/mohae/joefriday"
)

type Ticker struct {
	*joe.Ticker
	Data chan Info
}

// NewTicker returns a new Ticker containing a ticker channel, T,
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan Info)}
	go t.Run()
	return &t, nil
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

// tick runs on each tick.  When a done signal is received, it returns;
// closing the Ticker's channels.

// Info holds information about system memory.
type Info struct {
	Timestamp int64
	TotalRAM  uint64
	FreeRAM   uint64
	SharedRAM uint64
	BufferRAM uint64
	TotalSwap uint64
	FreeSwap  uint64
}

// Get gets the meminfo information.
func (m *Info) Get() error {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return err
	}
	m.Timestamp = time.Now().UTC().UnixNano()
	m.TotalRAM = sysinfo.Totalram
	m.FreeRAM = sysinfo.Freeram
	m.SharedRAM = sysinfo.Sharedram
	m.BufferRAM = sysinfo.Bufferram
	m.TotalSwap = sysinfo.Totalswap
	m.FreeSwap = sysinfo.Freeswap
	return nil
}

// Get gets the meminfo information.
func Get() (m Info, err error) {
	err = m.Get()
	return m, err
}
