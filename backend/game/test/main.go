package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	game "go_pakito/game"
)

func draw(s *game.StateMsg, scale float32) {
	rl.ClearBackground(rl.DarkGray)

	{
		ball_pos := s.PosBall
		radius := s.BallRadius * scale
		rl.DrawCircleSector(rl.Vector2{float32(ball_pos.X) * scale, float32(ball_pos.Y) * scale}, float32(radius), 0.0, 360.0, 40, rl.Green)
		fmt.Println("Ball pos:", ball_pos)
	}

	{
		ball_pos := s.PosP1
		radius := s.HandleRadius * scale
		rl.DrawCircleSector(rl.Vector2{float32(ball_pos.X) * scale, float32(ball_pos.Y) * scale}, float32(radius), 0.0, 360.0, 40, rl.White)
		fmt.Println("Hand pos:", ball_pos)
	}
	{
		ball_pos := s.PosP2
		radius := s.HandleRadius * scale
		rl.DrawCircleSector(rl.Vector2{float32(ball_pos.X) * scale, float32(ball_pos.Y) * scale}, float32(radius), 0.0, 360.0, 40, rl.White)
		fmt.Println("Hand pos:", ball_pos)
	}
}

func main() {
	output_height := 850
	state := game.Init()
	state.InitScene()

	scale := float32(output_height) / state.Field().Y
	output_res := state.Field().Scale(scale)

	rl.InitWindow(int32(output_res.X), int32(output_res.Y), "Air Hokey")
	defer rl.CloseWindow()

	running_game := game.RunGame(&state)

	// Ground body
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		select {
		case running_game.MousePos <- game.MousePosMsg{
			PlayerNum: 0,
			MousePos: state.Get_bounds_p(0).Clamp(
				game.Vec2(rl.GetMousePosition()).Scale(1 / scale)),
		}:
		default:
		}
		state_msg := <-running_game.RecvState
		draw(&state_msg, scale)
		rl.EndDrawing()
	}
}
