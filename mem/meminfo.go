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
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/SermoDigital/helpers"
	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/mem/flat"
)

const procMemInfo = "/proc/meminfo"

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
	bldr := fb.NewBuilder(0)
	return i.SerializeFlatBuilder(bldr)
}

func (i *Info) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, int64(i.Timestamp))
	flat.InfoAddMemTotal(bldr, int64(i.MemTotal))
	flat.InfoAddMemFree(bldr, int64(i.MemFree))
	flat.InfoAddMemAvailable(bldr, int64(i.MemAvailable))
	flat.InfoAddBuffers(bldr, int64(i.Buffers))
	flat.InfoAddCached(bldr, int64(i.Cached))
	flat.InfoAddSwapCached(bldr, int64(i.SwapCached))
	flat.InfoAddActive(bldr, int64(i.Active))
	flat.InfoAddInactive(bldr, int64(i.Inactive))
	flat.InfoAddSwapTotal(bldr, int64(i.SwapTotal))
	flat.InfoAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeInfoFlat deserializes bytes serialized with Flatbuffers from
// InfoFlat into *Info.
func DeserializeInfoFlat(p []byte) *Info {
	infoFlat := flat.GetRootAsInfo(p, 0)
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
	proc, err = os.Open(procMemInfo)
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
// If an error occurs while opening /proc/meminfo, the error will be sent
// to the errs channel and this func will exit.
//
// To stop processing and exit; send a signal on the done channel.  This
// will cause the function to stop the ticker, close the out channel and
// return.
func InfoFlatTicker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		l, i, pos, nameLen int
		v                  byte
		n                  uint64
	)
	// premake some temp slices
	val := make([]byte, 0, 32)
	// just reset the bldr at the end of every ticker
	bldr := fb.NewBuilder(0)
	f, err := os.Open(procMemInfo)
	if err != nil {
		errs <- joe.Error{Type: "cpu", Op: "InfoFlatTicker: open /proc/meminfo", Err: err}
		return
	}
	buf := bufio.NewReaderSize(f, 1536)
	// ticker
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// The current timestamp is always in UTC
			_, err = f.Seek(0, os.SEEK_SET)
			if err != nil {
				errs <- joe.Error{Type: "mem", Op: "seek byte 0: /proc/meminfo", Err: err}
				continue
			}
			bldr.Reset()
			buf.Reset(f)
			flat.InfoStart(bldr)
			flat.InfoAddTimestamp(bldr, time.Now().UTC().UnixNano())
			for l = 0; l < 16; l++ {
				line, err := buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errs <- joe.Error{Type: "mem", Op: "read output bytes", Err: err}
					break
				}
				if l > 7 && l < 14 {
					continue
				}
				// first grab the key name (everything up to the ':')
				for i, v = range line {
					if v == 0x3A {
						val = line[:i]
						break
					}
				}
				nameLen = len(val)
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
				n, err = helpers.ParseUint(val[nameLen:])
				if err != nil {
					errs <- joe.Error{Type: "mem", Op: fmt.Sprintf("convert %s", val[:nameLen]), Err: err}
					continue
				}
				v = val[0]
				if v == 'M' {
					v = val[3]
					if v == 'T' {
						flat.InfoAddMemTotal(bldr, int64(n))
					} else if v == 'F' {
						flat.InfoAddMemFree(bldr, int64(n))
					} else {
						flat.InfoAddMemAvailable(bldr, int64(n))
					}
				} else if v == 'S' {
					v = val[4]
					if v == 'C' {
						flat.InfoAddSwapCached(bldr, int64(n))
					} else if v == 'T' {
						flat.InfoAddSwapTotal(bldr, int64(n))
					} else if v == 'F' {
						flat.InfoAddSwapFree(bldr, int64(n))
					}
				} else if v == 'B' {
					flat.InfoAddBuffers(bldr, int64(n))
				} else if v == 'I' {
					flat.InfoAddInactive(bldr, int64(n))
				} else if v == 'C' {
					flat.InfoAddMemAvailable(bldr, int64(n))
				} else if v == 'A' {
					flat.InfoAddInactive(bldr, int64(n))
				}
			}
			bldr.Finish(flat.InfoEnd(bldr))
			inf := bldr.Bytes[bldr.Head():]
			out <- inf
		}
	}
}

// func (i *InfoFlat) String() string {
// 	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, i.Timestamp()).UTC(), i.MemTotal(), i.MemFree(), i.MemAvailable(), i.Active(), i.Inactive(), i.Cached(), i.Buffers(), i.SwapTotal(), i.SwapCached(), i.SwapFree())
// }
