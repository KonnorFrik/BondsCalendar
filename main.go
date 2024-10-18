package main

import (
	"fmt"
	// "os"
	"bonds_payment_calendar/bonds"
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

	AllBonds = bonds.BondsNew()
)

const (
	ExitKey         = 'q'
	HelpKey         = 'h'
	IncreaseYearKey = '>'
	DecreaseYearKey = '<'
	AppendBondsKey  = 'a'
	SaveBondsKey    = 's'
	LoadBondsKey    = 'l'

	// DefaultTimeout = 500
)

func PrintSimpleDate(obj time.Time) {
	fmt.Printf("%02d.%02d.%d\n", obj.Day(), obj.Month(), obj.Year())
}

func HelpWindow() {
	startY, startX := MaxY/4, MaxX/4
	height, width := MaxY/2, MaxX/2
	scr, err := goncurses.NewWindow(height, width, startY, startX)

	if err != nil {
		// fmt.Fprint(os.Stderr, "Can't create Help window:", err)
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

		input = scr.GetChar()
	}
}

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

			if graphY <= 0 {
				break
			}
		}

		x += offsetX
	}

	win.MovePrintf(countY, x-(offsetX/2), ":%d", result.PaymentCount)
	return result
}

func DrawInfoByYear(win *goncurses.Window, sizeX, sizeY int, yearInfo YearInfo) {
	var y int = 1
	win.MovePrint(0, sizeX/2-2, "|Info|")

	win.MovePrintf(y, sizeX/3, "Year: %d", yearInfo.Year)
	y++
	win.MovePrintf(y, 1, "Payments count: %d", yearInfo.PaymentCount)
	y++
}

func CreateBondsByUser() (*bonds.BondsData, error) {
	winHeight, winWidth := MaxY/2, MaxX/2
	winPosY, winPosX := MaxY/4, MaxX/4
	win, err := goncurses.NewWindow(winHeight, winWidth, winPosY, winPosX)

	if err != nil {
		return nil, err
	}

	defer win.Delete()
	goncurses.Echo(true)
	goncurses.Cursor(1)
	win.Box(0, 0)
	win.MovePrint(0, winWidth/2-6, "|Bonds Create|")

	var x int = 2
	var y int = 1

	question := "Name: "
	name, err := AskString(win, y, x, question, winWidth-2-len(question))

	if err != nil {
		return nil, err
	}

	y++
	question = "Coupons count: "
	couponCount, err := AskInt(win, y, x, question, winWidth-2-len(question))

	if err != nil {
		return nil, err
	}

	y++
	question = "Nearest pay day[dd.mm.yyyy]: "
	couponNearestPayDate, err := AskDate(win, y, x, question, winWidth-2-len(question), "02.01.2006")

	if err != nil {
		return nil, err
	}

	y++
	question = "Next pay day[dd.mm.yyyy]: "
	couponNextPayDate, err := AskDate(win, y, x, question, winWidth-2-len(question), "02.01.2006")

	if err != nil {
		return nil, err
	}

	goncurses.Echo(false)
	goncurses.Cursor(0)

	result := bonds.BondsDataNew()
	result.Name = name
	result.CouponCount = couponCount
	result.CouponPeriod = bonds.CouponPeriodCreate(couponNearestPayDate, couponNextPayDate)
	result.CouponNearPayDate = couponNearestPayDate
	return result, nil
}

// TODO:
// [ ] add more info in BoundsData
// [ ] by key show list of all bonds in window (format: "index | Name | ...")
// [ ]      in list window by key delete choosen bonds

// TODO:
// [ ] create different windows
//     [x] for graph
//     [ ] for info by year
//         payments count
//         payments in roubles
//     [ ] for log?

// TODO:
// [1/2] by key save bonds data in json in default path
//      for now - ask for path and save into it

// TODO:
// [ ] create terminal window (and maybe struct for wrap window)
//      ask all input through terminal
//      print all error, output to terminal

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
	infoHeight, infoWidth := MaxY-1, (MaxX/3)*1
	infoPosY, infoPosX := 0, (MaxX/3)*2
	info, err := goncurses.NewWindow(infoHeight, infoWidth, infoPosY, infoPosX)

	if err != nil {
		panic(err)
	}

	// data := bonds.BondsDataNew()
	// data.Name = "test";
	// data.CouponCount = 2;
	// data.CouponNearPayDate = bonds.CouponPayDay(2024, 10, 31);
	// data.CouponPeriod = bonds.CouponPeriodCreate(data.CouponNearPayDate, bonds.CouponPayDay(2025, 5, 2));
	//
	// AllBonds.Append(data)
	// AllBonds.Append(data)

	var year int = CurrentYear
	var graphOffsetX int = (mainWidth - 6) / 12
	// fmt.Fprintf(os.Stderr, "graphOffsetX: %d\n", graphOffsetX)

	var input goncurses.Key
	var loop bool = true

	for loop {
		main.Clear()
		main.Box(0, 0)
		info.Box(0, 0)

		switch input {
		case ExitKey:
			{
				tmp := PopUpText(mainHeight/2-3, mainWidth/2-5, "Really exit?[y/n]")

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
					PopUpText(MaxY/2, MaxX/2-5, err.Error())

				} else {
					AllBonds.Append(data)
				}

			}

		case SaveBondsKey:
            {
                msg := "Filename for save: "
                x := MaxX / 2 - len(msg) - 6

                goncurses.Echo(true)
                goncurses.Cursor(1)
                filename, err := PopUpAskString(MaxY / 2 - 1, x, msg, len(msg) + 2 + 10)
                goncurses.Echo(false)
                goncurses.Cursor(0)

                if err != nil {
                    PopUpText(MaxY / 2, MaxX / 2, err.Error())
                    continue
                }

                err = AllBonds.SaveToFile(filename)

                if err != nil {
                    PopUpText(MaxY / 2, MaxX / 2, err.Error())
                }
            }

        case LoadBondsKey:
            {
                msg := "Filename for load: "
                x := MaxX / 2 - len(msg) - 6
                goncurses.Echo(true)
                goncurses.Cursor(1)
                filename, err := PopUpAskString(MaxY / 2 - 1, x, msg, len(msg) + 2 + 10)
                goncurses.Echo(false)
                goncurses.Cursor(0)

                if err != nil {
                    PopUpText(MaxY / 2, MaxX / 2, err.Error())
                    continue
                }

                err = AllBonds.LoadFromFile(filename)

                if err != nil {
                    PopUpText(MaxY / 2, MaxX / 2, err.Error())
                }

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
		input = stdscr.GetChar()
	}
}
