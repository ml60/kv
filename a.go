package main

import (
	"fmt"
	"time"
)

func main() {
	for i:=0; i<10; i++ {
		fmt.Print("\r", time.Now())
		time.Sleep(1*time.Second)
	}
	fmt.Println()
}

