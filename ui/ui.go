package ui

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/nsf/termbox-go"
)

var (
	modkernel32         = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleTitle = modkernel32.NewProc("SetConsoleTitleW")
)

func Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetOutputMode(termbox.Output256)
	return nil
}

func ClearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func PrintMenu(menuItems []string, selectedItem int) {
	for i, item := range menuItems {
		x, y := 2, i+2
		fg, bg := termbox.ColorWhite, termbox.ColorBlack
		if i == selectedItem {
			fg, bg = termbox.ColorBlack, termbox.ColorWhite
		}
		for _, char := range item {
			termbox.SetCell(x, y, char, fg, bg)
			x++
		}
	}
	termbox.Flush()
}

func PrintPlayfield(playfield *[20][30]int) {
	for i := 0; i < len(playfield); i++ {
		for j := 0; j < len(playfield[i]); j++ {
			x, y := j*2, i
			switch playfield[i][j] {
			case 0:
				termbox.SetCell(x, y, '█', termbox.ColorMagenta, termbox.ColorMagenta)
				termbox.SetCell(x+1, y, '█', termbox.ColorMagenta, termbox.ColorMagenta)
			default:
				termbox.SetCell(x, y, '█', termbox.ColorWhite, termbox.ColorBlack)
				termbox.SetCell(x+1, y, '█', termbox.ColorWhite, termbox.ColorBlack)
			}
		}
	}
	termbox.Flush()
}

func PrintInfoOnScreenAtXY(text string, x int, y int) {
	w, h := termbox.Size()
	textLen := len([]rune(text))

	if x < 0 {
		x = 0
	} else if x > w-textLen {
		x = w - textLen
	}

	if y < 0 {
		y = 0
	} else if y > h-1 {
		y = h - 1
	}

	for _, c := range text {
		termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
	termbox.Flush()
}

func PrintInfoOnScreen(text string) {
	w, h := termbox.Size()
	x, y := w/2-len(text)/2, h/2

	for _, c := range text {
		termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
	termbox.Flush()
}

func PrintDebugInfo(title string) {
	titlePtr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ret, _, err := procSetConsoleTitle.Call(uintptr(unsafe.Pointer(titlePtr)))
	if ret == 0 {
		fmt.Println("Error:", err)
	}
}
