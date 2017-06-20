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

// Package version handles Flatbuffer based processing of a platform's kernel
// and version information: /proc/version. Instead of returning a Go struct, it
// returns Flatbuffer serialized bytes. A function to deserialize the
// Flatbuffer serialized bytes into a kernel.Kernel struct is provided. After
// the first use, the flatbuffer builder is reused.
//
// Note: the package name is version and not the final element of the import
// path (flat). 
package version

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	v "github.com/mohae/joefriday/platform/version"
	"github.com/mohae/joefriday/platform/version/flat/flat"
)

// Profiler is used to process the version information, /proc/version, using
// Flatbuffers.
type Profiler struct {
	*v.Profiler
	*fb.Builder
}

// Initializes and returns a version information profiler that utilizes
// FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := v.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current version information as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current version information as Flatbuffer serialized bytes
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

// Serialize serializes version information using Flatbuffers.
func (prof *Profiler) Serialize(inf *v.Info) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	os := prof.Builder.CreateString(inf.OS)
	version := prof.Builder.CreateString(inf.Version)
	compileUser := prof.Builder.CreateString(inf.CompileUser)
	gcc := prof.Builder.CreateString(inf.GCC)
	osgcc := prof.Builder.CreateString(inf.OSGCC)
	typ := prof.Builder.CreateString(inf.Type)
	compileDate := prof.Builder.CreateString(inf.CompileDate)
	arch := prof.Builder.CreateString(inf.Arch)
	flat.InfoStart(prof.Builder)
	flat.InfoAddOS(prof.Builder, os)
	flat.InfoAddVersion(prof.Builder, version)
	flat.InfoAddCompileUser(prof.Builder, compileUser)
	flat.InfoAddGCC(prof.Builder, gcc)
	flat.InfoAddOSGCC(prof.Builder, osgcc)
	flat.InfoAddType(prof.Builder, typ)
	flat.InfoAddCompileDate(prof.Builder, compileDate)
	flat.InfoAddArch(prof.Builder, arch)
	prof.Builder.Finish(flat.InfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes version information using Flatbuffers with the
// package's global Profiler.
func Serialize(inf *v.Info) (p []byte, err error) {
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
// as kernel.Kernel.
func Deserialize(p []byte) *v.Info {
	flatInf := flat.GetRootAsInfo(p, 0)
	var inf v.Info
	inf.OS = string(flatInf.OS())
	inf.Version = string(flatInf.Version())
	inf.CompileUser = string(flatInf.CompileUser())
	inf.GCC = string(flatInf.GCC())
	inf.OSGCC = string(flatInf.OSGCC())
	inf.Type = string(flatInf.Type())
	inf.CompileDate = string(flatInf.CompileDate())
	inf.Arch = string(flatInf.Arch())
	return &inf
}
