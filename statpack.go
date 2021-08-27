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
	var size int64
	if d.Info != nil {
		size = d.Info.Size()
	}
	s.total.update(size)
	if d.Error != nil {
		s.errors.update(size)
	}
	if d.Attr == nil {
		s.untracked.update(size)
	} else {
		key := filemeta.ToHashKey(d.Attr.Hash)
		if !s.hashes[key] {
			s.hashes[key] = true
			s.unique.update(size)
		}
	}
	if d.Hashed {
		s.hashed.update(size)
	}
	if d.Changed {
		s.changed.update(size)
	}
	if d.VerifyFailed {
		s.failed.update(size)
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
