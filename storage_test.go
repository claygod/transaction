package transaction

// Core
// Storage/sections test
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "testing"

func TestStorageNew(t *testing.T) {
	s := newStorage()
	if s == nil {
		t.Error("Error creating `Storage`")
	}
	if uint64(len(s)) != storageNumber {
		t.Error("Error in the number of sections")
	}
}

func TestStorageAdd(t *testing.T) {
	s := newStorage()
	if !s.addUnit(1) {
		t.Error("Error adding a unit")
	}
	if s.addUnit(1) {
		t.Error("Failed to add unit again")
	}
}

func TestStorageGet(t *testing.T) {
	s := newStorage()
	s.addUnit(1)
	if s.getUnit(1) == nil {
		t.Error("No unit found (has been added)")
	}
	if s.getUnit(2) != nil {
		t.Error("Found a non-existent unit")
	}
}

func TestStorageDel(t *testing.T) {
	s := newStorage()
	s.addUnit(1)
	if _, ok := s.delUnit(1); !ok {
		t.Error("Unable to delete unit")
	}
	if _, ok := s.delUnit(1); ok {
		t.Error("Repeated deletion of the same unit!")
	}
}

func TestStorageId(t *testing.T) {
	s := newStorage()
	if s.id(1) != 1 {
		t.Error("The lower bits are incorrectly recalculated")
	}
	if s.id(1048575) != 65535 {
		t.Error("Improperly recalculated high-order bits")
	}
}

func TestSectionNew(t *testing.T) {
	if newSection() == nil {
		t.Error("Error creating `Section`")
	}
}

func TestSectionAdd(t *testing.T) {
	s := newSection()
	if !s.addUnit(1) {
		t.Error("Error adding a unit")
	}
	if s.addUnit(1) {
		t.Error("Failed to add unit again")
	}
}

func TestSectionGet(t *testing.T) {
	s := newSection()
	s.addUnit(1)
	if s.getUnit(1) == nil {
		t.Error("No unit found (has been added)")
	}
	if s.getUnit(2) != nil {
		t.Error("Found a non-existent unit")
	}
}

func TestSectionDel(t *testing.T) {
	s := newSection()
	s.addUnit(1)
	if _, ok := s.delUnit(1); !ok {
		t.Error("Unable to delete unit")
	}
	if _, ok := s.delUnit(1); ok {
		t.Error("Repeated deletion of the same unit!")
	}
}
