package platform

import (
	"bufio"
	"io"
	"os"
	"strings"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/platform/flat"
)

// Kernel holds information about the kernel; this is /proc/version output
type Kernel struct {
	Version     string `json:"version"`
	CompileUser string `json:"compile_user"`
	GCC         string `json:"gcc"`
	OSGCC       string `json:"os_gcc"`
	Type        string `json:"type"`
	CompileDate string `json:"compile_date"`
	Arch        string `json:"arch"`
}

// SerializeFlat serializes Kernel using Flatbuffers.
func (k *Kernel) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return k.SerializeFlatBuilder(bldr)
}

// SerializeFlatBuilder uses the passed flat.Builder to serialize Kernel.
// The builder is expected to be in a usable state.
func (k *Kernel) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	version := bldr.CreateString(k.Version)
	compileUser := bldr.CreateString(k.CompileUser)
	gcc := bldr.CreateString(k.GCC)
	osgcc := bldr.CreateString(k.OSGCC)
	typ := bldr.CreateString(k.Type)
	compileDate := bldr.CreateString(k.CompileDate)
	arch := bldr.CreateString(k.Arch)
	flat.KernelStart(bldr)
	flat.KernelAddVersion(bldr, version)
	flat.KernelAddCompileUser(bldr, compileUser)
	flat.KernelAddGCC(bldr, gcc)
	flat.KernelAddOSGCC(bldr, osgcc)
	flat.KernelAddType(bldr, typ)
	flat.KernelAddCompileDate(bldr, compileDate)
	flat.KernelAddArch(bldr, arch)
	bldr.Finish(flat.KernelEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeKernelFlat deserializes the bytes into Kernel using Flatbuffers.
func DeserializeKernelFlat(p []byte) *Kernel {
	flatKernel := flat.GetRootAsKernel(p, 0)
	var kernel Kernel
	kernel.Version = string(flatKernel.Version())
	kernel.CompileUser = string(flatKernel.CompileUser())
	kernel.GCC = string(flatKernel.GCC())
	kernel.OSGCC = string(flatKernel.OSGCC())
	kernel.Type = string(flatKernel.Type())
	kernel.CompileDate = string(flatKernel.CompileDate())
	kernel.Arch = string(flatKernel.Arch())
	return &kernel
}

// GetKernel populates Kernel with /proc/version information.
func GetKernel() (*Kernel, error) {
	var i, pos, pos2 int
	var v byte
	var kernel Kernel
	f, err := os.Open("/proc/version")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, joe.Error{Type: "platform", Op: "read /proc/verion", Err: err}
		}
		// The version is everything up to the first '(', 0x28, - 1 byte
		for i, v = range line {
			if v == 0x28 {
				kernel.Version = string(line[:i-1])
				pos = i + 1
				break
			}
		}

		// The last part of the version should be the arch
		parts := strings.Split(kernel.Version, "-")
		kernel.Arch = parts[len(parts)-1]
		// The CompileUser is everything up to the next ')', 0x29
		for i, v = range line[pos:] {
			if v == 0x29 {
				kernel.CompileUser = string(line[pos : pos+i])
				pos += i + 3
				break
			}
		}

		var inOSGCC bool
		// GCC info; this may include os specific gcc info
		for i, v = range line[pos:] {
			if v == 0x28 {
				inOSGCC = true
				kernel.GCC = string(line[pos : pos+i-1])
				pos2 = i + pos + 1
				continue
			}
			if v == 0x29 {
				if inOSGCC {
					kernel.OSGCC = string(line[pos2 : pos+i])
					inOSGCC = false
					continue
				}
				pos, pos2 = pos+i+2, pos
				break
			}
		}
		// Check if GCC is empty, this happens if there wasn't an OSGCC value
		if kernel.GCC == "" {
			kernel.GCC = string(line[pos2 : pos-1])
		}
		// Get the type information, everything up to '('
		for i, v = range line[pos:] {
			if v == 0x28 {
				kernel.Type = string(line[pos : pos+i-1])
				pos += i + 1
				break
			}
		}
		// The rest is the compile date.
		kernel.CompileDate = string(line[pos : len(line)-2])
	}
	return &kernel, nil
}

// GetKernelFlat returns the Flatbuffer serialized Kernel information.
func GetKernelFlat() ([]byte, error) {
	k, err := GetKernel()
	if err != nil {
		return nil, err
	}
	return k.SerializeFlat(), nil
}

// Release holds information about the release.  The source depends on the
// OS.  Currently only Debian and Redhat families are supported.
type Release struct {
	ID           string `json:"id"`
	IDLike       string `json:"id_like"`
	PrettyName   string `json:"pretty_name"`
	Version      string `json:"version"`
	VersionID    string `json:"version_id"`
	HomeURL      string `json:"home_url"`
	BugReportURL string `json:"bug_report_url"`
}

// SerializeFlat serializes Release using Flatbuffers.
func (r *Release) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return r.SerializeFlatBuilder(bldr)
}

// SerializeFlatBuilder uses the passed flat.Builder to serialize Release.
// The builder is expected to be in a usable state.
func (r *Release) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	id := bldr.CreateString(r.ID)
	idLike := bldr.CreateString(r.IDLike)
	prettyName := bldr.CreateString(r.PrettyName)
	version := bldr.CreateString(r.Version)
	versionID := bldr.CreateString(r.VersionID)
	homeURL := bldr.CreateString(r.HomeURL)
	bugReportURL := bldr.CreateString(r.BugReportURL)
	flat.ReleaseStart(bldr)
	flat.ReleaseAddID(bldr, id)
	flat.ReleaseAddIDLike(bldr, idLike)
	flat.ReleaseAddPrettyName(bldr, prettyName)
	flat.ReleaseAddVersion(bldr, version)
	flat.ReleaseAddVersionID(bldr, versionID)
	flat.ReleaseAddHomeURL(bldr, homeURL)
	flat.ReleaseAddBugReportURL(bldr, bugReportURL)
	bldr.Finish(flat.ReleaseEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeReleaseFlat deserializes bytes into Release using Flatbuffers.
func DeserializeReleaseFlat(p []byte) *Release {
	releaseFlat := flat.GetRootAsRelease(p, 0)
	var release Release
	release.ID = string(releaseFlat.ID())
	release.IDLike = string(releaseFlat.IDLike())
	release.PrettyName = string(releaseFlat.PrettyName())
	release.Version = string(releaseFlat.Version())
	release.VersionID = string(releaseFlat.VersionID())
	release.HomeURL = string(releaseFlat.HomeURL())
	release.BugReportURL = string(releaseFlat.BugReportURL())
	return &release
}

// GetRelease populates release: the source depends on the OS
func GetRelease() (*Release, error) {
	var i int
	var v byte
	var release Release
	var err error
	var key string
	val := make([]byte, 32)
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, joe.Error{Type: "platform", Op: "read /proc/verion", Err: err}
		}
		// The key is everything up to '='; 0x3D
		for i, v = range line {
			if v == 0x3D {
				key = string(line[:i-1])
				val = line[i+1 : len(line)-1]
				break
			}
		}
		// See if the value is quoted; remove quotes if it is
		if val[0] == 0x22 {
			val = val[1 : len(val)-1]
		}
		if key == "ID" {
			release.ID = string(val)
			continue
		}
		if key == "ID_LIKE" {
			release.IDLike = string(val)
			continue
		}
		if key == "PRETTY_NAME" {
			release.PrettyName = string(val)
			continue
		}
		if key == "VERSION" {
			release.Version = string(val)
			continue
		}
		if key == "VERSION_ID" {
			release.VersionID = string(val)
			continue
		}
		if key == "HOME_URL" {
			release.HomeURL = string(val)
			continue
		}
		if key == "BUG_REPORT_URL" {
			release.BugReportURL = string(val)
			continue
		}
	}
	return &release, nil
}

// GetReleaseFlat returns Flatbuffer serialized Release information.
func GetReleaseFlat() ([]byte, error) {
	r, err := GetRelease()
	if err != nil {
		return nil, err
	}
	return r.SerializeFlat(), nil
}
