package main

import (
	"sync"

	"github.com/ddirect/filemeta"
	"github.com/ddirect/format"
)

const (
	flagFailed = 1 << iota
)

type countSize struct {
	count, size int64
}

func (c *countSize) update(size int64) {
	c.count++
	c.size += size
}

func (c *countSize) appendRow(t format.Table, name string) format.Table {
	if c.count > 0 {
		return append(t, format.TableRow{name, c.count, format.Size(c.size)})
	}
	return t
}

type statPack struct {
	total, untracked, hashed, changed, failed countSize
	mutex                                     sync.Mutex
}

func (s *statPack) update(d *filemeta.Data) {
	s.updateX(d, 0)
}

func (s *statPack) updateX(d *filemeta.Data, flags int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	size := d.Info.Size()
	s.total.update(size)
	if d.Attr == nil {
		s.untracked.update(size)
	}
	if d.Hashed {
		s.hashed.update(size)
	}
	if d.Changed {
		s.changed.update(size)
	}
	if flags&flagFailed != 0 {
		s.failed.update(size)
	}
}

func (s *statPack) toTable() format.Table {
	a := format.Table{format.TableRow{".", "count", "size"}}
	a = s.total.appendRow(a, "total")
	a = s.changed.appendRow(a, "changed")
	a = s.untracked.appendRow(a, "untracked")
	a = s.hashed.appendRow(a, "hashed")
	a = s.failed.appendRow(a, "failed")
	return a
}
