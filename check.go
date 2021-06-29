package main

func handlePanic(err *error) {
	if r := recover(); r != nil {
		*err = r.(error)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkI(n int, err error) int {
	if err != nil {
		panic(err)
	}
	return n
}

func checkI64(n int64, err error) int64 {
	if err != nil {
		panic(err)
	}
	return n
}

func checkS(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}
