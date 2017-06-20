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

// Package cpustats handles Flatbuffer based processing of kernel activity,
// /proc/stat. Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes. A function to deserialize the Flatbuffer serialized bytes
// into a cpustats.Stats struct is provided.  After the first use, the
// flatbuffer builder is re-used.
//
// Note: the package name is cpustats and not the final element of the import
// path (flat). 
package cpustats

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	stats "github.com/mohae/joefriday/cpu/cpustats"
	"github.com/mohae/joefriday/cpu/cpustats/flat/flat"
)

// Profiler is used to process the stats, /proc/stat, as Flatbuffers
// serialized bytes.
type Profiler struct {
	*stats.Profiler
	*fb.Builder
}

// Initialized and returns a new stats Profiler that uses Flatbuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := stats.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current Stats as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	stts, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(stts), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current Stats as Flatbuffer serialized bytes using the
// package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	} else {
		std.Builder.Reset()
	}

	return std.Get()
}

// Serialize serializes the Stats using Flatbuffers.
func (prof *Profiler) Serialize(stts *stats.Stats) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	statsF := make([]fb.UOffsetT, len(stts.CPU))
	ids := make([]fb.UOffsetT, len(stts.CPU))
	for i := 0; i < len(ids); i++ {
		ids[i] = prof.Builder.CreateString(stts.CPU[i].ID)
	}
	for i := 0; i < len(statsF); i++ {
		flat.StatStart(prof.Builder)
		flat.StatAddID(prof.Builder, ids[i])
		flat.StatAddUser(prof.Builder, stts.CPU[i].User)
		flat.StatAddNice(prof.Builder, stts.CPU[i].Nice)
		flat.StatAddSystem(prof.Builder, stts.CPU[i].System)
		flat.StatAddIdle(prof.Builder, stts.CPU[i].Idle)
		flat.StatAddIOWait(prof.Builder, stts.CPU[i].IOWait)
		flat.StatAddIRQ(prof.Builder, stts.CPU[i].IRQ)
		flat.StatAddSoftIRQ(prof.Builder, stts.CPU[i].SoftIRQ)
		flat.StatAddSteal(prof.Builder, stts.CPU[i].Steal)
		flat.StatAddQuest(prof.Builder, stts.CPU[i].Quest)
		flat.StatAddQuestNice(prof.Builder, stts.CPU[i].QuestNice)
		statsF[i] = flat.StatEnd(prof.Builder)
	}
	flat.StatsStartCPUVector(prof.Builder, len(statsF))
	for i := len(statsF) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(statsF[i])
	}
	statsV := prof.Builder.EndVector(len(statsF))
	flat.StatsStart(prof.Builder)
	flat.StatsAddClkTck(prof.Builder, stts.ClkTck)
	flat.StatsAddTimestamp(prof.Builder, stts.Timestamp)
	flat.StatsAddCtxt(prof.Builder, stts.Ctxt)
	flat.StatsAddBTime(prof.Builder, stts.BTime)
	flat.StatsAddProcesses(prof.Builder, stts.Processes)
	flat.StatsAddCPU(prof.Builder, statsV)
	prof.Builder.Finish(flat.StatsEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize the Stats using the package global Profiler.
func Serialize(stts *stats.Stats) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(stts), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as a stats.Stats.
func Deserialize(p []byte) *stats.Stats {
	stts := &stats.Stats{}
	statF := &flat.Stat{}
	statsFlat := flat.GetRootAsStats(p, 0)
	stts.ClkTck = statsFlat.ClkTck()
	stts.Timestamp = statsFlat.Timestamp()
	stts.Ctxt = statsFlat.Ctxt()
	stts.BTime = statsFlat.BTime()
	stts.Processes = statsFlat.Processes()
	len := statsFlat.CPULength()
	stts.CPU = make([]stats.Stat, len)
	for i := 0; i < len; i++ {
		var stat stats.Stat
		if statsFlat.CPU(statF, i) {
			stat.ID = string(statF.ID())
			stat.User = statF.User()
			stat.Nice = statF.Nice()
			stat.System = statF.System()
			stat.Idle = statF.Idle()
			stat.IOWait = statF.IOWait()
			stat.IRQ = statF.IRQ()
			stat.SoftIRQ = statF.SoftIRQ()
			stat.Steal = statF.Steal()
			stat.Quest = statF.Quest()
			stat.QuestNice = statF.QuestNice()
		}
		stts.CPU[i] = stat
	}
	return stts
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			p, err := t.Get()
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
