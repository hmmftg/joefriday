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

// Package flat handles Flatbuffer based processing of a platform's kernel
// information: /proc/version.  Instead of returning a Go struct, it returns
// Flatbuffer serialized bytes.  A function to deserialize the Flatbuffer
// serialized bytes into a kernel.Kernel struct is provided.  After the first
// use, the flatbuffer builder is reused.
package flat

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/platform/kernel"
)

// Profiler is used to process the kernel information, /proc/version, using
// Flatbuffers.
type Profiler struct {
	*kernel.Profiler
	*fb.Builder
}

// Initializes and returns a kernel information profiler that utilizes
// FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := kernel.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current kernel information as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current kernel information as Flatbuffer serialized bytes
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

// Serialize serializes kernel information using Flatbuffers.
func (prof *Profiler) Serialize(k *kernel.Kernel) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	os := prof.Builder.CreateString(k.OS)
	version := prof.Builder.CreateString(k.Version)
	compileUser := prof.Builder.CreateString(k.CompileUser)
	gcc := prof.Builder.CreateString(k.GCC)
	osgcc := prof.Builder.CreateString(k.OSGCC)
	typ := prof.Builder.CreateString(k.Type)
	compileDate := prof.Builder.CreateString(k.CompileDate)
	arch := prof.Builder.CreateString(k.Arch)
	KernelStart(prof.Builder)
	KernelAddOS(prof.Builder, os)
	KernelAddVersion(prof.Builder, version)
	KernelAddCompileUser(prof.Builder, compileUser)
	KernelAddGCC(prof.Builder, gcc)
	KernelAddOSGCC(prof.Builder, osgcc)
	KernelAddType(prof.Builder, typ)
	KernelAddCompileDate(prof.Builder, compileDate)
	KernelAddArch(prof.Builder, arch)
	prof.Builder.Finish(KernelEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes kernel information using Flatbuffers with the
// package's global Profiler.
func Serialize(k *kernel.Kernel) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(k), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as kernel.Kernel.
func Deserialize(p []byte) *kernel.Kernel {
	flatK := GetRootAsKernel(p, 0)
	var k kernel.Kernel
	k.OS = string(flatK.OS())
	k.Version = string(flatK.Version())
	k.CompileUser = string(flatK.CompileUser())
	k.GCC = string(flatK.GCC())
	k.OSGCC = string(flatK.OSGCC())
	k.Type = string(flatK.Type())
	k.CompileDate = string(flatK.CompileDate())
	k.Arch = string(flatK.Arch())
	return &k
}
