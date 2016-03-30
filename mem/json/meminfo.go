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

// Package mem gets and processes /proc/meminfo, returning the data in the
// appropriate format.
package json

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"time"

	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/mem"
)

type InfoProfiler struct {
	Info mem.InfoProfiler
}

func NewInfoProfiler() (proc *InfoProfiler, err error) {
	f, err := os.Open(mem.ProcMemInfo)
	if err != nil {
		return nil, err
	}
	return &InfoProfiler{Info: mem.InfoProfiler{Proc: joe.Proc{File: f, Buf: bufio.NewReader(f)}, Val: make([]byte, 0, 32)}}, nil
}

// Get returns some of the results of /proc/meminfo.
func (prof *InfoProfiler) Get() (p []byte, err error) {
	prof.Info.Proc.Reset()
	inf, err := prof.Info.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *InfoProfiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// GetInfo get's the current meminfo.
func GetInfo() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewInfoProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Ticker gathers the meminfo on a ticker, whose interval is defined by the
// received duration, and sends the results to the channel.  The output is
// JSON serialized bytes of mem.Info.  Any error encountered during
// processing is sent to the error channel; processing will continue.
//
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
func (prof *InfoProfiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan mem.Info)
	defer close(outCh)
	go prof.Info.Ticker(interval, outCh, done, errs)
	for {
		select {
		case inf, ok := <-outCh:
			if !ok {
				return
			}
			b, err := prof.Serialize(&inf)
			if err != nil {
				errs <- err
				continue
			}
			out <- b
		}
	}
}

// InfoTicker gathers the meminfo on a ticker, whose interval is defined
// by the received duration, and sends the results to the channel.  The
// output is the JSON serialized bytes of Info.  Any error encountered
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
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	p, err := NewInfoProfiler()
	if err != nil {
		errs <- err
		return
	}
	p.Ticker(interval, out, done, errs)
}

// Serialize mem.Info as JSON
func (prof *InfoProfiler) Serialize(inf *mem.Info) ([]byte, error) {
	return json.Marshal(inf)
}

// UnmarshalInfo unmarshals JSON into *Info.
func UnmarshalInfo(p []byte) (*mem.Info, error) {
	info := &mem.Info{}
	err := json.Unmarshal(p, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}
