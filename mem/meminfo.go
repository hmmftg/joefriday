// Copyright 2016 The JoeFriday authors.
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
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
)

type Info struct {
	Timestamp    int64
	MemTotal     int
	MemFree      int
	MemAvailable int
	Buffers      int
	Cached       int
	SwapCached   int
	Active       int
	Inactive     int
	SwapTotal    int
	SwapFree     int
}

// Serialize serializes the Info using flatbuffers.
func (i *Info) Serialize() []byte {
	bldr := fb.NewBuilder(0)
	DataStart(bldr)
	DataAddTimestamp(bldr, int64(i.Timestamp))
	DataAddMemTotal(bldr, int64(i.MemTotal))
	DataAddMemFree(bldr, int64(i.MemFree))
	DataAddMemAvailable(bldr, int64(i.MemAvailable))
	DataAddBuffers(bldr, int64(i.Buffers))
	DataAddCached(bldr, int64(i.Cached))
	DataAddSwapCached(bldr, int64(i.SwapCached))
	DataAddActive(bldr, int64(i.Active))
	DataAddInactive(bldr, int64(i.Inactive))
	DataAddSwapTotal(bldr, int64(i.SwapTotal))
	DataAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(DataEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// Deserialize deserializes bytes representing flatbuffers serialized Data
// into *Info.  If the bytes are not from flatbuffers serialization of
// Data, it is a programmer error and a panic will occur.
func Deserialize(p []byte) *Info {
	data := GetRootAsData(p, 0)
	info := &Info{}
	info.Timestamp = data.Timestamp()
	info.MemTotal = int(data.MemTotal())
	info.MemFree = int(data.MemFree())
	info.MemAvailable = int(data.MemAvailable())
	info.Buffers = int(data.Buffers())
	info.Cached = int(data.Cached())
	info.SwapCached = int(data.SwapCached())
	info.Active = int(data.Active())
	info.Inactive = int(data.Inactive())
	info.SwapTotal = int(data.SwapTotal())
	info.SwapFree = int(data.SwapFree())
	return info
}

// GetInfo returns some of the results of /proc/meminfo.
func GetInfo() (*Info, error) {
	var out bytes.Buffer
	var l, i int
	var name string
	var err error
	var v byte
	line := make([]byte, 0, 50)
	t := time.Now().UTC().UnixNano()
	err = meminfo(&out)
	if err != nil {
		return nil, err
	}
	inf := &Info{Timestamp: t}
	var pos int
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		l++
		if l > 8 || l < 15 {
			continue
		}
		line, err = out.ReadBytes(joe.LF)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
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
			if v == 0x20 || v == joe.LF || v == joe.CR {
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
			inf.MemTotal = i
			continue
		}
		if name == "MemFree" {
			inf.MemFree = i
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = i
			continue
		}
		if name == "Buffers" {
			inf.Buffers = i
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = i
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = i
			continue
		}
		if name == "Active" {
			inf.Active = i
			continue
		}
		if name == "Inactive" {
			inf.Inactive = i
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = i
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = i
			continue
		}
	}
	return inf, nil
}

// GetData returns the current meminfo as flatbuffer serialized bytes.
// TODO: Benchmark to see if we should just create the flatbuffers w/o
// doing the intermediate step of to the data structure.
func GetData() ([]byte, error) {
	inf, err := GetInfo()
	if err != nil {
		return nil, err
	}
	return inf.Serialize(), nil
}

func meminfo(buff *bytes.Buffer) error {
	cmd := exec.Command("cat", "/proc/meminfo")
	cmd.Stdout = buff
	return cmd.Run()
}
