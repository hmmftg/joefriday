package bench

// These are implimentations for bench purposes.

// GetInfoR accepts a *bufio.Reader and returns some of the results of
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/mem"
)

var bldr = fb.NewBuilder(0)
var buf bufio.Reader

type MemInfo struct {
	Timestamp    int64
	MemTotal     int
	MemFree      int
	MemAvailable int
	Buffers      int
	Cached       int
	SwapCached   int
	Active       int
	Inactive     int
	SwapTotal    int
	SwapFree     int
}

// Serialize serializes the MemInfo using flatbuffers.
func (i *MemInfo) Serialize() []byte {
	bldrL := fb.NewBuilder(0)
	mem.DataStart(bldrL)
	mem.DataAddTimestamp(bldrL, int64(i.Timestamp))
	mem.DataAddMemTotal(bldrL, int64(i.MemTotal))
	mem.DataAddMemFree(bldrL, int64(i.MemFree))
	mem.DataAddMemAvailable(bldrL, int64(i.MemAvailable))
	mem.DataAddBuffers(bldrL, int64(i.Buffers))
	mem.DataAddCached(bldrL, int64(i.Cached))
	mem.DataAddSwapCached(bldrL, int64(i.SwapCached))
	mem.DataAddActive(bldrL, int64(i.Active))
	mem.DataAddInactive(bldrL, int64(i.Inactive))
	mem.DataAddSwapTotal(bldrL, int64(i.SwapTotal))
	mem.DataAddSwapFree(bldrL, int64(i.SwapFree))
	bldrL.Finish(mem.DataEnd(bldrL))
	return bldrL.Bytes[bldrL.Head():]
}

// BldrSerialize serializes the MemInfo using flatbuffers: the builder is
// reused.
func (i *MemInfo) BldrSerialize() []byte {
	bldr.Reset()
	mem.DataStart(bldr)
	mem.DataAddTimestamp(bldr, int64(i.Timestamp))
	mem.DataAddMemTotal(bldr, int64(i.MemTotal))
	mem.DataAddMemFree(bldr, int64(i.MemFree))
	mem.DataAddMemAvailable(bldr, int64(i.MemAvailable))
	mem.DataAddBuffers(bldr, int64(i.Buffers))
	mem.DataAddCached(bldr, int64(i.Cached))
	mem.DataAddSwapCached(bldr, int64(i.SwapCached))
	mem.DataAddActive(bldr, int64(i.Active))
	mem.DataAddInactive(bldr, int64(i.Inactive))
	mem.DataAddSwapTotal(bldr, int64(i.SwapTotal))
	mem.DataAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(mem.DataEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// Deserialize deserializes bytes representing flatbuffers serialized Data
// into *Info.  If the bytes are not from flatbuffers serialization of
// Data, it is a programmer error and a panic will occur.
func Deserialize(p []byte) *MemInfo {
	data := mem.GetRootAsData(p, 0)
	info := &MemInfo{}
	info.Timestamp = data.Timestamp()
	info.MemTotal = int(data.MemTotal())
	info.MemFree = int(data.MemFree())
	info.MemAvailable = int(data.MemAvailable())
	info.Buffers = int(data.Buffers())
	info.Cached = int(data.Cached())
	info.SwapCached = int(data.SwapCached())
	info.Active = int(data.Active())
	info.Inactive = int(data.Inactive())
	info.SwapTotal = int(data.SwapTotal())
	info.SwapFree = int(data.SwapFree())
	return info
}

// cat /proc/meminfo.  This is mainly here for benchmark purposes.
// GetMemInfoCat returns some of the results of /proc/meminfo.
func GetMemInfoCat() (*MemInfo, error) {
	var out bytes.Buffer
	var l, i int
	var name string
	var err error
	var v byte
	t := time.Now().UTC().UnixNano()
	err = meminfo(&out)
	if err != nil {
		return nil, err
	}
	inf := &MemInfo{Timestamp: t}
	var pos int
	line := make([]byte, 0, 50)
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		line, err = out.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l > 8 && l < 15 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		name = string(val[:])
		val = val[:0]
		// skip all spaces
		for i, v = range line[pos:] {
			if v != 0x20 {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range line[pos:] {
			if v == 0x20 || v == '\r' {
				break
			}
			val = append(val, v)
		}
		// any conversion error results in 0
		i, err = strconv.Atoi(string(val[:]))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", name, err)
		}
		val = val[:0]
		if name == "MemTotal" {
			inf.MemTotal = i
			continue
		}
		if name == "MemFree" {
			inf.MemFree = i
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = i
			continue
		}
		if name == "Buffers" {
			inf.Buffers = i
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = i
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = i
			continue
		}
		if name == "Active" {
			inf.Active = i
			continue
		}
		if name == "Inactive" {
			inf.Inactive = i
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = i
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = i
			continue
		}
	}
	return inf, nil
}

// GetDataCat returns the current meminfo as flatbuffer serialized bytes.
func GetMemDataCat() ([]byte, error) {
	inf, err := GetMemInfoCat()
	if err != nil {
		return nil, err
	}
	return inf.Serialize(), nil
}

// GetMemDataCatReuseBldr reuses the Builder.
func GetMemDataCatReuseBldr() ([]byte, error) {
	inf, err := GetMemInfoCat()
	if err != nil {
		return nil, err
	}
	return inf.BldrSerialize(), nil
}

func meminfo(buff *bytes.Buffer) error {
	cmd := exec.Command("cat", "/proc/meminfo")
	cmd.Stdout = buff
	return cmd.Run()
}

// GetMemInfoRead returns some of the results of /proc/meminfo.
func GetMemInfoRead() (*MemInfo, error) {
	var l, i int
	var name string
	var err error
	var v byte
	t := time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	inf := &MemInfo{Timestamp: t}
	var pos int
	line := make([]byte, 0, 50)
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		line, err = buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l > 8 && l < 15 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		name = string(val[:])
		val = val[:0]
		// skip all spaces
		for i, v = range line[pos:] {
			if v != 0x20 {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range line[pos:] {
			if v == 0x20 || v == '\r' {
				break
			}
			val = append(val, v)
		}
		// any conversion error results in 0
		i, err = strconv.Atoi(string(val[:]))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", name, err)
		}
		val = val[:0]
		if name == "MemTotal" {
			inf.MemTotal = i
			continue
		}
		if name == "MemFree" {
			inf.MemFree = i
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = i
			continue
		}
		if name == "Buffers" {
			inf.Buffers = i
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = i
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = i
			continue
		}
		if name == "Active" {
			inf.Active = i
			continue
		}
		if name == "Inactive" {
			inf.Inactive = i
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = i
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = i
			continue
		}
	}
	return inf, nil
}

// GetMemDataRead returns the current meminfo as flatbuffer serialized bytes.
func GetMemDataRead() ([]byte, error) {
	inf, err := GetMemInfoRead()
	if err != nil {
		return nil, err
	}
	return inf.Serialize(), nil
}

// GetMemDataReadReuseBldr reuses the Builder.
func GetMemDataReadReuseBldr() ([]byte, error) {
	inf, err := GetMemInfoRead()
	if err != nil {
		return nil, err
	}
	return inf.BldrSerialize(), nil
}

// GetInfoReadReuseR returns some of the results of /proc/meminfo.
func GetMemInfoReadReuseR() (*MemInfo, error) {
	var l, i int
	var name string
	var err error
	var v byte
	t := time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf.Reset(f)
	inf := &MemInfo{Timestamp: t}
	var pos int
	line := make([]byte, 0, 50)
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		line, err = buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l > 8 && l < 15 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		name = string(val[:])
		val = val[:0]
		// skip all spaces
		for i, v = range line[pos:] {
			if v != 0x20 {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range line[pos:] {
			if v == 0x20 || v == '\r' {
				break
			}
			val = append(val, v)
		}
		// any conversion error results in 0
		i, err = strconv.Atoi(string(val[:]))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", name, err)
		}
		val = val[:0]
		if name == "MemTotal" {
			inf.MemTotal = i
			continue
		}
		if name == "MemFree" {
			inf.MemFree = i
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = i
			continue
		}
		if name == "Buffers" {
			inf.Buffers = i
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = i
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = i
			continue
		}
		if name == "Active" {
			inf.Active = i
			continue
		}
		if name == "Inactive" {
			inf.Inactive = i
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = i
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = i
			continue
		}
	}
	return inf, nil
}

// GetMemDataReadReuseR returns the current meminfo as flatbuffer serialized bytes.
func GetMemDataReadReuseR() ([]byte, error) {
	inf, err := GetMemInfoReadReuseR()
	if err != nil {
		return nil, err
	}
	return inf.Serialize(), nil
}

// GetMemDataReuseRReuseBldr reuses the Builder.
func GetMemDataReuseRReuseBldr() ([]byte, error) {
	inf, err := GetMemInfoReadReuseR()
	if err != nil {
		return nil, err
	}
	return inf.BldrSerialize(), nil
}
