package main

import (
	"github.com/gbin/goncurses"
)

/*
Register key with window to return
Register key for special input processing
Own input processing func related to window
Own drawing func
*/

type SpecialInputFunc func()

type FSMWindow struct {
    Window *goncurses.Window
    SizeY, SizeX int
    PosY, posX int
    Title string

    nextWindow map[goncurses.Key]*FSMWindow
    specialInput map[goncurses.Key]SpecialInputFunc
}

func FSMWindowNew(sizeY, sizeX, posY, posX int) (*FSMWindow, error) {
    var err error
    obj := new(FSMWindow)
    obj.Window, err = goncurses.NewWindow(sizeY, sizeX, posY, posX)

    if err != nil {
        return nil, err
    }

    obj.nextWindow = make(map[goncurses.Key]*FSMWindow)
    obj.specialInput = make(map[goncurses.Key]SpecialInputFunc)

    return obj, nil
}

func (self *FSMWindow) SetTitle(title string) {
    self.Title = title
}

func (self *FSMWindow) DrawBox() {
    self.Window.Clear()
    self.Window.Box(0, 0)

    if self.Title != "" {
        self.Window.MovePrint(0, self.SizeX / 2 - (len(self.Title) / 2), self.Title)
    }
}
