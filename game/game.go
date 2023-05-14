package game

import (
	"console-snake/ui"
	"strconv"
	"time"

	"math/rand"

	"github.com/nsf/termbox-go"
)

type State int
type UserInput int
type SnakeDirection int
type GameState int

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
	alive          bool
}

func Start(menu func()) {
	ui.ClearScreen()

	game := Game{
		playfield:      new([20][30]int),
		snakeLength:    5,
		snakePosition:  Stack{},
		fruitPosition:  []int{5, 6},
		userInput:      NoInput,
		snakeDirection: RightDirection,
		alive:          true}

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

	for game.alive {
		// game logic
		gameTick(&game)
		// draw the current state of 2-dimensional array with colors
		ui.PrintPlayfield(game.playfield)

		// var str string
		// for i := range game.snakePosition.items {
		// 	str += strconv.Itoa(game.snakePosition.items[i][0]) + " " + strconv.Itoa(game.snakePosition.items[i][1]) + " "
		// }
		ui.PrintDebugInfo(strconv.FormatBool(game.alive))

		// wait for a short period of time and reset the timer for the next tick
		<-tickTimer.C
		tickTimer.Reset(tickDuration)

		if game.alive {
			handleUserInput(&game, eventChan)
		}
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

	// check newHead for collision
	checkCollision(game, newHead)

	// if newHead is on fruit position then we eat the fruit and respawn it
	if newHead[0] == game.fruitPosition[0] && newHead[1] == game.fruitPosition[1] {
		game.snakePosition.Push(newHead)
		respawnFruit(game)
	} else {
		// remove snake last block (tail), add new block (newHead)
		game.snakePosition.Push(newHead)
		game.snakePosition.Pop()
	}

	// clear gamefield

	for y := 0; y < len(game.playfield); y++ {
		for x := 0; x < len(game.playfield[y]); x++ {
			if game.playfield[y][x] == 1 {
				game.playfield[y][x] = 0
			}
		}
	}

	// draw snake

	for i := range game.snakePosition.items {
		pos := game.snakePosition.items[i]
		game.playfield[pos[0]][pos[1]] = 1
	}

	// draw fruit

	pos := game.fruitPosition
	game.playfield[pos[0]][pos[1]] = 2
}

func checkCollision(game *Game, newHead []int) {
	// check if we collide with border, if so then get our head to the opposite side (teleport)
	if newHead[0] > 19 {
		newHead[0] = 0
	}
	if newHead[0] < 0 {
		newHead[0] = 19
	}
	if newHead[1] > 29 {
		newHead[1] = 0
	}
	if newHead[1] < 0 {
		newHead[1] = 29
	}
	// if snake collides with its body then it's game over
	for i := range game.snakePosition.items {
		if newHead[0] == game.snakePosition.items[i][0] && newHead[1] == game.snakePosition.items[i][1] {
			gameOver(game)
			break
		}
	}
}

func gameOver(game *Game) {
	ui.PrintDebugInfo("game over")
	game.alive = false
}

func respawnFruit(game *Game) {

	newFruitPos := []int{rand.Intn(19), rand.Intn(29)}

	for i := range game.snakePosition.items {
		if newFruitPos[0] == game.snakePosition.items[i][0] && newFruitPos[1] == game.snakePosition.items[i][1] {
			respawnFruit(game)
			return
		}
	}

	game.fruitPosition = newFruitPos
}

func handleUserInput(game *Game, eventChan chan termbox.Event) {
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
				game.alive = false
			default:
				game.userInput = NoInput
			}
		}
	default:
		// no event waiting, continue the game
	}
}
