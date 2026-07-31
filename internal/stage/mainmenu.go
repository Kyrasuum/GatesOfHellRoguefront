package stage

import (
	"fmt"
	"log"

	"roguefront/pkg/app"

	"github.com/sqweek/dialog"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type MainMenu struct {
	panel int

	profile string
	game    string
	mods    string
}

func (s *MainMenu) Init() error {
	if s == nil {
		return fmt.Errorf("Invalid stage")
	}

	s.panel = 0

	return nil
}

func (s *MainMenu) OnResize(w int32, h int32) {
	if s == nil {
		return
	}
}

func (s *MainMenu) Render() {
	if s == nil {
		return
	}

	rl.ClearBackground(rl.Black)

	if s.panel == 1 {
		gui.SetState(gui.STATE_DISABLED)
	} else {
		gui.SetState(gui.STATE_NORMAL)
	}
	nwbtn := gui.Button(rl.NewRectangle(10, 10, 100, 40), "New Game")
	if nwbtn {
		s.panel = 1
	}
	if s.panel == 2 {
		gui.SetState(gui.STATE_DISABLED)
	} else {
		gui.SetState(gui.STATE_NORMAL)
	}
	ldbtn := gui.Button(rl.NewRectangle(10, 50, 100, 40), "Load Game")
	if ldbtn {
		s.panel = 2
	}
	gui.SetState(gui.STATE_NORMAL)
	exbtn := gui.Button(rl.NewRectangle(10, 90, 100, 40), "Exit")
	if exbtn {
		app.CurApp.Exit()
	}

	if s.panel == 1 {
		gui.Panel(rl.NewRectangle(110, 10, 400, 120), "New Game")
		cnbtn := gui.Button(rl.NewRectangle(488, 12, 20, 20), "X")
		if cnbtn {
			s.panel = 0
		}
		gui.Label(rl.NewRectangle(115, 35, 100, 20), "User Profile Path:")
		gui.TextBox(rl.NewRectangle(215, 35, 290, 20), &s.profile, 16, true)
		prfile := gui.Button(rl.NewRectangle(485, 35, 20, 20), "...")
		if prfile {
			file, err := dialog.Directory().
				Title("Select User Profile Path for Call to Arms Gates of Hell Ostfront").Browse()
			if err == nil {
				s.profile = file
			}
		}
		gui.Label(rl.NewRectangle(115, 60, 100, 20), "Game Directory:")
		gui.TextBox(rl.NewRectangle(215, 60, 290, 20), &s.game, 16, true)
		gmfile := gui.Button(rl.NewRectangle(485, 60, 20, 20), "...")
		if gmfile {
			file, err := dialog.Directory().
				Title("Select Game Directory for Call to Arms Gates of Hell Ostfront").Browse()
			if err == nil {
				s.game = file
			}
		}
		gui.Label(rl.NewRectangle(115, 85, 100, 20), "Mods Directory:")
		gui.TextBox(rl.NewRectangle(215, 85, 290, 20), &s.mods, 16, true)
		mdfile := gui.Button(rl.NewRectangle(485, 85, 20, 20), "...")
		if mdfile {
			file, err := dialog.Directory().
				Title("Select Mods Directory for Call to Arms Gates of Hell Ostfront").Browse()
			if err == nil {
				s.mods = file
			}
		}
		stbtn := gui.Button(rl.NewRectangle(252, 108, 100, 20), "Start Game")
		if stbtn {
			game := Game{}
			err := game.Init()
			if err != nil {
				log.Printf("%+v", err)
				return
			}
			app.CurApp.SetStage(&game)
		}
	}
	if s.panel == 2 && !ldbtn {
		file, err := dialog.File().
			Filter("Json Files", "json").
			Title("Load Previous Roguefront Campaign").Load()
		if err == nil {
			file = file
			// TO DO: write loading campaigns
		}
		s.panel = 0
	}
}

func (s *MainMenu) Update(dt float32) {
	if s == nil {
		return
	}
}

func (s *MainMenu) OnInput(dt float32) {
	if s == nil {
		return
	}
}

func (s *MainMenu) OnAdd() {
	if s == nil {
		return
	}
}

func (s *MainMenu) OnRemove() {
	if s == nil {
		return
	}
}
