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

// Package flat handles Flatbuffer based processing of memmory information
// using syscall.  Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes.  A function to deserialize the Flatbuffer serialized
// bytes into a mem.Info struct is provided.  After the first use, the
// flatbuffer builder is reused.
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

// Get returns the current loadavg as Flatbuffer serialized bytes.
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
	p := builder.Bytes[builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
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

// Ticker delivers mem.Info as Flatbuffers serialized bytes at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
}

// NewTicker returns a new Ticker continaing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte)}
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
			p, err := Get()
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
