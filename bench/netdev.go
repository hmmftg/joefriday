package bench

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/mohae/joefriday/net"
)

// predeclare some vars
var iData net.Iface
var fieldNum, fieldVal int
var bs []byte
var inf = &net.Info{Interfaces: make([]net.Iface, 0, 4)}
var t int64

func init() {
	tmp := bytes.NewBuffer(bs)
	buf = bufio.NewReaderSize(tmp, 4096)
	tmp = nil
}

func EmulateNetDevDataTicker() ([]byte, error) {
	t = time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf.Reset(f)
	inf.Interfaces = inf.Interfaces[:0]
	// there's always at least 2 interfaces (I think)
	inf.Timestamp = t
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("/proc/mem/dev: read output bytes: %s", err)
		}
		l++
		if l < 3 {
			continue
		}

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
				return nil, fmt.Errorf("/proc/net/dev ticker: %s: %s", iData.Name, err)
			}
			val = val[:0]
			if fieldNum == 1 {
				iData.RBytes = int64(fieldVal)
				continue
			}
			if fieldNum == 2 {
				iData.RPackets = int64(fieldVal)
				continue
			}
			if fieldNum == 3 {
				iData.RErrs = int64(fieldVal)
				continue
			}
			if fieldNum == 4 {
				iData.RDrop = int64(fieldVal)
				continue
			}
			if fieldNum == 5 {
				iData.RFIFO = int64(fieldVal)
				continue
			}
			if fieldNum == 6 {
				iData.RFrame = int64(fieldVal)
				continue
			}
			if fieldNum == 7 {
				iData.RCompressed = int64(fieldVal)
				continue
			}
			if fieldNum == 8 {
				iData.RMulticast = int64(fieldVal)
				continue
			}
			if fieldNum == 9 {
				iData.TBytes = int64(fieldVal)
				continue
			}
			if fieldNum == 10 {
				iData.TPackets = int64(fieldVal)
				continue
			}
			if fieldNum == 11 {
				iData.TErrs = int64(fieldVal)
				continue
			}
			if fieldNum == 12 {
				iData.TDrop = int64(fieldVal)
				continue
			}
			if fieldNum == 13 {
				iData.TFIFO = int64(fieldVal)
				continue
			}
			if fieldNum == 14 {
				iData.TColls = int64(fieldVal)
				continue
			}
			if fieldNum == 15 {
				iData.TCarrier = int64(fieldVal)
				continue
			}
			if fieldNum == 16 {
				iData.TCompressed = int64(fieldVal)
				break
			}
		}
		inf.Interfaces = append(inf.Interfaces, iData)
	}
	f.Close()
	l = 0
	data := net.Serialize(inf, bldr)
	return data, nil
}
