// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package polynomial

// #include <stdlib.h>
import "C" //a COMMENT has to be places just at the right place for this to work. Excellent language design!
import (
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"unsafe"
)

// Memory management for polynomials
// WARNING: This is not thread safe TODO: Make sure that is not a problem
// TODO: There is a lot of "unsafe" memory management here and needs to be vetted thoroughly

type enormousArray = [1 << 32]fr.Element // semantic necessity

type sizedPool struct {
	maxN  int
	pool  sync.Pool
	stats poolStats
}

type inUseData struct {
	allocatedFor []uintptr
	pool         *sizedPool
}

type Pool struct {
	//lock     sync.Mutex
	inUse    map[unsafe.Pointer]inUseData
	subPools []sizedPool
}

func (p *sizedPool) get(n int) unsafe.Pointer {
	p.stats.maake(n)
	return p.pool.Get().(unsafe.Pointer)
}

func (p *sizedPool) put(ptr unsafe.Pointer) {
	p.stats.dump()
	p.pool.Put(ptr)
}

func NewPool(maxN ...int) (pool Pool) {

	sort.Ints(maxN)
	pool = Pool{
		inUse:    make(map[unsafe.Pointer]inUseData),
		subPools: make([]sizedPool, len(maxN)),
	}

	for i := range pool.subPools {
		subPool := &pool.subPools[i]
		subPool.maxN = maxN[i]
		subPool.pool = sync.Pool{
			New: func() interface{} {
				subPool.stats.Allocated++
				return C.malloc(C.ulong(8 * fr.Bytes * subPool.maxN))
			},
		}
	}
	return
}

func (p *Pool) findCorrespondingPool(n int) *sizedPool {
	poolI := 0
	for poolI < len(p.subPools) && n > p.subPools[poolI].maxN {
		poolI++
	}
	return &p.subPools[poolI] // out of bounds error here would mean that n is too large
}

func (p *Pool) Make(n int) []fr.Element {
	pool := p.findCorrespondingPool(n)
	ptr := pool.get(n)
	p.addInUse(ptr, pool)
	return (*enormousArray)(ptr)[:n]
}

// Dump dumps a set of polynomials into the pool
func (p *Pool) Dump(slices ...[]fr.Element) {
	for _, slice := range slices {
		ptr := getDataPointer(slice)
		if metadata, ok := p.inUse[ptr]; ok {
			delete(p.inUse, ptr)
			metadata.pool.put(ptr)
		} else {
			panic("attempting to dump a slice not created by the pool")
		}
	}
}

func (p *Pool) addInUse(ptr unsafe.Pointer, pool *sizedPool) {
	pcs := make([]uintptr, 2)
	n := runtime.Callers(3, pcs)

	if prevPcs, ok := p.inUse[ptr]; ok { // TODO: remove if unnecessary for security
		panic(fmt.Errorf("re-allocated non-dumped slice, previously allocated at %v", runtime.CallersFrames(prevPcs.allocatedFor)))
	}
	p.inUse[ptr] = inUseData{
		allocatedFor: pcs[:n],
		pool:         pool,
	}
}

func printFrame(frame runtime.Frame) {
	fmt.Printf("\t%s line %d, function %s\n", frame.File, frame.Line, frame.Function)
}

func (p *Pool) printInUse() {
	fmt.Println("slices never dumped allocated at:")
	for _, pcs := range p.inUse {
		fmt.Println("-------------------------")

		var frame runtime.Frame
		frames := runtime.CallersFrames(pcs.allocatedFor)
		more := true
		for more {
			frame, more = frames.Next()
			printFrame(frame)
		}
	}
}

type poolStats struct {
	Used          int
	Allocated     int
	ReuseRate     float64
	InUse         int
	GreatestNUsed int
	SmallestNUsed int
}

type poolsStats struct {
	SubPools []poolStats
	InUse    int
}

func (s *poolStats) maake(n int) {
	s.Used++
	s.InUse++
	if n > s.GreatestNUsed {
		s.GreatestNUsed = n
	}
	if s.SmallestNUsed == 0 || s.SmallestNUsed > n {
		s.SmallestNUsed = n
	}
}

func (s *poolStats) dump() {
	s.InUse--
}

func (s *poolStats) finalize() {
	s.ReuseRate = float64(s.Used) / float64(s.Allocated)
}

func getDataPointer(slice []fr.Element) unsafe.Pointer {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	return unsafe.Pointer(header.Data)
}

func (p *Pool) PrintPoolStats() {
	InUse := 0
	subStats := make([]poolStats, len(p.subPools))
	for i := range p.subPools {
		subPool := &p.subPools[i]
		subPool.stats.finalize()
		subStats[i] = subPool.stats
		InUse += subPool.stats.InUse
	}

	poolsStats := poolsStats{
		SubPools: subStats,
		InUse:    InUse,
	}
	serialized, _ := json.MarshalIndent(poolsStats, "", "  ")
	fmt.Println(string(serialized))
	p.printInUse()
}

func (p *Pool) Clone(slice []fr.Element) []fr.Element {
	res := p.Make(len(slice))
	copy(res, slice)
	return res
}

func (p *Pool) Free() {
	for ptr := range p.inUse {
		C.free(ptr)
	}
}
