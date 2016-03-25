package cpu

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/cpu/flat"
)

const procStat = "/proc/stat"

var CLK_TCK int16 // the ticks per clock cycle

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
	v, err := strconv.Atoi(string(b[:len(b)-1]))
	if err != nil {
		return joe.Error{Type: "cpu", Op: "processing conf CLK_TCK output", Err: err}
	}
	CLK_TCK = int16(v)
	return nil
}

// Stats holds the /proc/stat information
type Stats struct {
	ClkTck    int16  `json:"clk_tck"`
	Timestamp int64  `json:"timestamp"`
	Ctxt      int64  `json:"ctxt"`
	BTime     int64  `json:"btime"`
	Processes int64  `json:"processes"`
	CPU       []Stat `json:"cpu"`
}

// Stat is for capturing the CPU information of /proc/stat.
type Stat struct {
	ID        string `json:"ID"`
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
func (s *Stats) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return s.SerializeFlatBuilder(bldr)
}

// SerializeFlat serializes Stats into Flatbuffer serialized bytes.
func (s *Stats) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	stats := make([]fb.UOffsetT, len(s.CPU))
	ids := make([]fb.UOffsetT, len(s.CPU))
	for i := 0; i < len(ids); i++ {
		ids[i] = bldr.CreateString(s.CPU[i].ID)
	}
	for i := 0; i < len(stats); i++ {
		flat.StatStart(bldr)
		flat.StatAddID(bldr, ids[i])
		flat.StatAddUser(bldr, s.CPU[i].User)
		flat.StatAddNice(bldr, s.CPU[i].Nice)
		flat.StatAddSystem(bldr, s.CPU[i].System)
		flat.StatAddIdle(bldr, s.CPU[i].Idle)
		flat.StatAddIOWait(bldr, s.CPU[i].IOWait)
		flat.StatAddIRQ(bldr, s.CPU[i].IRQ)
		flat.StatAddSoftIRQ(bldr, s.CPU[i].SoftIRQ)
		flat.StatAddSteal(bldr, s.CPU[i].Steal)
		flat.StatAddQuest(bldr, s.CPU[i].Quest)
		flat.StatAddQuestNice(bldr, s.CPU[i].QuestNice)
		stats[i] = flat.StatEnd(bldr)
	}
	flat.StatsStartCPUVector(bldr, len(stats))
	for i := len(stats) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(stats[i])
	}
	statsV := bldr.EndVector(len(stats))
	flat.StatsStart(bldr)
	flat.StatsAddClkTck(bldr, s.ClkTck)
	flat.StatsAddTimestamp(bldr, s.Timestamp)
	flat.StatsAddCtxt(bldr, s.Ctxt)
	flat.StatsAddBTime(bldr, s.BTime)
	flat.StatsAddProcesses(bldr, s.Processes)
	flat.StatsAddCPU(bldr, statsV)
	bldr.Finish(flat.StatsEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeStatsFlat deserializes Flatbuffer serialized bytes into Stats.
func DeserializeStatsFlat(p []byte) Stats {
	var stats Stats
	statF := &flat.Stat{}
	statsFlat := flat.GetRootAsStats(p, 0)
	stats.ClkTck = statsFlat.ClkTck()
	stats.Timestamp = statsFlat.Timestamp()
	stats.Ctxt = statsFlat.Ctxt()
	stats.BTime = statsFlat.BTime()
	stats.Processes = statsFlat.Processes()
	len := statsFlat.CPULength()
	stats.CPU = make([]Stat, len)
	for i := 0; i < len; i++ {
		var stat Stat
		if statsFlat.CPU(statF, i) {
			stat.ID = string(statF.ID())
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
		stats.CPU[i] = stat
	}
	return stats
}

func init() {
	var err error
	proc, err = os.Open(procStat)
	if err != nil {
		log.Fatalln(err)
	}
	buf = bufio.NewReader(proc)
}

var proc *os.File
var buf *bufio.Reader

// GetStats gets the output of /proc/stat.
func GetStats() (Stats, error) {
	var stats Stats
	if CLK_TCK == 0 {
		err := Init()
		if err != nil {
			return stats, err
		}
	}
	_, err := proc.Seek(0, os.SEEK_SET)
	if err != nil {
		return stats, err
	}
	buf.Reset(proc)

	stats.ClkTck = CLK_TCK
	stats.CPU = make([]Stat, 0, 2)
	stats.Timestamp = time.Now().UTC().UnixNano()

	var (
		name                     string
		i, j, pos, val, fieldNum int
		v                        byte
		stop                     bool
	)

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
			stat := Stat{ID: name}
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
					if fieldNum < 4 {
						if fieldNum == 1 {
							stat.User = int64(val)
						} else if fieldNum == 2 {
							stat.Nice = int64(val)
						} else if fieldNum == 3 {
							stat.System = int64(val)
						}
					} else if fieldNum < 7 {
						if fieldNum == 4 {
							stat.Idle = int64(val)
						} else if fieldNum == 5 {
							stat.IOWait = int64(val)
						} else if fieldNum == 6 {
							stat.IRQ = int64(val)
						}
					} else if fieldNum < 10 {
						if fieldNum == 7 {
							stat.SoftIRQ = int64(val)
						} else if fieldNum == 8 {
							stat.Steal = int64(val)
						} else if fieldNum == 9 {
							stat.Quest = int64(val)
						}
					} else if fieldNum == 10 {
						stat.QuestNice = int64(val)
					}
				}
			}
			stats.CPU = append(stats.CPU, stat)
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

// GetStatsFlat gets /proc/stat as Flatbuffer serialized bytes.
func GetStatsFlat() ([]byte, error) {
	s, err := GetStats()
	if err != nil {
		return nil, err
	}
	return s.SerializeFlat(), nil
}

// Utilization holds information about cpu utilization.
type Utilization struct {
	Timestamp int64 `json:"timestamp"`
	// time since last reboot, in seconds
	BTimeDelta int32 `json:"btime_delta"`
	// context switches since last snapshot
	CtxtDelta int64 `json:"ctxt_delta"`
	// current number of Processes
	Processes int32 `json:"processes"`
	// cpu specific utilization information
	CPU []Util `json:"cpu"`
}

// Util holds utilization information for a CPU.
type Util struct {
	ID        string  `json:"id"`
	Usage     float32 `json:"total"`
	User      float32 `json:"user"`
	Nice      float32 `json:"nice"`
	System    float32 `json:"system"`
	Idle      float32 `json:"idle"`
	IOWait    float32 `json:"io_wait"`
	IRQ       float32 `json:"irq"`
	SoftIRQ   float32 `json:"soft_irq"`
	Steal     float32 `json:"steal"`
	Quest     float32 `json:"quest"`
	QuestNice float32 `json:"quest_nice"`
}

// SerializeFlat serializes Utilization into Flatbuffer serialized bytes.
func (u *Utilization) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return u.SerializeFlatBuilder(bldr)
}

// SerializeFlatBuilder serializes Utilization into Flatbuffer serialized
// bytes using the received builder.  It is assumed that the passed builder
// is in a usable state.
func (u *Utilization) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	utils := make([]fb.UOffsetT, len(u.CPU))
	ids := make([]fb.UOffsetT, len(u.CPU))
	for i := 0; i < len(ids); i++ {
		ids[i] = bldr.CreateString(u.CPU[i].ID)
	}
	for i := 0; i < len(utils); i++ {
		flat.UtilStart(bldr)
		flat.UtilAddID(bldr, ids[i])
		flat.UtilAddUsage(bldr, u.CPU[i].Usage)
		flat.UtilAddUser(bldr, u.CPU[i].User)
		flat.UtilAddNice(bldr, u.CPU[i].Nice)
		flat.UtilAddSystem(bldr, u.CPU[i].System)
		flat.UtilAddIdle(bldr, u.CPU[i].Idle)
		flat.UtilAddIOWait(bldr, u.CPU[i].IOWait)
		utils[i] = flat.UtilEnd(bldr)
	}
	flat.UtilizationStartCPUVector(bldr, len(utils))
	for i := len(utils) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(utils[i])
	}
	utilsV := bldr.EndVector(len(utils))
	flat.UtilizationStart(bldr)
	flat.UtilizationAddTimestamp(bldr, u.Timestamp)
	flat.UtilizationAddBTimeDelta(bldr, u.BTimeDelta)
	flat.UtilizationAddCtxtDelta(bldr, u.CtxtDelta)
	flat.UtilizationAddProcesses(bldr, u.Processes)
	flat.UtilizationAddCPU(bldr, utilsV)
	bldr.Finish(flat.UtilizationEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeUtilizationFlat deserializes Flatbuffer serialized bytes into
// Utilization.
func DeserializeUtilizationFlat(p []byte) Utilization {
	var u Utilization
	uF := &flat.Util{}
	flatUtil := flat.GetRootAsUtilization(p, 0)
	u.Timestamp = flatUtil.Timestamp()
	u.CtxtDelta = flatUtil.CtxtDelta()
	u.BTimeDelta = flatUtil.BTimeDelta()
	u.Processes = flatUtil.Processes()
	len := flatUtil.CPULength()
	u.CPU = make([]Util, len)
	for i := 0; i < len; i++ {
		var util Util
		if flatUtil.CPU(uF, i) {
			util.ID = string(uF.ID())
			util.Usage = uF.Usage()
			util.User = uF.User()
			util.Nice = uF.Nice()
			util.System = uF.System()
			util.Idle = uF.Idle()
			util.IOWait = uF.IOWait()
		}
		u.CPU[i] = util
	}
	return u
}

// GetUtilization returns the cpu utilization.  Utilization calculations
// requires two pieces of data.  This func gets a snapshot of /proc/stat,
// sleeps for a second, takes another snapshot and calcualtes the utilization
// from the two snapshots.  If ongoing utilitzation information is desired,
// the UtilizationTicker should be used; it's better suited for ongoing
// utilization information being; using less cpu cycles and generating less
// garbage.
func GetUtilization() (Utilization, error) {
	stat1, err := GetStats()
	if err != nil {
		return Utilization{}, err
	}
	time.Sleep(time.Second)
	stat2, err := GetStats()
	if err != nil {
		return Utilization{}, err
	}

	return calculateUtilization(stat1, stat2), nil
}

// GetUtilizationFlat returns CPU Utilization informaton as Flatbuffer
// serialized bytes.
func GetUtilizationFlat() ([]byte, error) {
	u, err := GetUtilization()
	if err != nil {
		return nil, err
	}
	return u.SerializeFlat(), nil
}

// UtilizationTicker processes CPU utilization information on a ticker.  The
// generated utilization data is sent to the outCh.  Any errors encountered
// are sent to the errCh.  Processing ends when either a done signal is
// received or the done channel is closed.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func UtilizationTicker(interval time.Duration, outCh chan Utilization, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(outCh)
	// predeclare some vars
	var (
		i, j, pos, val, fieldNum int
		v                        byte
		name                     string
		stop                     bool
		prior                    Stats
	)
	// first get stats as the baseline
	cur, err := GetStats()
	if err != nil {
		errs <- err
	}
	// ticker
tick:
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			prior.Ctxt = cur.Ctxt
			prior.BTime = cur.BTime
			prior.Processes = cur.Processes
			if len(prior.CPU) != len(cur.CPU) {
				prior.CPU = make([]Stat, len(cur.CPU))
			}
			copy(prior.CPU, cur.CPU)
			cur.Timestamp = time.Now().UTC().UnixNano()
			_, err := proc.Seek(0, os.SEEK_SET)
			if err != nil {
				errs <- joe.Error{Type: "cpu", Op: "utilization ticker: seek /proc/stat", Err: err}
				continue tick
			}
			buf.Reset(proc)
			cur.CPU = cur.CPU[:0]
			// read each line until eof
			for {
				line, err := buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errs <- joe.Error{Type: "cpu stat", Op: "reading /proc/stat output", Err: err}
					break
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
					stat := Stat{ID: name}
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
								errs <- joe.Error{Type: "cpu stat", Op: "convert cpu data", Err: err}
								continue
							}
							j = pos + i + 1
							if fieldNum < 4 {
								if fieldNum == 1 {
									stat.User = int64(val)
								} else if fieldNum == 2 {
									stat.Nice = int64(val)
								} else if fieldNum == 3 {
									stat.System = int64(val)
								}
							} else if fieldNum < 7 {
								if fieldNum == 4 {
									stat.Idle = int64(val)
								} else if fieldNum == 5 {
									stat.IOWait = int64(val)
								} else if fieldNum == 6 {
									stat.IRQ = int64(val)
								}
							} else if fieldNum < 10 {
								if fieldNum == 7 {
									stat.SoftIRQ = int64(val)
								} else if fieldNum == 8 {
									stat.Steal = int64(val)
								} else if fieldNum == 9 {
									stat.Quest = int64(val)
								}
							} else if fieldNum == 10 {
								stat.QuestNice = int64(val)
							}
						}
					}
					cur.CPU = append(cur.CPU, stat)
					stop = false
					continue
				}
				if name == "ctxt" {
					// rest of the line is the data
					val, err = strconv.Atoi(string(line[pos : len(line)-1]))
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert ctxt data", Err: err}
					}
					cur.Ctxt = int64(val)
					continue
				}
				if name == "btime" {
					// rest of the line is the data
					val, err = strconv.Atoi(string(line[pos : len(line)-1]))
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert btime data", Err: err}
					}
					cur.BTime = int64(val)
					continue
				}
				if name == "processes" {
					// rest of the line is the data
					val, err = strconv.Atoi(string(line[pos : len(line)-1]))
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert processes data", Err: err}
					}
					cur.Processes = int64(val)
					continue
				}
			}
			outCh <- calculateUtilization(prior, cur)
		}
	}
}

// UtilizationTickerFlat processes CPU utilization information on a ticker
// The generated utilization data serialized with flatbuffers and is sent to
// the outCh.  Any errors encountered are sent to the errCh.  Processing ends
// when either a done signal is received or the done channel is closed.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func UtilizationTickerFlat(interval time.Duration, outCh chan []byte, done chan struct{}, errs chan error) {
	out := make(chan Utilization)
	defer close(outCh)
	go UtilizationTicker(interval, out, done, errs)
	bldr := fb.NewBuilder(0)
	var u Utilization
	for {
		select {
		case <-done:
			return
		case u = <-out:
			bldr.Reset()
			outCh <- u.SerializeFlatBuilder(bldr)
		}
	}
}

// usage = ()(Δuser + Δnice + Δsystem)/(Δuser+Δnice+Δsystem+Δidle)) * CLK_TCK
func calculateUtilization(s1, s2 Stats) Utilization {
	u := Utilization{
		Timestamp:  s2.Timestamp,
		BTimeDelta: int32(s2.Timestamp/1000000000 - s2.BTime),
		CtxtDelta:  s2.Ctxt - s1.Ctxt,
		Processes:  int32(s2.Processes),
		CPU:        make([]Util, len(s2.CPU)),
	}
	var dUser, dNice, dSys, dIdle, tot float32
	// Rest of the calculations are per core
	for i := 0; i < len(s2.CPU); i++ {
		v := Util{ID: s2.CPU[i].ID}
		dUser = float32(s2.CPU[i].User - s1.CPU[i].User)
		dNice = float32(s2.CPU[i].Nice - s1.CPU[i].Nice)
		dSys = float32(s2.CPU[i].System - s1.CPU[i].System)
		dIdle = float32(s2.CPU[i].Idle - s1.CPU[i].Idle)
		tot = dUser + dNice + dSys + dIdle
		v.Usage = (dUser + dNice + dSys) / tot * float32(s2.ClkTck)
		v.User = dUser / tot * float32(s2.ClkTck)
		v.Nice = dNice / tot * float32(s2.ClkTck)
		v.System = dSys / tot * float32(s2.ClkTck)
		v.Idle = dIdle / tot * float32(s2.ClkTck)
		v.IOWait = float32(s2.CPU[i].IOWait-s1.CPU[i].IOWait) / tot * float32(s2.ClkTck)
		u.CPU[i] = v
	}
	return u
}
