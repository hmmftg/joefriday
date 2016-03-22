package platform

import (
	"bufio"
	"io"
	"os"
	"strings"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
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

// Serialize serializes Kernel using Flatbuffers.
func (k Kernel) Serialize() []byte {
	bldr := fb.NewBuilder(0)
	version := bldr.CreateString(k.Version)
	compileUser := bldr.CreateString(k.CompileUser)
	gcc := bldr.CreateString(k.GCC)
	osgcc := bldr.CreateString(k.OSGCC)
	typ := bldr.CreateString(k.Type)
	compileDate := bldr.CreateString(k.CompileDate)
	arch := bldr.CreateString(k.Arch)
	KernelFBStart(bldr)
	KernelFBAddVersion(bldr, version)
	KernelFBAddCompileUser(bldr, compileUser)
	KernelFBAddGCC(bldr, gcc)
	KernelFBAddOSGCC(bldr, osgcc)
	KernelFBAddType(bldr, typ)
	KernelFBAddCompileDate(bldr, compileDate)
	KernelFBAddArch(bldr, arch)
	bldr.Finish(KernelFBEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeKernel deserializes the bytes into Kernel using Flatbuffers.
func DeserializeKernel(p []byte) Kernel {
	kernelFB := GetRootAsKernelFB(p, 0)
	var kernel Kernel
	kernel.Version = string(kernelFB.Version())
	kernel.CompileUser = string(kernelFB.CompileUser())
	kernel.GCC = string(kernelFB.GCC())
	kernel.OSGCC = string(kernelFB.OSGCC())
	kernel.Type = string(kernelFB.Type())
	kernel.CompileDate = string(kernelFB.CompileDate())
	kernel.Arch = string(kernelFB.Arch())
	return kernel
}

// GetKernel populates Kernel with /proc/version information.
func GetKernel() (Kernel, error) {
	var i, pos, pos2 int
	var v byte
	var kernel Kernel
	f, err := os.Open("/proc/version")
	if err != nil {
		return kernel, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return kernel, joe.Error{Type: "platform", Op: "read /proc/verion", Err: err}
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
	return kernel, nil
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

// Serialize serializes Release using Flatbuffers.
func (r Release) Serialize() []byte {
	bldr := fb.NewBuilder(0)
	id := bldr.CreateString(r.ID)
	idLike := bldr.CreateString(r.IDLike)
	prettyName := bldr.CreateString(r.PrettyName)
	version := bldr.CreateString(r.Version)
	versionID := bldr.CreateString(r.VersionID)
	homeURL := bldr.CreateString(r.HomeURL)
	bugReportURL := bldr.CreateString(r.BugReportURL)
	ReleaseFBStart(bldr)
	ReleaseFBAddID(bldr, id)
	ReleaseFBAddIDLike(bldr, idLike)
	ReleaseFBAddPrettyName(bldr, prettyName)
	ReleaseFBAddVersion(bldr, version)
	ReleaseFBAddVersionID(bldr, versionID)
	ReleaseFBAddHomeURL(bldr, homeURL)
	ReleaseFBAddBugReportURL(bldr, bugReportURL)
	bldr.Finish(ReleaseFBEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeRelease deserializes bytes into Release using Flatbuffers.
func DeserializeRelease(p []byte) Release {
	releaseFB := GetRootAsReleaseFB(p, 0)
	var release Release
	release.ID = string(releaseFB.ID())
	release.IDLike = string(releaseFB.IDLike())
	release.PrettyName = string(releaseFB.PrettyName())
	release.Version = string(releaseFB.Version())
	release.VersionID = string(releaseFB.VersionID())
	release.HomeURL = string(releaseFB.HomeURL())
	release.BugReportURL = string(releaseFB.BugReportURL())
	return release
}

// GetRelease populates release: the source depends on the OS
func GetRelease() (Release, error) {
	var i int
	var v byte
	var release Release
	var err error
	var key string
	val := make([]byte, 32)
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return release, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return release, joe.Error{Type: "platform", Op: "read /proc/verion", Err: err}
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
	return release, nil
}
