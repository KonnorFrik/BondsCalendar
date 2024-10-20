package main

import (
	"bonds_payment_calendar/bonds"
	"bonds_payment_calendar/terminal"
	"fmt"
	"strconv"
	"time"

	"github.com/gbin/goncurses"
)

/* Bonds info for one year */
type YearInfo struct {
	Year         int
	PaymentCount int
}

var (
	MaxX int
	MaxY int

	CurrentYear = time.Now().Year()
	Terminal    = terminal.TerminalNew()
	AllBonds    = bonds.BondsNew()
)

const (
	ExitKey           = 'q'
	HelpKey           = 'h'
	IncreaseYearKey   = '>'
	DecreaseYearKey   = '<'
	AppendBondsKey    = 'a'
	SaveBondsKey      = 's'
	LoadBondsKey      = 'l'
	ListBondsKey      = 'v'
	ScrollUpKey       = 'w'
	ScrollDownKey     = 's'
	StartOfCommandKey = ':'

	DefaultDateLayout = "02.01.2006"
)

func CommandHelp(args []string) error {
	if len(args) == 0 {
		helpCommand, _ := CommandTable["help"]
		Terminal.Print(helpCommand.Info)
		return nil
	}

	err := IsCommandExist(args[0])

	if err != nil {
		return err
	}

	commandStruct, _ := CommandTable[args[0]]
	Terminal.Print("\t" + args[0])
	Terminal.Print(commandStruct.Info)

	return err
}

func CommandList(args []string) error {
	DrawListBonds(AllBonds, MaxY-1, MaxX/3*2, 0, 0)
	return nil
}

func CommandLoad(args []string) error {
	var filename string

	if len(args) == 0 {
		input, err := Terminal.AskString("Filename for load:")

		if err != nil {
			return err
		}

		filename = input

	} else {
		filename = args[0]
	}

	err := AllBonds.LoadFromFile(filename)

	if err != nil {
		return err
	}

	Terminal.Print(fmt.Sprintf("Loaded: %d bonds", len(AllBonds.Bonds)))

	return err
}

func CommandSave(args []string) error {
	var filename string

	if len(args) == 0 {
		input, err := Terminal.AskString("Filename for save:")

		if err != nil {
			return err
		}

		filename = input

	} else {
		filename = args[0]
	}

	err := AllBonds.SaveToFile(filename)

	if err != nil {
		return err
	}

	Terminal.Print(fmt.Sprintf("Saved: %d bonds", len(AllBonds.Bonds)))

	return err
}

func CommandDelete(args []string) error {
	var index int
	var err error

	if len(args) == 0 {
		tmp, err := Terminal.AskInt("Index for delete:")

		if err != nil {
			return err
		}

		index = tmp

	} else {
		tmp, err := strconv.Atoi(args[0])

		if err != nil {
			return err
		}

		index = tmp
	}

	AllBonds.Bonds, err = SliceRemoveByIndex(AllBonds.Bonds, index)
	return err
}

/* Draw a window with own input processing */
func HelpWindow() {
	startY, startX := MaxY/4, MaxX/4
	height, width := MaxY/2, MaxX/2
	scr, err := goncurses.NewWindow(height, width, startY, startX)

	if err != nil {
		panic(err)
	}

	defer scr.Delete()
	var input goncurses.Key

	for input != ExitKey {
		var x int = 1
		var y int = 0
		scr.Clear()
		scr.Box(0, 0)

		scr.MovePrint(y, width/2-2, "|Help|")
		y++

		// Info for main window
		scr.MovePrintf(y, x, "For main programm")
		x += 2
		y++
		scr.MovePrintf(y, x, "Exit key: %c", ExitKey)
		y++
		scr.MovePrintf(y, x, "Next year: %c", IncreaseYearKey)
		y++
		scr.MovePrintf(y, x, "Previous year: %c", DecreaseYearKey)
		y++
		scr.MovePrintf(y, x, "Append Bonds: %c", AppendBondsKey)
		y++
		scr.MovePrintf(y, x, "Save Bonds: %c", SaveBondsKey)
		y++
		scr.MovePrintf(y, x, "Load Bonds: %c", LoadBondsKey)
		y++
		scr.MovePrintf(y, x, "List Bonds: %c", ListBondsKey)
		y++

		// Info for bondsList
		x = width / 3
		y = 1
		scr.MovePrintf(y, x, "For bonds list")
		y++
		x += 2
		scr.MovePrintf(y, x, "Scroll Up: %c", ScrollUpKey)
		y++
		scr.MovePrintf(y, x, "Scroll Down: %c", ScrollDownKey)
		y++

		scr.Refresh()
		input = scr.GetChar()
	}
}

/* Draw graph of payments for given year */
func DrawGraphByYear(obj *bonds.Bonds, year int, win *goncurses.Window, sizeX, sizeY, offsetX int) YearInfo {
	result := YearInfo{
		Year:         year,
		PaymentCount: 0,
	}

	var x int = 1
	var monthY int = sizeY - 1
	var countY int = sizeY - 3
	win.MovePrint(monthY, x, "M")
	win.MovePrint(countY, x, "C")
	x += 3

	for m := 1; m < 13; m++ {
		win.MovePrintf(monthY, x, "%02d", m)
		payCount := obj.PayCountByYearMonth(year, m)
		result.PaymentCount += payCount
		win.MovePrintf(countY, x, "%2d", payCount)
		var graphY int = countY - 1

		for count := 0; count < payCount; count++ {
			win.MovePrint(graphY, x+1, "+")
			graphY--

			if graphY < 1 {
				break
			}
		}

		x += offsetX
	}

	win.MovePrintf(countY, x-(offsetX/2), ":%d", result.PaymentCount)
	return result
}

/* Draw info about payments for given year */
func DrawInfoByYear(win *goncurses.Window, sizeX, sizeY int, yearInfo YearInfo) {
	var y int = 1
	win.MovePrint(0, sizeX/2-2, "|Info|")
	win.MovePrintf(y, sizeX/3, "Year: %d", yearInfo.Year)
	y++
	win.MovePrintf(y, 1, "Payments count: %d", yearInfo.PaymentCount)
	y++
}

/* Draw a list of all bonds in own window, with own input processing */
func DrawListBonds(bondsArr *bonds.Bonds, sizeY, sizeX, posY, posX int) {
	win, err := goncurses.NewWindow(sizeY, sizeX, posY, posX)

	if err != nil {
		Terminal.Print(err.Error())
		return
	}

	defer win.Delete()
	bondsTable := make([]string, 0, len(bondsArr.Bonds))

	for id, obj := range bondsArr.Bonds {
		tmp := fmt.Sprintf("%d. Name:'%s' Coupon remaining:'%d', Near payday:(%02d.%02d.%d)", id, obj.Name, obj.CouponCount, obj.CouponNearPayDate.Day(), obj.CouponNearPayDate.Month(), obj.CouponNearPayDate.Year())
		bondsTable = append(bondsTable, tmp)
	}

	var startInd int = 0

	printSlice := func(startInd int) int {
		var endInd int = min(len(bondsTable)-1, sizeY-1)

		if startInd >= len(bondsTable) {
			startInd = len(bondsTable) - 1
		}

		if startInd < 0 {
			startInd = 0
		}

		var x int = 1
		var y int = 1

		for ind := startInd; ind <= endInd; ind++ {
			win.MovePrint(y, x, bondsTable[ind])
			y++
		}

		return startInd
	}

	var input goncurses.Key

	for input != ExitKey {
		win.Clear()
		win.Box(0, 0)
		win.MovePrint(0, sizeX/2-6, "|Bonds List|")
		win.Refresh()

		switch input {
		case ScrollUpKey:
			// startInd = printSlice(startInd - 1)
			startInd -= 1

		case ScrollDownKey:
			// startInd = printSlice(startInd + 1)
			startInd += 1
		}

		startInd = printSlice(startInd)
		input = win.GetChar()
	}
}

/* Ask user for bonds params and create new one */
func CreateBondsByUser() (*bonds.BondsData, error) {
	Terminal.Print("***Bonds Create***")
	question := "Name: "
	name, err := Terminal.AskString(question)

	if err != nil {
		return nil, err
	}

	question = "Coupons count: "
	couponCount, err := Terminal.AskInt(question)

	if err != nil {
		return nil, err
	}

	question = "Nearest pay day[dd.mm.yyyy]: "
	couponNearestPayDate, err := Terminal.AskDate(question, DefaultDateLayout)

	if err != nil {
		return nil, err
	}

	question = "Next pay day[dd.mm.yyyy]: "
	couponNextPayDate, err := Terminal.AskDate(question, DefaultDateLayout)

	Terminal.Print("******************")

	if err != nil {
		return nil, err
	}

	result := bonds.BondsDataNew()
	result.Name = name
	result.CouponCount = couponCount
	result.CouponPeriod = bonds.CouponPeriodCreate(couponNearestPayDate, couponNextPayDate)
	result.CouponNearPayDate = couponNearestPayDate
	return result, nil
}

// TODO:
// [ ] add more info in BoundsData
// [x] by key show list of all bonds in window (format: "index | Name | ...")
// [ ]      add command 'delete <index>'

// TODO:
// [ ] Create a small system for hold windows data/functions and call them
//      usefull for substitute window
//      implement this as finite state machine
//      at register pass key and window for return into it by key pressing
//          Win.AddWay(key, window) -> add a way into self.WayMap -> at each key pressing check is key in WayMap and return window if yes
//      Each window wrap must implement function for input commads ':<command>' and execute them
// [ ] create different windows
//     [x] for graph
//     [1/2] for info by year
//         payments count
//         payments in roubles

// TODO:
// [1/2] by key save bonds data in json in default path
// [ ] Create way to init default path and save into it
//      for now - ask for path and save into it

func main() {
	stdscr, err := goncurses.Init()
	goncurses.Echo(false)
	goncurses.Cursor(0)

	if err != nil {
		panic(err)
	}

	defer goncurses.End()
	MaxY, MaxX = goncurses.StdScr().MaxYX()
	mainHeight, mainWidth := MaxY-1, (MaxX/3)*2
	mainPosY, mainPosX := 0, 0
	main, err := goncurses.NewWindow(mainHeight, mainWidth, mainPosY, mainPosX)

	if err != nil {
		panic(err)
	}

	defer main.Delete()
	infoHeight, infoWidth := MaxY/2, (MaxX/3)*1
	infoPosY, infoPosX := 0, (MaxX/3)*2
	info, err := goncurses.NewWindow(infoHeight, infoWidth, infoPosY, infoPosX)

	if err != nil {
		panic(err)
	}

	err = Terminal.Init(terminal.TerminalSettings{
		Title:         "|Terminal|",
		SizeY:         infoHeight - 1,
		SizeX:         infoWidth,
		PosX:          infoPosX,
		PosY:          infoHeight,
		DefaultEcho:   false,
		DefaultCursor: 0,
	})

	if err != nil {
		panic(err)
	}

	defer Terminal.Delete()
	Terminal.Print("Inited successfully. Type ':help' for info")
	var year int = CurrentYear
	var graphOffsetX int = (mainWidth - 6) / 12
	var input goncurses.Key
	var loop bool = true

	for loop {
		main.Clear()
		main.Box(0, 0)
		info.Box(0, 0)

		switch input {
		case ExitKey:
			{
				tmp := Terminal.AskChar("Really exit?[y/n]")

				if tmp == 'y' || tmp == 'Y' {
					loop = false
					continue
				}
			}

		case HelpKey:
			HelpWindow()

		case IncreaseYearKey:
			year++

		case DecreaseYearKey:
			year--

			if year < CurrentYear {
				year = CurrentYear
			}

		case AppendBondsKey:
			{
				data, err := CreateBondsByUser()

				if err != nil {
					Terminal.Print(err.Error())

				} else {
					AllBonds.Append(data)
				}

			}

		case SaveBondsKey:
			{
				msg := "Filename for save: "
				filename, err := Terminal.AskString(msg)

				if err != nil {
					Terminal.Print(err.Error())
					continue
				}

				err = AllBonds.SaveToFile(filename)

				if err != nil {
					Terminal.Print(err.Error())
				}
			}

		case LoadBondsKey:
			{
				msg := "Filename for load: "
				filename, err := Terminal.AskString(msg)

				if err != nil {
					Terminal.Print(err.Error())
					continue
				}

				err = AllBonds.LoadFromFile(filename)

				if err != nil {
					Terminal.Print(err.Error())
				}

			}

		case ListBondsKey:
			DrawListBonds(AllBonds, mainHeight, mainWidth, mainPosY, mainPosX)

		case StartOfCommandKey:
			{
				command, err := Terminal.AskString("")

				if err != nil {
					Terminal.Print(err.Error())
					continue
				}

				ExecuteCommand(command)
			}
		}

		yearInfo := DrawGraphByYear(AllBonds, year, main, MaxX, MaxY-2, graphOffsetX)
		DrawInfoByYear(info, infoWidth, infoHeight, yearInfo)

		stdscr.MovePrintf(MaxY-1, 0, "Help:%c ", HelpKey)
		stdscr.Printf("Exit:%c ", ExitKey)
		stdscr.Printf("Prev year:%c ", DecreaseYearKey)
		stdscr.Printf("Next year:%c ", IncreaseYearKey)

		stdscr.Refresh()
		main.Refresh()
		info.Refresh()
		Terminal.Refresh()

		input = stdscr.GetChar()

	}
}

func init() {
	RegisterCommand("help", Command{"':help <command>'-Show info about commands", CommandHelp})
	RegisterCommand("list", Command{"Show list of all bonds", CommandList})
	RegisterCommand("save", Command{"'save <file>' - Save bonds info into file", CommandSave})
	RegisterCommand("load", Command{"'load <file>' - Load bonds info from file", CommandLoad})
	RegisterCommand("delete", Command{"'delete <index>' - Delete bonds info from list", CommandDelete})
}
