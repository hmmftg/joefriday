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

// Package structs defines the datastructures for net/info.
package structs

// Info holds the information for all of the network interfaces.
type Info struct {
	Timestamp  int64       `json:"timestamp"`
	Interfaces []Interface `json:"interfaces"`
}

// Interface contains information for a given network interface.
type Interface struct {
	Name        string `json:"name"`
	RBytes      int64  `json:"receive_bytes"`
	RPackets    int64  `json:"receive_packets"`
	RErrs       int64  `json:"receive_errs"`
	RDrop       int64  `json:"receive_drop"`
	RFIFO       int64  `json:"recieve_fifo"`
	RFrame      int64  `json:"receive_frame"`
	RCompressed int64  `json:"receive_compressed"`
	RMulticast  int64  `json:"receive_multicast"`
	TBytes      int64  `json:"transmit_bytes"`
	TPackets    int64  `json:"transmit_packets"`
	TErrs       int64  `json:"transmit_errs"`
	TDrop       int64  `json:"transmit_drop"`
	TFIFO       int64  `json:"transmit_fifo"`
	TColls      int64  `json:"transmit_colls"`
	TCarrier    int64  `json:"transmit_carrier"`
	TCompressed int64  `json:"transmit_compressed"`
}
