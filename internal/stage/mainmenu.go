package stage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"roguefront/pkg/app"
	"roguefront/pkg/data"
	// "github.com/sqweek/dialog"
)

var ()

type MainMenu struct {
	Profile  string `json:"profile"`
	Game     string `json:"game"`
	Workshop string `json:"workshop"`

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

	modlist []data.Mod `json:"-"`

	// profilemode  bool         `json:"-"`
	// gamemode     bool         `json:"-"`
	// workshopmode bool         `json:"-"`
	// panel        bool         `json:"-"`
	// ddnation     bool         `json:"-"`
	// ddenemy      bool         `json:"-"`
	// ddlevel      bool         `json:"-"`
	// ddresources  bool         `json:"-"`
	// scroll       rl.Vector2   `json:"-"`
	// view         rl.Rectangle `json:"-"`
}

func (s *MainMenu) Init() error {
	if s == nil {
		return fmt.Errorf("Invalid stage")
	}

	s.Profile = ""
	s.Game = ""
	s.Workshop = ""

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
	s.modlist = []data.Mod{}

	// s.profilemode = false
	// s.gamemode = false
	// s.workshopmode = false
	// s.panel = false
	// s.ddnation = false
	// s.ddenemy = false
	// s.ddlevel = false
	// s.ddresources = false
	// s.scroll = rl.NewVector2(0, 0)
	// s.view = rl.NewRectangle(0, 0, 0, 0)

	err := s.LoadSettings()
	if err != nil {
		return err
	}
	return s.DetectSettings()
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

	// rl.ClearBackground(rl.Black)
	//
	// height := float32(app.CurApp.GetHeight())
	// width := float32(app.CurApp.GetWidth())
	//
	//	if s.panel {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// // main buttons
	// nwbtn := gui.Button(rl.NewRectangle(0, height-120, 100, 40), "New Game")
	//
	//	if nwbtn {
	//		s.panel = true
	//		s.UpdateMods()
	//		s.UpdateNations()
	//		s.UpdateLevels()
	//		s.UpdateResources()
	//	}
	//
	// ldbtn := gui.Button(rl.NewRectangle(0, height-80, 100, 40), "Load Game")
	//
	//	if ldbtn {
	//		file, err := dialog.File().
	//			Filter("Json Files", "json").
	//			Title("Load Previous Roguefront Campaign").Load()
	//		if err == nil {
	//			s.LoadGame(file)
	//		}
	//	}
	//
	// exbtn := gui.Button(rl.NewRectangle(0, height-40, 100, 40), "Exit")
	//
	//	if exbtn {
	//		app.CurApp.Exit()
	//	}
	//
	// // settings panel
	// gui.Panel(rl.NewRectangle(100, height-120, width-100, 120), "Settings")
	// gui.Label(rl.NewRectangle(105, height-95, 100, 20), "User Profile Path:")
	//
	//	if gui.TextBox(rl.NewRectangle(205, height-95, width-209, 20), &s.Profile, 260, s.profilemode) {
	//		s.profilemode = !s.profilemode
	//	}
	//
	// prfile := gui.Button(rl.NewRectangle(width-22, height-95, 20, 20), "...")
	//
	//	if prfile {
	//		file, err := dialog.Directory().
	//			Title("Select User Profile Path for Call to Arms Gates of Hell Ostfront (Documents\\My Games\\gates of hell\\profiles\\#######)").Browse()
	//		if err == nil {
	//			s.Profile = file
	//		}
	//	}
	//
	// gui.Label(rl.NewRectangle(105, height-70, 100, 20), "Game Directory:")
	//
	//	if gui.TextBox(rl.NewRectangle(205, height-70, width-209, 20), &s.Game, 260, s.gamemode) {
	//		s.gamemode = !s.gamemode
	//	}
	//
	// gmfile := gui.Button(rl.NewRectangle(width-22, height-70, 20, 20), "...")
	//
	//	if gmfile {
	//		file, err := dialog.Directory().
	//			Title("Select Game Directory for Call to Arms Gates of Hell Ostfront").Browse()
	//		if err == nil {
	//			s.Game = file
	//		}
	//	}
	//
	// gui.Label(rl.NewRectangle(105, height-45, 100, 20), "Mods Directory:")
	//
	//	if gui.TextBox(rl.NewRectangle(205, height-45, width-209, 20), &s.Workshop, 260, s.workshopmode) {
	//		s.workshopmode = !s.workshopmode
	//	}
	//
	// mdfile := gui.Button(rl.NewRectangle(width-22, height-45, 20, 20), "...")
	//
	//	if mdfile {
	//		file, err := dialog.Directory().
	//			Title("Select Mods Directory for Call to Arms Gates of Hell Ostfront").Browse()
	//		if err == nil {
	//			s.Workshop = file
	//		}
	//	}
	//
	// // new game panel
	//
	//	if s.panel {
	//		gui.SetState(gui.STATE_NORMAL)
	//		gui.Panel(rl.NewRectangle((width-500)/2, (height-400)/2, 500, 400), "New Game")
	//		cnbtn := gui.Button(rl.NewRectangle((width-500)/2+478, (height-400)/2+2, 20, 20), "X")
	//		if cnbtn {
	//			s.panel = false
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+370, 80, 20), "Campaign Name:")
	//		gui.TextBox(rl.NewRectangle((width-500)/2+95, (height-400)/2+370, 285, 20), &s.Name, 260, true)
	//		stbtn := gui.Button(rl.NewRectangle((width-500)/2+390, (height-400)/2+370, 100, 20), "Start Game")
	//		if stbtn {
	//			s.NewGame()
	//		}
	//
	//		gui.ScrollPanel(rl.NewRectangle((width-500)/2+10, (height-400)/2+187, 480, 175), "", rl.NewRectangle(0, 0, 465, float32(len(s.modlist)*25)), &s.scroll, &s.view)
	//		rl.BeginScissorMode(int32(s.view.X), int32(s.view.Y), int32(s.view.Width), int32(s.view.Height))
	//		for i, mod := range s.modlist {
	//			if gui.CheckBox(rl.NewRectangle((width-500)/2+20, (height-400)/2+192+float32(25*i+5)+s.scroll.Y, 10, 10), mod.Name, &s.modlist[i].Enabled) {
	//				s.UpdateNations()
	//				s.UpdateLevels()
	//				s.UpdateResources()
	//			}
	//		}
	//		rl.EndScissorMode()
	//
	//		gui.Line(rl.NewRectangle((width-500)/2+10, (height-400)/2+167, 480, 20), "Mods")
	//		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+147, 80, 20), "Airborne:")
	//		if gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+152, 10, 10), "  ", &s.Airborne) {
	//			s.UpdateNations()
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+122, 80, 20), "Scorched Earth:")
	//		if gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+127, 10, 10), "  ", &s.Scorched) {
	//			s.UpdateNations()
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+97, 80, 20), "Talvisota:")
	//		if gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+102, 10, 10), "  ", &s.Talvisota) {
	//			s.UpdateNations()
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+72, 80, 20), "Liberation:")
	//		if gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+77, 10, 10), "  ", &s.Liberation) {
	//			s.UpdateNations()
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+360, (height-400)/2+47, 80, 20), "Finest Hour:")
	//		if gui.CheckBox(rl.NewRectangle((width-500)/2+440, (height-400)/2+52, 10, 10), "  ", &s.Finest) {
	//			s.UpdateNations()
	//		}
	//		gui.Line(rl.NewRectangle((width-500)/2+360, (height-400)/2+27, 90, 20), "DLC / Maps")
	//
	//		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+147, 80, 20), "Fog of War:")
	//		gui.CheckBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+152, 10, 10), "  ", &s.Fow)
	//		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+122, 80, 20), "Resources:")
	//		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+122, 220, 20), s.resources, &s.Reslvl, s.ddresources) {
	//			s.ddresources = !s.ddresources
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+92, 80, 20), "Difficulty:")
	//		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+92, 220, 20), s.levels, &s.Difficulty, s.ddlevel) {
	//			s.ddlevel = !s.ddlevel
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+62, 80, 20), "Enemy Nation:")
	//		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+62, 220, 20), s.nations, &s.Enemy, s.ddenemy) {
	//			s.ddenemy = !s.ddenemy
	//		}
	//		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+32, 80, 20), "Player Nation:")
	//		if gui.DropdownBox(rl.NewRectangle((width-500)/2+90, (height-400)/2+32, 220, 20), s.nations, &s.Nation, s.ddnation) {
	//			s.ddnation = !s.ddnation
	//		}
	//	}
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

func (s *MainMenu) NewGame() {
	if s == nil {
		return
	}

	game := Game{}
	err := game.Init()
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	game.Profile = s.Profile
	game.Game = s.Game
	game.Workshop = s.Workshop

	game.Nations = s.nationlist

	game.State.Status.Name = s.Name
	game.State.Status.Army = s.nationlist[s.Nation]
	game.State.Status.EnemyArmy = s.nationlist[s.Enemy]
	game.State.Status.Difficulty = strings.ToLower(s.levlist[s.Difficulty])
	game.State.Status.Resources = int(s.Reslvl)
	for _, mod := range s.modlist {
		if mod.Enabled {
			game.Mods = append(game.Mods, mod)
			game.State.Status.Mods = append(game.State.Status.Mods, fmt.Sprintf("%s:0", mod.Id))
		}
	}
	if s.Fow {
		game.State.Status.FogOfWar = "fog_realistic"
	} else {
		game.State.Status.FogOfWar = "fog_off"
	}

	game.Money = int64(s.Reslvl*300) + 300

	err = game.Populate()
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	game.ReRoll()

	s.SaveSettings()
	app.CurApp.SetStage(&game)
}

func (s *MainMenu) LoadGame(file string) {
	if s == nil {
		return
	}

	game := Game{}
	err := game.Init()
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	err = game.LoadGame(file)
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	s.SaveSettings()
	app.CurApp.SetStage(&game)
}

func (s *MainMenu) LoadSettings() error {
	if s == nil {
		return fmt.Errorf("MainMenu is nil")
	}
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
	if s == nil {
		return fmt.Errorf("MainMenu is nil")
	}
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

func (s *MainMenu) DetectSettings() error {
	if s == nil {
		return fmt.Errorf("MainMenu is nil")
	}

	if s.Profile == "" {
		currentUser, err := user.Current()
		if err == nil {
			username := strings.Split(currentUser.Username, "\\")
			if len(username) > 1 {
				name := username[1]
				path := fmt.Sprintf("C:\\Users\\%s\\Documents\\My Games\\gates of hell\\profiles", name)
				info, err := os.Stat(path)
				if err == nil && info.IsDir() {
					entries, err := os.ReadDir(path)
					if err == nil {
						for _, entry := range entries {
							if entry.IsDir() {
								s.Profile = path + "\\" + entry.Name() + "\\campaign"
								break
							}
						}
					}
				}
			}
		}
	}

	if s.Game == "" {
		info, err := os.Stat("C:\\Program Files (x86)\\Steam\\steamapps\\common\\Call to Arms - Gates of Hell")
		if err == nil && info.IsDir() {
			s.Game = "C:\\Program Files (x86)\\Steam\\steamapps\\common\\Call to Arms - Gates of Hell"
		}
	}

	if s.Workshop == "" {
		info, err := os.Stat("C:\\Program Files (x86)\\Steam\\steamapps\\workshop\\content\\400750")
		if err == nil && info.IsDir() {
			s.Workshop = "C:\\Program Files (x86)\\Steam\\steamapps\\workshop\\content\\400750"
		}
	}

	return nil
}

func (s *MainMenu) UpdateNations() {
	nations := ""
	nationlist := []string{}

	// pretty names for nations
	namelist := map[string]string{
		"ger":  "Germany",
		"rus":  "Soviet Union",
		"fin":  "Finland",
		"usa":  "United States (GoH)",
		"eng":  "Commonwealth (GoH)",
		"aus":  "Austria",
		"fra":  "France",
		"pol":  "Poland",
		"hun":  "Hungary",
		"jap":  "Japan",
		"ita":  "Italy",
		"usaf": "United States",
		"uk":   "Commonwealth",
	}

	nset := data.FindNations(s.Game, s.modlist)
	for _, nation := range nset {
		if nation == "eng" && !s.Finest {
			continue
		}
		if nation == "usa" && !s.Liberation {
			continue
		}
		if nation == "fin" && !s.Talvisota {
			continue
		}
		nationlist = append(s.nationlist, nation)
		name := nation
		if n, ok := namelist[nation]; ok {
			name = n
		}
		nations = nations + fmt.Sprintf("#%02d#%s;", len(nationlist), name)
	}

	if len(nations) > 0 {
		nations = nations[:len(nations)-1]
		s.nations = nations
		s.nationlist = nationlist
		if s.Nation >= int32(len(s.nationlist)) {
			s.Nation = 0
		}
		if s.Enemy >= int32(len(s.nationlist)) {
			s.Enemy = 0
		}
	}
}

func (s *MainMenu) UpdateLevels() {
	levels := ""
	levlist := []string{}

	lset := data.FindLevels(s.Game, s.modlist)
	for _, level := range lset {
		levlist = append(levlist, level)
		levels = levels + fmt.Sprintf("#%02d#%s;", len(levlist), level)
	}

	if len(levels) > 0 {
		levels = levels[:len(levels)-1]
		s.levels = levels
		s.levlist = levlist
		if s.Difficulty >= int32(len(s.levlist)) {
			s.Difficulty = 0
		}
	}
}

func (s *MainMenu) UpdateResources() {
	resources := ""

	rset := data.FindResources(s.Game, s.modlist)
	for i, level := range rset {
		resources = resources + fmt.Sprintf("#%02d#%s;", i+1, level)
	}

	if len(resources) > 0 {
		resources = resources[:len(resources)-1]
		s.resources = resources
		if s.Reslvl >= int32(len(s.resources)) {
			s.Reslvl = 0
		}
	}
}

func (s *MainMenu) UpdateMods() (err error) {
	if len(s.modlist) > 0 {
		return nil
	}

	s.modlist, err = data.FindMods(s.Workshop)
	return err
}
