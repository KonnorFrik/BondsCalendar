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

type SpecialInputFunc func() bool // Custom func for process input, returning false mean stop to work and exit
type CustomDraw func()  // Custom func from draw anything

type FSMWindow struct {
    Window *goncurses.Window
    SizeY, SizeX int
    PosY, posX int
    Title string

    drawFunc CustomDraw
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
    obj.SizeX = sizeX
    obj.SizeY = sizeY
    obj.posX = posX
    obj.PosY = posY

    return obj, nil
}

func (self *FSMWindow) FreeWindow() {
    self.Window.Delete()
}

func (self *FSMWindow) SetTitle(title string) *FSMWindow {
    self.Title = title
    return self
}

func (self *FSMWindow) SetCustomDraw(function CustomDraw) *FSMWindow {
    self.drawFunc = function
    return self
}

func (self *FSMWindow) RegisterNextWindow(key goncurses.Key, window *FSMWindow) *FSMWindow {
    self.nextWindow[key] = window
    return self
}

func (self *FSMWindow) RegisterInput(key goncurses.Key, function SpecialInputFunc) *FSMWindow {
    self.specialInput[key] = function
    return self
}

/* Call clear and then custom draw function. Call Draw before DrawBox */
func (self *FSMWindow) Draw() *FSMWindow {
    self.Window.Clear()

    if self.drawFunc != nil {
        self.drawFunc()
    }

    return self
}

/* Draw box and title if setted */
func (self *FSMWindow) DrawBox() *FSMWindow {
    self.Window.Box(0, 0)

    if self.Title != "" {
        self.Window.MovePrint(0, self.SizeX / 2 - (len(self.Title) / 2), self.Title)
    }

    self.Window.Refresh()
    return self
}

/*
Return next window or nil and status if this fsm window can continue work 
true - Yes, can
false - No
*/
func (self *FSMWindow) Input() (*FSMWindow, bool) {
    inputKey := self.Window.GetChar()
    var status bool = true
    nextWin, exist := self.nextWindow[inputKey]

    if exist {
        self.Window.Clear()
        return nextWin, status
    }

    inputFunc, exist := self.specialInput[inputKey]

    if exist {
        status = inputFunc()
    }

    return self, status
}
