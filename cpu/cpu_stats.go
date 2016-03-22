package cpu

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
)

var CLK_TCK int // the ticks per clock cycle

// Init: set's the CLK_TCK.
func Init() error {
	var out bytes.Buffer
	cmd := exec.Command("getconf", "CLK_TCK")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return joe.Error{Type: "cpu", Op: "get conf CLK_TCK", Err: err}
	}
	b, err := out.ReadBytes('\n')
	if err != nil {
		return joe.Error{Type: "cpu", Op: "read conf CLK_TCK output", Err: err}
	}
	CLK_TCK, err = strconv.Atoi(string(b[:len(b)-1]))
	if err != nil {
		return joe.Error{Type: "cpu", Op: "processing conf CLK_TCK output", Err: err}
	}
	return nil
}

type Stats struct {
	ClkTck    int16  `json:"clk_tck"`
	Timestamp int64  `json:"timestamp"`
	Ctxt      int64  `json:"ctxt"`
	BTime     int64  `json:"btime"`
	Processes int64  `json:"processes"`
	CPUs      []Stat `json:"cpus"`
}

// Stat is for capturing the output of /proc/stat.
type Stat struct {
	CPU       string `json:"CPU"`
	User      int64  `json:"user"`
	Nice      int64  `json:"nice"`
	System    int64  `json:"system"`
	Idle      int64  `json:"idle"`
	IOWait    int64  `json:"io_wait"`
	IRQ       int64  `json:"irq"`
	SoftIRQ   int64  `json:"soft_irq"`
	Steal     int64  `json:"steal"`
	Quest     int64  `json:"quest"`
	QuestNice int64  `json:"quest_nice"`
}

// SerializeFlat serializes Stats into Flatbuffer serialized bytes.
func (s Stats) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	stats := make([]fb.UOffsetT, len(s.CPUs))
	cpus := make([]fb.UOffsetT, len(s.CPUs))
	for i := 0; i < len(cpus); i++ {
		cpus[i] = bldr.CreateString(s.CPUs[i].CPU)
	}
	for i := 0; i < len(stats); i++ {
		StatFlatStart(bldr)
		StatFlatAddCPU(bldr, cpus[i])
		StatFlatAddUser(bldr, s.CPUs[i].User)
		StatFlatAddNice(bldr, s.CPUs[i].Nice)
		StatFlatAddSystem(bldr, s.CPUs[i].System)
		StatFlatAddIdle(bldr, s.CPUs[i].Idle)
		StatFlatAddIOWait(bldr, s.CPUs[i].IOWait)
		StatFlatAddIRQ(bldr, s.CPUs[i].IRQ)
		StatFlatAddSoftIRQ(bldr, s.CPUs[i].SoftIRQ)
		StatFlatAddSteal(bldr, s.CPUs[i].Steal)
		StatFlatAddQuest(bldr, s.CPUs[i].Quest)
		StatFlatAddQuestNice(bldr, s.CPUs[i].QuestNice)
		stats[i] = StatFlatEnd(bldr)
	}
	StatsFlatStartCPUsVector(bldr, len(stats))
	for i := len(stats) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(stats[i])
	}
	statsV := bldr.EndVector(len(stats))
	StatsFlatStart(bldr)
	StatsFlatAddClkTck(bldr, s.ClkTck)
	StatsFlatAddTimestamp(bldr, s.Timestamp)
	StatsFlatAddCtxt(bldr, s.Ctxt)
	StatsFlatAddBTime(bldr, s.BTime)
	StatsFlatAddProcesses(bldr, s.Processes)
	StatsFlatAddCPUs(bldr, statsV)
	bldr.Finish(StatsFlatEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeStatsFlat deserializes Flatbuffer serialized bytes into Stats.
func DeserializeStatsFlat(p []byte) Stats {
	var stats Stats
	statF := &StatFlat{}
	data := GetRootAsStatsFlat(p, 0)
	stats.ClkTck = data.ClkTck()
	stats.Timestamp = data.Timestamp()
	stats.Ctxt = data.Ctxt()
	stats.BTime = data.BTime()
	stats.Processes = data.Processes()
	len := data.CPUsLength()
	stats.CPUs = make([]Stat, len)
	for i := 0; i < len; i++ {
		var stat Stat
		if data.CPUs(statF, i) {
			stat.CPU = string(statF.CPU())
			stat.User = statF.User()
			stat.Nice = statF.Nice()
			stat.System = statF.System()
			stat.Idle = statF.Idle()
			stat.IOWait = statF.IOWait()
			stat.IRQ = statF.IRQ()
			stat.SoftIRQ = statF.SoftIRQ()
			stat.Steal = statF.Steal()
			stat.Quest = statF.Quest()
			stat.QuestNice = statF.QuestNice()
		}
		stats.CPUs[i] = stat
	}
	return stats
}

// GetStats gets the output of /proc/stat.
func GetStats() (Stats, error) {
	stats := Stats{Timestamp: time.Now().UTC().UnixNano(), CPUs: []Stat{}}
	f, err := os.Open("/proc/stat")
	if err != nil {
		return stats, err
	}
	defer f.Close()

	var name string
	var i, j, pos, val, fieldNum int
	var v byte
	var stop bool

	buf := bufio.NewReader(f)
	// read each line until eof
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return stats, joe.Error{Type: "cpu stat", Op: "reading /proc/stat output", Err: err}
		}
		// Get everything up to the first space, this is the key.  Not all keys are processed.
		for i, v = range line {
			if v == 0x20 {
				name = string(line[:i])
				pos = i + 1
				break
			}
		}
		// skip the intr line
		if name == "intr" {
			continue
		}
		if name[:3] == "cpu" {
			j = 0
			// skip over any remaining spaces
			for i, v = range line[pos:] {
				if v != 0x20 {
					break
				}
				j++
			}
			stat := Stat{CPU: name}
			fieldNum = 0
			pos, j = j+pos, j+pos
			// space is the field separator
			for i, v = range line[pos:] {
				if v == '\n' {
					stop = true
				}
				if v == 0x20 || stop {
					fieldNum++
					val, err = strconv.Atoi(string(line[j : pos+i]))
					if err != nil {
						return stats, joe.Error{Type: "cpu stat", Op: "convert cpu data", Err: err}
					}
					j = pos + i + 1
					if fieldNum == 1 {
						stat.User = int64(val)
						continue
					}
					if fieldNum == 2 {
						stat.Nice = int64(val)
						continue
					}
					if fieldNum == 3 {
						stat.System = int64(val)
						continue
					}
					if fieldNum == 4 {
						stat.Idle = int64(val)
						continue
					}
					if fieldNum == 5 {
						stat.IOWait = int64(val)
						continue
					}
					if fieldNum == 6 {
						stat.IRQ = int64(val)
						continue
					}
					if fieldNum == 7 {
						stat.SoftIRQ = int64(val)
						continue
					}
					if fieldNum == 8 {
						stat.Steal = int64(val)
						continue
					}
					if fieldNum == 9 {
						stat.Quest = int64(val)
						continue
					}
					if fieldNum == 10 {
						stat.QuestNice = int64(val)
						continue
					}
				}
			}
			stats.CPUs = append(stats.CPUs, stat)
			stop = false
			continue
		}
		if name == "ctxt" {
			// rest of the line is the data
			val, err = strconv.Atoi(string(line[pos : len(line)-1]))
			if err != nil {
				return stats, joe.Error{Type: "cpu stat", Op: "convert ctxt data", Err: err}
			}
			stats.Ctxt = int64(val)
			continue
		}
		if name == "btime" {
			// rest of the line is the data
			val, err = strconv.Atoi(string(line[pos : len(line)-1]))
			if err != nil {
				return stats, joe.Error{Type: "cpu stat", Op: "convert btime data", Err: err}
			}
			stats.BTime = int64(val)
			continue
		}
		if name == "processes" {
			// rest of the line is the data
			val, err = strconv.Atoi(string(line[pos : len(line)-1]))
			if err != nil {
				return stats, joe.Error{Type: "cpu stat", Op: "convert processes data", Err: err}
			}
			stats.Processes = int64(val)
			continue
		}
	}
	return stats, nil
}
