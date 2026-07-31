package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func UnexpectedVal(t *testing.T, expect int, got int) {
	t.Error(fmt.Sprintf("Got unexpected value from main: expected %v, got %v", expect, got))
}

func TestSchemaMainFunc(t *testing.T) {
	ch := make(chan int)
	go func() {
		os.Args = []string{"test", "--schema", "all"}
		main()
	}()
	time.Sleep(1 * time.Second)
	close(ch)
	output := <-ch
	if output != 0 {
		UnexpectedVal(t, 0, output)
	}
}
