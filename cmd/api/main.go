package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	// Get the Unix timestamp in seconds
	unixTimestamp := now.Unix()
	fmt.Println(unixTimestamp)
}
