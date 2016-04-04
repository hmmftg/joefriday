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

package flat

// Package flat handles Flatbuffer based processing of CPU utilization
// information.  Instead of returning a Go struct, it returns Flatbuffer
// serialized bytes.  A function to deserialize the Flatbuffer serialized
// bytes into a utilization.Utilization struct is provided.  After the first
// use, the flatbuffer builder is reused.
import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/cpu/utilization"
)

// Profiler is used to process the cpu utilization information using
// Flatbuffers.
type Profiler struct {
	*utilization.Profiler
	*fb.Builder
}

// Initializes and returns a cpu utilization profiler that uses FlatBuffers.
func New() (prof *Profiler, err error) {
	p, err := utilization.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

func (prof *Profiler) reset() {
	prof.Builder.Reset()
}

// Get returns the current cpu utilization as Flatbuffer serialized bytes.
// Utilization calculations requires two pieces of data.  This func gets a
// snapshot of /proc/stat, sleeps for a second, takes another snapshot and
// calcualtes the utilization from the two snapshots.  If ongoing utilitzation
// information is desired, the Ticker should be used; it's better suited for
// ongoing utilization information.
//
// TODO: should this be changed so that this calculates utilization since
// last time the stats were obtained.  If there aren't pre-existing stats
// it would get current utilization (which may be a separate method (or
// should be?))
func (prof *Profiler) Get() (p []byte, err error) {
	prof.reset()
	u, err := prof.Profiler.Get()
	return prof.Serialize(u), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current cpu utilization as Flatbuffer serialized bytes
// using the package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Ticker processes CPU utilization information on a ticker.  The generated
// utilization data is sent to the out channel.  Any errors encountered are
// sent to the errs.  Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan *utilization.Utilization)
	defer close(outCh)
	go prof.Profiler.Ticker(interval, outCh, done, errs)
	for {
		select {
		case u, ok := <-outCh:
			if !ok {
				return
			}
			out <- prof.Serialize(u)
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

// Serialize cpu utilization using Flatbuffers.
func (prof *Profiler) Serialize(u *utilization.Utilization) []byte {
	utils := make([]fb.UOffsetT, len(u.CPU))
	ids := make([]fb.UOffsetT, len(u.CPU))
	for i := 0; i < len(ids); i++ {
		ids[i] = prof.Builder.CreateString(u.CPU[i].ID)
	}
	for i := 0; i < len(utils); i++ {
		UtilStart(prof.Builder)
		UtilAddID(prof.Builder, ids[i])
		UtilAddUsage(prof.Builder, u.CPU[i].Usage)
		UtilAddUser(prof.Builder, u.CPU[i].User)
		UtilAddNice(prof.Builder, u.CPU[i].Nice)
		UtilAddSystem(prof.Builder, u.CPU[i].System)
		UtilAddIdle(prof.Builder, u.CPU[i].Idle)
		UtilAddIOWait(prof.Builder, u.CPU[i].IOWait)
		utils[i] = UtilEnd(prof.Builder)
	}
	UtilizationStartCPUVector(prof.Builder, len(utils))
	for i := len(utils) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(utils[i])
	}
	utilsV := prof.Builder.EndVector(len(utils))
	UtilizationStart(prof.Builder)
	UtilizationAddTimestamp(prof.Builder, u.Timestamp)
	UtilizationAddBTimeDelta(prof.Builder, u.BTimeDelta)
	UtilizationAddCtxtDelta(prof.Builder, u.CtxtDelta)
	UtilizationAddProcesses(prof.Builder, u.Processes)
	UtilizationAddCPU(prof.Builder, utilsV)
	prof.Builder.Finish(UtilizationEnd(prof.Builder))
	return prof.Builder.Bytes[prof.Builder.Head():]
}

// Serialize the Utilization using the package global Profiler.
func Serialize(u *utilization.Utilization) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	} else {
		std.reset()
	}
	return std.Serialize(u), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as a utilization.Utilization.
func Deserialize(p []byte) *utilization.Utilization {
	u := &utilization.Utilization{}
	uF := &Util{}
	flatUtil := GetRootAsUtilization(p, 0)
	u.Timestamp = flatUtil.Timestamp()
	u.CtxtDelta = flatUtil.CtxtDelta()
	u.BTimeDelta = flatUtil.BTimeDelta()
	u.Processes = flatUtil.Processes()
	len := flatUtil.CPULength()
	u.CPU = make([]utilization.Util, len)
	for i := 0; i < len; i++ {
		var util utilization.Util
		if flatUtil.CPU(uF, i) {
			util.ID = string(uF.ID())
			util.Usage = uF.Usage()
			util.User = uF.User()
			util.Nice = uF.Nice()
			util.System = uF.System()
			util.Idle = uF.Idle()
			util.IOWait = uF.IOWait()
		}
		u.CPU[i] = util
	}
	return u
}
