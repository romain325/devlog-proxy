package main

import (
	"fmt"
	"golang.org/x/exp/rand"
	"time"
)

func getHex(num int) string {
	hex := fmt.Sprintf("%x", num)
	if len(hex) == 1 {
		hex = "0" + hex
	}
	return hex
}

func getRngHexColor() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	color := []int{rand.Intn(255), rand.Intn(255), rand.Intn(255)}
	hex := "#" + getHex(color[0]) + getHex(color[1]) + getHex(color[2])
	return hex
}
