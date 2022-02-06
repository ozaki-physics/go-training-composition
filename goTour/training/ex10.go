package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.\n", os)
	}

	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
	// 以下と同じこと
	if t := time.Now(); t.Hour() < 12 {
		fmt.Println("Good morning!")
	} else if t.Hour() < 17 {
		fmt.Println("Good afternoon.")
	} else {
		fmt.Println("Good evening.")
	}
}
