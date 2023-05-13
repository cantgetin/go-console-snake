package game

import (
	"console-snake/ui"
	"time"

	"github.com/nsf/termbox-go"
)

type State int
type UserInput int
type SnakeDirection int

// User input enum
const (
	LeftInput  UserInput = 0
	RightInput           = 1
	DownInput            = 2
	UpInput              = 3
	NoInput              = 4
)

const (
	LeftDirection  SnakeDirection = 0
	RightDirection                = 1
	DownDirection                 = 2
	UpDirection                   = 3
)

type Game struct {
	playfield      *[10][10]int
	snakeLength    int
	snakePosition  [][2]int
	userInput      UserInput
	fruitPosition  []int
	snakeDirection SnakeDirection
}

func Start(menu func()) {
	ui.ClearScreen()

	alive := true

	game := Game{
		playfield:     new([10][10]int),
		snakeLength:   5,
		snakePosition: [][2]int{{2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5}},
		fruitPosition: []int{5, 6},
		userInput:     NoInput}

	eventChan := make(chan termbox.Event)
	go func() {
		for {
			event := termbox.PollEvent()
			eventChan <- event
		}
	}()

	tickDuration := 500 * time.Millisecond
	tickTimer := time.NewTimer(tickDuration)

	for alive {
		// game logic
		gameTick(&game)
		// draw the current state of 2-dimensional array with colors
		ui.PrintPlayfield(game.playfield)

		// wait for a short period of time and reset the timer for the next tick
		<-tickTimer.C
		tickTimer.Reset(tickDuration)

		//ui.PrintDebugInfo("y:" + strconv.Itoa(game.blockPosition[0]) + " x: " + strconv.Itoa(game.blockPosition[1]))
		alive = handleUserInput(&game, eventChan)
	}
	menu()
}

func gameTick(game *Game) {
	// move snake 1 unit to the direction it's facing
	// add block next to the head of snake
	for i := range game.snakePosition {
		pos := game.snakePosition[i]
		game.playfield[pos[0]][pos[1]] = 1
	}

	blockToAppend := game.snakePosition[len(game.snakePosition)-1]

	switch game.snakeDirection {
	case LeftDirection:
		blockToAppend[0] -= 1
	case RightDirection:
		blockToAppend[0] += 1
	case DownDirection:
		blockToAppend[1] -= 1
	case UpDirection:
		blockToAppend[0] += 1
	}

	game.snakePosition = append(game.snakePosition, blockToAppend)
}

func handleUserInput(game *Game, eventChan chan termbox.Event) bool {
	// check for user input
	select {
	case event := <-eventChan:
		for len(eventChan) > 0 {
			<-eventChan
		}

		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeyArrowLeft:
				game.userInput = LeftInput
				game.snakeDirection = LeftDirection
			case termbox.KeyArrowRight:
				game.userInput = RightInput
				game.snakeDirection = RightDirection
			case termbox.KeyArrowDown:
				game.userInput = DownInput
				game.snakeDirection = DownDirection
			case termbox.KeyArrowUp:
				game.userInput = UpInput
				game.snakeDirection = UpDirection
			case termbox.KeyEsc:
				return false
			default:
				game.userInput = NoInput
			}
		}
	default:
		// no event waiting, continue the game
	}
	return true
}
