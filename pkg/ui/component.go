package ui

import ()

var ()

type Component interface {
	Init() error
	Render()
	Update(dt float64)
	OnInput()
	OnAdd()
	OnRemove()
	OnResize(w int32, h int32)
}
