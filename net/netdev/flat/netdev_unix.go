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

// Package netdev handles Flatbuffer based processing of network device
// information; /proc/net/dev. Instead of returning a Go struct, it returns
// Flatbuffer serialized bytes. A function to deserialize the Flatbuffer
// serialized bytes into a structs.Info struct is provided. After the first
// use, the flatbuffer builder is reused.
//
// Note: the package name is netdev and not the final element of the import
// path (flat). 
package netdev

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	dev "github.com/mohae/joefriday/net/netdev"
	"github.com/mohae/joefriday/net/structs"
	"github.com/mohae/joefriday/net/structs/flat"
)

// Profiler is used to process the /proc/net/dev file using Flatbuffers.
type Profiler struct {
	*dev.Profiler
	*fb.Builder
}

// Initializes and returns a net info profiler that utilizes FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := dev.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current network information as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current network information as Flatbuffer serialized bytes
// using the package's global Profiler.
func Get() (p []byte, err error) {
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

// Serialize serializes Info using Flatbuffers.
func (prof *Profiler) Serialize(inf *structs.DevInfo) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	devs := make([]fb.UOffsetT, len(inf.Devices))
	names := make([]fb.UOffsetT, len(inf.Devices))
	for i := 0; i < len(inf.Devices); i++ {
		names[i] = prof.Builder.CreateString(inf.Devices[i].Name)
	}
	for i := 0; i < len(inf.Devices); i++ {
		flat.DeviceStart(prof.Builder)
		flat.DeviceAddName(prof.Builder, names[i])
		flat.DeviceAddRBytes(prof.Builder, inf.Devices[i].RBytes)
		flat.DeviceAddRPackets(prof.Builder, inf.Devices[i].RPackets)
		flat.DeviceAddRErrs(prof.Builder, inf.Devices[i].RErrs)
		flat.DeviceAddRDrop(prof.Builder, inf.Devices[i].RDrop)
		flat.DeviceAddRFIFO(prof.Builder, inf.Devices[i].RFIFO)
		flat.DeviceAddRFrame(prof.Builder, inf.Devices[i].RFrame)
		flat.DeviceAddRCompressed(prof.Builder, inf.Devices[i].RCompressed)
		flat.DeviceAddRMulticast(prof.Builder, inf.Devices[i].RMulticast)
		flat.DeviceAddTBytes(prof.Builder, inf.Devices[i].TBytes)
		flat.DeviceAddTPackets(prof.Builder, inf.Devices[i].TPackets)
		flat.DeviceAddTErrs(prof.Builder, inf.Devices[i].TErrs)
		flat.DeviceAddTDrop(prof.Builder, inf.Devices[i].TDrop)
		flat.DeviceAddTFIFO(prof.Builder, inf.Devices[i].TFIFO)
		flat.DeviceAddTColls(prof.Builder, inf.Devices[i].TColls)
		flat.DeviceAddTCarrier(prof.Builder, inf.Devices[i].TCarrier)
		flat.DeviceAddTCompressed(prof.Builder, inf.Devices[i].TCompressed)
		devs[i] = flat.DeviceEnd(prof.Builder)
	}
	flat.DevInfoStartDevicesVector(prof.Builder, len(devs))
	for i := len(inf.Devices) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(devs[i])
	}
	devsV := prof.Builder.EndVector(len(devs))
	flat.DevInfoStart(prof.Builder)
	flat.DevInfoAddTimestamp(prof.Builder, inf.Timestamp)
	flat.DevInfoAddDevices(prof.Builder, devsV)
	prof.Builder.Finish(flat.DevInfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes strcts.DevInfo using Flatbuffers with the package global
// Profiler.
func Serialize(inf *structs.DevInfo) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as structs.DevInfo.
func Deserialize(p []byte) *structs.DevInfo {
	devInfo := flat.GetRootAsDevInfo(p, 0)
	// get the # of interfaces
	dLen := devInfo.DevicesLength()
	info := &structs.DevInfo{Timestamp: devInfo.Timestamp(), Devices: make([]structs.Device, dLen)}
	fDev := &flat.Device{}
	sDev := structs.Device{}
	for i := 0; i < dLen; i++ {
		if devInfo.Devices(fDev, i) {
			sDev.Name = string(fDev.Name())
			sDev.RBytes = fDev.RBytes()
			sDev.RPackets = fDev.RPackets()
			sDev.RErrs = fDev.RErrs()
			sDev.RDrop = fDev.RDrop()
			sDev.RFIFO = fDev.RFIFO()
			sDev.RFrame = fDev.RFrame()
			sDev.RCompressed = fDev.RCompressed()
			sDev.RMulticast = fDev.RMulticast()
			sDev.TBytes = fDev.TBytes()
			sDev.TPackets = fDev.TPackets()
			sDev.TErrs = fDev.TErrs()
			sDev.TDrop = fDev.TDrop()
			sDev.TFIFO = fDev.TFIFO()
			sDev.TColls = fDev.TColls()
			sDev.TCarrier = fDev.TCarrier()
			sDev.TCompressed = fDev.TCompressed()
		}
		info.Devices[i] = sDev
	}
	return info
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
