package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ddirect/check"
	"github.com/ddirect/filemeta"
)

func formatTime(x time.Time) string {
	return x.Format("2006-01-02 15:04:05.000")
}

var nullHash = make([]byte, 32)

func listCore(fileName string) {
	data, err := filemeta.Get(fileName)
	check.E(err)
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

func handle(err error) int {
	flags := 0
	if err != nil {
		fmt.Println(err)
		flags = flagError
	}
	return flags
}

func fetch(fetchFunc filemeta.FetchFunc) (func(string), func()) {
	var s statPack
	return func(fileName string) {
			data, err := fetchFunc(fileName)
			s.updateX(&data, handle(err))
		}, func() {
			fmt.Print(s.toTable())
		}
}
