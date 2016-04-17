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

// Package flat handles Flatbuffer based processing of the syscalls uptime
// information.  Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes.  A function to deserialize the Flatbuffer serialized
// bytes into a uptime.Uptime struct is provided.  After the first use,
// the flatbuffer builder is reused.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	uptime "github.com/mohae/joefriday/sysinfo/uptime"
)

var builder = fb.NewBuilder(0)
var mu sync.Mutex

// Get returns the current uptime as Flatbuffer serialized bytes.
func Get() (p []byte, err error) {
	var u uptime.Uptime
	err = u.Get()
	if err != nil {
		return nil, err
	}
	return Serialize(&u), nil
}

// Serialize uptime.Uptime using Flatbuffers.
func Serialize(u *uptime.Uptime) []byte {
	mu.Lock()
	defer mu.Unlock()
	// ensure the Builder is in a usable state.
	builder.Reset()
	UptimeStart(builder)
	UptimeAddTimestamp(builder, u.Timestamp)
	UptimeAddUptime(builder, u.Uptime)
	builder.Finish(UptimeEnd(builder))
	return builder.Bytes[builder.Head():]
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as mem.Uptime.
func Deserialize(p []byte) *uptime.Uptime {
	uF := GetRootAsUptime(p, 0)
	var u uptime.Uptime
	u.Timestamp = uF.Timestamp()
	u.Uptime = uF.Uptime()
	return &u
}

type Ticker struct {
	*joe.Ticker
	Data chan []byte
}

// NewTicker returns a new Uptime Ticker that uses Flatbuffers to serialize
// the data.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte)}
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
			p, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- p
		}
	}
}
