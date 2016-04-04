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

// Package flat handles Flatbuffer based processing of CPU stats.  Instead
// of returning a Go struct, it returns Flatbuffer serialized bytes.
// A function to deserialize the Flatbuffer serialized bytes into a
// facts.Facts struct is provided.  After the first use, the flatbuffer
// builder is reused.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/cpu/stats"
)

// Profiler is used to process the /proc/stat file, as stats, using
// Flatbuffers.
type Profiler struct {
	*stats.Profiler
	*fb.Builder
}

// Initialized a new stats Profiler that utilizes Flatbuffers.
func New() (prof *Profiler, err error) {
	p, err := stats.New()
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
		std, err = New()
		if err != nil {
			return nil, err
		}
	} else {
		std.Builder.Reset()
	}

	return std.Get()
}

// Ticker processes CPU utilization information on a ticker.  The generated
// utilization data is sent to the outCh.  Any errors encountered are sent
// to the errCh.  Processing ends when either a done signal is received or
// the done channel is closed.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan *stats.Stats)
	defer close(out)
	go prof.Profiler.Ticker(interval, outCh, done, errs)
	for {
		select {
		case s, ok := <-outCh:
			if !ok {
				return
			}
			out <- prof.Serialize(s)
		}
	}
}

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// Serialize serializes the Stats using Flatbuffers.
func (prof *Profiler) Serialize(stts *stats.Stats) []byte {
	// ensure the Builder is in a usable state.
	std.Builder.Reset()
	statsF := make([]fb.UOffsetT, len(stts.CPU))
	ids := make([]fb.UOffsetT, len(stts.CPU))
	for i := 0; i < len(ids); i++ {
		ids[i] = prof.Builder.CreateString(stts.CPU[i].ID)
	}
	for i := 0; i < len(statsF); i++ {
		StatStart(prof.Builder)
		StatAddID(prof.Builder, ids[i])
		StatAddUser(prof.Builder, stts.CPU[i].User)
		StatAddNice(prof.Builder, stts.CPU[i].Nice)
		StatAddSystem(prof.Builder, stts.CPU[i].System)
		StatAddIdle(prof.Builder, stts.CPU[i].Idle)
		StatAddIOWait(prof.Builder, stts.CPU[i].IOWait)
		StatAddIRQ(prof.Builder, stts.CPU[i].IRQ)
		StatAddSoftIRQ(prof.Builder, stts.CPU[i].SoftIRQ)
		StatAddSteal(prof.Builder, stts.CPU[i].Steal)
		StatAddQuest(prof.Builder, stts.CPU[i].Quest)
		StatAddQuestNice(prof.Builder, stts.CPU[i].QuestNice)
		statsF[i] = StatEnd(prof.Builder)
	}
	StatsStartCPUVector(prof.Builder, len(statsF))
	for i := len(statsF) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(statsF[i])
	}
	statsV := prof.Builder.EndVector(len(statsF))
	StatsStart(prof.Builder)
	StatsAddClkTck(prof.Builder, stts.ClkTck)
	StatsAddTimestamp(prof.Builder, stts.Timestamp)
	StatsAddCtxt(prof.Builder, stts.Ctxt)
	StatsAddBTime(prof.Builder, stts.BTime)
	StatsAddProcesses(prof.Builder, stts.Processes)
	StatsAddCPU(prof.Builder, statsV)
	prof.Builder.Finish(StatsEnd(prof.Builder))
	return prof.Builder.Bytes[prof.Builder.Head():]
}

// Serialize the Stats using the package global Profiler.
func Serialize(stts *stats.Stats) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
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
	statF := &Stat{}
	statsFlat := GetRootAsStats(p, 0)
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
