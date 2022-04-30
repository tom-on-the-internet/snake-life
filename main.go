package main

import (
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type game struct {
	score  int
	snake  *snake
	mice   []position
	maxPos position
}

type snake struct {
	body []position
}

type position [2]int

func main() {
	game := newGame()
	game.beforeGame()

	mouseTurn := false

	for {
		mouseTurn = !mouseTurn
		if rand.Intn(100) > 90 {
			game.mice = append(game.mice, wallPosition())
		}
		maxX, maxY := getSize()
		game.maxPos[0] = maxX
		game.maxPos[1] = maxY

		// each mouse should move
		// randomly if snake not near
		// away if snake is near

		// mice should move first, if their turn

		// calculate new head position
		if mouseTurn {
			game.moveMice()
		}

		game.moveSnake()

		// ate mouse check
		// remember he could eat another mouse on the way
		// ateFood := positionsAreSame(game.food, newHeadPos)
		if false {
			game.score++
			// game.placeNewFood()
		} else {
			// game.snake.body = game.snake.body[:len(game.snake.body)-1]
		}

		game.draw()
	}
}

func newGame() *game {
	rand.Seed(time.Now().UnixNano())

	snake := newSnake()

	game := &game{
		score: 0,
		snake: snake,
	}

	return game
}

func positionsAreSame(a, b position) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func randomPosition() position {
	width, height := getSize()
	x := rand.Intn(width) + 1
	y := rand.Intn(height) + 2

	return [2]int{x, y}
}

func moveTowardPosition(pos, destPos position) position {
	if pos[0] < destPos[0] {
		pos[0]++
	} else if pos[0] > destPos[0] {
		pos[0]--
	}

	if pos[1] < destPos[1] {
		pos[1]++
	} else if pos[1] > destPos[1] {
		pos[1]--
	}

	return pos
}

func moveAwayFromPosition(pos, destPos position) position {
	if pos[0] < destPos[0] {
		pos[0]--
	} else if pos[0] > destPos[0] {
		pos[0]++
	}

	if pos[1] < destPos[1] {
		pos[1]--
	} else if pos[1] > destPos[1] {
		pos[1]++
	}

	return pos
}

func wallPosition() position {
	maxX, maxY := getSize()

	switch rand.Intn(4) {
	case 0:
		// top
		return position{rand.Intn(maxX), 1}
	case 1:
		// bottom
		return position{rand.Intn(maxX), maxY}
	case 2:
		// left
		return position{1, rand.Intn(maxY)}
	default:
		// right
		return position{maxX, rand.Intn(maxY)}
	}
}

func newSnake() *snake {
	maxX, maxY := getSize()
	pos := position{maxX / 2, maxY / 2}

	return &snake{
		body: []position{pos},
	}
}

func (g *game) moveMice() {
	mice := []position{}
	head := g.snake.body[0]

	for _, mouse := range g.mice {
		curDist := ((head[0] - mouse[0]) * (head[0] - mouse[0])) + ((head[1] - mouse[1]) * (head[1] - mouse[1]))
		if curDist > 100 && rand.Intn(8) != 7 {
			mouse = moveTowardPosition(mouse, randomPosition())
		} else {
			mouse = moveAwayFromPosition(mouse, head)
		}

		if mouse[0] == 0 {
			mouse[0] = 1
		}
		if mouse[1] == 0 {
			mouse[1] = 1
		}
		if mouse[0] > g.maxPos[0] {
			mouse[0] = g.maxPos[0] - 1
		}
		if mouse[1] > g.maxPos[1] {
			mouse[1] = g.maxPos[1] - 1
		}
		mice = append(mice, mouse)
	}

	// each mouse should move randomly
	// if it is near the snake, move away
	g.mice = mice
}

func (g *game) moveSnake() {
	// if no mice, move randomly
	var destPos position

	if len(g.mice) > 0 {
		var closest position

		dist := (g.maxPos[0] * g.maxPos[0]) + (g.maxPos[1] * g.maxPos[1])

		head := g.snake.body[0]

		for _, mouse := range g.mice {
			curDist := ((head[0] - mouse[0]) * (head[0] - mouse[0])) + ((head[1] - mouse[1]) * (head[1] - mouse[1]))
			if curDist < dist {
				dist = curDist
				closest = mouse
			}
		}

		destPos = closest
	} else {
		destPos = randomPosition()
	}

	curPos := g.snake.body[0]
	newPos := moveTowardPosition(curPos, destPos)

	g.snake.body = append([]position{newPos}, g.snake.body...)
	g.snake.body = g.snake.body[:len(g.snake.body)-1]

	mice := []position{}

	for _, mouse := range g.mice {
		if positionsAreSame(g.snake.body[0], mouse) {
			var newLastPos position

			head := g.snake.body[0]
			lastPos := g.snake.body[len(g.snake.body)-1]

			if head[0] > lastPos[0] {
				newLastPos[0] = lastPos[0] - 1
			} else {
				newLastPos[0] = lastPos[0] + 1
			}

			if head[1] > lastPos[1] {
				newLastPos[1] = lastPos[1] - 1
			} else {
				newLastPos[1] = lastPos[1] + 1
			}

			g.snake.body = append(g.snake.body, newLastPos)
			g.score++
		} else {
			mice = append(mice, mouse)
		}
	}

	g.mice = mice
}

func (g *game) draw() {
	clear()
	maxX, _ := getSize()

	status := "mice eaten: " + strconv.Itoa(g.score)
	statusXPos := maxX/2 - len(status)/2

	moveCursor(position{statusXPos, 0})
	draw(status)

	for _, mouse := range g.mice {
		moveCursor(mouse)
		draw("🐁")
	}

	for i := len(g.snake.body) - 1; i >= 0; i-- {
		pos := g.snake.body[i]

		moveCursor(pos)

		if i == 0 {
			draw("🟢")
		} else {
			green()
			draw("O")
			resetColor()
		}
	}

	render()
	time.Sleep(time.Millisecond * 60)
}

func (g *game) beforeGame() {
	hideCursor()

	// handle CTRL C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			g.over()
		}
	}()
}

func (g *game) over() {
	clear()
	showCursor()

	moveCursor(position{1, 1})
	draw("game over. the snake ate " + strconv.Itoa(g.score) + " mice.")

	render()

	os.Exit(0)
}
