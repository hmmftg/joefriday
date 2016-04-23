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

// Package mem gets and processes mem info: information for the /proc/meminfo
// file.
package mem

import (
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/meminfo"

// Info holds the mem info information.
type Info struct {
	Timestamp         int64 `json:"timestamp"`
	Active            int64 `json:"active"`
	ActiveAnon        int64 `json:"active_anon"`
	ActiveFile        int64 `json:"active_file"`
	AnonHugePages     int64 `json:"anon_huge_pages"`
	AnonPages         int64 `json:"anon_pages"`
	Bounce            int64 `json:"bounce"`
	Buffers           int64 `json:"buffers"`
	Cached            int64 `json:"cached"`
	CommitLimit       int64 `json:"commit_limit"`
	CommittedAS       int64 `json:"commited_as"`
	DirectMap4K       int64 `json:"direct_map_4k"`
	DirectMap2M       int64 `json:"direct_map_2m"`
	Dirty             int64 `json:"dirty"`
	HardwareCorrupted int64 `json:"hardware_corrupted"`
	HugePagesFree     int64 `json:"huge_pages_free"`
	HugePagesRsvd     int64 `json:"huge_pages_rsvd"`
	HugePagesSize     int64 `json:"huge_pages_size"`
	HugePagesSurp     int64 `json:"huge_pages_surp"`
	HugePagesTotal    int64 `json:"huge_pages_total"`
	Inactive          int64 `json:"inactive"`
	InactiveAnon      int64 `json:"inactive_anon"`
	InactiveFile      int64 `json:"inactive_file"`
	KernelStack       int64 `json:"kernel_stack"`
	Mapped            int64 `json:"mapped"`
	MemAvailable      int64 `json:"mem_available"`
	MemFree           int64 `json:"mem_free"`
	MemTotal          int64 `json:"mem_total"`
	Mlocked           int64 `json:"mlocked"`
	NFSUnstable       int64 `json:"nfs_unstable"`
	PageTables        int64 `json:"page_tables"`
	Shmem             int64 `json:"shmem"`
	Slab              int64 `json:"slab"`
	SReclaimable      int64 `json:"s_reclaimable"`
	SUnreclaim        int64 `json:"s_unreclaim"`
	SwapCached        int64 `json:"swap_cached"`
	SwapFree          int64 `json:"swap_free"`
	SwapTotal         int64 `json:"swap_total"`
	Unevictable       int64 `json:"unevictable"`
	VmallocChunk      int64 `json:"vmalloc_chunk"`
	VmallocTotal      int64 `json:"vmalloc_total"`
	VmallocUsed       int64 `json:"vmalloc_used"`
	Writeback         int64 `json:"writeback"`
	WritebackTmp      int64 `json:"writeback_tmp"`
}

// Profiler is used to process the /proc/meminfo file.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.New(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get returns the current meminfo.
func (prof *Profiler) Get() (inf *Info, err error) {
	var (
		i, pos, nameLen int
		v               byte
		n               uint64
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	inf = &Info{}
	inf.Timestamp = time.Now().UTC().UnixNano()
	for {
		prof.Val = prof.Val[:0]
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return inf, &joe.ReadError{Err: err}
		}
		// first grab the key name (everything up to the ':')
		for i, v = range prof.Line {
			if v == ':' {
				pos = i + 1
				break
			}
			prof.Val = append(prof.Val, v)
		}
		nameLen = len(prof.Val)

		// skip all spaces
		for i, v = range prof.Line[pos:] {
			if v != ' ' {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range prof.Line[pos:] {
			if v == ' ' || v == '\n' {
				break
			}
			prof.Val = append(prof.Val, v)
		}
		// any conversion error results in 0
		n, err = helpers.ParseUint(prof.Val[nameLen:])
		if err != nil {
			return inf, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
		}

		v = prof.Val[0]
		// evaluate the key
		if v == 'A' {
			if prof.Val[5] == 'e' {
				if nameLen == 6 {
					inf.Active = int64(n)
					continue
				}
				if prof.Val[7] == 'a' {
					inf.ActiveAnon = int64(n)
					continue
				}
				inf.ActiveFile = int64(n)
				continue
			}
			if nameLen == 9 {
				inf.AnonPages = int64(n)
				continue
			}
			inf.AnonHugePages = int64(n)
			continue
		}
		if v == 'C' {
			if nameLen == 6 {
				inf.Cached = int64(n)
				continue
			}
			if nameLen == 11 {
				inf.CommitLimit = int64(n)
				continue
			}
			inf.CommittedAS = int64(n)
			continue
		}
		if v == 'D' {
			if nameLen == 5 {
				inf.Dirty = int64(n)
				continue
			}
			if prof.Val[10] == 'k' {
				inf.DirectMap4K = int64(n)
				continue
			}
			inf.DirectMap2M = int64(n)
			continue
		}
		if v == 'H' {
			if nameLen == 14 {
				if prof.Val[10] == 'F' {
					inf.HugePagesFree = int64(n)
					continue
				}
				if prof.Val[10] == 'R' {
					inf.HugePagesRsvd = int64(n)
					continue
				}
				inf.HugePagesSurp = int64(n)
			}
			if prof.Val[1] == 'a' {
				inf.HardwareCorrupted = int64(n)
				continue
			}
			if prof.Val[9] == 'i' {
				inf.HugePagesSize = int64(n)
				continue
			}
			inf.HugePagesTotal = int64(n)
			continue
		}
		if v == 'I' {
			if nameLen == 8 {
				inf.Inactive = int64(n)
				continue
			}
			if prof.Val[9] == 'a' {
				inf.InactiveAnon = int64(n)
				continue
			}
			inf.InactiveFile = int64(n)
		}
		if v == 'M' {
			v = prof.Val[3]
			if nameLen < 8 {
				if v == 'p' {
					inf.Mapped = int64(n)
					continue
				}
				if v == 'F' {
					inf.MemFree = int64(n)
					continue
				}
				inf.Mlocked = int64(n)
				continue
			}
			if v == 'A' {
				inf.MemAvailable = int64(n)
				continue
			}
			inf.MemTotal = int64(n)
			continue
		}
		if v == 'S' {
			v = prof.Val[1]
			if v == 'w' {
				if prof.Val[4] == 'C' {
					inf.SwapCached = int64(n)
					continue
				}
				if prof.Val[4] == 'F' {
					inf.SwapFree = int64(n)
					continue
				}
				inf.SwapTotal = int64(n)
				continue
			}
			if v == 'h' {
				inf.Shmem = int64(n)
				continue
			}
			if v == 'l' {
				inf.Slab = int64(n)
				continue
			}
			if v == 'R' {
				inf.SReclaimable = int64(n)
				continue
			}
			inf.SUnreclaim = int64(n)
			continue
		}
		if v == 'V' {
			if prof.Val[8] == 'C' {
				inf.VmallocChunk = int64(n)
				continue
			}
			if prof.Val[8] == 'T' {
				inf.VmallocTotal = int64(n)
				continue
			}
			inf.VmallocUsed = int64(n)
			continue
		}
		if v == 'W' {
			if nameLen == 9 {
				inf.Writeback = int64(n)
				continue
			}
			inf.WritebackTmp = int64(n)
			continue
		}
		if v == 'B' {
			if nameLen == 6 {
				inf.Bounce = int64(n)
				continue
			}
			inf.Buffers = int64(n)
			continue
		}
		if v == 'K' {
			inf.KernelStack = int64(n)
			continue
		}
		if v == 'N' {
			inf.NFSUnstable = int64(n)
			continue
		}
		if v == 'P' {
			inf.PageTables = int64(n)
		}
		inf.Unevictable = int64(n)
	}
	return inf, nil
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current meminfo using the package's global Profiler.
func Get() (inf *Info, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan Info
	*Profiler
}

// NewTicker returns a new Ticker continaing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan Info), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// predeclare some vars
	var (
		i, pos, nameLen int
		v               byte
		n               uint64
		err             error
		inf             Info
	)
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			err = t.Profiler.Reset()
			if err != nil {
				t.Errs <- err
				continue
			}
			inf.Timestamp = time.Now().UTC().UnixNano()
			for {
				t.Val = t.Val[:0]
				t.Line, err = t.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errs <- &joe.ReadError{Err: err}
				}
				// first grab the key name (everything up to the ':')
				for i, v = range t.Line {
					if v == ':' {
						pos = i + 1
						break
					}
					t.Val = append(t.Val, v)
				}
				nameLen = len(t.Val)

				// skip all spaces
				for i, v = range t.Line[pos:] {
					if v != ' ' {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range t.Line[pos:] {
					if v == ' ' || v == '\n' {
						break
					}
					t.Val = append(t.Val, v)
				}
				// any conversion error results in 0
				n, err = helpers.ParseUint(t.Val[nameLen:])
				if err != nil {
					t.Errs <- &joe.ParseError{Info: string(t.Val[:nameLen]), Err: err}
				}

				v = t.Val[0]
				// evaluate the key
				if v == 'A' {
					if t.Val[5] == 'e' {
						if nameLen == 6 {
							inf.Active = int64(n)
							continue
						}
						if t.Val[7] == 'a' {
							inf.ActiveAnon = int64(n)
							continue
						}
						inf.ActiveFile = int64(n)
						continue
					}
					if nameLen == 9 {
						inf.AnonPages = int64(n)
						continue
					}
					inf.AnonHugePages = int64(n)
					continue
				}
				if v == 'C' {
					if nameLen == 6 {
						inf.Cached = int64(n)
						continue
					}
					if nameLen == 11 {
						inf.CommitLimit = int64(n)
						continue
					}
					inf.CommittedAS = int64(n)
					continue
				}
				if v == 'D' {
					if nameLen == 5 {
						inf.Dirty = int64(n)
						continue
					}
					if t.Val[10] == 'k' {
						inf.DirectMap4K = int64(n)
						continue
					}
					inf.DirectMap2M = int64(n)
					continue
				}
				if v == 'H' {
					if nameLen == 14 {
						if t.Val[10] == 'F' {
							inf.HugePagesFree = int64(n)
							continue
						}
						if t.Val[10] == 'R' {
							inf.HugePagesRsvd = int64(n)
							continue
						}
						inf.HugePagesSurp = int64(n)
					}
					if t.Val[1] == 'a' {
						inf.HardwareCorrupted = int64(n)
						continue
					}
					if t.Val[9] == 'i' {
						inf.HugePagesSize = int64(n)
						continue
					}
					inf.HugePagesTotal = int64(n)
					continue
				}
				if v == 'I' {
					if nameLen == 8 {
						inf.Inactive = int64(n)
						continue
					}
					if t.Val[9] == 'a' {
						inf.InactiveAnon = int64(n)
						continue
					}
					inf.InactiveFile = int64(n)
				}
				if v == 'M' {
					v = t.Val[3]
					if nameLen < 8 {
						if v == 'p' {
							inf.Mapped = int64(n)
							continue
						}
						if v == 'F' {
							inf.MemFree = int64(n)
							continue
						}
						inf.Mlocked = int64(n)
						continue
					}
					if v == 'A' {
						inf.MemAvailable = int64(n)
						continue
					}
					inf.MemTotal = int64(n)
					continue
				}
				if v == 'S' {
					v = t.Val[1]
					if v == 'w' {
						if t.Val[4] == 'C' {
							inf.SwapCached = int64(n)
							continue
						}
						if t.Val[4] == 'F' {
							inf.SwapFree = int64(n)
							continue
						}
						inf.SwapTotal = int64(n)
						continue
					}
					if v == 'h' {
						inf.Shmem = int64(n)
						continue
					}
					if v == 'l' {
						inf.Slab = int64(n)
						continue
					}
					if v == 'R' {
						inf.SReclaimable = int64(n)
						continue
					}
					inf.SUnreclaim = int64(n)
					continue
				}
				if v == 'V' {
					if t.Val[8] == 'C' {
						inf.VmallocChunk = int64(n)
						continue
					}
					if t.Val[8] == 'T' {
						inf.VmallocTotal = int64(n)
						continue
					}
					inf.VmallocUsed = int64(n)
					continue
				}
				if v == 'W' {
					if nameLen == 9 {
						inf.Writeback = int64(n)
						continue
					}
					inf.WritebackTmp = int64(n)
					continue
				}
				if v == 'B' {
					if nameLen == 6 {
						inf.Bounce = int64(n)
						continue
					}
					inf.Buffers = int64(n)
					continue
				}
				if v == 'K' {
					inf.KernelStack = int64(n)
					continue
				}
				if v == 'N' {
					inf.NFSUnstable = int64(n)
					continue
				}
				if v == 'P' {
					inf.PageTables = int64(n)
				}
				inf.Unevictable = int64(n)
			}
			t.Data <- inf
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
