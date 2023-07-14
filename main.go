package main

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	screenWidth      = 640
	screenHeight     = 480
	playAreaTopLeftY = 50
	playAreaWidth    = screenWidth
	playAreaHeight   = screenHeight - playAreaTopLeftY - 1 // Adjusted for bottom border
	fontPath         = "valorax.ttf"
	fontSize         = 20
	scoreTextX       = 10
	scoreTextY       = 10
	gameOverTextX    = screenWidth / 2
	gameOverTextY    = screenHeight / 2
)

// Point represents a 2D point
type Point struct {
	X int32
	Y int32
}

// Game represents the game state
type Game struct {
	snake       []sdl.Rect
	food        sdl.Rect
	direction   Point
	gameOver    bool
	score       int
	updateCount int
}

// newGame creates a new game
func newGame() Game {
	// Initialize the snake
	snake := []sdl.Rect{
		{X: 10, Y: playAreaTopLeftY + 10, W: 10, H: 10},
		{X: 20, Y: playAreaTopLeftY + 10, W: 10, H: 10},
		{X: 30, Y: playAreaTopLeftY + 10, W: 10, H: 10},
	}

	// Initialize the food
	food := generateFoodPosition(snake)

	// Initialize the game state
	game := Game{
		snake:       snake,
		food:        food,
		direction:   Point{X: 10, Y: 0},
		gameOver:    false,
		score:       0,
		updateCount: 0,
	}

	return game
}

func generateFoodPosition(snakePoints []sdl.Rect) sdl.Rect {
	for {
		// Generate a random position for the food within the play area
		foodPos := sdl.Rect{
			X: rand.Int31n(screenWidth-10) + 5,
			Y: rand.Int31n(playAreaHeight-10) + playAreaTopLeftY + 5,
			W: 10,
			H: 10,
		}

		// Check if the food position is not on the snake
		if !isPointInSlice(foodPos, convertRectsToPoints(snakePoints)) {
			return foodPos
		}
	}
}

// isPointInSlice checks if a point is in a slice of points
func isPointInSlice(point sdl.Rect, points []Point) bool {
	for _, p := range points {
		if p.X == point.X && p.Y == point.Y {
			return true
		}
	}
	return false
}

// convertRectsToPoints converts a slice of sdl.Rect values to a slice of Point values
func convertRectsToPoints(rects []sdl.Rect) []Point {
	points := make([]Point, len(rects))
	for i, rect := range rects {
		points[i] = Point{X: rect.X, Y: rect.Y}
	}
	return points
}

// clearScreen clears the screen
func clearScreen(renderer *sdl.Renderer, gameOver bool) {
	// Clear the screen with black color
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	if !gameOver {
		// Draw white borders on all four sides
		renderer.SetDrawColor(255, 255, 255, 255)
		topBorderRect := sdl.Rect{X: 0, Y: playAreaTopLeftY - 1, W: playAreaWidth, H: 1}
		bottomBorderRect := sdl.Rect{X: 0, Y: playAreaTopLeftY + playAreaHeight, W: playAreaWidth, H: 1}
		leftBorderRect := sdl.Rect{X: 0, Y: playAreaTopLeftY, W: 1, H: playAreaHeight}
		rightBorderRect := sdl.Rect{X: playAreaWidth - 1, Y: playAreaTopLeftY, W: 1, H: playAreaHeight}
		renderer.FillRect(&topBorderRect)
		renderer.FillRect(&bottomBorderRect)
		renderer.FillRect(&leftBorderRect)
		renderer.FillRect(&rightBorderRect)
	} else {
		// Clear the play area space when game is over
		renderer.SetDrawColor(0, 0, 0, 255)
		playAreaRect := sdl.Rect{X: 0, Y: playAreaTopLeftY, W: playAreaWidth, H: playAreaHeight}
		renderer.FillRect(&playAreaRect)
	}
}

// renderScore renders the score on the screen
func renderScore(renderer *sdl.Renderer, font *ttf.Font, score int) {
	scoreText := fmt.Sprintf("Score: %d", score)
	surface, err := font.RenderUTF8Solid(scoreText, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	_, _, textureWidth, textureHeight, err := texture.Query()
	if err != nil {
		panic(err)
	}

	renderer.Copy(texture, nil, &sdl.Rect{X: scoreTextX, Y: scoreTextY, W: textureWidth, H: textureHeight})
}

// renderGameOver renders the game over screen
func renderGameOver(renderer *sdl.Renderer, font *ttf.Font, score int) {
	gameOverText := "Game Over"
	fontColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}

	surface, err := font.RenderUTF8Solid(gameOverText, fontColor)
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	_, _, textureWidth, textureHeight, err := texture.Query()
	if err != nil {
		panic(err)
	}

	gameOverTextXPos := gameOverTextX - textureWidth/2
	gameOverTextYPos := gameOverTextY - textureHeight/2

	renderer.Copy(texture, nil, &sdl.Rect{X: gameOverTextXPos, Y: gameOverTextYPos, W: textureWidth, H: textureHeight})

	scoreText := fmt.Sprintf("Score: %d", score)
	scoreSurface, err := font.RenderUTF8Solid(scoreText, fontColor)
	if err != nil {
		panic(err)
	}
	defer scoreSurface.Free()

	scoreTexture, err := renderer.CreateTextureFromSurface(scoreSurface)
	if err != nil {
		panic(err)
	}
	defer scoreTexture.Destroy()

	_, _, scoreTextureWidth, scoreTextureHeight, err := scoreTexture.Query()
	if err != nil {
		panic(err)
	}

	scoreTextXPos := gameOverTextX - scoreTextureWidth/2
	scoreTextYPos := gameOverTextY + textureHeight/2 + 10

	renderer.Copy(scoreTexture, nil, &sdl.Rect{X: scoreTextXPos, Y: scoreTextYPos, W: scoreTextureWidth, H: scoreTextureHeight})
}

// gameLoop runs the game loop
func gameLoop(renderer *sdl.Renderer, font *ttf.Font, game *Game) {
	// Set the initial food position
	foodPos := generateFoodPosition(game.snake)
	game.food = foodPos

	// Set the initial direction
	game.direction = Point{X: 10, Y: 0}

	// Set the initial update count
	game.updateCount = 0

	// Run the game loop
	for !game.gameOver {
		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				game.gameOver = true
			case *sdl.KeyboardEvent:
				keyEvent := event.(*sdl.KeyboardEvent)
				if keyEvent.Type == sdl.KEYDOWN {
					switch keyEvent.Keysym.Sym {
					case sdl.K_UP:
						if game.direction.Y != 10 {
							game.direction = Point{X: 0, Y: -10}
						}
					case sdl.K_DOWN:
						if game.direction.Y != -10 {
							game.direction = Point{X: 0, Y: 10}
						}
					case sdl.K_LEFT:
						if game.direction.X != 10 {
							game.direction = Point{X: -10, Y: 0}
						}
					case sdl.K_RIGHT:
						if game.direction.X != -10 {
							game.direction = Point{X: 10, Y: 0}
						}
					}
				}
			}
		}

		// Update the game state
		game.updateCount++
		if game.updateCount == 10 {
			game.updateCount = 0

			// Move the snake
			head := game.snake[len(game.snake)-1]
			newHead := sdl.Rect{
				X: head.X + game.direction.X,
				Y: head.Y + game.direction.Y,
				W: 10,
				H: 10,
			}

			// Check if the snake hit the wall
			if newHead.X < 0 || newHead.X >= screenWidth || newHead.Y < playAreaTopLeftY || newHead.Y >= screenHeight {
				game.gameOver = true
			}

			// Check if the snake hit itself
			for _, rect := range game.snake[:len(game.snake)-1] {
				if newHead.X == rect.X && newHead.Y == rect.Y {
					game.gameOver = true
				}
			}

			// Check if the snake ate the food
			if newHead.HasIntersection(&game.food) {
				game.score++

				// Generate a new segment for the snake at the same position as the head
				newSegment := sdl.Rect{
					X: newHead.X,
					Y: newHead.Y,
					W: 10,
					H: 10,
				}
				game.snake = append(game.snake, newSegment)

				// Generate a new food position
				foodPos := generateFoodPosition(game.snake)
				game.food = foodPos
			} else {
				// Move the snake
				game.snake = append(game.snake[1:], newHead)
			}
		}

		// Clear the screen
		clearScreen(renderer, game.gameOver)

		// Draw the snake
		renderer.SetDrawColor(255, 255, 255, 255)
		for _, rect := range game.snake {
			renderer.FillRect(&rect)
		}

		// Draw the food
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.FillRect(&game.food)

		// Render the score
		renderScore(renderer, font, game.score)

		// Update the screen
		renderer.Present()

		// Delay to control the frame rate
		sdl.Delay(16)
	}

	// Render the game over screen
	clearScreen(renderer, game.gameOver)
	renderGameOver(renderer, font, game.score)
	renderer.Present()

	// Delay to show the game over screen for a few seconds
	//sdl.Delay(3000)

	// pause the game
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
	}
}

func main() {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Initialize SDL_ttf
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer ttf.Quit()

	// Create a window and renderer
	window, err := sdl.CreateWindow("Snake", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// Create a font
	font, err := ttf.OpenFont(fontPath, fontSize)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	// Create the game
	game := newGame()

	// Run the game loop
	gameLoop(renderer, font, &game)
}
