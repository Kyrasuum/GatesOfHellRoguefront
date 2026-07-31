package stage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"roguefront/pkg/app"

	"github.com/sqweek/dialog"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type MainMenu struct {
	Profile string `json:"profile"`
	Game    string `json:"game"`
	Mods    string `json:"mods"`

	Name       string `json:"name"`
	Nation     int32  `json:"nation"`
	Enemy      int32  `json:"enemy"`
	Difficulty int32  `json:"difficulty"`
	Reslvl     int32  `json:"reslvl"`
	Fow        bool   `json:"fow"`

	Finest     bool `json:"finest"`
	Liberation bool `json:"liberation"`
	Talvisota  bool `json:"talvisota"`
	Scorched   bool `json:"scorched"`
	Airborne   bool `json:"airborne"`

	nations    string   `json:"-"`
	nationlist []string `json:"-"`
	levels     string   `json:"-"`
	levlist    []string `json:"-"`
	resources  string   `json:"-"`
	reslist    []int    `json:"-"`

	modlist []string `json:"-"`
	chkmods []bool   `json:"-"`

	panel       bool         `json:"-"`
	ddnation    bool         `json:"-"`
	ddenemy     bool         `json:"-"`
	ddlevel     bool         `json:"-"`
	ddresources bool         `json:"-"`
	scroll      rl.Vector2   `json:"-"`
	view        rl.Rectangle `json:"-"`
}

func (s *MainMenu) Init() error {
	if s == nil {
		return fmt.Errorf("Invalid stage")
	}

	s.Profile = ""
	s.Game = ""
	s.Mods = ""

	s.Name = "New Campaign"
	s.Nation = 0
	s.Enemy = 0
	s.Difficulty = 1
	s.Reslvl = 1
	s.Fow = true

	s.Finest = false
	s.Liberation = false
	s.Talvisota = false
	s.Scorched = false
	s.Airborne = false

	s.nations = "#01#Germany;#02#Soviets"
	s.nationlist = []string{"ger", "sov"}
	s.levels = "#01#Easy;#02#Normal;#03#Hard;#04#Heroic"
	s.levlist = []string{"easy", "normal", "hard", "heroic"}
	s.resources = "#01#Low;#02#Standard;#03#High"
	s.reslist = []int{1, 2, 3}
	s.modlist = []string{}
	s.chkmods = []bool{}

	s.panel = false
	s.ddnation = false
	s.ddenemy = false
	s.ddlevel = false
	s.ddresources = false
	s.scroll = rl.NewVector2(0, 0)
	s.view = rl.NewRectangle(0, 0, 0, 0)

	return s.LoadSettings()
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

	height := float32(app.CurApp.GetHeight())
	width := float32(app.CurApp.GetWidth())

	if s.panel {
		gui.SetState(gui.STATE_DISABLED)
	} else {
		gui.SetState(gui.STATE_NORMAL)
	}

	// main buttons
	nwbtn := gui.Button(rl.NewRectangle(0, height-120, 100, 40), "New Game")
	if nwbtn {
		s.panel = true
		s.UpdateMods()
	}
	ldbtn := gui.Button(rl.NewRectangle(0, height-80, 100, 40), "Load Game")
	if ldbtn {
		file, err := dialog.File().
			Filter("Json Files", "json").
			Title("Load Previous Roguefront Campaign").Load()
		if err == nil {
			s.LoadGame(file)
		}
	}
	exbtn := gui.Button(rl.NewRectangle(0, height-40, 100, 40), "Exit")
	if exbtn {
		app.CurApp.Exit()
	}

	// settings panel
	gui.Panel(rl.NewRectangle(100, height-120, width-100, 120), "Settings")
	gui.Label(rl.NewRectangle(105, height-95, 100, 20), "User Profile Path:")
	gui.TextBox(rl.NewRectangle(205, height-95, width-209, 20), &s.Profile, 16, true)
	prfile := gui.Button(rl.NewRectangle(width-22, height-95, 20, 20), "...")
	if prfile {
		file, err := dialog.Directory().
			Title("Select User Profile Path for Call to Arms Gates of Hell Ostfront").Browse()
		if err == nil {
			s.Profile = file
		}
	}
	gui.Label(rl.NewRectangle(105, height-70, 100, 20), "Game Directory:")
	gui.TextBox(rl.NewRectangle(205, height-70, width-209, 20), &s.Game, 16, true)
	gmfile := gui.Button(rl.NewRectangle(width-22, height-70, 20, 20), "...")
	if gmfile {
		file, err := dialog.Directory().
			Title("Select Game Directory for Call to Arms Gates of Hell Ostfront").Browse()
		if err == nil {
			s.Game = file
		}
	}
	gui.Label(rl.NewRectangle(105, height-45, 100, 20), "Mods Directory:")
	gui.TextBox(rl.NewRectangle(205, height-45, width-209, 20), &s.Mods, 16, true)
	mdfile := gui.Button(rl.NewRectangle(width-22, height-45, 20, 20), "...")
	if mdfile {
		file, err := dialog.Directory().
			Title("Select Mods Directory for Call to Arms Gates of Hell Ostfront").Browse()
		if err == nil {
			s.Mods = file
		}
	}

	// new game panel
	if s.panel {
		s.UpdateNations()
		s.UpdateLevels()
		s.UpdateResources()

		gui.SetState(gui.STATE_NORMAL)
		gui.Panel(rl.NewRectangle((width-500)/2, (height-400)/2, 500, 400), "New Game")
		cnbtn := gui.Button(rl.NewRectangle((width-500)/2+478, (height-400)/2+2, 20, 20), "X")
		if cnbtn {
			s.panel = false
		}
		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+370, 80, 20), "Campaign Name:")
		gui.TextBox(rl.NewRectangle((width-500)/2+95, (height-400)/2+370, 285, 20), &s.Name, 16, true)
		stbtn := gui.Button(rl.NewRectangle((width-500)/2+390, (height-400)/2+370, 100, 20), "Start Game")
		if stbtn {
			s.NewGame()
		}

		gui.ScrollPanel(rl.NewRectangle((width-500)/2+10, (height-400)/2+187, 480, 175), "", rl.NewRectangle(0, 0, 465, float32(len(s.modlist)*25)), &s.scroll, &s.view)
		rl.BeginScissorMode(int32(s.view.X), int32(s.view.Y), int32(s.view.Width), int32(s.view.Height))
		for i, mod := range s.modlist {
			gui.Label(rl.NewRectangle((width-500)/2+20, (height-400)/2+187+float32(25*i+5)+s.scroll.Y, 80, 20), fmt.Sprintf("%s %d", mod, i))
		}
		rl.EndScissorMode()

		gui.Line(rl.NewRectangle((width-500)/2+10, (height-400)/2+167, 480, 20), "Mods")
		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+147, 80, 20), "Airborne:")
		gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+152, 10, 10), "  ", &s.Airborne)
		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+122, 80, 20), "Scorched Earth:")
		gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+127, 10, 10), "  ", &s.Scorched)
		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+97, 80, 20), "Talvisota:")
		gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+102, 10, 10), "  ", &s.Talvisota)
		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+72, 80, 20), "Liberation:")
		gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+77, 10, 10), "  ", &s.Liberation)
		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+47, 80, 20), "Finest Hour:")
		gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+52, 10, 10), "  ", &s.Finest)
		gui.Line(rl.NewRectangle((width-500)/2+360, (height-400)/2+27, 90, 20), "DLC / Maps")

		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+147, 80, 20), "Fog of War:")
		gui.CheckBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+152, 10, 10), "  ", &s.Fow)
		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+122, 80, 20), "Resources:")
		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+122, 220, 20), s.resources, &s.Reslvl, s.ddresources) {
			s.ddresources = !s.ddresources
		}
		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+92, 80, 20), "Difficulty:")
		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+92, 220, 20), s.levels, &s.Difficulty, s.ddlevel) {
			s.ddlevel = !s.ddlevel
		}
		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+62, 80, 20), "Enemy Nation:")
		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+62, 220, 20), s.nations, &s.Enemy, s.ddenemy) {
			s.ddenemy = !s.ddenemy
		}
		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+32, 80, 20), "Player Nation:")
		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+32, 220, 20), s.nations, &s.Nation, s.ddnation) {
			s.ddnation = !s.ddnation
		}
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

func (s *MainMenu) UpdateNations() {
	s.nations = "#01#Germany;#02#Soviets"
	s.nationlist = []string{"ger", "sov"}

	if s.Finest {
		s.nationlist = append(s.nationlist, "eng")
		s.nations = s.nations + fmt.Sprintf(";#0%d#Commonwealth (GOH)", len(s.nationlist))
	}
	if s.Liberation {
		s.nationlist = append(s.nationlist, "usa")
		s.nations = s.nations + fmt.Sprintf(";#0%d#United States (GOH)", len(s.nationlist))
	}
	if s.Talvisota {
		s.nationlist = append(s.nationlist, "fin")
		s.nations = s.nations + fmt.Sprintf(";#0%d#Finland", len(s.nationlist))
	}
	if s.Nation >= int32(len(s.nationlist)) {
		s.Nation = 0
	}
	if s.Enemy >= int32(len(s.nationlist)) {
		s.Enemy = 0
	}
}

func (s *MainMenu) UpdateLevels() {
	//TODO: go through mods for modded levels

	if s.Difficulty >= int32(len(s.levlist)) {
		s.Difficulty = 0
	}
}

func (s *MainMenu) UpdateResources() {
	//TODO: go through mods for modded resources

	if s.Reslvl >= int32(len(s.reslist)) {
		s.Reslvl = 0
	}
}

func (s *MainMenu) UpdateMods() {
	//TODO: go through mod directory for all mods
}

func (s *MainMenu) NewGame() {
	if s == nil {
		return
	}

	game := Game{
		profile: s.Profile,
		game:    s.Game,
		mods:    s.Mods,
	}
	err := game.Init()
	game.state.Status.Name = s.Name
	game.state.Status.Army = s.nationlist[s.Nation]
	game.state.Status.EnemyArmy = s.nationlist[s.Enemy]
	game.state.Status.Difficulty = s.levlist[s.Difficulty]
	game.state.Status.Resources = s.reslist[s.Reslvl]
	if s.Fow {
		game.state.Status.FogOfWar = "fog_realistic"
	} else {
		game.state.Status.FogOfWar = "fog_realistic"
	}
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	s.SaveSettings()
	app.CurApp.SetStage(&game)
}

func (s *MainMenu) LoadGame(file string) {
	//TODO: implement loading game logic

	s.SaveSettings()
}

func (s *MainMenu) LoadSettings() error {
	// Get the absolute path of the running executable
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	// Extract the directory from the executable path
	exPath := filepath.Dir(ex)

	// Load the settings file
	data, err := os.ReadFile(exPath + "/settings.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	// Unmarshall the settings
	return json.Unmarshal(data, s)
}

func (s *MainMenu) SaveSettings() error {
	// Marshall the settings
	data, err := json.Marshal(*s)
	if err != nil {
		return err
	}

	// Get the absolute path of the running executable
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	// Extract the directory from the executable path
	exPath := filepath.Dir(ex)

	// Save the settings file
	return os.WriteFile(exPath+"/settings.json", data, 0644)
}
