package test

import (
	"fmt"
	"testing"
)

func TestDemo(t *testing.T) {
	var a int = -3
	var b uint = uint(a)
	fmt.Println(a)
	fmt.Println(b)
}
