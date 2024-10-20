/*
Wrap for ncurses Window
Always boxes
Features:
	print strings in terminal window
	get inputs from user with converting
*/
package terminal

import (
	"strconv"
	"time"

	"github.com/gbin/goncurses"
)

type TerminalSettings struct {
	SizeX, SizeY  int
	PosX, PosY    int
	DefaultEcho   bool   // Default value for your programm
	DefaultCursor byte   // Default value for your programm
	Title         string // Title for terminal

	// clearPosX, clearPosY int // position for set cursor in and delete line
	printPosX, printPosY int // position for set cursor in and print line
	inputPosX, inputPosY int // position for set cursor in and ask for input
	maxInput             int // maximum length of user input, defines as SizeX - 4
}

type Terminal struct {
	Window   *goncurses.Window
	Settings TerminalSettings
}

// TODO:
// fix terminal input
// terminal must allow input big string without interrupting
// user type string in terminal -> if length too big -> print part, save part in buffer, and continue input getting

/* Allocate new terminal, without initializing */
func TerminalNew() *Terminal {
	obj := new(Terminal)
	return obj
}

/* Init new ncurses window with given parameters */
func (self *Terminal) Init(settings TerminalSettings) error {
	var err error
	self.Window, err = goncurses.NewWindow(settings.SizeY, settings.SizeX, settings.PosY, settings.PosX)

	if err != nil {
		return err
	}

	// settings.clearPosX = 1
	// settings.clearPosY = 1
	settings.printPosX = 1
	settings.printPosY = settings.SizeY - 3
	settings.inputPosX = 1
	settings.inputPosY = settings.SizeY - 2
	settings.maxInput = settings.SizeX - 4

	self.Settings = settings
	self.Window.ScrollOk(true)
	self.Refresh()

	return err
}

/* Call .Delete for window for free memory */
func (self *Terminal) Delete() *Terminal {
	self.Window.ScrollOk(false)
	self.Window.Delete()
	self.Settings = TerminalSettings{}
	return self
}

/* Draw box, title and refresh changes */
func (self *Terminal) Refresh() {
	self.Window.Box(0, 0)
	self.Window.MovePrint(0, self.Settings.SizeX/2-5, self.Settings.Title)
	self.Window.Refresh()
}

/*
Clear all printed content
Return error from one stage:
  - Clearing
  - Drawing box
*/
func (self *Terminal) Clear() error {
	err := self.Window.Clear()

	if err != nil {
		return err
	}

	// err = self.Window.Box(0, 0)
	// self.Window.Refresh()
    self.Refresh()

	return err
}

/*
Print given msg to terminal. Split msg if it can't fit in one terminal line (Terminal.SizeX)
*/
func (self *Terminal) Print(msg string) *Terminal {
	if len(msg) >= (self.Settings.SizeX - 2) {
		// print splitted
		var startInd int
		var endInd int = self.Settings.SizeX - 2

		for endInd < len(msg) {
			self.printScrolled(msg[startInd:endInd])

			startInd = endInd
			endInd += self.Settings.SizeX - 2
		}

		self.printScrolled(msg[startInd:])

	} else {
		self.printScrolled(msg)
	}

	self.Refresh()

	return self
}

/*
Show message and ask for input string.
Function return after a pressing 'enter'
*/
func (self *Terminal) AskString(question string) (string, error) {
	goncurses.Echo(true)
	goncurses.Cursor(1)

	result, err := self.askInput(question)
	self.Print(result)

	goncurses.Echo(self.Settings.DefaultEcho)
	goncurses.Cursor(self.Settings.DefaultCursor)
	return result, err
}

/*
Show message and ask for input char.
Function return after a pressin any key
*/
func (self *Terminal) AskChar(question string) goncurses.Key {
	goncurses.Echo(true)
	goncurses.Cursor(1)

	result := self.askChar(question)
	self.Print(goncurses.KeyString(result))

	goncurses.Echo(self.Settings.DefaultEcho)
	goncurses.Cursor(self.Settings.DefaultCursor)
	return result
}

/*
Show message, ask input string and convert it to int.
Function return after a pressing 'enter'
*/
func (self *Terminal) AskInt(question string) (int, error) {
	var result int
	goncurses.Echo(true)
	goncurses.Cursor(1)

	input, err := self.askInput(question)
	self.Print(input)

	if err == nil {
		result, err = strconv.Atoi(input)
	}

	goncurses.Echo(self.Settings.DefaultEcho)
	goncurses.Cursor(self.Settings.DefaultCursor)
	return result, err
}

/*
Show message, ask input string and convert it to time.Time.
Function return after a pressing 'enter'
*/
func (self *Terminal) AskDate(question, layout string) (time.Time, error) {
	var result time.Time
	goncurses.Echo(true)
	goncurses.Cursor(1)

	input, err := self.askInput(question)
	self.Print(input)

	if err == nil {
		result, err = time.Parse(layout, input)
	}

	goncurses.Echo(self.Settings.DefaultEcho)
	goncurses.Cursor(self.Settings.DefaultCursor)
	return result, err
}

/* Help function - scroll terminal window and print only one line, which must fit in */
func (self *Terminal) printScrolled(oneLineMsg string) {
	// Clear input field from HLine option 1
	// self.Window.Move(self.Settings.SizeY - 1, 0)
	// self.Window.ClearToEOL()

	self.Window.Scroll(1)
	self.Window.MovePrint(self.Settings.printPosY, self.Settings.printPosX, oneLineMsg)
	self.Window.ClearToEOL()

	// Clear input field from HLine option 2
	self.Window.Move(self.Settings.printPosY+1, self.Settings.printPosX)
	self.Window.ClearToEOL()
}

/* Help function - Correctly process a scrolling, printing a request prompt, and get string from user */
func (self *Terminal) askInput(question string) (string, error) {
	self.printScrolled(question)
	self.Refresh()
	self.Window.MovePrint(self.Settings.inputPosY, self.Settings.inputPosX, "> ")
	result, err := self.Window.GetString(self.Settings.maxInput)

	return result, err
}

/* Help function - Correctly process a scrolling, printing a request prompt, and get pressed key from user */
func (self *Terminal) askChar(question string) goncurses.Key {
	self.printScrolled(question)
	self.Refresh()
	self.Window.MovePrint(self.Settings.inputPosY, self.Settings.inputPosX, "> ")
	result := self.Window.GetChar()

	return result
}
