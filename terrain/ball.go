package terrain

import (
	"math"

	config "github.com/Guidjy/wireframe/config"
	. "github.com/gen2brain/raylib-go/raylib"
)

const acceleration float32 = 50000.0
const friction float32 = 0.98
const maxVelocity = 25000.0

type Ball struct {
	Pos Vector3

	velocity Vector2

	Radius float32
}

func (ball *Ball) init() {
	ball.Pos = Vector3Zero()
	ball.velocity = Vector2Zero()
	ball.Radius = 50
}

func (ball *Ball) move(dir Vector2) {
	dt := GetFrameTime()

	if dir.X != 0 || dir.Y != 0 {
		ball.velocity = ball.velocity.Add(dir.Scale(acceleration)).Scale(dt)
	}

	ball.velocity = ball.velocity.Scale(friction)

	if ball.velocity.X > maxVelocity {
		ball.velocity.X = maxVelocity
	}
	if ball.velocity.X < -maxVelocity {
		ball.velocity.X = -maxVelocity
	}
	if ball.velocity.Y > maxVelocity {
		ball.velocity.Y = maxVelocity
	}
	if ball.velocity.Y < -maxVelocity {
		ball.velocity.Y = -maxVelocity
	}

	ball.Pos.X += ball.velocity.X * dt
	ball.Pos.Z += ball.velocity.Y * dt

	// keeps ball in bounds
	hw := float32(config.TerrainWidth/2 - 250)
	if ball.Pos.X < -hw {
		ball.Pos.X = -hw
	}
	if ball.Pos.X > hw {
		ball.Pos.X = hw
	}
	if ball.Pos.Z < -hw {
		ball.Pos.Z = -hw
	}
	if ball.Pos.Z > hw {
		ball.Pos.Z = hw
	}
}

func (ball *Ball) handleKeyboardInput() {
	dir := Vector2Zero()

	if IsKeyDown(KeyUp) {
		dir.Y = 1
	}
	if IsKeyDown(KeyDown) {
		dir.Y = -1
	}
	if IsKeyDown(KeyLeft) {
		dir.X = 1
	}
	if IsKeyDown(KeyRight) {
		dir.X = -1
	}

	ball.move(dir)
}

func (ball Ball) render() {
	rings := 10
	slices := 10

	vertices := make([][]Vector3, rings+1)
	for i := range vertices {
		vertices[i] = make([]Vector3, slices+1)
	}

	// Creates vertices
	for i := 0; i <= rings; i++ {
		// sweeps vertically
		phi := math.Pi * float64(i) / float64(rings)

		for j := 0; j <= slices; j++ {
			// sweeps horizontally
			theta := 2.0 * math.Pi * float64(j) / float64(slices)

			x := float32(float64(ball.Radius) * math.Sin(phi) * math.Cos(theta))
			y := float32(float64(ball.Radius) * math.Cos(phi))
			z := float32(float64(phi) * math.Sin(theta))
			z = float32(float64(ball.Radius) * math.Sin(phi) * math.Sin(theta))

			vertices[i][j] = Vector3{X: x, Y: y, Z: z}.Add(ball.Pos)
		}
	}

	// Draws edges
	for i := 0; i < rings; i++ {
		for j := 0; j < slices; j++ {
			drawEdge(vertices[i][j], vertices[i+1][j], Red)
			drawEdge(vertices[i][j], vertices[i][j+1], Red)
		}
	}
}

func (ball *Ball) Update() {
	ball.handleKeyboardInput()

	ball.render()
}
