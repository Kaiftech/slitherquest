package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 500
	screenHeight = 500
	cellSize     = 25
	gridCols     = screenWidth / cellSize
	gridRows     = screenHeight / cellSize
	startSpeed   = 10 // Starting speed (lower value = faster)
)

type Coordinate struct {
	X, Y int
}

type Snake struct {
	Body      []Coordinate
	Direction Direction
}

type Direction int

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

var (
	snake        Snake
	food         Coordinate
	score        int
	speed        int
	gameTick     int
	gameState    GameState
	restartDelay int
)

type GameState int

const (
	Start GameState = iota
	Playing
	GameOver
)

type SlitherQuest struct{}

func (g *SlitherQuest) Update() error {
	switch gameState {
		case Start:
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				gameState = Playing
				InitGame()
			}
		case Playing:
			HandleInput()

			gameTick++
			if gameTick >= speed {
				MoveSnake()
				gameTick = 0
			}

			// Check for collision with the snake itself
			if CheckCollisionWithSelf() || CheckBoundaryCollision() {
				gameState = GameOver
				restartDelay = 60 // A delay before allowing a restart
			}

		case GameOver:
			restartDelay--
			if restartDelay <= 0 && ebiten.IsKeyPressed(ebiten.KeySpace) {
				gameState = Start
			}
	}
	return nil
}

func (g *SlitherQuest) Draw(screen *ebiten.Image) {
	switch gameState {
		case Start:
			ebitenutil.DebugPrintAt(screen, "Press SPACE to Start", screenWidth/2-70, screenHeight/2-20)
		case Playing:
			// Draw the Snake
			for _, segment := range snake.Body {
				DrawSquare(screen, segment.X*cellSize, segment.Y*cellSize, color.RGBA{0, 255, 0, 255})
			}
			// Draw the Food
			DrawSquare(screen, food.X*cellSize, food.Y*cellSize, color.RGBA{255, 0, 0, 255})
			// Draw the Score
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", score), 10, 10)
		case GameOver:
			ebitenutil.DebugPrintAt(screen, "Game Over", screenWidth/2-40, screenHeight/2-20)
			ebitenutil.DebugPrintAt(screen, "Press SPACE to Restart", screenWidth/2-80, screenHeight/2+10)
	}
}

func (g *SlitherQuest) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func InitGame() {
	// Initialize the snake and food positions
	snake = Snake{
		Body:      []Coordinate{{1, 1}, {1, 2}, {1, 3}},
		Direction: RIGHT,
	}
	food = Coordinate{rand.Intn(gridCols), rand.Intn(gridRows)}
	score = 0
	speed = startSpeed
	gameTick = 0
}

func HandleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && snake.Direction != DOWN {
		snake.Direction = UP
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && snake.Direction != LEFT {
		snake.Direction = RIGHT
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && snake.Direction != UP {
		snake.Direction = DOWN
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && snake.Direction != RIGHT {
		snake.Direction = LEFT
	}
}

func MoveSnake() {
	head := snake.Body[len(snake.Body)-1]
	var newHead Coordinate

	switch snake.Direction {
		case UP:
			newHead = Coordinate{head.X, head.Y - 1}
		case RIGHT:
			newHead = Coordinate{head.X + 1, head.Y}
		case DOWN:
			newHead = Coordinate{head.X, head.Y + 1}
		case LEFT:
			newHead = Coordinate{head.X - 1, head.Y}
	}

	// Check if the snake eats the food
	if newHead == food {
		score++
		food = Coordinate{rand.Intn(gridCols), rand.Intn(gridRows)}
		speed = max(1, startSpeed-score/2) // Increase speed as score increases
	} else {
		snake.Body = snake.Body[1:] // Remove tail if no food is eaten
	}

	snake.Body = append(snake.Body, newHead) // Add new head to the snake body
}

func CheckCollisionWithSelf() bool {
	head := snake.Body[len(snake.Body)-1]
	for _, segment := range snake.Body[:len(snake.Body)-1] {
		if segment == head {
			return true
		}
	}
	return false
}

func CheckBoundaryCollision() bool {
	head := snake.Body[len(snake.Body)-1]
	if head.X < 0 || head.X >= gridCols || head.Y < 0 || head.Y >= gridRows {
		return true
	}
	return false
}

func DrawSquare(screen *ebiten.Image, x, y int, c color.Color) {
	ebitenutil.DrawRect(screen, float64(x), float64(y), cellSize, cellSize, c)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed for random food generation

	InitGame()

	// Create the game instance
	game := &SlitherQuest{}

	// Set up the window and run the game loop
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("SlitherQuest")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
