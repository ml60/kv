package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

type (
	Kvdb struct {
		base  [4096]*Lev1
		count [3]int
	}
	Lev1 [32]*Lev2
	Lev2 [32]*int32
)

var bit = []uint{0, 12, 17}
var mask = []int32{0x0fff, 0x1f, 0x1f}

func (db *Kvdb) insert(val int32) {
	p1 := &db.base[val&mask[0]]
	if *p1 == nil {
		*p1 = new(Lev1)
		db.count[0]++
	}

	p2 := &(*p1)[val>>bit[1]&mask[1]]
	if *p2 == nil {
		*p2 = new(Lev2)
		db.count[1]++
	}

	p3 := &(*p2)[val>>bit[2]&mask[2]]
	if *p3 == nil {
		*p3 = new(int32)
		db.count[2]++
	}

	**p3 = val
}

func (db *Kvdb) search(key int32) int32 {
	p1 := db.base[key&mask[0]]
	if p1 == nil {
		return -1
	}

	p2 := p1[key>>bit[1]&mask[1]]
	if p2 == nil {
		return -2
	}

	p3 := p2[key>>bit[2]&mask[2]]
	if p3 == nil {
		return -3
	}

	return *p3
}

func (db *Kvdb) stats() {
	tot := 12592
	fmt.Printf("base:           = %d\n", tot)

	mbyte := []int{256, 256, 4}
	for i := 0; i < 3; i++ {
		ptot := db.count[i] * mbyte[i]
		tot += ptot
		fmt.Printf("lev%d: %d * %d = %d\n", i, db.count[i], mbyte[i], ptot)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("alloc: %d memcalc: %d diff: %d\n",
		m.Alloc, tot, m.Alloc-uint64(tot))
}

func main() {
	n := 0
	if len(os.Args) > 1 {
		n, _ = strconv.Atoi(os.Args[1])
	}

	db := new(Kvdb)
	for i := 0; i < n; i++ {
		db.insert(int32(i * 3767 % 4191304))
	}

	fmt.Printf("key:%d val:%d\n", 7534, db.search(7534))
	fmt.Printf("key:%d val:%d\n", 7535, db.search(7535))
	db.stats()
}
