package stage

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"roguefront/pkg/app"

	"roguefront/res"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Game struct {
	profile  string
	game     string
	workshop string

	state res.Save

	paused bool
}

// initialize game object
func (g *Game) Init() error {
	if g == nil {
		return fmt.Errorf("Invalid stage")
	}

	g.paused = false

	g.state = res.Save{
		Campaign: &res.Campaign{
			Soldiers:    []*res.Soldier{},
			Inventories: []*res.Inventory{},
			Squads:      []*res.Squad{},
		},
		Status: &res.Status{
			Mods:        []string{},
			Timestamp:   0,
			Seed:        0,
			Name:        "new campaign",
			Army:        "ger",
			EnemyArmy:   "ger",
			Difficulty:  "easy",
			FogOfWar:    "fog_realistic",
			Resources:   1,
			Region:      "ostfront",
			PlayedGames: 0,
			WonGames:    0,
			Landscape:   "wood",
			Map:         "multi/dcg_narofominsk:campaign_capture_the_flag:4x4",
			Texmod:      "camo",
			Risk:        "standard",
			Gamemode:    "campaign_capture_the_flag",
		},
	}

	return nil
}

// handle resize event
func (g *Game) OnResize(w int32, h int32) {
	if g == nil {
		return
	}
}

// render hook
func (g *Game) Render() {
	if g == nil {
		return
	}

	rl.ClearBackground(rl.Black)

	width := float32(app.CurApp.GetWidth())
	height := float32(app.CurApp.GetHeight())

	gui.SetState(gui.STATE_NORMAL)
	gui.Panel(rl.NewRectangle(0, 0, width, 20), "")
	menu := gui.Button(rl.NewRectangle(width-50, 0, 50, 20), "Menu")
	if menu {
		g.paused = true
	}

	if g.paused {
		gui.Panel(rl.NewRectangle((width-300)/2, (height-200)/2, 300, 200), "Menu")
		cnbtn := gui.Button(rl.NewRectangle((width-300)/2+278, (height-200)/2+2, 20, 20), "X")
		if cnbtn {
			g.paused = false
		}
		svbtn := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+35, 250, 30), "Save")
		if svbtn {
			g.SaveGame()
		}
		exbtn1 := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+75, 250, 30), "Exit to Main Menu")
		if exbtn1 {
			menu := MainMenu{}
			err := menu.Init()
			if err != nil {
				log.Printf("%+v", err)
				return
			}
			app.CurApp.SetStage(&menu)
		}
		exbtn2 := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+115, 250, 30), "Exit to Desktop")
		if exbtn2 {
			app.CurApp.Exit()
		}
		clbtn := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+155, 250, 30), "Cancel")
		if clbtn {
			g.paused = false
		}
	}

	return
}

// handle update cycle
func (g *Game) Update(dt float32) {
	if g == nil {
		return
	}
}

// handle player input
func (g *Game) OnInput(dt float32) {
	if g == nil {
		return
	}
}

// handle add event
func (g *Game) OnAdd() {
	if g == nil {
		return
	}
}

// handle remove event
func (g *Game) OnRemove() {
	if g == nil {
		return
	}
}

func (g *Game) PopulateMaps() {}

func (g *Game) PopulateEnemies() {}

func (g *Game) PopulateUnits() {}

func (g *Game) RollLandscape() {
	// roll landscape
	landscapes := []string{
		"wood",
		"winter",
		"desert",
	}
	g.state.Status.Landscape = landscapes[rand.IntN(len(landscapes))]
}

func (g *Game) RollRegion() {
}

func (g *Game) RollMap() {
}

func (g *Game) RollEnemy() {
}

func (g *Game) SaveGame() {
	g.state.Status.Seed = int64(rand.Int32())
	g.state.Status.Timestamp = time.Now().Unix()
	g.state.Write(g.profile)
}

func (g *Game) LoadGame() {}
