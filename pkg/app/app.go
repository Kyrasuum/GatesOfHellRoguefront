package app

import (
	"roguefront/pkg/stage"
)

var (
	CurApp App
)

type App interface {
	GetWidth() int32
	GetHeight() int32
	Start(debug bool) error
	Exit()
	SetStage(nextStage stage.Stage)
	GetStage() stage.Stage
}
