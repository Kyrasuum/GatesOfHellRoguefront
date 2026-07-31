package ui

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Text struct {
	pos rl.Vector2

	text string
	font rl.Font

	fontSize float32
	spacing  float32

	textCol color.RGBA
}

func (t *Text) defaults() {
	if t == nil {
		return
	}

	t.pos = rl.Vector2{0, 0}

	t.text = ""
	t.font = rlx.GetFontDefault()

	t.fontSize = 16
	t.spacing = 1

	t.textCol = color.RGBA{255, 255, 255, 255}
}

func (t *Text) Init() error {
	if t == nil {
		return fmt.Errorf("Invalid text")
	}

	t.defaults()
	return nil
}

func (t *Text) SetPosition(p rl.Vector2) {
	if t == nil {
		return
	}

	t.pos = p
}

func (t *Text) SetText(s string) {
	if t == nil {
		return
	}

	t.text = s
}

func (t *Text) SetFont(f rl.Font) {
	if t == nil {
		return
	}

	t.font = f
}

func (t *Text) SetFontSize(f float32) {
	if t == nil {
		return
	}

	t.fontSize = f
}

func (t *Text) SetTextColor(c color.RGBA) {
	if t == nil {
		return
	}

	t.textCol = c
}

func (t *Text) Render() {
	if t == nil {
		return
	}

	rlx.DrawTextEx(t.font, t.text, t.pos, t.fontSize, t.spacing, t.textCol)
}

func (t *Text) Update(dt float64) {
	if t == nil {
		return
	}
}

func (t *Text) OnInput() {
	if t == nil {
		return
	}
}

func (t *Text) OnAdd() {
	if t == nil {
		return
	}
}

func (t *Text) OnRemove() {
	if t == nil {
		return
	}
}

func (t *Text) OnResize(w int32, h int32) {
	if t == nil {
		return
	}
}
