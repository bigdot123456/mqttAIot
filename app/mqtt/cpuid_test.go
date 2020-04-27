package main

import (
	"fmt"
	cpuid "github.com/jeek120/cpuid"
	"testing"
)

func TestCPUID(t *testing.T) {

	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 0)
	fmt.Printf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])
}
