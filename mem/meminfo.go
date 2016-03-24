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
package mem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	flat "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
)

type Info struct {
	Timestamp    int64 `json:"timestamp"`
	MemTotal     int64 `json:"mem_total"`
	MemFree      int64 `json:"mem_free"`
	MemAvailable int64 `json:"mem_available"`
	Buffers      int64 `json:"buffers"`
	Cached       int64 `json:"cached"`
	SwapCached   int64 `json:"swap_cached"`
	Active       int64 `json:"active"`
	Inactive     int64 `json:"inactive"`
	SwapTotal    int64 `json:"swap_total"`
	SwapFree     int64 `json:"swap_free"`
}

// Serialize serializes the Info using flatbuffers.
func (i *Info) SerializeFlat() []byte {
	bldr := flat.NewBuilder(0)
	return i.SerializeFlatBuilder(bldr)
}

func (i *Info) SerializeFlatBuilder(bldr *flat.Builder) []byte {
	InfoFlatStart(bldr)
	InfoFlatAddTimestamp(bldr, int64(i.Timestamp))
	InfoFlatAddMemTotal(bldr, int64(i.MemTotal))
	InfoFlatAddMemFree(bldr, int64(i.MemFree))
	InfoFlatAddMemAvailable(bldr, int64(i.MemAvailable))
	InfoFlatAddBuffers(bldr, int64(i.Buffers))
	InfoFlatAddCached(bldr, int64(i.Cached))
	InfoFlatAddSwapCached(bldr, int64(i.SwapCached))
	InfoFlatAddActive(bldr, int64(i.Active))
	InfoFlatAddInactive(bldr, int64(i.Inactive))
	InfoFlatAddSwapTotal(bldr, int64(i.SwapTotal))
	InfoFlatAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(InfoFlatEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeInfoFlat deserializes bytes serialized with Flatbuffers from
// InfoFlat into *Info.
func DeserializeInfoFlat(p []byte) *Info {
	infoFlat := GetRootAsInfoFlat(p, 0)
	info := &Info{}
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

func (i *Info) String() string {
	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, i.Timestamp).UTC(), i.MemTotal, i.MemFree, i.MemAvailable, i.Active, i.Inactive, i.Cached, i.Buffers, i.SwapTotal, i.SwapCached, i.SwapFree)
}

// GetInfo returns some of the results of /proc/meminfo.
func GetInfo() (*Info, error) {
	var l, i int
	var name string
	var v byte
	t := time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	inf := &Info{Timestamp: t}
	var pos int
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l > 8 && l < 15 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		name = string(val[:])
		val = val[:0]
		// skip all spaces
		for i, v = range line[pos:] {
			if v != 0x20 {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range line[pos:] {
			if v == 0x20 || v == '\r' {
				break
			}
			val = append(val, v)
		}
		// any conversion error results in 0
		i, err = strconv.Atoi(string(val[:]))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", name, err)
		}
		val = val[:0]
		if name == "MemTotal" {
			inf.MemTotal = int64(i)
			continue
		}
		if name == "MemFree" {
			inf.MemFree = int64(i)
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "Buffers" {
			inf.Buffers = int64(i)
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = int64(i)
			continue
		}
		if name == "Active" {
			inf.Active = int64(i)
			continue
		}
		if name == "Inactive" {
			inf.Inactive = int64(i)
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = int64(i)
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = int64(i)
			continue
		}
	}
	return inf, nil
}

// GetInfoFlat returns the current meminfo as flatbuffer serialized bytes.
func GetInfoFlat() ([]byte, error) {
	inf, err := GetInfo()
	if err != nil {
		return nil, err
	}
	return inf.SerializeFlat(), nil
}

// InfoFlatTicker gathers the meminfo on a ticker, whose interval is defined
// by the received duration, and sends the results to the channel.  The
// output is a Flatbuffers serialization of InfoFlat.  Any error encountered
// during processing is sent to the error channel; processing will continue.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
//
// This pre-allocates the builder and everything other than the []byte that
// gets sent to the out channel to reduce allocations, as this is expected
// to be both a frequent and a long-running process.
func InfoFlatTicker(interval time.Duration, outCh chan []byte, done chan struct{}, errCh chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(outCh)
	// predeclare some vars
	var l, i, pos int
	var t int64
	var v byte
	var name string
	// premake some temp slices
	val := make([]byte, 0, 32)
	// just reset the bldr at the end of every ticker
	bldr := flat.NewBuilder(0)
	// Some hoops to jump through to ensure we don't get a ErrBufferFull.
	var bs []byte
	tmp := bytes.NewBuffer(bs)
	buf := bufio.NewReaderSize(tmp, 1536)
	tmp = nil
	// ticker
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// The current timestamp is always in UTC
			t = time.Now().UTC().UnixNano()
			f, err := os.Open("/proc/meminfo")
			if err != nil {
				errCh <- joe.Error{Type: "mem", Op: "open /proc/meminfo", Err: err}
				continue
			}
			buf.Reset(f)
			InfoFlatStart(bldr)
			InfoFlatAddTimestamp(bldr, t)
			for {
				if l == 16 {
					break
				}
				line, err := buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errCh <- joe.Error{Type: "mem", Op: "read command results", Err: err}
					break
				}
				l++
				if l > 8 && l < 15 {
					continue
				}
				// first grab the key name (everything up to the ':')
				for i, v = range line {
					if v == 0x3A {
						pos = i + 1
						break
					}
					val = append(val, v)
				}
				name = string(val[:])
				val = val[:0]
				// skip all spaces
				for i, v = range line[pos:] {
					if v != 0x20 {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range line[pos:] {
					if v == 0x20 || v == '\r' {
						break
					}
					val = append(val, v)
				}
				// any conversion error results in 0
				i, err = strconv.Atoi(string(val[:]))
				if err != nil {
					errCh <- joe.Error{Type: "mem", Op: "convert to int", Err: err}
					continue
				}
				val = val[:0]
				if name == "MemTotal" {
					InfoFlatAddMemTotal(bldr, int64(i))
					continue
				}
				if name == "MemFree" {
					InfoFlatAddMemFree(bldr, int64(i))
					continue
				}
				if name == "MemAvailable" {
					InfoFlatAddMemAvailable(bldr, int64(i))
					continue
				}
				if name == "Buffers" {
					InfoFlatAddBuffers(bldr, int64(i))
					continue
				}
				if name == "Cached" {
					InfoFlatAddMemAvailable(bldr, int64(i))
					continue
				}
				if name == "SwapCached" {
					InfoFlatAddSwapCached(bldr, int64(i))
					continue
				}
				if name == "Active" {
					InfoFlatAddActive(bldr, int64(i))
					continue
				}
				if name == "Inactive" {
					InfoFlatAddInactive(bldr, int64(i))
					continue
				}
				if name == "SwapTotal" {
					InfoFlatAddSwapTotal(bldr, int64(i))
					continue
				}
				if name == "SwapFree" {
					InfoFlatAddSwapFree(bldr, int64(i))
					continue
				}
			}
			f.Close()
			bldr.Finish(InfoFlatEnd(bldr))
			inf := bldr.Bytes[bldr.Head():]
			outCh <- inf
			bldr.Reset()
			l = 0
		}
	}
}

func (i *InfoFlat) String() string {
	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, i.Timestamp()).UTC(), i.MemTotal(), i.MemFree(), i.MemAvailable(), i.Active(), i.Inactive(), i.Cached(), i.Buffers(), i.SwapTotal(), i.SwapCached(), i.SwapFree())
}
