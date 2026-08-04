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

	"github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
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

	namelist   map[string]string
	nations    []giu.Widget `json:"-"`
	enemies    []giu.Widget `json:"-"`
	nationlist []string     `json:"-"`
	levels     []giu.Widget `json:"-"`
	levlist    []string     `json:"-"`
	resources  []giu.Widget `json:"-"`

	modlist []data.Mod `json:"-"`
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

	// pretty names for nations
	s.namelist = map[string]string{
		"ger":  "Germany",
		"rus":  "Soviet Union",
		"sov":  "Soviet Union",
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

	s.nationlist = []string{"ger", "rus"}
	s.nations = []giu.Widget{
		giu.Selectable("Germany").OnClick(func() { s.Nation = 0 }),
		giu.Selectable("Soviet Union").OnClick(func() { s.Nation = 1 }),
	}
	s.enemies = []giu.Widget{
		giu.Selectable("Germany").OnClick(func() { s.Enemy = 0 }),
		giu.Selectable("Soviet Union").OnClick(func() { s.Enemy = 1 }),
	}
	s.levlist = []string{"easy", "normal", "hard", "heroic"}
	s.levels = []giu.Widget{
		giu.Selectable("Easy").OnClick(func() { s.Difficulty = 0 }),
		giu.Selectable("Normal").OnClick(func() { s.Difficulty = 1 }),
		giu.Selectable("Hard").OnClick(func() { s.Difficulty = 2 }),
		giu.Selectable("Heroic").OnClick(func() { s.Difficulty = 3 }),
	}
	s.resources = []giu.Widget{
		giu.Selectable("Low").OnClick(func() { s.Reslvl = 0 }),
		giu.Selectable("Medium").OnClick(func() { s.Reslvl = 1 }),
		giu.Selectable("High").OnClick(func() { s.Reslvl = 2 }),
	}
	s.modlist = []data.Mod{}

	err := s.LoadSettings()
	if err != nil {
		return err
	}
	return s.DetectSettings()
}

func (s *MainMenu) Build() {
	if s == nil {
		return
	}

	// height := float32(app.CurApp.GetHeight())
	width := float32(app.CurApp.GetWidth())
	giu.Column(
		giu.Row(
			// main menu
			giu.Dummy(0, 0),
			giu.Column(
				giu.Dummy(0, 0),
				giu.Button("New Game").Size(100, 27).OnClick(func() {
					s.UpdateMods()
					s.UpdateNations()
					s.UpdateLevels()
					s.UpdateResources()
					giu.OpenPopup("New Game")
				}),
				// new game
				giu.PopupModal("New Game").Layout(
					giu.Row(
						giu.Dummy(0, 0),
						giu.Column(
							giu.Label("------------------- Settings: -------------------"),
							giu.Buttonf("Player Nation: %s", s.namelist[s.nationlist[s.Nation]]).Size(200, 20).OnClick(func() {
								giu.OpenPopup("Player")
							}),
							giu.Popup("Player").Layout(
								s.nations...,
							),
							giu.Buttonf("Enemy Nation: %s", s.namelist[s.nationlist[s.Enemy]]).Size(200, 20).OnClick(func() {
								giu.OpenPopup("Enemy")
							}),
							giu.Popup("Enemy").Layout(
								s.enemies...,
							),
							giu.Buttonf("Difficulty: %s", s.levlist[s.Difficulty]).Size(200, 20).OnClick(func() {
								giu.OpenPopup("Difficulty")
							}),
							giu.Popup("Difficulty").Layout(
								s.levels...,
							),
							giu.Buttonf("Resources: %d", s.Reslvl).Size(200, 20).OnClick(func() {
								giu.OpenPopup("Resources")
							}),
							giu.Popup("Resources").Layout(
								s.resources...,
							),
							giu.Checkbox("Fog of War:", &s.Fow),
							giu.Dummy(0, 0),
						),
						giu.Dummy(0, 0),
						giu.Column(
							giu.Label("---- DLC / Maps: ----"),
							giu.Checkbox("Airborne", &s.Airborne),
							giu.Checkbox("Scorched Earth", &s.Scorched),
							giu.Checkbox("Talvisota", &s.Talvisota),
							giu.Checkbox("Liberation", &s.Liberation),
							giu.Checkbox("Finest Hour", &s.Finest),
							giu.Dummy(0, 0),
						),
						giu.Dummy(0, 0),
					),
					giu.Row(
						giu.Dummy(0, 0),
						giu.Label("------------------------------------ Mods: ------------------------------------"),
						giu.Dummy(0, 0),
					),
					giu.Row(
						giu.Label("       Name:"),
						giu.InputText(&s.Name),
					),
					giu.Row(
						giu.Dummy(50, 0),
						giu.Button("Start Game").Size(100, 20).OnClick(func() {
							giu.CloseCurrentPopup()
							s.NewGame()
						}),
						giu.Dummy(0, 0),
						giu.Button("Cancel").Size(100, 20).OnClick(func() {
							giu.CloseCurrentPopup()
						}),
						giu.Dummy(50, 0),
					),
				),
				giu.Button("Load Game").Size(100, 27).OnClick(func() {
					file, err := dialog.File().
						Filter("Json Files", "json").
						Title("Load Previous Roguefront Campaign").Load()
					if err == nil {
						s.LoadGame(file)
					}
				}),
				giu.Button("Exit").Size(100, 27).OnClick(func() {
					app.CurApp.Exit()
				}),
				giu.Dummy(0, 0),
			),
			giu.Dummy(0, 0),
			// settings panel
			giu.Column(
				giu.Dummy(0, 0),
				giu.Label("Settings:"),
				giu.Row(
					giu.Dummy(0, 0),
					giu.Column(
						giu.Dummy(0, 1),
						giu.Label("User Profile Path:"),
						giu.Dummy(0, 2),
						giu.Label("Game Directory:"),
						giu.Dummy(0, 2),
						giu.Label("Mods Directory:"),
						giu.Dummy(0, 1),
					),
					giu.Column(
						giu.Dummy(0, 0),
						giu.InputText(&s.Profile).Size(width-320),
						giu.InputText(&s.Game).Size(width-320),
						giu.InputText(&s.Workshop).Size(width-320),
						giu.Dummy(0, 0),
					),
					giu.Column(
						giu.Dummy(0, 0),
						giu.Button("Browse").Size(60, 20).OnClick(func() {
							file, err := dialog.Directory().
								Title("Select User Profile Path for Call to Arms Gates of Hell Ostfront (Documents\\My Games\\gates of hell\\profiles\\#######)").Browse()
							if err == nil {
								s.Profile = file
							}
						}),
						giu.Button("Browse").Size(60, 20).OnClick(func() {
							file, err := dialog.Directory().
								Title("Select Game Directory for Call to Arms Gates of Hell Ostfront").Browse()
							if err == nil {
								s.Game = file
							}
						}),
						giu.Button("Browse").Size(60, 20).OnClick(func() {
							file, err := dialog.Directory().
								Title("Select Mods Directory for Call to Arms Gates of Hell Ostfront").Browse()
							if err == nil {
								s.Workshop = file
							}
						}),
						giu.Dummy(0, 0),
					),
				),
				giu.Dummy(0, 0),
			),
			giu.Dummy(100, 100),
		),
	).Build()
}

func (s *MainMenu) OnResize(w int, h int) {
	if s == nil {
		return
	}
}

func (s *MainMenu) Update(dt int64) {
	if s == nil {
		return
	}
}

func (s *MainMenu) OnInput(dt int64) {
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
	nations := []giu.Widget{}
	enemies := []giu.Widget{}
	nationlist := []string{}

	nset := data.FindNations(s.Game, s.modlist)
	for i, nation := range nset {
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
		if n, ok := s.namelist[nation]; ok {
			name = n
		} else {
			s.namelist[nation] = nation
		}
		nations = append(nations, giu.Selectable(name).OnClick(func() { s.Nation = int32(i) }))
		enemies = append(enemies, giu.Selectable(name).OnClick(func() { s.Enemy = int32(i) }))
	}

	if len(nations) > 0 {
		s.nations = nations
		s.enemies = enemies
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
	levels := []giu.Widget{}
	levlist := []string{}

	lset := data.FindLevels(s.Game, s.modlist)
	for i, level := range lset {
		levlist = append(levlist, level)
		levels = append(levels, giu.Selectable(level).OnClick(func() { s.Difficulty = int32(i) }))
	}

	if len(levels) > 0 {
		s.levels = levels
		s.levlist = levlist
		if s.Difficulty >= int32(len(s.levlist)) {
			s.Difficulty = 0
		}
	}
}

func (s *MainMenu) UpdateResources() {
	resources := []giu.Widget{}

	rset := data.FindResources(s.Game, s.modlist)
	for i, level := range rset {
		resources = append(resources, giu.Selectable(level).OnClick(func() { s.Reslvl = int32(i) }))
	}

	if len(resources) > 0 {
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
