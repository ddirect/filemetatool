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

var nullHash = make([]byte, filemeta.HashSize)

func listCore(fileName string) {
	data := filemeta.Get(fileName)
	check.E(data.Error)
	hash := "<changed>"
	if !data.Changed {
		hash = hex.EncodeToString(data.Hash)
	}
	fmt.Printf("%64s%20d  %s  %s\n", hash, data.Size, formatTime(data.GetModTime()), fileName)
}
