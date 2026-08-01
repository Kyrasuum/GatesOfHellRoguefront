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
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"

	"roguefront/pkg/app"
	"roguefront/res"

	"github.com/sqweek/dialog"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
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
	reslist    []int    `json:"-"`

	modname []string `json:"-"`
	modlist []string `json:"-"`
	chkmods []bool   `json:"-"`

	profilemode  bool         `json:"-"`
	gamemode     bool         `json:"-"`
	workshopmode bool         `json:"-"`
	panel        bool         `json:"-"`
	ddnation     bool         `json:"-"`
	ddenemy      bool         `json:"-"`
	ddlevel      bool         `json:"-"`
	ddresources  bool         `json:"-"`
	scroll       rl.Vector2   `json:"-"`
	view         rl.Rectangle `json:"-"`
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
	s.reslist = []int{1, 2, 3}
	s.modname = []string{}
	s.modlist = []string{}
	s.chkmods = []bool{}

	s.profilemode = false
	s.gamemode = false
	s.workshopmode = false
	s.panel = false
	s.ddnation = false
	s.ddenemy = false
	s.ddlevel = false
	s.ddresources = false
	s.scroll = rl.NewVector2(0, 0)
	s.view = rl.NewRectangle(0, 0, 0, 0)

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
		s.UpdateNations()
		s.UpdateLevels()
		s.UpdateResources()
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
	if gui.TextBox(rl.NewRectangle(205, height-95, width-209, 20), &s.Profile, 260, s.profilemode) {
		s.profilemode = !s.profilemode
	}
	prfile := gui.Button(rl.NewRectangle(width-22, height-95, 20, 20), "...")
	if prfile {
		file, err := dialog.Directory().
			Title("Select User Profile Path for Call to Arms Gates of Hell Ostfront (Documents\\My Games\\gates of hell\\profiles\\#######)").Browse()
		if err == nil {
			s.Profile = file
		}
	}
	gui.Label(rl.NewRectangle(105, height-70, 100, 20), "Game Directory:")
	if gui.TextBox(rl.NewRectangle(205, height-70, width-209, 20), &s.Game, 260, s.gamemode) {
		s.gamemode = !s.gamemode
	}
	gmfile := gui.Button(rl.NewRectangle(width-22, height-70, 20, 20), "...")
	if gmfile {
		file, err := dialog.Directory().
			Title("Select Game Directory for Call to Arms Gates of Hell Ostfront").Browse()
		if err == nil {
			s.Game = file
		}
	}
	gui.Label(rl.NewRectangle(105, height-45, 100, 20), "Mods Directory:")
	if gui.TextBox(rl.NewRectangle(205, height-45, width-209, 20), &s.Workshop, 260, s.workshopmode) {
		s.workshopmode = !s.workshopmode
	}
	mdfile := gui.Button(rl.NewRectangle(width-22, height-45, 20, 20), "...")
	if mdfile {
		file, err := dialog.Directory().
			Title("Select Mods Directory for Call to Arms Gates of Hell Ostfront").Browse()
		if err == nil {
			s.Workshop = file
		}
	}

	// new game panel
	if s.panel {

		gui.SetState(gui.STATE_NORMAL)
		gui.Panel(rl.NewRectangle((width-500)/2, (height-400)/2, 500, 400), "New Game")
		cnbtn := gui.Button(rl.NewRectangle((width-500)/2+478, (height-400)/2+2, 20, 20), "X")
		if cnbtn {
			s.panel = false
		}
		gui.Label(rl.NewRectangle((width-500)/2+10, (height-400)/2+370, 80, 20), "Campaign Name:")
		gui.TextBox(rl.NewRectangle((width-500)/2+95, (height-400)/2+370, 285, 20), &s.Name, 260, true)
		stbtn := gui.Button(rl.NewRectangle((width-500)/2+390, (height-400)/2+370, 100, 20), "Start Game")
		if stbtn {
			s.NewGame()
		}

		gui.ScrollPanel(rl.NewRectangle((width-500)/2+10, (height-400)/2+187, 480, 175), "", rl.NewRectangle(0, 0, 465, float32(len(s.modlist)*25)), &s.scroll, &s.view)
		rl.BeginScissorMode(int32(s.view.X), int32(s.view.Y), int32(s.view.Width), int32(s.view.Height))
		for i, mod := range s.modname {
			gui.CheckBox(rl.NewRectangle((width-500)/2+20, (height-400)/2+192+float32(25*i+5)+s.scroll.Y, 10, 10), mod, &s.chkmods[i])
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

	LoadNations := func(path string) {
		r, err := zip.OpenReader(path)
		if err != nil {
			return
		}
		defer r.Close()

		for _, f := range r.File {
			switch {
			case strings.HasPrefix(f.Name, "set/multiplayer/armies/"):
				if strings.HasSuffix(f.Name, ".set") {
					split := strings.Split(f.Name, "/")
					name := split[len(split)-1]
					name = name[:len(name)-4]
					if !slices.Contains(s.nationlist[:], name) {
						s.nationlist = append(s.nationlist, name)
						s.nations = s.nations + fmt.Sprintf(";#0%d#"+name, len(s.nationlist))
					}
				}
			}
		}
	}

	for i, chk := range s.chkmods {
		if chk {
			paks, err := os.ReadDir(s.Workshop + "/" + s.modlist[i] + "/resource")
			if err != nil {
				log.Printf("%+v", err)
				continue
			}
			for _, pak := range paks {
				if !pak.IsDir() && strings.HasSuffix(pak.Name(), ".pak") {
					LoadNations(s.Workshop + "/" + s.modlist[i] + "/resource/" + pak.Name())
				}
			}

			// load loose files
			entries, err := os.ReadDir(s.Workshop + "/" + s.modlist[i])
			if err != nil {
				log.Printf("%+v", err)
				return
			}
			for _, entry := range entries {
				name := entry.Name()
				if !entry.IsDir() && strings.HasSuffix(name, ".set") {
					if !slices.Contains(s.nationlist[:], name[:3]) {
						s.nationlist = append(s.nationlist, name[:3])
						s.nations = s.nations + fmt.Sprintf(";#0%d#"+name[:3], len(s.nationlist))
					}
				}
			}
		}
	}
}

func (s *MainMenu) UpdateLevels() {
	ReadLevel := func(data []byte) {
		scanner := bufio.NewScanner(bytes.NewReader(data))

		levels := ""
		levlist := []string{}
		for scanner.Scan() {
			line := scanner.Text()

			if line == "" {
				continue
			}

			fmt.Sprintf("%s\n", line)
			switch {
			case strings.HasPrefix(line, "{GameModes"):
			case strings.HasPrefix(line, "{ResupplyCost"):
			case strings.HasPrefix(line, "{"):
				levels = levels + fmt.Sprintf("%s;", strings.TrimSpace(line)[1:])
				levlist = append(levlist, strings.TrimSpace(line)[1:])
			}
		}
		if len(levels) > 0 {
			levels = levels[:len(levels)-1]
		}
		s.levels = levels
		s.levlist = levlist
	}

	UpdateLevel := func(path string) {
		r, err := zip.OpenReader(path)
		if err != nil {
			return
		}
		defer r.Close()

		for _, f := range r.File {
			switch f.Name {
			case "set/dynamic_campaign.set":
				rc, err := f.Open()
				if err != nil {
					return
				}

				data, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					return
				}

				ReadLevel(data)
			}
		}
	}
	UpdateLevel(s.Game + "/resource/gamelogic.pak")
	for i, chk := range s.chkmods {
		if chk {
			UpdateLevel(s.Workshop + "/" + s.modlist[i] + "/resource/gamelogic.pak")
			data, err := os.ReadFile(s.Workshop + "/" + s.modlist[i] + "/resource/set/dynamic_campaign.set")
			if err == nil {
				ReadLevel(data)
			}
		}
	}

	if s.Difficulty >= int32(len(s.levlist)) {
		s.Difficulty = 0
	}
}

func (s *MainMenu) UpdateResources() {
	ReadResource := func(data []byte) {
		scanner := bufio.NewScanner(bytes.NewReader(data))

		for scanner.Scan() {
			line := scanner.Text()

			if line == "" {
				continue
			}

			switch {
			//find resources section
			case strings.HasPrefix(line, "\t{Resources"):
				resources := ""
				reslist := []int{}
				for scanner.Scan() {
					line := scanner.Text()

					if line == "" {
						continue
					}

					switch {
					case strings.HasPrefix(line, "\t\t{"):
						resources = resources + fmt.Sprintf("%s;", strings.TrimSpace(line)[1:])
						reslist = append(reslist, len(resources))
					case strings.HasPrefix(line, "\t}"):
						if len(resources) > 0 {
							resources = resources[:len(resources)-1]
						}
						s.resources = resources
						s.reslist = reslist
						return
					}
				}
			}
		}
	}

	UpdateResource := func(path string) {
		r, err := zip.OpenReader(path)
		if err != nil {
			return
		}
		defer r.Close()

		for _, f := range r.File {
			switch f.Name {
			case "set/dynamic_campaign.set":
				rc, err := f.Open()
				if err != nil {
					return
				}

				data, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					return
				}

				ReadResource(data)
			}
		}
	}
	UpdateResource(s.Game + "/resource/gamelogic.pak")
	for i, chk := range s.chkmods {
		if chk {
			UpdateResource(s.Workshop + "/" + s.modlist[i] + "/resource/gamelogic.pak")
			data, err := os.ReadFile(s.Workshop + "/" + s.modlist[i] + "/resource/set/dynamic_campaign.set")
			if err == nil {
				ReadResource(data)
			}
		}
	}

	if s.Reslvl >= int32(len(s.reslist)) {
		s.Reslvl = 0
	}
}

func (s *MainMenu) UpdateMods() {
	if len(s.modlist) > 0 {
		return
	}

	entries, err := os.ReadDir(s.Workshop)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	s.modname = []string{}
	s.modlist = []string{}
	s.chkmods = []bool{}

	for _, entry := range entries {
		// Filter out standard files, keeping only folders for mods
		if entry.IsDir() {
			data, err := os.ReadFile(s.Workshop + "/" + entry.Name() + "/mod.info")
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				log.Printf("%+v", err)
				continue
			}

			scanner := bufio.NewScanner(bytes.NewReader(data))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())

				if line == "" {
					continue
				}

				switch {
				case strings.HasPrefix(line, "{name"):
					s.modlist = append(s.modlist, entry.Name())
					s.chkmods = append(s.chkmods, false)
					s.modname = append(s.modname, res.ParseString(line))
					break
				}
			}
		}
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
	for i, chk := range s.chkmods {
		if chk {
			game.Mods = append(game.Mods, s.modlist[i])
			game.State.Status.Mods = append(game.State.Status.Mods, fmt.Sprintf("%s:0", s.modlist[i]))
		}
	}
	if s.Fow {
		game.State.Status.FogOfWar = "fog_realistic"
	} else {
		game.State.Status.FogOfWar = "fog_off"
	}
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
