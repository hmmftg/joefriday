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
	"log"
	"os"
	"strconv"
	"time"

	"github.com/EricLagergren/joefriday/mem/meminfo"
	"github.com/SermoDigital/helpers"

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
	meminfo.InfoFlatStart(bldr)
	meminfo.InfoFlatAddTimestamp(bldr, int64(i.Timestamp))
	meminfo.InfoFlatAddMemTotal(bldr, int64(i.MemTotal))
	meminfo.InfoFlatAddMemFree(bldr, int64(i.MemFree))
	meminfo.InfoFlatAddMemAvailable(bldr, int64(i.MemAvailable))
	meminfo.InfoFlatAddBuffers(bldr, int64(i.Buffers))
	meminfo.InfoFlatAddCached(bldr, int64(i.Cached))
	meminfo.InfoFlatAddSwapCached(bldr, int64(i.SwapCached))
	meminfo.InfoFlatAddActive(bldr, int64(i.Active))
	meminfo.InfoFlatAddInactive(bldr, int64(i.Inactive))
	meminfo.InfoFlatAddSwapTotal(bldr, int64(i.SwapTotal))
	meminfo.InfoFlatAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(meminfo.InfoFlatEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeInfoFlat deserializes bytes serialized with Flatbuffers from
// InfoFlat into *Info.
func DeserializeInfoFlat(p []byte) *Info {
	infoFlat := meminfo.GetRootAsInfoFlat(p, 0)
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

func init() {
	var err error
	proc, err = os.Open("/proc/meminfo")
	if err != nil {
		log.Fatalln(err)
	}
	buf = bufio.NewReader(proc)
}

var proc *os.File
var buf *bufio.Reader
var val = make([]byte, 0, 32)

// GetInfo returns some of the results of /proc/meminfo.
func GetInfo() (inf Info, err error) {
	_, err = proc.Seek(0, os.SEEK_SET)
	if err != nil {
		return inf, err
	}
	buf.Reset(proc)
	var (
		i       int
		v       byte
		pos     int
		nameLen int
	)
	for l := 0; l < 16; l++ {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return inf, fmt.Errorf("error reading output bytes: %s", err)
		}
		if l > 8 && l < 14 {
			continue
		}

		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == ':' {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		nameLen = len(val)

		// skip all spaces
		for i, v = range line[pos:] {
			if v != ' ' {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range line[pos:] {
			if v == ' ' || v == '\r' {
				break
			}
			val = append(val, v)
		}
		// any conversion error results in 0
		n, err := helpers.ParseUint(val[nameLen:])
		if err != nil {
			return inf, fmt.Errorf("%s: %s", val[:nameLen], err)
		}

		v = val[0]

		// Forgive me.
		if v == 'M' {
			v = val[3]
			if v == 'T' {
				inf.MemTotal = int64(n)
			} else if v == 'F' {
				inf.MemFree = int64(n)
			} else {
				inf.MemAvailable = int64(n)
			}
		} else if v == 'S' {
			v = val[4]
			if v == 'C' {
				inf.SwapCached = int64(n)
			} else if v == 'T' {
				inf.SwapTotal = int64(n)
			} else if v == 'F' {
				inf.SwapFree = int64(n)
			}
		} else if v == 'B' {
			inf.Buffers = int64(n)
		} else if v == 'I' {
			inf.Inactive = int64(n)
		} else if v == 'C' {
			inf.Cached = int64(n)
		} else if v == 'A' {
			inf.Active = int64(n)
		}
		val = val[:0]
	}
	inf.Timestamp = time.Now().UTC().UnixNano()
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
			meminfo.InfoFlatStart(bldr)
			meminfo.InfoFlatAddTimestamp(bldr, t)
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
					meminfo.InfoFlatAddMemTotal(bldr, int64(i))
					continue
				}
				if name == "MemFree" {
					meminfo.InfoFlatAddMemFree(bldr, int64(i))
					continue
				}
				if name == "MemAvailable" {
					meminfo.InfoFlatAddMemAvailable(bldr, int64(i))
					continue
				}
				if name == "Buffers" {
					meminfo.InfoFlatAddBuffers(bldr, int64(i))
					continue
				}
				if name == "Cached" {
					meminfo.InfoFlatAddMemAvailable(bldr, int64(i))
					continue
				}
				if name == "SwapCached" {
					meminfo.InfoFlatAddSwapCached(bldr, int64(i))
					continue
				}
				if name == "Active" {
					meminfo.InfoFlatAddActive(bldr, int64(i))
					continue
				}
				if name == "Inactive" {
					meminfo.InfoFlatAddInactive(bldr, int64(i))
					continue
				}
				if name == "SwapTotal" {
					meminfo.InfoFlatAddSwapTotal(bldr, int64(i))
					continue
				}
				if name == "SwapFree" {
					meminfo.InfoFlatAddSwapFree(bldr, int64(i))
					continue
				}
			}
			f.Close()
			bldr.Finish(meminfo.InfoFlatEnd(bldr))
			inf := bldr.Bytes[bldr.Head():]
			outCh <- inf
			bldr.Reset()
			l = 0
		}
	}
}

// func (i *InfoFlat) String() string {
// 	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, i.Timestamp()).UTC(), i.MemTotal(), i.MemFree(), i.MemAvailable(), i.Active(), i.Inactive(), i.Cached(), i.Buffers(), i.SwapTotal(), i.SwapCached(), i.SwapFree())
// }
