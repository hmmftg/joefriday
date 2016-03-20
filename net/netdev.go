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

// Package net gets and processes /proc/net/dev, returning the data in the
// appropriate format.
package net

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
)

type Info struct {
	Timestamp  int64
	Interfaces []Iface
}

// Iface: contains information for a given network interface; names as
// such to prevent collision with the Flatbuffers struct.
type Iface struct {
	Name string
	RCum Received
	TCum Transmitted
}

// Received: data related to receive; named as such to prevent collision
// with the Flatbuffers struct.
type Received struct {
	Bytes      int64
	Packets    int64
	Errs       int64
	Drop       int64
	FIFO       int64
	Frame      int64
	Compressed int64
	Multicast  int64
}

// Transmitted: data related to transmit; named as such to prevent collision
// with the Flatbuffers struct.
type Transmitted struct {
	Bytes      int64
	Packets    int64
	Errs       int64
	Drop       int64
	FIFO       int64
	Colls      int64
	Carrier    int64
	Compressed int64
}

// Serialize serializes the Info using flatbuffers.
func (i *Info) Serialize() []byte {
	ifaces := make([]fb.UOffsetT, 0, len(i.Interfaces))

	bldr := fb.NewBuilder(0)
	for _, v := range i.Interfaces {
		ReceiveStart(bldr)
		ReceiveAddBytes(bldr, v.RCum.Bytes)
		ReceiveAddPackets(bldr, v.RCum.Packets)
		ReceiveAddErrs(bldr, v.RCum.Errs)
		ReceiveAddDrop(bldr, v.RCum.Drop)
		ReceiveAddFIFO(bldr, v.RCum.FIFO)
		ReceiveAddFrame(bldr, v.RCum.Frame)
		ReceiveAddCompressed(bldr, v.RCum.Compressed)
		ReceiveAddMulticast(bldr, v.RCum.Multicast)
		r := ReceiveEnd(bldr)
		TransmitStart(bldr)
		TransmitAddBytes(bldr, v.TCum.Bytes)
		TransmitAddPackets(bldr, v.TCum.Packets)
		TransmitAddErrs(bldr, v.TCum.Errs)
		TransmitAddDrop(bldr, v.TCum.Drop)
		TransmitAddFIFO(bldr, v.TCum.FIFO)
		TransmitAddColls(bldr, v.TCum.Colls)
		TransmitAddCarrier(bldr, v.TCum.Carrier)
		TransmitAddCompressed(bldr, v.TCum.Compressed)
		t := ReceiveEnd(bldr)
		n := bldr.CreateString(v.Name)
		IFaceStart(bldr)
		IFaceAddName(bldr, n)
		IFaceAddRCum(bldr, r)
		IFaceAddTCum(bldr, t)
		ifaces = append(ifaces, IFaceEnd(bldr))
	}
	DataStartInterfacesVector(bldr, len(ifaces))
	for i := len(i.Interfaces) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(ifaces[i])
	}
	bldr.EndVector(len(ifaces))
	DataStart(bldr)
	DataAddTimestamp(bldr, i.Timestamp)
	bldr.Finish(DataEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// Deserialize deserializes bytes representing flatbuffers serialized Data
// into *Info.  If the bytes are not from flatbuffers serialization of
// Data, it is a programmer error and a panic will occur.
func Deserialize(p []byte) *Info {
	data := GetRootAsData(p, 0)
	// get the # of interfaces
	iLen := data.InterfacesLength()
	info := &Info{Timestamp: data.Timestamp(), Interfaces: make([]Iface, 0, iLen)}
	iFace := &IFace{}
	iface := Iface{}
	r := &Receive{}
	t := &Transmit{}
	for i := 0; i < iLen; i++ {
		if data.Interfaces(iFace, i) {
			r = iFace.RCum(r)
			iface.RCum.Bytes = r.Bytes()
			iface.RCum.Packets = r.Packets()
			iface.RCum.Errs = r.Errs()
			iface.RCum.Drop = r.Drop()
			iface.RCum.FIFO = r.FIFO()
			iface.RCum.Frame = r.Frame()
			iface.RCum.Compressed = r.Compressed()
			iface.RCum.Multicast = r.Multicast()
			t = iFace.TCum(t)
			iface.TCum.Bytes = t.Bytes()
			iface.TCum.Packets = t.Packets()
			iface.TCum.Errs = t.Errs()
			iface.TCum.Drop = t.Drop()
			iface.TCum.FIFO = t.FIFO()
			iface.TCum.Colls = t.Colls()
			iface.TCum.Carrier = t.Carrier()
			iface.TCum.Compressed = t.Compressed()
		}
		info.Interfaces = append(info.Interfaces, iface)
	}
	return info
}

func (i Info) String() string {
	var buf bytes.Buffer
	buf.WriteString(time.Unix(0, i.Timestamp).UTC().String())
	buf.WriteByte('\n')
	for _, v := range i.Interfaces {
		buf.WriteString(joe.Column(8, v.Name))
		buf.WriteString(joe.Int64Column(22, v.RCum.Bytes))
		buf.WriteString(joe.Int64Column(22, v.RCum.Packets))
		buf.WriteString(joe.Int64Column(22, v.RCum.Errs))
		buf.WriteString(joe.Int64Column(22, v.RCum.Drop))
		buf.WriteString(joe.Int64Column(22, v.RCum.FIFO))
		buf.WriteString(joe.Int64Column(22, v.RCum.Frame))
		buf.WriteString(joe.Int64Column(22, v.RCum.Compressed))
		buf.WriteString(joe.Int64Column(22, v.RCum.Multicast))
		buf.WriteString(joe.Int64Column(22, v.TCum.Bytes))
		buf.WriteString(joe.Int64Column(22, v.TCum.Packets))
		buf.WriteString(joe.Int64Column(22, v.TCum.Errs))
		buf.WriteString(joe.Int64Column(22, v.TCum.Drop))
		buf.WriteString(joe.Int64Column(22, v.TCum.FIFO))
		buf.WriteString(joe.Int64Column(22, v.TCum.Colls))
		buf.WriteString(joe.Int64Column(22, v.TCum.Carrier))
		buf.WriteString(joe.Int64Column(22, v.TCum.Compressed))
		buf.WriteByte('\n')
	}
	return buf.String()
}

// GetInfo returns some of the results of /proc/meminfo.
func GetInfo() (*Info, error) {
	var l, i, pos, fieldNum, fieldVal int
	var v byte
	t := time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	// there's always at least 2 interfaces (I think)
	inf := &Info{Timestamp: t, Interfaces: make([]Iface, 0, 2)}
	val := make([]byte, 0, 32)
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l < 3 {
			continue
		}
		var iData Iface

		// first grab the interface name (everything up to the ':')
		for i, v = range line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		iData.Name = string(val[:])
		val = val[:0]
		fieldNum = 0
		// process the rest of the line
		for {
			fieldNum++
			// skip all spaces
			for i, v = range line[pos:] {
				if v != 0x20 {
					pos += i
					break
				}
			}

			// grab the numbers
			for i, v = range line[pos:] {
				if v == 0x20 || v == '\n' {
					pos += i
					break
				}
				val = append(val, v)
			}
			// any conversion error results in 0
			fieldVal, err = strconv.Atoi(string(val[:]))
			if err != nil {
				return nil, fmt.Errorf("%s: %s", iData.Name, err)
			}
			val = val[:0]
			if fieldNum == 1 {
				iData.RCum.Bytes = int64(fieldVal)
				continue
			}
			if fieldNum == 2 {
				iData.RCum.Packets = int64(fieldVal)
				continue
			}
			if fieldNum == 3 {
				iData.RCum.Errs = int64(fieldVal)
				continue
			}
			if fieldNum == 4 {
				iData.RCum.Drop = int64(fieldVal)
				continue
			}
			if fieldNum == 5 {
				iData.RCum.FIFO = int64(fieldVal)
				continue
			}
			if fieldNum == 6 {
				iData.RCum.Frame = int64(fieldVal)
				continue
			}
			if fieldNum == 7 {
				iData.RCum.Compressed = int64(fieldVal)
				continue
			}
			if fieldNum == 8 {
				iData.RCum.Multicast = int64(fieldVal)
				continue
			}
			if fieldNum == 9 {
				iData.RCum.Bytes = int64(fieldVal)
				continue
			}
			if fieldNum == 10 {
				iData.RCum.Packets = int64(fieldVal)
				continue
			}
			if fieldNum == 11 {
				iData.TCum.Errs = int64(fieldVal)
				continue
			}
			if fieldNum == 12 {
				iData.TCum.Drop = int64(fieldVal)
				continue
			}
			if fieldNum == 13 {
				iData.TCum.FIFO = int64(fieldVal)
				continue
			}
			if fieldNum == 14 {
				iData.TCum.Colls = int64(fieldVal)
				continue
			}
			if fieldNum == 15 {
				iData.TCum.Carrier = int64(fieldVal)
				continue
			}
			if fieldNum == 16 {
				iData.TCum.Compressed = int64(fieldVal)
				break
			}
		}
		inf.Interfaces = append(inf.Interfaces, iData)
	}
	return inf, nil
}

// GetData returns the current meminfo as flatbuffer serialized bytes.
func GetData() ([]byte, error) {
	inf, err := GetInfo()
	if err != nil {
		return nil, err
	}
	return inf.Serialize(), nil
}

/*
// DataTicker gathers the meminfo on a ticker, whose interval is defined by
// the received duration, and sends the results to the channel.  The output
// is Flatbuffers serialized Data.  Any error encountered during processing
// is sent to the error channel.  Processing will continue
//
// Either closing the done channel or sending struct{} to the done channel
// will result in function exit.  The out channel is closed on exit.
//
// This pre-allocates the builder and everything other than the []byte that
// gets sent to the out channel to reduce allocations, as this is expected
// to be both a frequent and a long-running process.  Doing so reduces
// byte allocations per tick just ~ 42%.
func DataTicker(interval time.Duration, outCh chan []byte, done chan struct{}, errCh chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(outCh)
	// predeclare some vars
	var l, i, pos int
	var t int64
	var v byte
	var name string
	// premake some temp slices
	val := make([]byte, 0, 32)
	// just reset the bldr at the end of every ticker
	bldr := fb.NewBuilder(0)
	// Some hopes to jump through to ensure we don't get a ErrBufferFull; which was
	// occuring with var buf bufio.Reader (which works in the bench code)
	var bs []byte
	tmp := bytes.NewBuffer(bs)
	buf := bufio.NewReaderSize(tmp, 1536)
	tmp = nil
	// ticker
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// The current timestamp is always in UTC
			t = time.Now().UTC().UnixNano()
			f, err := os.Open("/proc/meminfo")
			if err != nil {
				errCh <- joe.Error{Type: "mem", Op: "open /proc/meminfo", Err: err}
				continue
			}
			buf.Reset(f)
			DataStart(bldr)
			DataAddTimestamp(bldr, t)
			for {
				if l == 16 {
					break
				}
				line, err := buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errCh <- joe.Error{Type: "mem", Op: "read command results", Err: err}
					break
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
					errCh <- joe.Error{Type: "mem", Op: "convert to int", Err: err}
					continue
				}
				val = val[:0]
				if name == "MemTotal" {
					DataAddMemTotal(bldr, int64(i))
					continue
				}
				if name == "MemFree" {
					DataAddMemFree(bldr, int64(i))
					continue
				}
				if name == "MemAvailable" {
					DataAddMemAvailable(bldr, int64(i))
					continue
				}
				if name == "Buffers" {
					DataAddBuffers(bldr, int64(i))
					continue
				}
				if name == "Cached" {
					DataAddMemAvailable(bldr, int64(i))
					continue
				}
				if name == "SwapCached" {
					DataAddSwapCached(bldr, int64(i))
					continue
				}
				if name == "Active" {
					DataAddActive(bldr, int64(i))
					continue
				}
				if name == "Inactive" {
					DataAddInactive(bldr, int64(i))
					continue
				}
				if name == "SwapTotal" {
					DataAddSwapTotal(bldr, int64(i))
					continue
				}
				if name == "SwapFree" {
					DataAddSwapFree(bldr, int64(i))
					continue
				}
			}
			f.Close()
			bldr.Finish(DataEnd(bldr))
			data := bldr.Bytes[bldr.Head():]
			outCh <- data
			bldr.Reset()
			l = 0
		}
	}
}

func (d *Data) String() string {
	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, d.Timestamp()).UTC(), d.MemTotal(), d.MemFree(), d.MemAvailable(), d.Active(), d.Inactive(), d.Cached(), d.Buffers(), d.SwapTotal(), d.SwapCached(), d.SwapFree())
}
*/
