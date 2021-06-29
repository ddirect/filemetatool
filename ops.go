package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ddirect/filemeta"
)

func formatTime(x time.Time) string {
	return x.Format("2006-01-02 15:04:05.000")
}

var nullHash = make([]byte, 32)

func listCore(fileName string) {
	data, err := filemeta.Get(fileName)
	check(err)
	var hash string
	if data.Attr != nil {
		hash = hex.EncodeToString(data.Attr.Hash)
	} else {
		if data.Changed {
			hash = "<changed>"
		}
	}
	fmt.Printf("%64s%20d  %s  %s\n", hash, data.Info.Size(), formatTime(data.Info.ModTime()), fileName)
}

func fetch(fetchFunc func(string) (data filemeta.Data, errOut error)) (func(string), func()) {
	var s statPack
	return func(fileName string) {
			data, err := fetchFunc(fileName)
			check(err)
			s.update(&data)
		}, func() {
			printkv(s.toKv())
		}
}

func scrub() (func(string), func()) {
	var s statPack
	return func(fileName string) {
			data, err := filemeta.Get(fileName)
			check(err)
			flags := 0
			if data.Attr != nil {
				ok, err := data.Verify()
				check(err)
				if !ok {
					fmt.Printf("failed: %s\n", data.Path)
					flags = flagFailed
				}
			}
			s.updateX(&data, flags)
		}, func() {
			printkv(s.toKv())
		}
}
