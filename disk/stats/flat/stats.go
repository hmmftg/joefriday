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

// Package flat handles Flatbuffer based processing of disk stats.  Instead
// of returning a Go struct, it returns Flatbuffer serialized bytes.  A
// function to deserialize the Flatbuffer serialized bytes into a
// structs.Stats struct is provided.  After the first use, the flatbuffer
// builder is reused.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/disk/stats"
	"github.com/mohae/joefriday/disk/structs"
	"github.com/mohae/joefriday/disk/structs/flat"
)

// Profiler is used to process the /proc/stat file, as stats, using
// Flatbuffers.
type Profiler struct {
	*stats.Profiler
	*fb.Builder
}

// Initialized a new stats Profiler that utilizes Flatbuffers.
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
func (prof *Profiler) Serialize(stts *structs.Stats) []byte {
	// ensure the Builder is in a usable state.
	std.Builder.Reset()
	devF := make([]fb.UOffsetT, len(stts.Devices))
	names := make([]fb.UOffsetT, len(stts.Devices))
	for i := 0; i < len(names); i++ {
		names[i] = prof.Builder.CreateString(stts.Devices[i].Name)
	}
	for i := 0; i < len(devF); i++ {
		flat.DeviceStart(prof.Builder)
		flat.DeviceAddMajor(prof.Builder, stts.Devices[i].Major)
		flat.DeviceAddMinor(prof.Builder, stts.Devices[i].Minor)
		flat.DeviceAddName(prof.Builder, names[i])
		flat.DeviceAddReadsCompleted(prof.Builder, stts.Devices[i].ReadsCompleted)
		flat.DeviceAddReadsMerged(prof.Builder, stts.Devices[i].ReadsMerged)
		flat.DeviceAddReadSectors(prof.Builder, stts.Devices[i].ReadSectors)
		flat.DeviceAddReadingTime(prof.Builder, stts.Devices[i].ReadingTime)
		flat.DeviceAddWritesCompleted(prof.Builder, stts.Devices[i].WritesCompleted)
		flat.DeviceAddWritesMerged(prof.Builder, stts.Devices[i].WritesMerged)
		flat.DeviceAddWrittenSectors(prof.Builder, stts.Devices[i].WrittenSectors)
		flat.DeviceAddWritingTime(prof.Builder, stts.Devices[i].WritingTime)
		flat.DeviceAddIOInProgress(prof.Builder, stts.Devices[i].IOInProgress)
		flat.DeviceAddIOTime(prof.Builder, stts.Devices[i].IOTime)
		flat.DeviceAddWeightedIOTime(prof.Builder, stts.Devices[i].WeightedIOTime)
		devF[i] = flat.DeviceEnd(prof.Builder)
	}
	flat.StatsStartDevicesVector(prof.Builder, len(devF))
	for i := len(devF) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(devF[i])
	}
	devV := prof.Builder.EndVector(len(devF))
	flat.StatsStart(prof.Builder)
	flat.StatsAddTimestamp(prof.Builder, stts.Timestamp)
	flat.StatsAddDevices(prof.Builder, devV)
	prof.Builder.Finish(flat.StatsEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize the Stats using the package global Profiler.
func Serialize(stts *structs.Stats) (p []byte, err error) {
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
// as a structs.Stats.
func Deserialize(p []byte) *structs.Stats {
	stts := &structs.Stats{}
	devF := &flat.Device{}
	statsFlat := flat.GetRootAsStats(p, 0)
	stts.Timestamp = statsFlat.Timestamp()
	len := statsFlat.DevicesLength()
	stts.Devices = make([]structs.Device, len)
	for i := 0; i < len; i++ {
		var dev structs.Device
		if statsFlat.Devices(devF, i) {
			dev.Major = devF.Major()
			dev.Minor = devF.Minor()
			dev.Name = string(devF.Name())
			dev.ReadsCompleted = devF.ReadsCompleted()
			dev.ReadsMerged = devF.ReadsMerged()
			dev.ReadSectors = devF.ReadSectors()
			dev.ReadingTime = devF.ReadingTime()
			dev.WritesCompleted = devF.WritesCompleted()
			dev.WritesMerged = devF.WritesMerged()
			dev.WrittenSectors = devF.WrittenSectors()
			dev.WritingTime = devF.WritingTime()
			dev.IOInProgress = devF.IOInProgress()
			dev.IOTime = devF.IOTime()
			dev.WeightedIOTime = devF.WeightedIOTime()
		}
		stts.Devices[i] = dev
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
