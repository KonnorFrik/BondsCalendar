/*
A wrap for ncurses window
Register function:
    CustomDraw - for draw anything
    SpecialInputFunc - function binds with key, for processing this input key
Register a next windows:
    Same struct *FSMWindow binds with ncurses key. When you call FSMWindow.Input() you will recieve same *FSMWindow or next if pressed a special key for switch them

Recomended order of functions:
    .Draw() - for clear and draw your data
    .DrawBox() - for draw a box and title and refresh window
    .Input() - can block if ncurses.timeout not setted
        Also return a bool from SpecialInputFunc, which you can use as you want
*/
package main

import (
	"github.com/gbin/goncurses"
)

type SpecialInputFunc func() bool // Custom func for process input, returning false mean stop to work and exit
type CustomDraw func()            // Custom func from draw anything

type FSMWindow struct {
	Window       *goncurses.Window
	SizeY, SizeX int
	PosY, posX   int
	Title        string

	drawFunc     CustomDraw
	nextWindow   map[goncurses.Key]*FSMWindow
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

/* Bind a ncurses key with function for processing this key while press it */
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
		self.Window.MovePrint(0, self.SizeX/2-(len(self.Title)/2), self.Title)
	}

	self.Window.Refresh()
	return self
}

// TODO:
// change bool to int

/*
Return next window or same and bool status from custom input-proccessing function
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
