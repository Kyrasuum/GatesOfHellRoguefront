package stage

import ()

var ()

type Stage interface {
	Init() error
	OnResize(w int, h int)
	Build() // render call
	Update(dt int64)
	OnInput(dt int64)
	OnAdd()
	OnRemove()
}
