package main

import (
	"fmt"
)

type kv struct {
	key    string
	values []interface{}
}

func printkv(x []kv) {
	width0 := 0
	for _, i := range x {
		if len(i.key) > width0 {
			width0 = len(i.key)
		}
	}
	const widthN = 30
	for _, i := range x {
		fmt.Printf("%*s", width0, i.key)
		for _, v := range i.values {
			fmt.Printf("%*v", widthN, v)
		}
		fmt.Println()
	}
}
