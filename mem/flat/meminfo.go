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

// Package flat gets and processes /proc/meminfo using Flatbuffers.
package flat

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/SermoDigital/helpers"
	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/mem"
)

var std *InfoProfiler

// InfoProfilerFlat wraps InfoProfiler and provides a builder; enabling reuse.
type InfoProfiler struct {
	Info mem.InfoProfiler
	bldr *fb.Builder
}

func NewInfoProfiler() (proc *InfoProfiler, err error) {
	f, err := os.Open(mem.ProcMemInfo)
	if err != nil {
		return nil, err
	}
	return &InfoProfiler{Info: mem.InfoProfiler{Proc: joe.Proc{File: f, Buf: bufio.NewReader(f)}, Val: make([]byte, 0, 32)}, bldr: fb.NewBuilder(0)}, nil
}

func (p *InfoProfiler) reset() error {
	p.Info.Lock()
	p.bldr.Reset()
	p.Info.Unlock()
	return p.Info.Reset()
}

// Get returns the current meminfo as flatbuffer serialized bytes.
func (p *InfoProfiler) Get() ([]byte, error) {
	p.reset()
	inf, err := p.Info.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(inf), nil
}

// GetInfo get's the current meminfo.
func GetInfo() (p []byte, err error) {
	if std == nil {
		std, err = NewInfoProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Ticker gathers the meminfo on a ticker, whose interval is defined by
// the received duration, and sends the results to the channel.  The output
// is Flatbuffer serialized bytes of Info.  Any error encountered during
// processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
func (p *InfoProfiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		l, i, pos, nameLen int
		v                  byte
		n                  uint64
		err                error
	)
	// Lock now because the for loop unlocks to simplify unlock logic when
	// a continue occurs (instead of the tick completing.)
	p.Info.Lock()
	// ticker
Tick:
	for {
		p.Info.Unlock()
		select {
		case <-done:
			return
		case <-ticker.C:
			err = p.reset()
			p.Info.Lock()
			if err != nil {
				errs <- joe.Error{Type: "mem", Op: "seek byte 0: /proc/meminfo", Err: err}
				continue
			}
			InfoStart(p.bldr)
			InfoAddTimestamp(p.bldr, time.Now().UTC().UnixNano())
			for l = 0; l < 16; l++ {
				p.Info.Line, err = p.Info.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					// An error results in sending error message and stop processing of this tick.
					errs <- joe.Error{Type: "mem", Op: "read output bytes", Err: err}
					continue Tick
				}
				if l > 7 && l < 14 {
					continue
				}
				// first grab the key name (everything up to the ':')
				for i, v = range p.Info.Line {
					if v == 0x3A {
						p.Info.Val = p.Info.Line[:i]
						break
					}
				}
				nameLen = len(p.Info.Val)
				// skip all spaces
				for i, v = range p.Info.Line[pos:] {
					if v != 0x20 {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range p.Info.Line[pos:] {
					if v == 0x20 || v == '\n' {
						break
					}
					p.Info.Val = append(p.Info.Val, v)
				}
				// any conversion error results in 0
				n, err = helpers.ParseUint(p.Info.Val[nameLen:])
				if err != nil {
					errs <- joe.Error{Type: "mem", Op: fmt.Sprintf("convert %s", p.Info.Val[:nameLen]), Err: err}
					continue
				}
				v = p.Info.Val[0]
				if v == 'M' {
					v = p.Info.Val[3]
					if v == 'T' {
						InfoAddMemTotal(p.bldr, int64(n))
					} else if v == 'F' {
						InfoAddMemFree(p.bldr, int64(n))
					} else {
						InfoAddMemAvailable(p.bldr, int64(n))
					}
				} else if v == 'S' {
					v = p.Info.Val[4]
					if v == 'C' {
						InfoAddSwapCached(p.bldr, int64(n))
					} else if v == 'T' {
						InfoAddSwapTotal(p.bldr, int64(n))
					} else if v == 'F' {
						InfoAddSwapFree(p.bldr, int64(n))
					}
				} else if v == 'B' {
					InfoAddBuffers(p.bldr, int64(n))
				} else if v == 'I' {
					InfoAddInactive(p.bldr, int64(n))
				} else if v == 'C' {
					InfoAddMemAvailable(p.bldr, int64(n))
				} else if v == 'A' {
					InfoAddInactive(p.bldr, int64(n))
				}
			}
			p.bldr.Finish(InfoEnd(p.bldr))
			inf := p.bldr.Bytes[p.bldr.Head():]
			out <- inf
		}
	}
}

// TODO: should InfoTickerFlat use std or have a local proc?
// InfoTickerFlat gathers the meminfo on a ticker, whose interval is defined
// by the received duration, and sends the results to the channel.  The
// output is Flatbuffer serialized bytes of Info.  Any error encountered
// during processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
//
// This func uses a local InfoProfiler.  If an error occurs during the
// creation of the InfoProfiler, it will be sent to errs and exit.
func InfoTicker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	p, err := NewInfoProfiler()
	if err != nil {
		errs <- err
		return
	}
	p.Ticker(interval, out, done, errs)
}

func (prof *InfoProfiler) Serialize(inf *mem.Info) []byte {
	prof.Info.Lock()
	defer prof.Info.Unlock()
	InfoStart(prof.bldr)
	InfoAddTimestamp(prof.bldr, int64(inf.Timestamp))
	InfoAddMemTotal(prof.bldr, int64(inf.MemTotal))
	InfoAddMemFree(prof.bldr, int64(inf.MemFree))
	InfoAddMemAvailable(prof.bldr, int64(inf.MemAvailable))
	InfoAddBuffers(prof.bldr, int64(inf.Buffers))
	InfoAddCached(prof.bldr, int64(inf.Cached))
	InfoAddSwapCached(prof.bldr, int64(inf.SwapCached))
	InfoAddActive(prof.bldr, int64(inf.Active))
	InfoAddInactive(prof.bldr, int64(inf.Inactive))
	InfoAddSwapTotal(prof.bldr, int64(inf.SwapTotal))
	InfoAddSwapFree(prof.bldr, int64(inf.SwapFree))
	prof.bldr.Finish(InfoEnd(prof.bldr))
	return prof.bldr.Bytes[prof.bldr.Head():]
}

// DeserializeInfo deserializes bytes serialized with Flatbuffers from
// InfoFlat into *Info.
func DeserializeInfo(p []byte) *mem.Info {
	infoFlat := GetRootAsInfo(p, 0)
	info := &mem.Info{}
	info.Timestamp = infoFlat.Timestamp()
	info.MemTotal = infoFlat.MemTotal()
	info.MemFree = infoFlat.MemFree()
	info.MemAvailable = infoFlat.MemAvailable()
	info.Buffers = infoFlat.Buffers()
	info.Cached = infoFlat.Cached()
	info.SwapCached = infoFlat.SwapCached()
	info.Active = infoFlat.Active()
	info.Inactive = infoFlat.Inactive()
	info.SwapTotal = infoFlat.SwapTotal()
	info.SwapFree = infoFlat.SwapFree()
	return info
}
