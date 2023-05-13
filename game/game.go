package game

import (
	"console-snake/ui"
	"strconv"
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
	playfield      *[20][30]int
	snakeLength    int
	snakePosition  Stack
	userInput      UserInput
	fruitPosition  []int
	snakeDirection SnakeDirection
}

func Start(menu func()) {
	ui.ClearScreen()

	alive := true

	game := Game{
		playfield:      new([20][30]int),
		snakeLength:    5,
		snakePosition:  Stack{},
		fruitPosition:  []int{5, 6},
		userInput:      NoInput,
		snakeDirection: RightDirection}

	snakeStartPos := [][]int{{2, 3}, {2, 4}, {2, 5}, {2, 6}, {2, 7}}
	for i := range snakeStartPos {
		game.snakePosition.Push([]int{snakeStartPos[i][0], snakeStartPos[i][1]})
	}

	eventChan := make(chan termbox.Event)
	go func() {
		for {
			event := termbox.PollEvent()
			eventChan <- event
		}
	}()

	tickDuration := 250 * time.Millisecond
	tickTimer := time.NewTimer(tickDuration)

	for alive {
		// game logic
		gameTick(&game)
		// draw the current state of 2-dimensional array with colors
		ui.PrintPlayfield(game.playfield)

		var str string
		for i := range game.snakePosition.items {
			str += strconv.Itoa(game.snakePosition.items[i][0]) + " " + strconv.Itoa(game.snakePosition.items[i][1]) + " "
		}

		ui.PrintDebugInfo(str)

		// wait for a short period of time and reset the timer for the next tick
		<-tickTimer.C
		tickTimer.Reset(tickDuration)

		alive = handleUserInput(&game, eventChan)
	}
	menu()
}

func gameTick(game *Game) {
	// move head 1 block to the direction it's facing
	head := game.snakePosition.items[len(game.snakePosition.items)-1]

	newHead := []int{head[0], head[1]}

	switch game.snakeDirection {
	case LeftDirection:
		newHead[1] -= 1
	case RightDirection:
		newHead[1] += 1
	case DownDirection:
		newHead[0] += 1
	case UpDirection:
		newHead[0] -= 1
	}

	// snake previous blocks go to the next block's place

	game.snakePosition.Push(newHead)
	game.snakePosition.Pop()

	for y := 0; y < len(game.playfield); y++ {
		for x := 0; x < len(game.playfield[y]); x++ {
			if game.playfield[y][x] == 1 {
				game.playfield[y][x] = 0
			}
		}
	}

	for i := range game.snakePosition.items {
		pos := game.snakePosition.items[i]
		game.playfield[pos[0]][pos[1]] = 1
	}
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
