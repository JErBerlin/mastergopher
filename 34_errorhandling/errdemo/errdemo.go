package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

var errEmptyInput = errors.New("input cannot be empty")

type missingError struct {
	key string
}

func (e missingError) Error() string {
	return fmt.Sprintf("no entry found for %q", e.key)
}

func fetch(key string) (string, error) {
	db := map[string]string{"lang": "Go", "tool": "gopher", "pkg": "errors"}

	if key == "" {
		return "", errEmptyInput
	}
	if val, ok := db[key]; ok {
		return val, nil
	}
	return "", missingError{key}
}

func checkEmpty() {
	v, err := fetch("")
	if err != nil {
		if errors.Is(err, errEmptyInput) {
			fmt.Println("please provide a non-empty key")
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("value: %q\n", v)
}

func checkMissing() {
	key := "unknown"
	v, err := fetch(key)
	if err != nil {
		var me missingError
		if errors.As(err, &me) {
			fmt.Printf("missing: %s\n", err)
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("value: %q\n", v)
}

func nestedWalk(depth int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered at depth %d: %v\n", depth, r)
			// start custom stack printer
			pcs := make([]uintptr, 32)
			n := runtime.Callers(3, pcs) // skip Callers, this defer, and recover
			frames := runtime.CallersFrames(pcs[:n])
			for {
				frame, more := frames.Next()
				fmt.Printf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
				if !more {
					break
				}
			}
			// end custom stack printer
		}
	}()
	if depth == 5 {
		panic("reached max depth")
	}
	fmt.Println("depth", depth)
	nestedWalk(depth + 1)
}

func main() {
	checkEmpty()
	checkMissing()
	nestedWalk(1)
}
