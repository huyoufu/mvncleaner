package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPrintMetric(t *testing.T) {

}

func TestTimerTrack_NextStepNum(t *testing.T) {
	tt := new(TimerTrack)
	tt.NextStepNum()
	fmt.Println(tt.stepNum)
}
func TestTimerTrack_PrintBeautiful(t *testing.T) {
	DefaultTimerTrack.Start()
	time.Sleep(time.Second)
	DefaultTimerTrack.Step()
	DefaultTimerTrack.Step()
	time.Sleep(time.Second)
	DefaultTimerTrack.Step()
	DefaultTimerTrack.End()

	DefaultTimerTrack.PrintBeautiful()

}

func TestTimerTrack_Cost(t *testing.T) {
	DefaultTimerTrack.Start()
	time.Sleep(time.Second)
	DefaultTimerTrack.Step()
	DefaultTimerTrack.Step()
	time.Sleep(time.Second)
	DefaultTimerTrack.Step()
	DefaultTimerTrack.End()
	cost := DefaultTimerTrack.Cost()
	fmt.Printf("共花费了%d ms\n", cost)
}

func BenchmarkTimerTrack_NextStepNum(b *testing.B) {
	tt := new(TimerTrack)
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(tt *TimerTrack) {
			for i := 0; i < 10000; i++ {
				tt.NextStepNum()
			}
			wg.Done()
		}(tt)
	}
	wg.Wait()
	fmt.Printf("计数:%d\n", tt.stepNum)

}
