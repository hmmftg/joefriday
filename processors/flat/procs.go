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

// Package flat handles Flatbuffer based processing of Processor info. Instead
// of returning a Go struct, it returns Flatbuffer serialized bytes. A function
// to deserialize the Flatbuffer serialized bytes into a processors.Processors
// struct is provided.  After the first use, the flatbuffer builder is reused.
package flat

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/processors"
)

// Profiler is used to process the processors as Flatbuffers serialized bytes.
type Profiler struct {
	*processors.Profiler
	*fb.Builder
}

// Initializes and returns a processors piler that utilizes FlatBuffers.
func NewProfiler() (p *Profiler, err error) {
	prof, err := processors.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current processor info as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	procs, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(procs), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the current processor info as Flatbuffer serialized bytes using
// the package global Profiler.
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

// Serialize serializes Processors using Flatbuffers.
func (p *Profiler) Serialize(procs *processors.Processors) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(procs.Chips))
	for i, chip := range procs.Chips {
		uoffs[i] = p.SerializeChip(&chip)
	}
	ProcessorsStartChipsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	chips := p.Builder.EndVector(len(uoffs))
	ProcessorsStart(p.Builder)
	ProcessorsAddTimestamp(p.Builder, procs.Timestamp)
	ProcessorsAddCount(p.Builder, procs.Count)
	ProcessorsAddChips(p.Builder, chips)
	p.Builder.Finish(ProcessorsEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

func (p *Profiler) SerializeChip(c *processors.Chip) fb.UOffsetT {
	vendorID := p.Builder.CreateString(c.VendorID)
	cpuFamily := p.Builder.CreateString(c.CPUFamily)
	model := p.Builder.CreateString(c.Model)
	modelName := p.Builder.CreateString(c.ModelName)
	stepping := p.Builder.CreateString(c.Stepping)
	microcode := p.Builder.CreateString(c.Microcode)
	cacheSize := p.Builder.CreateString(c.CacheSize)
	uoffs := make([]fb.UOffsetT, len(c.Flags))
	for i, flag := range c.Flags {
		uoffs[i] = p.Builder.CreateString(flag)
	}
	ChipStartFlagsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	flags := p.Builder.EndVector(len(uoffs))
	ChipStart(p.Builder)
	ChipAddPhysicalID(p.Builder, c.PhysicalID)
	ChipAddVendorID(p.Builder, vendorID)
	ChipAddCPUFamily(p.Builder, cpuFamily)
	ChipAddModel(p.Builder, model)
	ChipAddModelName(p.Builder, modelName)
	ChipAddStepping(p.Builder, stepping)
	ChipAddMicrocode(p.Builder, microcode)
	ChipAddCPUMHz(p.Builder, c.CPUMHz)
	ChipAddCacheSize(p.Builder, cacheSize)
	ChipAddCPUCores(p.Builder, c.CPUCores)
	ChipAddFlags(p.Builder, flags)
	return ChipEnd(p.Builder)
}

// Serialize Facts using the package global Profiler.
func Serialize(fcts *processors.Processors) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(fcts), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as fact.Facts.
func Deserialize(p []byte) *processors.Processors {
	flatP := GetRootAsProcessors(p, 0)
	procs := &processors.Processors{}
	flatC := &Chip{}
	proc := processors.Chip{}
	procs.Timestamp = flatP.Timestamp()
	procs.Chips = make([]processors.Chip, flatP.ChipsLength())
	for i := 0; i < len(procs.Chips); i++ {
		if !flatP.Chips(flatC, i) {
			continue
		}
		proc.PhysicalID = flatC.PhysicalID()
		proc.VendorID = string(flatC.VendorID())
		proc.CPUFamily = string(flatC.CPUFamily())
		proc.Model = string(flatC.Model())
		proc.ModelName = string(flatC.ModelName())
		proc.Stepping = string(flatC.Stepping())
		proc.Microcode = string(flatC.Microcode())
		proc.CPUMHz = flatC.CPUMHz()
		proc.CacheSize = string(flatC.CacheSize())
		proc.CPUCores = flatC.CPUCores()
		proc.Flags = make([]string, flatC.FlagsLength())
		for i := 0; i < len(proc.Flags); i++ {
			proc.Flags[i] = string(flatC.Flags(i))
		}
		procs.Chips[i] = proc
	}
	return procs
}
