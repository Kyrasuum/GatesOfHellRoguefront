package stage

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"time"

	"roguefront/pkg/app"

	"roguefront/res"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Game struct {
	Profile  string `json:"profile"`
	Game     string `json:"game"`
	Workshop string `json:"workshop"`

	Maps    map[string][]string `json:"maps"`
	Mods    []string            `json:"mods"`
	Nations []string            `json:"nations"`

	State res.Save `json:"state"`

	paused bool `json:"-"`
}

func (g *Game) Default() {
}

// initialize game object
func (g *Game) Init() error {
	if g == nil {
		return fmt.Errorf("Invalid stage")
	}

	g.Profile = ""
	g.Game = ""
	g.Workshop = ""

	g.paused = false

	g.Maps = map[string][]string{}
	g.Mods = []string{}
	g.Nations = []string{}

	g.State = res.Save{
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

func (g *Game) Populate() error {
	g.PopulateMaps()
	g.PopulateUnits()
	g.PopulateSquads()
	g.PopulateItems()
	g.PopulateVehicles()

	ret := 0
	if len(g.Nations) < 1 {
		log.Printf("Failed to populate nations\n")
		ret = 1
	}
	if len(g.Maps) < 1 {
		log.Printf("Failed to populate maps\n")
		ret = 1
	}

	if ret != 0 {
		return fmt.Errorf("Failed to populate")
	}
	return nil
}

func (g *Game) PopulateMaps() {
	ReadMaps := func(data []byte) {
		scanner := bufio.NewScanner(bytes.NewReader(data))

		region := ""
		g.Maps = map[string][]string{}
		for scanner.Scan() {
			line := scanner.Text()

			if line == "" {
				continue
			}

			switch {
			//find resources section
			case strings.HasPrefix(line, "\t\t\t{"):
				if region != "" {
					g.Maps[region] = append(g.Maps[region], strings.Split(strings.TrimSpace(line)[2:], "\"")[0])
				}
			case strings.HasPrefix(line, "\t{"):
				region = strings.TrimSpace(line)[1:]
				g.Maps[region] = []string{}
			case strings.HasPrefix(line, "}"):
				return
			}
		}
	}

	UpdateMaps := func(path string) {
		r, err := zip.OpenReader(path)
		if err != nil {
			return
		}
		defer r.Close()

		for _, f := range r.File {
			switch f.Name {
			case "set/dynamic_campaign/map_points.set":
				rc, err := f.Open()
				if err != nil {
					return
				}

				data, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					return
				}

				ReadMaps(data)
			}
		}
	}

	UpdateMaps(g.Game + "/resource/gamelogic.pak")
	for _, mod := range g.Mods {
		UpdateMaps(g.Workshop + "/" + mod + "/resource/gamelogic.pak")
		data, err := os.ReadFile(g.Workshop + "/" + mod + "/resource/set/dynamic_campaign/map_points.set")
		if err == nil {
			ReadMaps(data)
		}
	}
}

func (g *Game) PopulateUnits() {}

func (g *Game) PopulateSquads() {}

func (g *Game) PopulateItems() {}

func (g *Game) PopulateVehicles() {}

func (g *Game) ReRoll() {
	g.RollLandscape()
	g.RollRegion()
	g.RollMap()
	g.RollEnemy()
}

func (g *Game) RollLandscape() {
	landscapes := []string{
		"wood",
		"winter",
		"desert",
	}
	g.State.Status.Landscape = landscapes[rand.IntN(len(landscapes))]
}

func (g *Game) RollRegion() {
	regions := slices.Collect(maps.Keys(g.Maps))
	g.State.Status.Region = regions[rand.IntN(len(regions))]
}

func (g *Game) RollMap() {
	maps := g.Maps[g.State.Status.Region]
	g.State.Status.Map = maps[rand.IntN(len(maps))]
}

func (g *Game) RollEnemy() {
	g.State.Status.EnemyArmy = g.Nations[rand.IntN(len(g.Nations))]
}

func (g *Game) SaveGame() {
	g.State.Status.Seed = int64(rand.Int32())
	g.State.Status.Timestamp = time.Now().Unix()
	g.State.Write(g.Profile)
}

func (g *Game) LoadGame(path string) error {
	// Load the settings file
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	// Unmarshall the campaign
	return json.Unmarshal(data, g)
}
