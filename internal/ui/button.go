package ui

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Button struct {
	pos    rl.Rectangle
	inset  rl.Vector2
	border float32

	hovered bool
	clicked bool

	text    string
	subtext string
	font    rl.Font

	fontSize float32
	spacing  float32

	textCol color.RGBA
	bgCol   color.RGBA
	bordCol color.RGBA

	onPress func()
}

func (b *Button) defaults() {
	if b == nil {
		return
	}

	b.pos = rl.Rectangle{0, 0, 0, 0}
	b.inset = rl.Vector2{0, 0}
	b.border = float32(3)

	b.hovered = false
	b.clicked = false

	b.text = ""
	b.subtext = ""
	b.font = rlx.GetFontDefault()

	b.fontSize = 16
	b.spacing = 1

	b.textCol = color.RGBA{255, 255, 255, 255}
	b.bgCol = color.RGBA{0, 0, 0, 255}
	b.bordCol = color.RGBA{255, 255, 255, 255}

	b.onPress = func() {}
}

func (b *Button) Init() error {
	if b == nil {
		return fmt.Errorf("Invalid button")
	}

	b.defaults()

	return nil
}

func (b *Button) SetOnPress(f func()) {
	if b == nil {
		return
	}

	b.onPress = f
}

func (b *Button) SetPosition(p rl.Rectangle) {
	if b == nil {
		return
	}

	b.pos = p
	b.calcText()
}

func (b *Button) SetInset(i rl.Vector2) {
	if b == nil {
		return
	}

	b.inset = i
	b.calcText()
}

func (b *Button) SetBorder(bo float32) {
	if b == nil {
		return
	}

	b.border = bo
	b.calcText()
}

func (b *Button) SetText(t string) {
	if b == nil {
		return
	}

	b.text = t
}

func (b *Button) SetFont(f rl.Font) {
	if b == nil {
		return
	}

	b.font = f
}

func (b *Button) SetFontSize(f float32) {
	if b == nil {
		return
	}

	b.fontSize = f
}

func (b *Button) SetTextColor(c color.RGBA) {
	if b == nil {
		return
	}

	b.textCol = c
}

func (b *Button) SetBGColor(c color.RGBA) {
	if b == nil {
		return
	}

	b.bgCol = c
}

func (b *Button) SetBorderColor(c color.RGBA) {
	if b == nil {
		return
	}

	b.bordCol = c
}

func (b *Button) calcText() {
	if b == nil {
		return
	}

	b.subtext = b.text[:0]
	for i := 1; 1 < len(b.text); i++ {
		w := rlx.MeasureTextEx(b.font, b.text[:i], b.fontSize, b.spacing)
		if w.X > b.pos.Width-b.inset.X-b.border {
			break
		}
		b.subtext = b.text[:i]
	}
}

func (b *Button) Render() {
	if b == nil {
		return
	}

	rlx.DrawRectangleLinesEx(b.pos, b.border, b.bordCol)
	rlx.DrawRectangleRec(b.pos, b.bgCol)
	rlx.DrawTextEx(b.font, b.subtext, rl.Vector2{(b.pos.X + b.inset.X + b.border), (b.pos.Y + b.inset.Y + b.border)}, b.fontSize, b.spacing, b.textCol)
}

func (b *Button) Update(dt float64) {
	if b == nil {
		return
	}
}

func (b *Button) OnInput() {
	if b == nil {
		return
	}

	mp := rlx.GetMousePosition()
	if rlx.CheckCollisionPointRec(mp, b.pos) {
		if rlx.IsMouseButtonPressed(rlx.MouseLeftButton) {
			if !b.clicked {
				b.onPress()
			}
			b.clicked = true
		} else {
			b.clicked = false
		}
		b.hovered = true
	} else {
		b.clicked = false
		b.hovered = false
	}
}

func (b *Button) OnResize(w int32, h int32) {
	if b == nil {
		return
	}
}

func (b *Button) OnAdd() {
	if b == nil {
		return
	}
}

func (b *Button) OnRemove() {
	if b == nil {
		return
	}
}
