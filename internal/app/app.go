package app

import (
	// "image/color"
	"log"
	"os"
	"time"

	"roguefront/internal/stage"
	"roguefront/pkg/config"

	App "roguefront/pkg/app"
	pub_stage "roguefront/pkg/stage"

	"github.com/AllenDang/giu"
)

var ()

type app struct {
	curStage pub_stage.Stage
	console  interface{}

	logicInterval int64
	drawInterval  int64

	width  int
	height int

	shutdown bool
	window   *giu.MasterWindow
}

// initialize app
func (a *app) init() error {
	a.curStage = nil
	a.console = nil

	a.width = 800
	a.height = 512

	a.logicInterval = 16
	a.drawInterval = 16

	a.window = nil

	App.CurApp = a

	return nil
}

// handle input
func (a *app) handleInput(dt int64) {
	if a.curStage != nil {
		a.curStage.OnInput(dt)
	}
}

// render cycle
func (a *app) Build() {
	if a.curStage != nil {
		a.curStage.Build()
	}
}

// handle resizing
func (a *app) onResize() {
	w, h := a.window.GetSize()

	//check for resize event
	if a.width != w || a.height != h {
		if a.curStage != nil {
			a.curStage.OnResize(w, h)
		}
		a.width = w
		a.height = h
	}
}

// update cycle
func (a *app) update(dt int64) {
	if a.curStage != nil {
		a.curStage.Update(dt)
	}
}

// get window width
func (a *app) GetWidth() int {
	return a.width
}

// get window height
func (a *app) GetHeight() int {
	return a.height
}

// main run loop for the app while running
func (a *app) run(debug bool) error {
	a.window = giu.NewMasterWindow(config.AppName, a.width, a.height, 0)
	a.window.SetTargetFPS(uint(time.Second / (time.Duration(a.drawInterval) * time.Millisecond)))

	if debug {
		file, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	menu := stage.MainMenu{}
	err := menu.Init()
	if err != nil {
		return err
	}
	a.SetStage(&menu)

	//render loop
	dt := time.Now().Unix()
	a.window.Run(func() {
		dt = time.Now().Unix() - dt
		//render

		giu.SingleWindow().Layout(
			a,
		)
		//logic
		a.onResize()
		a.update(dt)
		if giu.IsWindowFocused(giu.FocusedFlags(giu.FocusedFlagsAnyWindow)) {
			a.handleInput(dt)
		}
	})
	return nil
}

// set  the currently active stage
func (a *app) SetStage(nextStage pub_stage.Stage) {
	if a.curStage != nil {
		a.curStage.OnRemove()
	}
	if nextStage != nil {
		a.curStage = nextStage
		a.curStage.OnAdd()
	}
}

// get the currently active stage in the app
func (a *app) GetStage() pub_stage.Stage {
	return a.curStage
}

// Exit the application
func (a *app) Exit() {
	if a.curStage != nil {
		a.curStage.OnRemove()
	}
	a.window.Close()
}

// start the application
func (a *app) Start(debug bool) error {
	return a.run(debug)
}

// create a new app
func NewApp() *app {
	a := &app{}
	err := a.init()
	if err != nil {
		log.Printf("ERR: %+v\n", err)
		return nil
	}
	return a
}
