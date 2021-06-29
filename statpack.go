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

func (c *countSize) appendKv(a []kv, name string) []kv {
	if c.count > 0 {
		return append(a, kv{name, []interface{}{c.count, format.Size(c.size)}})
	}
	return a
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

func (s *statPack) toKv() []kv {
	a := []kv{kv{".", []interface{}{"count", "size"}}}
	a = s.total.appendKv(a, "total")
	a = s.changed.appendKv(a, "changed")
	a = s.untracked.appendKv(a, "untracked")
	a = s.hashed.appendKv(a, "hashed")
	a = s.failed.appendKv(a, "failed")
	return a
}
