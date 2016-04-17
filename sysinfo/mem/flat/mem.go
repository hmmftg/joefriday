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

// Package flat handles Flatbuffer based processing of /proc/meminfo.
// Instead of returning a Go struct, it returns Flatbuffer serialized bytes.
// A function to deserialize the Flatbuffer serialized bytes into a
// mem.Info struct is provided.  After the first use, the flatbuffer
// builder is reused.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/sysinfo/mem"
)

var builder = fb.NewBuilder(0)
var mu sync.Mutex

// Get returns the current meminfo as Flatbuffer serialized bytes using the
// package's global Profiler.
func Get() (p []byte, err error) {
	var inf mem.Info
	err = inf.Get()
	if err != nil {
		return nil, err
	}
	return Serialize(&inf), nil
}

// Serialize mem.Info using Flatbuffers.
func Serialize(inf *mem.Info) []byte {
	mu.Lock()
	defer mu.Unlock()
	// ensure the Builder is in a usable state.
	builder.Reset()
	InfoStart(builder)
	InfoAddTimestamp(builder, inf.Timestamp)
	InfoAddTotalRAM(builder, inf.TotalRAM)
	InfoAddFreeRAM(builder, inf.FreeRAM)
	InfoAddSharedRAM(builder, inf.SharedRAM)
	InfoAddBufferRAM(builder, inf.BufferRAM)
	InfoAddTotalSwap(builder, inf.TotalSwap)
	InfoAddFreeSwap(builder, inf.FreeSwap)
	builder.Finish(InfoEnd(builder))
	return builder.Bytes[builder.Head():]
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as mem.Info.
func Deserialize(p []byte) *mem.Info {
	infoFlat := GetRootAsInfo(p, 0)
	info := &mem.Info{}
	info.Timestamp = infoFlat.Timestamp()
	info.TotalRAM = infoFlat.TotalRAM()
	info.FreeRAM = infoFlat.FreeRAM()
	info.SharedRAM = infoFlat.SharedRAM()
	info.BufferRAM = infoFlat.BufferRAM()
	info.TotalSwap = infoFlat.TotalSwap()
	info.FreeSwap = infoFlat.FreeSwap()
	return info
}

type Ticker struct {
	*joe.Ticker
	Data chan []byte
}

// NewTicker returns a new Ticker containing a ticker channel, T,
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
			s, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- s
		}
	}
}
