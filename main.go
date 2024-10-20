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
	ScrollUpKey       = 'w'
	ScrollDownKey     = 's'
	StartOfCommandKey = ':'

	DefaultDateLayout = "02.01.2006"
)

func CommandHelp(args []string) error {
	if len(args) == 0 {
        for key, val := range CommandTable {
            if key == "help" {
                continue
            }

            Terminal.Print(val.Info)
        }

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
    err := DrawListBonds(AllBonds, MaxY-1, 0, 0)
	return err
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

func CommandNewBonds(args []string) error {
    data, err := CreateBondsByUser()

    if err != nil {
        return err

    } 

    AllBonds.Append(data)
    return nil
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

    var obj *bonds.BondsData
    var exist bool

    if index < len(AllBonds.Bonds) && index >= 0 {
        obj = AllBonds.Bonds[index]
        exist = true
    }

    AllBonds.Bonds, err = SliceRemoveByIndex(AllBonds.Bonds, index)

    if exist && err == nil {
        Terminal.Print(fmt.Sprintf("%d. %s - deleted", index, obj.Name))
    }

	return err
}

/*
Draw graph of payments for given year 
Called by main.Draw
*/
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

/*
Draw info about payments for given year 
Called by independent sub-window - info
*/
func DrawInfoByYear(win *goncurses.Window, sizeX, sizeY int, yearInfo YearInfo) {
	var y int = 1
	win.MovePrint(0, sizeX/2-2, "|Info|")
	win.MovePrintf(y, sizeX/3, "Year: %d", yearInfo.Year)
	y++
	win.MovePrintf(y, 1, "Payments count: %d", yearInfo.PaymentCount)
	y++
}

/* Draw a list of all bonds as scrollable pop up window */
func DrawListBonds(bondsArr *bonds.Bonds, sizeY, posY, posX int) error {
	bondsTable := make([]string, 0, len(bondsArr.Bonds))
    var format string = "%d. Name:'%s' Coupon remaining:'%d', Near payday:(%02d.%02d.%d), ~PeriodDays(%d)"

	for id, obj := range bondsArr.Bonds {

		tmp := fmt.Sprintf(
            format,
            id,
            obj.Name,
            obj.CouponCount,
            obj.CouponNearPayDate.Day(),
            obj.CouponNearPayDate.Month(),
            obj.CouponNearPayDate.Year(),
            obj.CouponPeriod,
        )
		bondsTable = append(bondsTable, tmp)
	}

    return PopUpScrollableList(bondsTable, "|Bonds List|", sizeY, posY, posX)
}

/* Ask user for bonds params and create new one */
func CreateBondsByUser() (*bonds.BondsData, error) {
	Terminal.Print("***Bonds Create***")
	name, err := Terminal.AskString("Name: ")

	if err != nil {
		return nil, err
	}

	couponCount, err := Terminal.AskInt("Coupons count: ")

	if err != nil {
		return nil, err
	}

	couponNearestPayDate, err := Terminal.AskDate("Nearest pay day[dd.mm.yyyy]: ", DefaultDateLayout)

	if err != nil {
		return nil, err
	}

	couponNextPayDate, err := Terminal.AskDate("Next pay day[dd.mm.yyyy](can be empty): ", DefaultDateLayout)
	Terminal.Print("******************")

	if err != nil {
        if couponCount == 1 {
            couponNextPayDate = time.Time{}
        } else {
            return nil, err
        }
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

// TODO:
// [1/2] by key save bonds data in json in default path
// [ ] Create way to init default path and save into it
//      for now - ask for path and save into it

func main() {
	var year int = CurrentYear

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
    main, err := FSMWindowNew(mainHeight, mainWidth, mainPosY, mainPosX)

	if err != nil {
		panic(err)
	}

	defer main.FreeWindow()
	var graphOffsetX int = (mainWidth - 6) / 12
    var focus *FSMWindow = main
    main.SetTitle("|Main|").RegisterInput(IncreaseYearKey, func() bool {
        year++
        return true
    }).RegisterInput(DecreaseYearKey, func() bool {
        year--

        if year < CurrentYear {
            year = CurrentYear
        }
        
        return true
    }).RegisterInput(ExitKey, func() bool {
        tmp := Terminal.AskChar("Really exit?[y/n]")

        if tmp == 'y' || tmp == 'Y' {
            return false
        }

        return true
    }).RegisterInput(StartOfCommandKey, func() bool {
        command, err := Terminal.AskString("")

        if err != nil {
            Terminal.Print(err.Error())
            return true
        }

        ExecuteCommand(command)
        return true
    }).RegisterInput(HelpKey, func() bool {
        arr := []string{
            fmt.Sprintf("  Keys for graph window"),
            fmt.Sprintf("%c - Exit from programm, or close sub-window", ExitKey),
            fmt.Sprintf("%c - Show next year info", IncreaseYearKey),
            fmt.Sprintf("%c - Show previous year info", DecreaseYearKey),
            fmt.Sprintf("%c - Start write command to terminal", StartOfCommandKey),
            fmt.Sprintf("%c - Show this window", HelpKey),
        }

        err := PopUpScrollableList(arr, "|Info|", main.SizeY, main.PosY, main.posX)

        if err != nil {
            Terminal.Print(err.Error())
        }

        return true
    })

    yearInfo := YearInfo{CurrentYear, 0}
    main.SetCustomDraw(func() {
		yearInfo = DrawGraphByYear(AllBonds, year, main.Window, MaxX, MaxY-2, graphOffsetX)
    })

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
	var loop bool = true

	for loop {
		info.Box(0, 0)

        focus.Draw() // this must be before DrawInfoByYear
		DrawInfoByYear(info, infoWidth, infoHeight, yearInfo)

		stdscr.MovePrintf(MaxY-1, 0, "Help:%c ", HelpKey)
		stdscr.Printf("Exit:%c ", ExitKey)
		stdscr.Printf("Prev year:%c ", DecreaseYearKey)
		stdscr.Printf("Next year:%c ", IncreaseYearKey)

		stdscr.Refresh()
		Terminal.Refresh()
        focus.DrawBox()
		info.Refresh()

        var isWork bool
        focus, isWork = focus.Input()

        if !isWork {
            loop = false
            continue
        }

	}
}

func init() {
	RegisterCommand("help", Command{"':help <command>' - Show info about commands", CommandHelp})
    RegisterCommand("list", Command{"':list' - Show list of all bonds", CommandList})
    RegisterCommand("save", Command{"':save <file>' - Save bonds info into file", CommandSave})
    RegisterCommand("load", Command{"':load <file>' - Load bonds info from file", CommandLoad})
    RegisterCommand("new", Command{"':new' - Create new bonds and append it into list", CommandNewBonds})
    RegisterCommand("delete", Command{"':delete <index>' - Delete bonds info from list", CommandDelete})
}
