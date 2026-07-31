package stage

import ()

var ()

type Stage interface {
	Init() error
	OnResize(w int32, h int32)
	Render()
	Update(dt float32)
	OnInput(dt float32)
	OnAdd()
	OnRemove()
}
