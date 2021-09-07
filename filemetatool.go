package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ddirect/filemeta"
)

func main() {
	var do string
	var probeThreads, hashThreads int
	var sync bool
	flag.StringVar(&do, "do", "", "list|refresh|stat|scrub|inspect")
	flag.IntVar(&probeThreads, "probe_threads", runtime.NumCPU(), "number of threads used to probe the metadata")
	flag.IntVar(&hashThreads, "hash_threads", 1, "number of threads used for hashing (always 1 for refresh)")
	flag.BoolVar(&sync, "sync", false, "use fully synchronous mode")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 || do == "" {
		flag.Usage()
		return
	}

	var op filemeta.Op
	switch do {
	case "list":
		walk(files, listCore)
		return
	case "refresh":
		op = filemeta.OpRefresh
	case "stat":
		op = filemeta.OpGet
	case "inspect":
		op = filemeta.OpInspect
	case "scrub":
		op = filemeta.OpVerify
	default:
		fmt.Fprintf(os.Stderr, "unknown operation '%s'\n", op)
		return
	}

	s := newStatPack()
	handleData := func(data *filemeta.Data) {
		if data.Error != nil {
			fmt.Println(data.Error)
		}
		s.update(data)
	}
	if sync {
		run, done := filemeta.SyncOperations(op)
		defer done()
		walk(files, func(path string) {
			data := run(path)
			handleData(&data)
		})
	} else {
		async := filemeta.AsyncOperations(op, probeThreads, hashThreads)

		go func() {
			walk(files, func(path string) {
				async.FileIn <- path
			})
			close(async.FileIn)
		}()

		for data := range async.DataOut {
			handleData(&data)
		}
	}
	fmt.Print(s.toTable())
}
