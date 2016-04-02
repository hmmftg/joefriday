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

// Package net gets and processes /proc/net/dev, returning the infoFlat in the
// appropriate format.
package flat

import (
	"sync"
	"time"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/net/info"
	"github.com/mohae/joefriday/net/structs"
	"github.com/mohae/joefriday/net/structs/flat"
)

type Profiler struct {
	Prof    *info.Profiler
	Builder *fb.Builder
}

func New() (prof *Profiler, err error) {
	p, err := info.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p, Builder: fb.NewBuilder(0)}, nil
}

func (prof *Profiler) Reset() error {
	prof.Prof.Lock()
	prof.Builder.Reset()
	prof.Prof.Unlock()
	return prof.Prof.Reset()
}

// Get returns the current meminfo as flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	prof.Reset()
	inf, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get get's the current meminfo.
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

func (prof *Profiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			info, err := prof.Get()
			if err != nil {
				errs <- err
				continue
			}
			out <- info
		}
	}
}

func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// Serialize serializes Info using Flatbuffers.
func (prof *Profiler) Serialize(inf *structs.Info) []byte {
	ifaces := make([]fb.UOffsetT, len(inf.Interfaces))
	names := make([]fb.UOffsetT, len(inf.Interfaces))
	for i := 0; i < len(inf.Interfaces); i++ {
		names[i] = prof.Builder.CreateString(inf.Interfaces[i].Name)
	}
	for i := 0; i < len(inf.Interfaces); i++ {
		flat.InterfaceStart(prof.Builder)
		flat.InterfaceAddName(prof.Builder, names[i])
		flat.InterfaceAddRBytes(prof.Builder, inf.Interfaces[i].RBytes)
		flat.InterfaceAddRPackets(prof.Builder, inf.Interfaces[i].RPackets)
		flat.InterfaceAddRErrs(prof.Builder, inf.Interfaces[i].RErrs)
		flat.InterfaceAddRDrop(prof.Builder, inf.Interfaces[i].RDrop)
		flat.InterfaceAddRFIFO(prof.Builder, inf.Interfaces[i].RFIFO)
		flat.InterfaceAddRFrame(prof.Builder, inf.Interfaces[i].RFrame)
		flat.InterfaceAddRCompressed(prof.Builder, inf.Interfaces[i].RCompressed)
		flat.InterfaceAddRMulticast(prof.Builder, inf.Interfaces[i].RMulticast)
		flat.InterfaceAddTBytes(prof.Builder, inf.Interfaces[i].TBytes)
		flat.InterfaceAddTPackets(prof.Builder, inf.Interfaces[i].TPackets)
		flat.InterfaceAddTErrs(prof.Builder, inf.Interfaces[i].TErrs)
		flat.InterfaceAddTDrop(prof.Builder, inf.Interfaces[i].TDrop)
		flat.InterfaceAddTFIFO(prof.Builder, inf.Interfaces[i].TFIFO)
		flat.InterfaceAddTColls(prof.Builder, inf.Interfaces[i].TColls)
		flat.InterfaceAddTCarrier(prof.Builder, inf.Interfaces[i].TCarrier)
		flat.InterfaceAddTCompressed(prof.Builder, inf.Interfaces[i].TCompressed)
		ifaces[i] = flat.InterfaceEnd(prof.Builder)
	}
	flat.InfoStartInterfacesVector(prof.Builder, len(ifaces))
	for i := len(inf.Interfaces) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(ifaces[i])
	}
	ifacesV := prof.Builder.EndVector(len(ifaces))
	flat.InfoStart(prof.Builder)
	flat.InfoAddTimestamp(prof.Builder, inf.Timestamp)
	flat.InfoAddInterfaces(prof.Builder, ifacesV)
	prof.Builder.Finish(flat.InfoEnd(prof.Builder))
	return prof.Builder.Bytes[prof.Builder.Head():]
}

// Deserialize deserializes struct.Info Flatbuffer serialized bytes.
func Deserialize(p []byte) *structs.Info {
	infoFlat := flat.GetRootAsInfo(p, 0)
	// get the # of interfaces
	iLen := infoFlat.InterfacesLength()
	info := &structs.Info{Timestamp: infoFlat.Timestamp(), Interfaces: make([]structs.Interface, iLen)}
	iFace := &flat.Interface{}
	iface := structs.Interface{}
	for i := 0; i < iLen; i++ {
		if infoFlat.Interfaces(iFace, i) {
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
		info.Interfaces[i] = iface
	}
	return info
}
