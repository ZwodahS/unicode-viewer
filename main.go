package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/zwodahs/tcell-region"
)

func min(v1, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

func max(v1, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func isControl(value int) bool {
	return (value >= 0x0000 && value <= 0x001F) || value == 0x007F || (value >= 0x0080 && value <= 0x009F)
}

func redraw(columnRegion, rowRegion, runeRegion *tcellr.Region, page, selectedX, selectedY int) {
	columnRegion.Fill(' ')
	rowRegion.Fill(' ')
	runeRegion.Fill(' ')
	runeStart := page * 16 * 16 * 2
	for i, c := runeStart, 0; i < runeStart+(32*32); i, c = i+1, c+1 {
		x := (c % 16)
		y := c / 16
		drawX := x * 2
		style := tcell.StyleDefault
		if x == selectedX && y == selectedY {
			style = style.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)
		}
		if isControl(i) {
			runeRegion.SetCell(drawX, y, ' ', style)
			continue
		}
		value := fmt.Sprintf("%U", i)[2:]
		quoted := "'\\u" + value + "'"
		c, err := strconv.Unquote(quoted)
		if err != nil {
			debug(err)
			panic(err)
		}
		for _, r := range c {
			if runewidth.RuneWidth(r) == 1 {
				runeRegion.SetCell(drawX, y, r, style)
				if x == selectedX && y == selectedY {
					runeRegion.SetCell(40, 3, r)
				}
			}
			break
		}
	}

	for i, c := runeStart/16, 0; c < 32; i, c = i+1, c+1 {
		value := fmt.Sprintf("%X", i)
		style := tcell.StyleDefault.Foreground(tcell.ColorRed)
		if c == selectedY {
			style = style.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)
		}
		rowRegion.SetTextRight(7, c, value, style)
	}

	columnRegion.SetText(0, 0, "0 1 2 3 4 5 6 7 8 9 A B C D E F", tcell.StyleDefault.Foreground(tcell.ColorRed))
	for i := 0; i < 16; i++ {
		style := tcell.StyleDefault.Foreground(tcell.ColorRed)
		if i == selectedX {
			style = style.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)
		}
		columnRegion.SetText(i*2, 0, fmt.Sprintf("%X", i), style)
	}

}

func _main() int {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	// set up tcell
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()
	screen.Fill(' ', tcell.StyleDefault)

	// set up tcell region
	err = tcellr.InitRegion()
	if err != nil {
		fmt.Println("Error Init regions")
		return 1
	}

	// set up tcell event polling
	eventChannel := make(chan tcell.Event)
	go func() {
		for {
			eventChannel <- screen.PollEvent()
		}
	}()
	mainRegion := tcellr.NewRegion(screen, 60, 60)
	columnRegion := mainRegion.NewRegion(32, 1)
	columnRegion.SetPosition(15, 0)
	runeRegion := mainRegion.NewRegion(62, 32)
	runeRegion.SetPosition(15, 1)
	rowRegion := mainRegion.NewRegion(8, 32)
	rowRegion.SetPosition(6, 1)

	x, y := 0, 0

	page := 0
	redraw(columnRegion, rowRegion, runeRegion, page, x, y)

	mainRegion.Draw(0, 0)

	// Core loop
	exit := false

loop:
	for !exit {
		select {
		case event := <-eventChannel: // handle event
			switch ev := event.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlL:
					screen.Fill(' ', tcell.StyleDefault)
					mainRegion.Fill(' ')
					runeRegion.Fill(' ')
				default:
					switch ev.Rune() {
					case '>':
						page += 1
						redraw(columnRegion, rowRegion, runeRegion, page, x, y)
					case '<':
						if page != 0 {
							page -= 1
							redraw(columnRegion, rowRegion, runeRegion, page, x, y)
						}
					case 'h':
						x = max(0, x-1)
						redraw(columnRegion, rowRegion, runeRegion, page, x, y)
					case 'j':
						y = min(31, y+1)
						redraw(columnRegion, rowRegion, runeRegion, page, x, y)
					case 'k':
						y = max(0, y-1)
						redraw(columnRegion, rowRegion, runeRegion, page, x, y)
					case 'l':
						x = min(15, x+1)
						redraw(columnRegion, rowRegion, runeRegion, page, x, y)
					default:
					}
				}
				if ev.Key() == tcell.KeyEscape {
					break loop
				}
			case *tcell.EventResize:
			}
			mainRegion.SetText(0, 0, "Page: "+strconv.Itoa(page))
			mainRegion.SetText(0, 1, fmt.Sprintf("%d %d", x, y))
			mainRegion.Draw(0, 0)
			screen.Show()
		}
	}

	// allow for defer methods to be run
	return 0
}

func main() { os.Exit(_main()) }
