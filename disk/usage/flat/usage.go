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

// Package flat handles Flatbuffer based processing of Disk usage.  Instead
// of returning a Go struct, it returns Flatbuffer serialized bytes.  A
// function to deserialize the Flatbuffer serialized bytes into a
// structs.Usage struct is provided.  After the first use, the flatbuffer
// builder is reused.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/disk/structs"
	"github.com/mohae/joefriday/disk/structs/flat"
	"github.com/mohae/joefriday/disk/usage"
)

// Profiler is used to process the disk usage; generating flatbuffers
// serialized bytes.
type Profiler struct {
	*usage.Profiler
	*fb.Builder
}

// Initialized a new disk usage Profiler that utilizes Flatbuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := usage.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current Usage as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	u, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(u), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current Usage as Flatbuffer serialized bytes using the
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

// Serialize serializes the Usage using Flatbuffers.
func (prof *Profiler) Serialize(u *structs.Usage) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	devF := make([]fb.UOffsetT, len(u.Devices))
	names := make([]fb.UOffsetT, len(u.Devices))
	for i := 0; i < len(names); i++ {
		names[i] = prof.Builder.CreateString(u.Devices[i].Name)
	}
	for i := 0; i < len(devF); i++ {
		flat.DeviceStart(prof.Builder)
		flat.DeviceAddMajor(prof.Builder, u.Devices[i].Major)
		flat.DeviceAddMinor(prof.Builder, u.Devices[i].Minor)
		flat.DeviceAddName(prof.Builder, names[i])
		flat.DeviceAddReadsCompleted(prof.Builder, u.Devices[i].ReadsCompleted)
		flat.DeviceAddReadsMerged(prof.Builder, u.Devices[i].ReadsMerged)
		flat.DeviceAddReadSectors(prof.Builder, u.Devices[i].ReadSectors)
		flat.DeviceAddReadingTime(prof.Builder, u.Devices[i].ReadingTime)
		flat.DeviceAddWritesCompleted(prof.Builder, u.Devices[i].WritesCompleted)
		flat.DeviceAddWritesMerged(prof.Builder, u.Devices[i].WritesMerged)
		flat.DeviceAddWrittenSectors(prof.Builder, u.Devices[i].WrittenSectors)
		flat.DeviceAddWritingTime(prof.Builder, u.Devices[i].WritingTime)
		flat.DeviceAddIOInProgress(prof.Builder, u.Devices[i].IOInProgress)
		flat.DeviceAddIOTime(prof.Builder, u.Devices[i].IOTime)
		flat.DeviceAddWeightedIOTime(prof.Builder, u.Devices[i].WeightedIOTime)
		devF[i] = flat.DeviceEnd(prof.Builder)
	}
	flat.UsageStartDevicesVector(prof.Builder, len(devF))
	for i := len(devF) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(devF[i])
	}
	devV := prof.Builder.EndVector(len(devF))
	flat.UsageStart(prof.Builder)
	flat.UsageAddTimestamp(prof.Builder, u.Timestamp)
	flat.UsageAddTimeDelta(prof.Builder, u.TimeDelta)
	flat.UsageAddDevices(prof.Builder, devV)
	prof.Builder.Finish(flat.UsageEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize the Usage using the package global Profiler.
func Serialize(u *structs.Usage) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(u), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as a structs.Usage.
func Deserialize(p []byte) *structs.Usage {
	u := &structs.Usage{}
	devF := &flat.Device{}
	uF := flat.GetRootAsUsage(p, 0)
	u.Timestamp = uF.Timestamp()
	u.TimeDelta = uF.TimeDelta()
	len := uF.DevicesLength()
	u.Devices = make([]structs.Device, len)
	for i := 0; i < len; i++ {
		var dev structs.Device
		if uF.Devices(devF, i) {
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
		u.Devices[i] = dev
	}
	return u
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
