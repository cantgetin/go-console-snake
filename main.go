package main

import (
	"console-snake/game"
	"console-snake/ui"

	"github.com/nsf/termbox-go"
)

func main() {
	menu()
}

func menu() {
	menuItems := []string{"Start", "About", "Exit"}
	selectedItem := 0
	alive := true

	ui.Init()
	ui.ClearScreen()
	ui.PrintMenu(menuItems, selectedItem)

	for alive {
		// wait for user input
		event := termbox.PollEvent()

		switch event.Key {
		case termbox.KeyArrowUp:
			if selectedItem > 0 {
				selectedItem--
			}

			ui.ClearScreen()
			ui.PrintMenu(menuItems, selectedItem)
		case termbox.KeyArrowDown:
			if selectedItem < len(menuItems)-1 {
				selectedItem++
			}

			ui.ClearScreen()
			ui.PrintMenu(menuItems, selectedItem)
		case termbox.KeyEnter:
			switch selectedItem {
			case 0:
				game.Start(menu)
			case 1:
				about()
			case 2:
				alive = false
			}
		}
	}
}

func about() {
	ui.ClearScreen()
	ui.PrintInfoOnScreen("Golang tetris by cantgetin, 2023")
	for {
		event := termbox.PollEvent()

		switch event.Key {
		case termbox.KeyEsc:
			menu()
		}
	}
}
