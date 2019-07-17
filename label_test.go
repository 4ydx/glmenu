package glmenu

import (
	"testing"

	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
)

// Test that inside point returns the expected screen coordinate point at the center of the object
func TestInsidePoint(t *testing.T) {
	// around center of screen
	f := &v41.Font{
		WindowHeight: 20,
		WindowWidth:  30,
	}
	l := Label{
		Menu: &Menu{
			Font: f,
		},
		Text: &v41.Text{
			Position:   mgl32.Vec2{0, 0},
			LowerLeft:  gltext.Point{X: -3, Y: -2},
			UpperRight: gltext.Point{X: +3, Y: +2},
		},
	}
	p := l.InsidePoint()
	if p.X != 15 || p.Y != 10 {
		t.Fatalf("bad point %+v", p)
	}

	// moved up the y axis by 6 units
	l.Text = &v41.Text{
		Position:   mgl32.Vec2{0, 6},
		LowerLeft:  gltext.Point{X: -3, Y: -2},
		UpperRight: gltext.Point{X: +3, Y: +2},
	}
	p = l.InsidePoint()
	if p.X != 15 || p.Y != 4 {
		t.Fatalf("bad point %+v", p)
	}

	// moved down the y axis by 6 units
	l.Text = &v41.Text{
		Position:   mgl32.Vec2{0, -6},
		LowerLeft:  gltext.Point{X: -3, Y: -2},
		UpperRight: gltext.Point{X: +3, Y: +2},
	}
	p = l.InsidePoint()
	if p.X != 15 || p.Y != 16 {
		t.Fatalf("bad point %+v", p)
	}

	// moved up the x axis by 6 units
	l.Text = &v41.Text{
		Position:   mgl32.Vec2{6, 0},
		LowerLeft:  gltext.Point{X: -3, Y: -2},
		UpperRight: gltext.Point{X: +3, Y: +2},
	}
	p = l.InsidePoint()
	if p.X != 21 || p.Y != 10 {
		t.Fatalf("bad point %+v", p)
	}

	// moved down the x axis by 6 units
	l.Text = &v41.Text{
		Position:   mgl32.Vec2{-6, 0},
		LowerLeft:  gltext.Point{X: -3, Y: -2},
		UpperRight: gltext.Point{X: +3, Y: +2},
	}
	p = l.InsidePoint()
	if p.X != 9 || p.Y != 10 {
		t.Fatalf("bad point %+v", p)
	}
}
