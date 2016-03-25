package bench

// These are implimentations for bench purposes.

// GetInfoR accepts a *bufio.Reader and returns some of the results of
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/SermoDigital/helpers"
	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/mem/flat"
)

var bldr = fb.NewBuilder(0)
var buf *bufio.Reader
var readB = make([]byte, 1536)

func init() {
	tmp := bytes.NewBuffer(readB)
	buf = bufio.NewReader(tmp)
	tmp = nil
}

type MemInfo struct {
	Timestamp    int64 `json:"timestamp"`
	MemTotal     int64 `json:"mem_total"`
	MemFree      int64 `json:"mem_free"`
	MemAvailable int64 `json:"mem_available"`
	Buffers      int64 `json:"buffers"`
	Cached       int64 `json:"cached"`
	SwapCached   int64 `json:"swap_cached"`
	Active       int64 `json:"active"`
	Inactive     int64 `json:"inactive"`
	SwapTotal    int64 `json:"swapt_total"`
	SwapFree     int64 `json:"swap_free"`
}

// Serialize serializes the MemInfo using flatbuffers.
func (i *MemInfo) SerializeFlat() []byte {
	bldrL := fb.NewBuilder(0)
	flat.InfoStart(bldrL)
	flat.InfoAddTimestamp(bldrL, int64(i.Timestamp))
	flat.InfoAddMemTotal(bldrL, int64(i.MemTotal))
	flat.InfoAddMemFree(bldrL, int64(i.MemFree))
	flat.InfoAddMemAvailable(bldrL, int64(i.MemAvailable))
	flat.InfoAddBuffers(bldrL, int64(i.Buffers))
	flat.InfoAddCached(bldrL, int64(i.Cached))
	flat.InfoAddSwapCached(bldrL, int64(i.SwapCached))
	flat.InfoAddActive(bldrL, int64(i.Active))
	flat.InfoAddInactive(bldrL, int64(i.Inactive))
	flat.InfoAddSwapTotal(bldrL, int64(i.SwapTotal))
	flat.InfoAddSwapFree(bldrL, int64(i.SwapFree))
	bldrL.Finish(flat.InfoEnd(bldrL))
	return bldrL.Bytes[bldrL.Head():]
}

// BldrSerialize serializes the MemInfo using flatbuffers: the builder is
// reused.
func (i *MemInfo) BldrSerializeFlat() []byte {
	bldr.Reset()
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, int64(i.Timestamp))
	flat.InfoAddMemTotal(bldr, int64(i.MemTotal))
	flat.InfoAddMemFree(bldr, int64(i.MemFree))
	flat.InfoAddMemAvailable(bldr, int64(i.MemAvailable))
	flat.InfoAddBuffers(bldr, int64(i.Buffers))
	flat.InfoAddCached(bldr, int64(i.Cached))
	flat.InfoAddSwapCached(bldr, int64(i.SwapCached))
	flat.InfoAddActive(bldr, int64(i.Active))
	flat.InfoAddInactive(bldr, int64(i.Inactive))
	flat.InfoAddSwapTotal(bldr, int64(i.SwapTotal))
	flat.InfoAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeFlat deserializes bytes representing flatbuffers serialized
// InfoFlat into *Info.
func DeserializeFlat(p []byte) *MemInfo {
	infoFlat := flat.GetRootAsInfo(p, 0)
	info := &MemInfo{}
	info.Timestamp = infoFlat.Timestamp()
	info.MemTotal = infoFlat.MemTotal()
	info.MemFree = infoFlat.MemFree()
	info.MemAvailable = infoFlat.MemAvailable()
	info.Buffers = infoFlat.Buffers()
	info.Cached = infoFlat.Cached()
	info.SwapCached = infoFlat.SwapCached()
	info.Active = infoFlat.Active()
	info.Inactive = infoFlat.Inactive()
	info.SwapTotal = infoFlat.SwapTotal()
	info.SwapFree = infoFlat.SwapFree()
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
			inf.MemTotal = int64(i)
			continue
		}
		if name == "MemFree" {
			inf.MemFree = int64(i)
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "Buffers" {
			inf.Buffers = int64(i)
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = int64(i)
			continue
		}
		if name == "Active" {
			inf.Active = int64(i)
			continue
		}
		if name == "Inactive" {
			inf.Inactive = int64(i)
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = int64(i)
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = int64(i)
			continue
		}
	}
	return inf, nil
}

func GetMemInfoCatToJSON() ([]byte, error) {
	inf, err := GetMemInfoCat()
	if err != nil {
		return nil, err
	}
	return json.Marshal(inf)
}

// GetDataCat returns the current meminfo as flatbuffer serialized bytes.
func GetMemDataCat() ([]byte, error) {
	inf, err := GetMemInfoCat()
	if err != nil {
		return nil, err
	}
	return inf.SerializeFlat(), nil
}

// GetMemDataCatReuseBldr reuses the Builder.
func GetMemDataCatReuseBldr() ([]byte, error) {
	inf, err := GetMemInfoCat()
	if err != nil {
		return nil, err
	}
	return inf.BldrSerializeFlat(), nil
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
	bf := bufio.NewReader(f)
	inf := &MemInfo{Timestamp: t}
	var pos int
	line := make([]byte, 0, 50)
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		line, err = bf.ReadSlice('\n')
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
			inf.MemTotal = int64(i)
			continue
		}
		if name == "MemFree" {
			inf.MemFree = int64(i)
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "Buffers" {
			inf.Buffers = int64(i)
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = int64(i)
			continue
		}
		if name == "Active" {
			inf.Active = int64(i)
			continue
		}
		if name == "Inactive" {
			inf.Inactive = int64(i)
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = int64(i)
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = int64(i)
			continue
		}
	}
	return inf, nil
}

func GetMemInfoReadToJSON() ([]byte, error) {
	inf, err := GetMemInfoRead()
	if err != nil {
		return nil, err
	}
	return json.Marshal(inf)
}

// GetMemDataRead returns the current meminfo as flatbuffer serialized bytes.
func GetMemDataRead() ([]byte, error) {
	inf, err := GetMemInfoRead()
	if err != nil {
		return nil, err
	}
	return inf.SerializeFlat(), nil
}

// GetMemDataReadReuseBldr reuses the Builder.
func GetMemDataReadReuseBldr() ([]byte, error) {
	inf, err := GetMemInfoRead()
	if err != nil {
		return nil, err
	}
	return inf.BldrSerializeFlat(), nil
}

// GetInfoReadReuseR returns some of the results of /proc/meminfo.
func GetMemInfoReadReuseR() (*MemInfo, error) {
	var l, i int
	var name string
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
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		line, err := buf.ReadSlice('\n')
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
			inf.MemTotal = int64(i)
			continue
		}
		if name == "MemFree" {
			inf.MemFree = int64(i)
			continue
		}
		if name == "MemAvailable" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "Buffers" {
			inf.Buffers = int64(i)
			continue
		}
		if name == "Cached" {
			inf.MemAvailable = int64(i)
			continue
		}
		if name == "SwapCached" {
			inf.SwapCached = int64(i)
			continue
		}
		if name == "Active" {
			inf.Active = int64(i)
			continue
		}
		if name == "Inactive" {
			inf.Inactive = int64(i)
			continue
		}
		if name == "SwapTotal" {
			inf.SwapTotal = int64(i)
			continue
		}
		if name == "SwapFree" {
			inf.SwapFree = int64(i)
			continue
		}
	}
	return inf, nil
}

func GetMemInfoReadReuseRToJSON() ([]byte, error) {
	inf, err := GetMemInfoReadReuseR()
	if err != nil {
		return nil, err
	}
	return json.Marshal(inf)
}

// GetMemDataReadReuseR returns the current meminfo as flatbuffer serialized bytes.
func GetMemDataReadReuseR() ([]byte, error) {
	inf, err := GetMemInfoReadReuseR()
	if err != nil {
		return nil, err
	}
	return inf.SerializeFlat(), nil
}

// GetMemDataReuseRReuseBldr reuses the Builder.
func GetMemDataReuseRReuseBldr() ([]byte, error) {
	inf, err := GetMemInfoReadReuseR()
	if err != nil {
		return nil, err
	}
	return inf.BldrSerializeFlat(), nil
}

// GetMemInfoToFlatbuffersReuseBldr returns some of the results of /proc/meminfo.
func GetMemInfoToFlatbuffersReuseBldr() ([]byte, error) {
	var l, i int
	var name string
	var v byte
	t := time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	bldr.Reset()
	defer f.Close()
	buf.Reset(f)
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, t)
	var pos int
	val := make([]byte, 0, 32)
	for {
		if l == 16 {
			break
		}
		line, err := buf.ReadSlice('\n')
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
			flat.InfoAddMemTotal(bldr, int64(i))
			continue
		}
		if name == "MemFree" {
			flat.InfoAddMemFree(bldr, int64(i))
			continue
		}
		if name == "MemAvailable" {
			flat.InfoAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "Buffers" {
			flat.InfoAddBuffers(bldr, int64(i))
			continue
		}
		if name == "Cached" {
			flat.InfoAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "SwapCached" {
			flat.InfoAddSwapCached(bldr, int64(i))
			continue
		}
		if name == "Active" {
			flat.InfoAddActive(bldr, int64(i))
			continue
		}
		if name == "Inactive" {
			flat.InfoAddInactive(bldr, int64(i))
			continue
		}
		if name == "SwapTotal" {
			flat.InfoAddSwapTotal(bldr, int64(i))
			continue
		}
		if name == "SwapFree" {
			flat.InfoAddSwapFree(bldr, int64(i))
			continue
		}
	}
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():], nil
}

var l, i, pos int
var v byte
var f *os.File
var err error
var line []byte
var name string
var val = make([]byte, 0, 20)

func GetMemInfoToFlatbuffersMinAllocs() ([]byte, error) {
	f, err = os.Open("/proc/meminfo")
	if err != nil {
		goto fclose
	}
	bldr.Reset()
	buf.Reset(f)
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, time.Now().UTC().UnixNano())

	for {
		if l == 16 {
			break
		}
		line, err = buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			err = fmt.Errorf("error reading output bytes: %s", err)
			goto fclose
		}
		l++
		if l > 8 && l < 15 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				name = string(line[:i])
				pos = i + 1
				break
			}
		}
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
				val = line[pos : pos+i]
				break
			}
		}
		// any conversion error results in 0
		i, err = strconv.Atoi(string(val))
		if err != nil {
			err = fmt.Errorf("%s: %s", name, err)
			goto fclose
		}
		if name == "MemTotal" {
			flat.InfoAddMemTotal(bldr, int64(i))
			continue
		}
		if name == "MemFree" {
			flat.InfoAddMemFree(bldr, int64(i))
			continue
		}
		if name == "MemAvailable" {
			flat.InfoAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "Buffers" {
			flat.InfoAddBuffers(bldr, int64(i))
			continue
		}
		if name == "Cached" {
			flat.InfoAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "SwapCached" {
			flat.InfoAddSwapCached(bldr, int64(i))
			continue
		}
		if name == "Active" {
			flat.InfoAddActive(bldr, int64(i))
			continue
		}
		if name == "Inactive" {
			flat.InfoAddInactive(bldr, int64(i))
			continue
		}
		if name == "SwapTotal" {
			flat.InfoAddSwapTotal(bldr, int64(i))
			continue
		}
		if name == "SwapFree" {
			flat.InfoAddSwapFree(bldr, int64(i))
			continue
		}
	}
fclose:
	f.Close()
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():], err
}

func GetMemInfoToFlatbuffersMinAllocsSeek(f *os.File) ([]byte, error) {
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	bldr.Reset()
	buf.Reset(f)
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, time.Now().UTC().UnixNano())

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
				name = string(line[:i])
				pos = i + 1
				break
			}
		}
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
				val = line[pos : pos+i]
				break
			}
		}
		// any conversion error results in 0
		i, err = strconv.Atoi(string(val))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", name, err)
		}
		if name == "MemTotal" {
			flat.InfoAddMemTotal(bldr, int64(i))
			continue
		}
		if name == "MemFree" {
			flat.InfoAddMemFree(bldr, int64(i))
			continue
		}
		if name == "MemAvailable" {
			flat.InfoAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "Buffers" {
			flat.InfoAddBuffers(bldr, int64(i))
			continue
		}
		if name == "Cached" {
			flat.InfoAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "SwapCached" {
			flat.InfoAddSwapCached(bldr, int64(i))
			continue
		}
		if name == "Active" {
			flat.InfoAddActive(bldr, int64(i))
			continue
		}
		if name == "Inactive" {
			flat.InfoAddInactive(bldr, int64(i))
			continue
		}
		if name == "SwapTotal" {
			flat.InfoAddSwapTotal(bldr, int64(i))
			continue
		}
		if name == "SwapFree" {
			flat.InfoAddSwapFree(bldr, int64(i))
			continue
		}
	}
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():], nil
}

var (
	nameLen int
	n       uint64
)

func GetMemInfoEmulateCurrentFlatTicker(f *os.File) ([]byte, error) {
	// The current timestamp is always in UTC
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	bldr.Reset()
	buf.Reset(f)
	flat.InfoStart(bldr)
	flat.InfoAddTimestamp(bldr, time.Now().UTC().UnixNano())
	for l = 0; l < 16; l++ {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if l > 7 && l < 14 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				val = line[:i]
				pos = i + 1
				break
			}
		}
		nameLen = len(val)

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
		n, err = helpers.ParseUint(val[nameLen:])
		if err != nil {
			return nil, err
		}
		v = val[0]
		if v == 'M' {
			v = val[3]
			if v == 'T' {
				flat.InfoAddMemTotal(bldr, int64(n))
			} else if v == 'F' {
				flat.InfoAddMemFree(bldr, int64(n))
			} else {
				flat.InfoAddMemAvailable(bldr, int64(n))
			}
		} else if v == 'S' {
			v = val[4]
			if v == 'C' {
				flat.InfoAddSwapCached(bldr, int64(n))
			} else if v == 'T' {
				flat.InfoAddSwapTotal(bldr, int64(n))
			} else if v == 'F' {
				flat.InfoAddSwapFree(bldr, int64(n))
			}
		} else if v == 'B' {
			flat.InfoAddBuffers(bldr, int64(n))
		} else if v == 'I' {
			flat.InfoAddInactive(bldr, int64(n))
		} else if v == 'C' {
			flat.InfoAddMemAvailable(bldr, int64(n))
		} else if v == 'A' {
			flat.InfoAddInactive(bldr, int64(n))
		}
	}
	bldr.Finish(flat.InfoEnd(bldr))
	return bldr.Bytes[bldr.Head():], nil
}

func GetMemInfoCurrent(proc *os.File) (inf MemInfo, err error) {
	_, err = proc.Seek(0, os.SEEK_SET)
	if err != nil {
		return inf, err
	}
	buf.Reset(proc)
	var (
		i       int
		v       byte
		pos     int
		nameLen int
	)
	for l := 0; l < 16; l++ {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return inf, fmt.Errorf("error reading output bytes: %s", err)
		}
		if l > 8 && l < 14 {
			continue
		}

		// first grab the key name (everything up to the ':')
		for i, v = range line {
			if v == ':' {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		nameLen = len(val)

		// skip all spaces
		for i, v = range line[pos:] {
			if v != ' ' {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range line[pos:] {
			if v == ' ' || v == '\r' {
				break
			}
			val = append(val, v)
		}
		// any conversion error results in 0
		n, err := helpers.ParseUint(val[nameLen:])
		if err != nil {
			return inf, fmt.Errorf("%s: %s", val[:nameLen], err)
		}

		v = val[0]

		// Reduce evaluations.
		if v == 'M' {
			v = val[3]
			if v == 'T' {
				inf.MemTotal = int64(n)
			} else if v == 'F' {
				inf.MemFree = int64(n)
			} else {
				inf.MemAvailable = int64(n)
			}
		} else if v == 'S' {
			v = val[4]
			if v == 'C' {
				inf.SwapCached = int64(n)
			} else if v == 'T' {
				inf.SwapTotal = int64(n)
			} else if v == 'F' {
				inf.SwapFree = int64(n)
			}
		} else if v == 'B' {
			inf.Buffers = int64(n)
		} else if v == 'I' {
			inf.Inactive = int64(n)
		} else if v == 'C' {
			inf.Cached = int64(n)
		} else if v == 'A' {
			inf.Active = int64(n)
		}
		val = val[:0]
	}
	inf.Timestamp = time.Now().UTC().UnixNano()
	return inf, nil
}
