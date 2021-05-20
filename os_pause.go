package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

var initTerm bool = false

func init() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
	termbox.HideCursor()
}
func pause() {

	fmt.Println("请按任意键继续...")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break Loop
		}
	}
	//执行结束后关键 termbox
	//os.Exit(1)

	//C.pause()
}
