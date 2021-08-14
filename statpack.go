package main

import (
	"sync"

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
		t.Append(name, c.count, format.Size(c.size))
	}
}

type statPack struct {
	total, untracked, hashed, changed, failed, errors countSize
	mutex                                             sync.Mutex
}

func (s *statPack) update(d *filemeta.Data) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	a.Append(".", "count", "size")
	s.total.appendRow(a, "total")
	s.changed.appendRow(a, "changed")
	s.untracked.appendRow(a, "untracked")
	s.hashed.appendRow(a, "hashed")
	s.failed.appendRow(a, "failed")
	s.errors.appendRow(a, "errors")
	return a
}
