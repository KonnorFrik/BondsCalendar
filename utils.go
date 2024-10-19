package main

import (
	"github.com/gbin/goncurses"
)
// TODO: make PopUpText working with multiply lines
//      iterate over given string and split by spaces and print each line
//      mb add scrolling

/* Show one line text at new boxed window, until any input given and return last input */
func PopUpText(y, x int, msg string) goncurses.Key {
	winHeight, winWidth := 3, len(msg)+2
	win, err := goncurses.NewWindow(winHeight, winWidth, y, x)

	if err != nil {
		return goncurses.Key(-1)
	}

	defer win.Delete()
	win.Box(0, 0)
	win.MovePrint(1, 1, msg)
	win.Refresh()
	input := win.GetChar()
	return input
}

/*
Create new window at y,x.
Show one line message and ask for input
New window height is always 3, width len(msg) + maxInputLen + 2(for box)
*/
func PopUpAskString(y, x int, msg string, maxInputLen int) (string, error) {
	winHeight, winWidth := 3, len(msg)+maxInputLen+2
	win, err := goncurses.NewWindow(winHeight, winWidth, y, x)

	if err != nil {
		return "", err
	}

	defer win.Delete()
	win.Box(0, 0)
	win.MovePrint(1, 1, msg)
	win.Refresh()

	result, err := win.GetString(maxInputLen)

	if err != nil {
		return "", err
	}

	return result, nil
}
