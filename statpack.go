package main

import (
	"github.com/ddirect/filemeta"
	"github.com/ddirect/format"
)

type countSize struct {
	count, size int64
}

func (c *countSize) update(size int64) {
	c.count++
	c.size += size
}

func (c *countSize) appendRow(t *format.Table, name string) {
	if c.count > 0 {
		t.AppendRow(name, c.count, format.Size(c.size))
	}
}

type statPack struct {
	total, unique, untracked, hashed, changed, failed, errors countSize
	hashes                                                    map[filemeta.HashKey]bool
}

func newStatPack() *statPack {
	return &statPack{hashes: make(map[filemeta.HashKey]bool)}
}

func (s *statPack) update(d *filemeta.Data) {
	s.total.update(d.Size)
	if d.Error != nil {
		s.errors.update(d.Size)
	}
	if d.Hash == nil {
		s.untracked.update(d.Size)
	} else {
		key := filemeta.ToHashKey(d.Hash)
		if !s.hashes[key] {
			s.hashes[key] = true
			s.unique.update(d.Size)
		}
	}
	if d.Hashed {
		s.hashed.update(d.Size)
	}
	if d.Changed {
		s.changed.update(d.Size)
	}
	if d.VerifyFailed {
		s.failed.update(d.Size)
	}
}

func (s *statPack) toTable() *format.Table {
	a := new(format.Table)
	a.AppendRow(".", "count", "size")
	s.total.appendRow(a, "total")
	s.unique.appendRow(a, "unique")
	s.changed.appendRow(a, "changed")
	s.untracked.appendRow(a, "untracked")
	s.hashed.appendRow(a, "hashed")
	s.failed.appendRow(a, "failed")
	s.errors.appendRow(a, "errors")
	return a
}
