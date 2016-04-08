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

// Package structs defines the data structures for disk.
package structs

// Stats holds the information for all of the network interfaces.
type Stats struct {
	Timestamp int64    `json:"timestamp"`
	Devices   []Device `json:"devices"`
}

// Device contains information for a given block device.
type Device struct {
	Major           uint32
	Minor           uint32
	Name            string `json:"name"`
	ReadsCompleted  uint64 `json:"reads_completed"`
	ReadsMerged     uint64 `json:"reads_merged"`
	ReadSectors     uint64 `json:'read_sectors'`
	ReadingTime     uint64 `json:"reading_time"`
	WritesCompleted uint64 `json:"writes_completed"`
	WritesMerged    uint64 `json:"writes_merged"`
	WrittenSectors  uint64 `json:"written_sectors"`
	WritingTime     uint64 `json:"writing_time"`
	IOInProgress    int32  `json:"io_in_progress"`
	IOTime          uint64 `json:"io_time"`
	WeightedIOTime  uint64 `json:"weighted_io_time"`
}

// Usage holds the information for all of the network interfaces.
type Usage struct {
	Timestamp int64    `json:"timestamp"`
	TimeDelta uint32   `json:"time_delta"`
	Devices   []Device `json:"devices"`
}
