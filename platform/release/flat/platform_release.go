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

// Package flat handles Flatbuffer based processing of a platform's OS
// information using /etc/os-release.  Instead of returning a Go struct, it
// returns Flatbuffer serialized bytes.  A function to deserialize the
// Flatbuffer serialized bytes into a release.Release struct is provided.
//  After the first use, the flatbuffer builder is reused.
package flat

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/platform/release"
)

// Profiler is used to process the os information, /etc/os-release using
// Flatbuffers.
type Profiler struct {
	*release.Profiler
	*fb.Builder
}

// Initializes and returns an OS information profiler that utilizes
// FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := release.NewProfiler()
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
func (prof *Profiler) Serialize(r *release.Release) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	id := prof.Builder.CreateString(r.ID)
	idLike := prof.Builder.CreateString(r.IDLike)
	prettyName := prof.Builder.CreateString(r.PrettyName)
	version := prof.Builder.CreateString(r.Version)
	versionID := prof.Builder.CreateString(r.VersionID)
	homeURL := prof.Builder.CreateString(r.HomeURL)
	bugReportURL := prof.Builder.CreateString(r.BugReportURL)
	ReleaseStart(prof.Builder)
	ReleaseAddID(prof.Builder, id)
	ReleaseAddIDLike(prof.Builder, idLike)
	ReleaseAddPrettyName(prof.Builder, prettyName)
	ReleaseAddVersion(prof.Builder, version)
	ReleaseAddVersionID(prof.Builder, versionID)
	ReleaseAddHomeURL(prof.Builder, homeURL)
	ReleaseAddBugReportURL(prof.Builder, bugReportURL)
	prof.Builder.Finish(ReleaseEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize serializes OS release information using Flatbuffers with the
// package's global Profiler.
func Serialize(r *release.Release) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(r), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as release.Release.
func Deserialize(p []byte) *release.Release {
	flatR := GetRootAsRelease(p, 0)
	var r release.Release
	r.ID = string(flatR.ID())
	r.IDLike = string(flatR.IDLike())
	r.HomeURL = string(flatR.HomeURL())
	r.PrettyName = string(flatR.PrettyName())
	r.Version = string(flatR.Version())
	r.VersionID = string(flatR.VersionID())
	r.BugReportURL = string(flatR.BugReportURL())
	return &r
}
