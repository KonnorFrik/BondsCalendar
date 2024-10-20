package main

import (
	"fmt"

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

/*
Create new window and show given strings as list with indices
Have own input processing
*/
func PopUpScrollableList(data []string, title string, sizeY, posY, posX int) error {
	// TODO:
	// pass MaxX MaxY of main window and spawn new subwindow with list at center
	if title == "" {
		title = "|List|"
	}

	var sizeX int = len(title) + 2

	for _, str := range data {
		sizeX = max(sizeX, len(str)+2)
	}

	win, err := goncurses.NewWindow(sizeY, sizeX, posY, posX)

	if err != nil {
		return err
	}

	var startInd int = 0

	printSlice := func(startInd int) int {
		var endInd int = min(len(data)-1, sizeY-1)

		if startInd >= len(data) {
			startInd = len(data) - 1
		}

		if startInd < 0 {
			startInd = 0
		}

		var x int = 1
		var y int = 1

		for ind := startInd; ind <= endInd; ind++ {
			win.MovePrint(y, x, data[ind])
			y++
		}

		return startInd
	}

	var input goncurses.Key

	for input != ExitKey {
		win.Clear()
		win.Box(0, 0)
		win.MovePrint(0, sizeX/2-(len(title)/2), title)
		win.Refresh()

		switch input {
		case ScrollUpKey:
			startInd -= 1

		case ScrollDownKey:
			startInd += 1
		}

		startInd = printSlice(startInd)
		input = win.GetChar()
	}

	return nil
}

/* Remove element from slice by it's index */
func SliceRemoveByIndex[T any](slc []T, index int) ([]T, error) {
	if index >= len(slc) {
		return slc, fmt.Errorf("Index: %d out of slice bounds with len: %d\n", index, len(slc))
	}

	// if len(slc) == 0 {
	//
	// }

	return append(slc[:index], slc[index+1:]...), nil
}
