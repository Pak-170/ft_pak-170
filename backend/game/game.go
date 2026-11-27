package game

import (
	"fmt"
	"time"

	box2d "github.com/bytearena/box2d"
)

type GameState struct {
	frame_num int64

	field       Vec2
	hand_radius float32

	world box2d.B2World

	ball *box2d.B2Body
	p1   *box2d.B2Body
	p2   *box2d.B2Body

	ground *box2d.B2Body
	drag1  *box2d.B2MouseJoint
	drag2  *box2d.B2MouseJoint

	next_mouse  Vec2
	next_mouse2 Vec2
}

// The messages should be:
// 1. State
//    Update for animations:
//  	  - Win / Loose (score change)
//  	  - Ball hit (for a little animations)

type StateMsg struct {
	Width        float32 `json:"width"`
	Height       float32 `json:"height"`
	PosP1        Vec2    `json:"pos_p1"`
	PosP2        Vec2    `json:"pos_p2"`
	PosBall      Vec2    `json:"pos_ball"`
	BallRadius   float32 `json:"ball_radius"`
	HandleRadius float32 `json:"HandleRadius"`
}

type RunningGame struct {
	RecvState chan StateMsg
	MousePos  chan MousePosMsg
}

type MousePosMsg struct {
	PlayerNum int  `json:"player_num"`
	MousePos  Vec2 `json:"mouse_pos"`
}

func Init() GameState {
	height := float32(10.0) // 10m
	width := height * 0.56

	return GameState{
		frame_num: 0,
		field:     Vec2{width, height},
		world:     box2d.MakeB2World(box2d.MakeB2Vec2(0.0, 0.0)),
	}
}

func (s *GameState) Curr_state() StateMsg {
	return StateMsg{
		Width:        s.field.X,
		Height:       s.field.Y,
		PosP1:        Vec2{float32(s.p1.GetPosition().X), float32(s.p1.GetPosition().Y)},
		PosP2:        Vec2{float32(s.p2.GetPosition().X), float32(s.p2.GetPosition().Y)},
		PosBall:      Vec2{float32(s.ball.GetPosition().X), float32(s.ball.GetPosition().Y)},
		HandleRadius: s.hand_radius,
		BallRadius:   float32(s.ball.M_fixtureList.M_shape.GetRadius()),
	}
}
func (s *GameState) Field() Vec2 {
	return s.field
}

func (s *GameState) InitScene() {
	println("Initializing scene")
	{
		bd := box2d.MakeB2BodyDef()
		ground := s.world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(0.0, float64(s.field.Y)), box2d.MakeB2Vec2(float64(s.field.X), float64(s.field.Y)))
		ground.CreateFixture(&shape, 0.0)

		shape.Set(box2d.MakeB2Vec2(0.0, 0.0), box2d.MakeB2Vec2(float64(s.field.X), 0.0))
		ground.CreateFixture(&shape, 0.0)

		shape.Set(box2d.MakeB2Vec2(0.0, 0.0), box2d.MakeB2Vec2(0.0, float64(s.field.Y)))
		ground.CreateFixture(&shape, 0.0)

		shape.Set(box2d.MakeB2Vec2(float64(s.field.X), 0.0), box2d.MakeB2Vec2(float64(s.field.X), float64(s.field.Y)))
		ground.CreateFixture(&shape, 0.0)
	}

	{
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(float64(s.field.X)/2.0, float64(s.field.Y)/2)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.AllowSleep = false

		body := s.world.CreateBody(&bd)

		shape := box2d.MakeB2CircleShape()
		shape.M_radius = 0.10 // 10 cm

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 10.0
		fd.Friction = 0.0
		fd.Restitution = 1.0
		body.CreateFixtureFromDef(&fd)
		body.M_linearVelocity = box2d.MakeB2Vec2(float64(10), float64(20))
		s.ball = body
	}
	{
		s.hand_radius = 0.40
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(float64(s.hand_radius*1.1), float64(s.field.Y-s.hand_radius*1.1))
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.AllowSleep = false

		body := s.world.CreateBody(&bd)

		shape := box2d.MakeB2CircleShape()
		shape.M_radius = float64(s.hand_radius)

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 0.1
		fd.Friction = 0.1
		fd.Restitution = 0.9
		body.CreateFixtureFromDef(&fd)
		s.p1 = body
	}
	{
		s.hand_radius = 0.40
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(float64(s.hand_radius*1.1), float64(s.hand_radius)*1.1)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.AllowSleep = false

		body := s.world.CreateBody(&bd)

		shape := box2d.MakeB2CircleShape()
		shape.M_radius = float64(s.hand_radius)

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 0.1
		fd.Friction = 0.1
		fd.Restitution = 0.9
		body.CreateFixtureFromDef(&fd)
		s.p2 = body
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := s.world.CreateBody(&bd)
		s.ground = ground
	}
	{
		join_def := box2d.MakeB2MouseJointDef()
		join_def.BodyA = s.ground
		join_def.BodyB = s.p1
		join_def.Target = s.p1.GetPosition()
		join_def.MaxForce = 10000.0 * join_def.BodyB.GetMass()
		join_def.FrequencyHz = 100.0
		join_def.DampingRatio = 0.2

		created := s.world.CreateJoint(&join_def)
		s.drag1 = created.(*box2d.B2MouseJoint)
	}
	{
		join_def := box2d.MakeB2MouseJointDef()
		join_def.BodyA = s.ground
		join_def.BodyB = s.p2
		join_def.Target = s.p2.GetPosition()
		join_def.MaxForce = 10000.0 * join_def.BodyB.GetMass()
		join_def.FrequencyHz = 100.0
		join_def.DampingRatio = 0.2

		created := s.world.CreateJoint(&join_def)
		s.drag2 = created.(*box2d.B2MouseJoint)
	}
}

func (s *GameState) Get_bounds_p(p int) Rect {
	if p == 0 {
		return Rect{
			Vec2{0, s.field.Y * 0.7},
			s.field,
		}.Erode(s.hand_radius)
	} else {
		return Rect{
			Vec2{0, 0},
			Vec2{s.field.X, s.field.Y * 0.3},
		}.Erode(s.hand_radius)
	}
}

func (s *GameState) update() {
	s.drag1.SetTarget(box2d.MakeB2Vec2(
		float64(s.next_mouse.X),
		float64(s.next_mouse.Y),
	))
	s.drag2.SetTarget(box2d.MakeB2Vec2(
		float64(s.next_mouse2.X),
		float64(s.next_mouse2.Y),
	))
	//mouse_pos := rl.GetMousePosition().Scale(s.scale)
	s.world.Step(1.0/60.0, 50, 50)
	s.frame_num++
}

func run_game_inner(state *GameState, c chan StateMsg, mouse_pos chan MousePosMsg) {
	// Ground body
	tick := make(chan int, 1000)
	go func() {
		for {
			time.Sleep(time.Duration((1000 / 60) * time.Millisecond))
			tick <- 1
		}
	}()
	fmt.Println("Started bg worker")

	state.update()
	for {
		select {
		case <-tick:
			state.next_mouse2 = state.Get_bounds_p(1).Clamp(Vec2{float32(state.ball.GetPosition().X), 0.0})
			state.update()
			select {
			case c <- state.Curr_state():
			default:
			}
		case mouse_pos_msg := <-mouse_pos:
			fmt.Println("Mouse pos: ", mouse_pos_msg)
			if mouse_pos_msg.PlayerNum == 0 {
				state.next_mouse = mouse_pos_msg.MousePos
			} else {
				state.next_mouse2 = mouse_pos_msg.MousePos
			}
		}
	}
}
func RunGame(state *GameState) RunningGame {
	running_game := RunningGame {
		RecvState: make(chan StateMsg, 2),
		MousePos: make(chan MousePosMsg, 2),
	}
	go func() {
		run_game_inner(state, running_game.RecvState, running_game.MousePos)
	}()
	return running_game
}
