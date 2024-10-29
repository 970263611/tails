package timed_task_tool

import (
	"fmt"
	"testing"
	"time"
)

func TestDelayAtFixedRate(t *testing.T) {
	DelayAtFixedRate(2, func() {
		fmt.Println("run TestDelayAtFixedRate")
	})
	time.Sleep(100 * time.Second)
}

func TestDelayWithFixedDelay(t *testing.T) {
	tt := DelayWithFixedDelay(2, func() {
		fmt.Println("run TestDelayWithFixedDelay")
	})
	time.Sleep(10 * time.Second)
	tt.Close()
}
