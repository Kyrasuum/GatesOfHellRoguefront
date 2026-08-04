package app

import (
	"roguefront/pkg/stage"
)

var (
	CurApp App
)

type App interface {
	GetWidth() int
	GetHeight() int
	Start(debug bool) error
	Build() // render call
	Exit()
	SetStage(nextStage stage.Stage)
	GetStage() stage.Stage
}
