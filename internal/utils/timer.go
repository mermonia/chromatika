package utils

import (
	"fmt"
	"time"
)

func Timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("[%s]: %s\n", name, time.Since(start))
	}
}
