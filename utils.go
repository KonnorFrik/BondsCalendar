package main

import (
    "strconv"
    "time"

	"github.com/gbin/goncurses"
)

/* Show message, ask input string and convert it to int */
func AskInt(win *goncurses.Window, y, x int, msg string, maxInputLen int) (int, error) {
    win.MovePrint(y, x, msg)
    win.Refresh()
    resultStr, err := win.GetString(maxInputLen)

    if err != nil {
        return 0, err
    }

    result, err := strconv.Atoi(resultStr)

    if err != nil {
        return 0, err
    }

    return result, nil
}

/* Show message and ask for input string */
func AskString(win *goncurses.Window, y, x int, msg string, maxInputLen int) (string, error) {
    win.MovePrint(y, x, msg)
    win.Refresh()

    result, err := win.GetString(maxInputLen)

    if err != nil {
        return "", err
    }

    return result, nil
}

/* Show message, ask input string and convert it to time.Time */
func AskDate(win *goncurses.Window, y, x int, msg string, maxInputLen int, layout string) (time.Time, error) {
    win.MovePrint(y, x, msg)
    win.Refresh()

    resultStr, err := win.GetString(maxInputLen)

    if err != nil {
        return time.Time{}, err
    }

    result, err := time.Parse(layout, resultStr)

    if err != nil {
        return time.Time{}, err
    }

    return result, nil
}

/* Show one line text at new boxed window, until any input given and return last input */
func PopUpText(y, x int, msg string) goncurses.Key {
    winHeight, winWidth := 3, len(msg) + 2
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
