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

// Package release handles Flatbuffer based processing of a platform's OS
// information using /etc/os-release.  Instead of returning a Go struct, it
// returns Flatbuffer serialized bytes.  A function to deserialize the
// Flatbuffer serialized bytes into a release.Release struct is provided.
//  After the first use, the flatbuffer builder is reused.
//
// Note: the package name is release and not the final element of the import
// path (flat). 
package release

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	r "github.com/mohae/joefriday/system/release"
	"github.com/mohae/joefriday/system/release/flat/structs"
)

// Profiler is used to process the os information, /etc/os-release using
// Flatbuffers.
type Profiler struct {
	*r.Profiler
	*fb.Builder
}

// Initializes and returns an OS information profiler that utilizes
// FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := r.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current OS release information as Flatbuffer serialized
// bytes.
func (prof *Profiler) Get() ([]byte, error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current OS release information as Flatbuffer serialized
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

// Serialize serializes OS release information using Flatbuffers.
func (prof *Profiler) Serialize(os *r.OS) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	name := prof.Builder.CreateString(os.Name)
	id := prof.Builder.CreateString(os.ID)
	idLike := prof.Builder.CreateString(os.IDLike)
	prettyName := prof.Builder.CreateString(os.PrettyName)
	version := prof.Builder.CreateString(os.Version)
	versionID := prof.Builder.CreateString(os.VersionID)
	homeURL := prof.Builder.CreateString(os.HomeURL)
	bugReportURL := prof.Builder.CreateString(os.BugReportURL)
	structs.OSStart(prof.Builder)
	structs.OSAddName(prof.Builder, name)
	structs.OSAddID(prof.Builder, id)
	structs.OSAddIDLike(prof.Builder, idLike)
	structs.OSAddPrettyName(prof.Builder, prettyName)
	structs.OSAddVersion(prof.Builder, version)
	structs.OSAddVersionID(prof.Builder, versionID)
	structs.OSAddHomeURL(prof.Builder, homeURL)
	structs.OSAddBugReportURL(prof.Builder, bugReportURL)
	prof.Builder.Finish(structs.OSEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes OS release information using Flatbuffers with the
// package's global Profiler.
func Serialize(os *r.OS) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(os), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as release.Release.
func Deserialize(p []byte) *r.OS {
	flatOS := structs.GetRootAsOS(p, 0)
	var os r.OS
	os.Name = string(flatOS.Name())
	os.ID = string(flatOS.ID())
	os.IDLike = string(flatOS.IDLike())
	os.HomeURL = string(flatOS.HomeURL())
	os.PrettyName = string(flatOS.PrettyName())
	os.Version = string(flatOS.Version())
	os.VersionID = string(flatOS.VersionID())
	os.BugReportURL = string(flatOS.BugReportURL())
	return &os
}
