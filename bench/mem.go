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

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/mem"
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
	mem.InfoFlatStart(bldrL)
	mem.InfoFlatAddTimestamp(bldrL, int64(i.Timestamp))
	mem.InfoFlatAddMemTotal(bldrL, int64(i.MemTotal))
	mem.InfoFlatAddMemFree(bldrL, int64(i.MemFree))
	mem.InfoFlatAddMemAvailable(bldrL, int64(i.MemAvailable))
	mem.InfoFlatAddBuffers(bldrL, int64(i.Buffers))
	mem.InfoFlatAddCached(bldrL, int64(i.Cached))
	mem.InfoFlatAddSwapCached(bldrL, int64(i.SwapCached))
	mem.InfoFlatAddActive(bldrL, int64(i.Active))
	mem.InfoFlatAddInactive(bldrL, int64(i.Inactive))
	mem.InfoFlatAddSwapTotal(bldrL, int64(i.SwapTotal))
	mem.InfoFlatAddSwapFree(bldrL, int64(i.SwapFree))
	bldrL.Finish(mem.InfoFlatEnd(bldrL))
	return bldrL.Bytes[bldrL.Head():]
}

// BldrSerialize serializes the MemInfo using flatbuffers: the builder is
// reused.
func (i *MemInfo) BldrSerializeFlat() []byte {
	bldr.Reset()
	mem.InfoFlatStart(bldr)
	mem.InfoFlatAddTimestamp(bldr, int64(i.Timestamp))
	mem.InfoFlatAddMemTotal(bldr, int64(i.MemTotal))
	mem.InfoFlatAddMemFree(bldr, int64(i.MemFree))
	mem.InfoFlatAddMemAvailable(bldr, int64(i.MemAvailable))
	mem.InfoFlatAddBuffers(bldr, int64(i.Buffers))
	mem.InfoFlatAddCached(bldr, int64(i.Cached))
	mem.InfoFlatAddSwapCached(bldr, int64(i.SwapCached))
	mem.InfoFlatAddActive(bldr, int64(i.Active))
	mem.InfoFlatAddInactive(bldr, int64(i.Inactive))
	mem.InfoFlatAddSwapTotal(bldr, int64(i.SwapTotal))
	mem.InfoFlatAddSwapFree(bldr, int64(i.SwapFree))
	bldr.Finish(mem.InfoFlatEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeFlat deserializes bytes representing flatbuffers serialized
// InfoFlat into *Info.
func DeserializeFlat(p []byte) *MemInfo {
	infoFlat := mem.GetRootAsInfoFlat(p, 0)
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
	mem.InfoFlatStart(bldr)
	mem.InfoFlatAddTimestamp(bldr, t)
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
			mem.InfoFlatAddMemTotal(bldr, int64(i))
			continue
		}
		if name == "MemFree" {
			mem.InfoFlatAddMemFree(bldr, int64(i))
			continue
		}
		if name == "MemAvailable" {
			mem.InfoFlatAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "Buffers" {
			mem.InfoFlatAddBuffers(bldr, int64(i))
			continue
		}
		if name == "Cached" {
			mem.InfoFlatAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "SwapCached" {
			mem.InfoFlatAddSwapCached(bldr, int64(i))
			continue
		}
		if name == "Active" {
			mem.InfoFlatAddActive(bldr, int64(i))
			continue
		}
		if name == "Inactive" {
			mem.InfoFlatAddInactive(bldr, int64(i))
			continue
		}
		if name == "SwapTotal" {
			mem.InfoFlatAddSwapTotal(bldr, int64(i))
			continue
		}
		if name == "SwapFree" {
			mem.InfoFlatAddSwapFree(bldr, int64(i))
			continue
		}
	}
	bldr.Finish(mem.InfoFlatEnd(bldr))
	return bldr.Bytes[bldr.Head():], nil
}

var l, i, pos int
var name string
var v byte
var t int64
var val = make([]byte, 0, 32)

func GetMemInfoToFlatbuffersMinAllocs() ([]byte, error) {
	t = time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	bldr.Reset()
	defer f.Close()
	buf.Reset(f)
	mem.InfoFlatStart(bldr)
	mem.InfoFlatAddTimestamp(bldr, t)

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
			mem.InfoFlatAddMemTotal(bldr, int64(i))
			continue
		}
		if name == "MemFree" {
			mem.InfoFlatAddMemFree(bldr, int64(i))
			continue
		}
		if name == "MemAvailable" {
			mem.InfoFlatAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "Buffers" {
			mem.InfoFlatAddBuffers(bldr, int64(i))
			continue
		}
		if name == "Cached" {
			mem.InfoFlatAddMemAvailable(bldr, int64(i))
			continue
		}
		if name == "SwapCached" {
			mem.InfoFlatAddSwapCached(bldr, int64(i))
			continue
		}
		if name == "Active" {
			mem.InfoFlatAddActive(bldr, int64(i))
			continue
		}
		if name == "Inactive" {
			mem.InfoFlatAddInactive(bldr, int64(i))
			continue
		}
		if name == "SwapTotal" {
			mem.InfoFlatAddSwapTotal(bldr, int64(i))
			continue
		}
		if name == "SwapFree" {
			mem.InfoFlatAddSwapFree(bldr, int64(i))
			continue
		}
	}
	bldr.Finish(mem.InfoFlatEnd(bldr))
	return bldr.Bytes[bldr.Head():], nil
}
