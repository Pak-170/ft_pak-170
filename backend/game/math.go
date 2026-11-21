package game

type Vec2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Rect struct {
	Positon Vec2 `json:"position"`
	Size    Vec2 `json:"size"`
}

func (r Rect) Clamp(point Vec2) Vec2 {
	return Vec2{
		min(r.Positon.X+r.Size.X, max(r.Positon.X, point.X)),
		min(r.Positon.Y+r.Size.Y, max(r.Positon.Y, point.Y)),
	}
}

// / Reduce the rectangle size by `amount` on each side. The center stays in the same place
func (r Rect) Erode(amount float32) Rect {
	r.Positon.X += amount
	r.Positon.Y += amount
	r.Size.X -= amount
	r.Size.Y -= amount
	return r
}

func (v Vec2) Scale(scale float32) Vec2 {
	v.X *= scale
	v.Y *= scale
	return v
}
