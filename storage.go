package transactor

// Transactor
// Storage
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"log"
	//"runtime"
	"sync"
	//"sync/atomic"
)

const storageDegree uint64 = 16
const storageNumber uint64 = 1 << storageDegree
const storageShift uint64 = 64 - storageDegree

type Storage struct {
	data [storageNumber]*Section
}

// newStorage - create new Storage
func newStorage() *Storage {
	s := &Storage{}
	for i := uint64(0); i < storageNumber; i++ {
		s.data[i] = newSection()
	}
	return s
}

func (s *Storage) addUnit(id int64) bool {
	section := s.data[(uint64(id)<<storageShift)>>storageShift]
	return section.addUnit(id)
}

func (s *Storage) getUnit(id int64) *Unit {
	return s.data[(uint64(id)<<storageShift)>>storageShift].getUnit(id)
}

func (s *Storage) delUnit(id int64) (*Unit, bool) {
	return s.data[(uint64(id)<<storageShift)>>storageShift].delUnit(id)
}

func (s *Storage) id(id int64) uint64 {
	return (uint64(id) << storageShift) >> storageShift
}

type Section struct {
	sync.Mutex
	data map[int64]*Unit
}

// newSection - create new Section
func newSection() *Section {
	s := &Section{
		data: make(map[int64]*Unit),
	}
	return s
}

/*
func (s *Section) lock() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.CompareAndSwapInt64(&s.hasp, 0, 1) {
			return true
		}
		runtime.Gosched()
	}
	return false
}
*/

func (s *Section) addUnit(id int64) bool {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.data[id]; !ok {
		s.data[id] = newUnit()
		return true
	}
	return false
}

func (s *Section) getUnit(id int64) *Unit {
	if u, ok := s.data[id]; ok {
		return u
	}
	return nil
}

func (s *Section) delUnit(id int64) (*Unit, bool) {
	s.Lock()
	defer s.Unlock()

	if u, ok := s.data[id]; ok {
		delete(s.data, id)
		return u, true
	}
	return nil, false
}
