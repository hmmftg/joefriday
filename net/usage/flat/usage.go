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

// Package flat handles Flatbuffer based processing of network interface
//usage.  Interface usage is calculated by taking the difference in two
// /proc/net/dev snapshots and reflect bytes received and transmitted since
// the prior snapshot.  Instead of returning a Go struct, it returns
// Flatbuffer serialized bytes.  A function to deserialize the Flatbuffer
// serialized bytes into a structs.Usage struct.  After the first use, the
// flatbuffer builder is reused.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/net/structs"
	"github.com/mohae/joefriday/net/structs/flat"
	"github.com/mohae/joefriday/net/usage"
)

// Profiler is used to process the network interface usage using Flatbuffers.
type Profiler struct {
	*usage.Profiler
	*fb.Builder
}

// Initializes and returns a network interface usage profiler that utilizes
// FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := usage.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current network interface usage as Flatbuffer serialized
// bytes.
// TODO: should this be changed so that this calculates usage since the last
// time the network uo was obtained.  If there aren't pre-existing uo
// it would get current usage (which may be a separate method (or should be?))
func (prof *Profiler) Get() (p []byte, err error) {
	u, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(u), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network interface usage as Flatbuffer serialized
// bytes using the package's global Profiler.
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

// Serialize serializes Usage using Flatbuffers.
func (prof *Profiler) Serialize(u *structs.Usage) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	ifaces := make([]fb.UOffsetT, len(u.Interfaces))
	names := make([]fb.UOffsetT, len(u.Interfaces))
	for i := 0; i < len(u.Interfaces); i++ {
		names[i] = prof.Builder.CreateString(u.Interfaces[i].Name)
	}
	for i := 0; i < len(u.Interfaces); i++ {
		flat.InterfaceStart(prof.Builder)
		flat.InterfaceAddName(prof.Builder, names[i])
		flat.InterfaceAddRBytes(prof.Builder, u.Interfaces[i].RBytes)
		flat.InterfaceAddRPackets(prof.Builder, u.Interfaces[i].RPackets)
		flat.InterfaceAddRErrs(prof.Builder, u.Interfaces[i].RErrs)
		flat.InterfaceAddRDrop(prof.Builder, u.Interfaces[i].RDrop)
		flat.InterfaceAddRFIFO(prof.Builder, u.Interfaces[i].RFIFO)
		flat.InterfaceAddRFrame(prof.Builder, u.Interfaces[i].RFrame)
		flat.InterfaceAddRCompressed(prof.Builder, u.Interfaces[i].RCompressed)
		flat.InterfaceAddRMulticast(prof.Builder, u.Interfaces[i].RMulticast)
		flat.InterfaceAddTBytes(prof.Builder, u.Interfaces[i].TBytes)
		flat.InterfaceAddTPackets(prof.Builder, u.Interfaces[i].TPackets)
		flat.InterfaceAddTErrs(prof.Builder, u.Interfaces[i].TErrs)
		flat.InterfaceAddTDrop(prof.Builder, u.Interfaces[i].TDrop)
		flat.InterfaceAddTFIFO(prof.Builder, u.Interfaces[i].TFIFO)
		flat.InterfaceAddTColls(prof.Builder, u.Interfaces[i].TColls)
		flat.InterfaceAddTCarrier(prof.Builder, u.Interfaces[i].TCarrier)
		flat.InterfaceAddTCompressed(prof.Builder, u.Interfaces[i].TCompressed)
		ifaces[i] = flat.InterfaceEnd(prof.Builder)
	}
	flat.UsageStartInterfacesVector(prof.Builder, len(ifaces))
	for i := len(u.Interfaces) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(ifaces[i])
	}
	ifacesV := prof.Builder.EndVector(len(ifaces))
	flat.UsageStart(prof.Builder)
	flat.UsageAddTimestamp(prof.Builder, u.Timestamp)
	flat.UsageAddTimeDelta(prof.Builder, u.TimeDelta)
	flat.UsageAddInterfaces(prof.Builder, ifacesV)
	prof.Builder.Finish(flat.UsageEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes Usage using Flatbuffers with the package global
// Profiler.
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
// as structs.Usage.
func Deserialize(p []byte) *structs.Usage {
	uFlat := flat.GetRootAsUsage(p, 0)
	// get the # of interfaces
	iLen := uFlat.InterfacesLength()
	u := &structs.Usage{
		Timestamp:  uFlat.Timestamp(),
		TimeDelta:  uFlat.TimeDelta(),
		Interfaces: make([]structs.Interface, iLen),
	}
	iFace := &flat.Interface{}
	iface := structs.Interface{}
	for i := 0; i < iLen; i++ {
		if uFlat.Interfaces(iFace, i) {
			iface.Name = string(iFace.Name())
			iface.RBytes = iFace.RBytes()
			iface.RPackets = iFace.RPackets()
			iface.RErrs = iFace.RErrs()
			iface.RDrop = iFace.RDrop()
			iface.RFIFO = iFace.RFIFO()
			iface.RFrame = iFace.RFrame()
			iface.RCompressed = iFace.RCompressed()
			iface.RMulticast = iFace.RMulticast()
			iface.TBytes = iFace.TBytes()
			iface.TPackets = iFace.TPackets()
			iface.TErrs = iFace.TErrs()
			iface.TDrop = iFace.TDrop()
			iface.TFIFO = iFace.TFIFO()
			iface.TColls = iFace.TColls()
			iface.TCarrier = iFace.TCarrier()
			iface.TCompressed = iFace.TCompressed()
		}
		u.Interfaces[i] = iface
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
