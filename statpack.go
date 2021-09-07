package main

import (
	"github.com/ddirect/filemeta"
	"github.com/ddirect/format"
	"github.com/ddirect/sys"
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
	total, uniqueFiles, uniqueHashes, untracked, hashed, changed, failed, errors countSize
	files                                                                        map[sys.FileKey]struct{}
	hashes                                                                       map[filemeta.HashKey]struct{}
}

func newStatPack() *statPack {
	return &statPack{
		files:  make(map[sys.FileKey]struct{}),
		hashes: make(map[filemeta.HashKey]struct{}),
	}
}

func (s *statPack) update(d *filemeta.Data) {
	s.total.update(d.Size)
	if d.Error != nil {
		s.errors.update(d.Size)
	}
	if d.Hash == nil {
		s.untracked.update(d.Size)
	} else {
		if _, ok := s.files[d.FileKey]; !ok {
			s.files[d.FileKey] = struct{}{}
			s.uniqueFiles.update(d.Size)
		}
		key := filemeta.ToHashKey(d.Hash)
		if _, ok := s.hashes[key]; !ok {
			s.hashes[key] = struct{}{}
			s.uniqueHashes.update(d.Size)
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
	s.uniqueFiles.appendRow(a, "unique files")
	s.uniqueHashes.appendRow(a, "unique hashes")
	s.changed.appendRow(a, "changed")
	s.untracked.appendRow(a, "untracked")
	s.hashed.appendRow(a, "hashed")
	s.failed.appendRow(a, "failed")
	s.errors.appendRow(a, "errors")
	return a
}
