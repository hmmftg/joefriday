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
	r "github.com/mohae/joefriday/platform/release"
	"github.com/mohae/joefriday/platform/release/flat/structs"
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
func (prof *Profiler) Serialize(inf *r.Info) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	name := prof.Builder.CreateString(inf.Name)
	id := prof.Builder.CreateString(inf.ID)
	idLike := prof.Builder.CreateString(inf.IDLike)
	prettyName := prof.Builder.CreateString(inf.PrettyName)
	version := prof.Builder.CreateString(inf.Version)
	versionID := prof.Builder.CreateString(inf.VersionID)
	homeURL := prof.Builder.CreateString(inf.HomeURL)
	bugReportURL := prof.Builder.CreateString(inf.BugReportURL)
	structs.InfoStart(prof.Builder)
	structs.InfoAddName(prof.Builder, name)
	structs.InfoAddID(prof.Builder, id)
	structs.InfoAddIDLike(prof.Builder, idLike)
	structs.InfoAddPrettyName(prof.Builder, prettyName)
	structs.InfoAddVersion(prof.Builder, version)
	structs.InfoAddVersionID(prof.Builder, versionID)
	structs.InfoAddHomeURL(prof.Builder, homeURL)
	structs.InfoAddBugReportURL(prof.Builder, bugReportURL)
	prof.Builder.Finish(structs.InfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes OS release information using Flatbuffers with the
// package's global Profiler.
func Serialize(inf *r.Info) (p []byte, err error) {
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
// as release.Release.
func Deserialize(p []byte) *r.Info {
	flatInf := structs.GetRootAsInfo(p, 0)
	var inf r.Info
	inf.Name = string(flatInf.Name())
	inf.ID = string(flatInf.ID())
	inf.IDLike = string(flatInf.IDLike())
	inf.HomeURL = string(flatInf.HomeURL())
	inf.PrettyName = string(flatInf.PrettyName())
	inf.Version = string(flatInf.Version())
	inf.VersionID = string(flatInf.VersionID())
	inf.BugReportURL = string(flatInf.BugReportURL())
	return &inf
}
