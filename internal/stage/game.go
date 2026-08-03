package stage

import (
	"encoding/json"
	"errors"
	"fmt"
	// "log"
	// rand "math/rand/v2"
	"os"
	"path/filepath"
	// "slices"
	"time"

	// "roguefront/pkg/app"
	"roguefront/pkg/data"
	// "roguefront/pkg/rng"

	"github.com/iancoleman/strcase"
)

var ()

type Game struct {
	Profile  string `json:"profile"`
	Game     string `json:"game"`
	Workshop string `json:"workshop"`

	Money int64 `json:"money"`

	Maps    map[string][]string `json:"maps"`
	Mods    []data.Mod          `json:"mods"`
	Nations []string            `json:"nations"`

	Squads   map[string]map[string]data.Squad `json:"sqd"`
	Infantry map[string]*data.Infantry        `json:"inf"`
	Vehicles map[string]data.Vehicle          `json:"veh"`
	Items    map[string][]data.Item           `json:"itm"`

	Shop Store `json:"store"`

	State *data.Save `json:"state"`

	// paused    bool         `json:"-"`
	// preview   bool         `json:"-"`
	// tab       int          `json:"-"`
	// ownscroll rl.Vector2   `json:"-"`
	// ownview   rl.Rectangle `json:"-"`
	// strscroll rl.Vector2   `json:"-"`
	// strview   rl.Rectangle `json:"-"`
	// examine   []examinable `josn:"-"`
}

// initialize game object
func (g *Game) Init() error {
	if g == nil {
		return fmt.Errorf("Invalid stage")
	}

	g.Profile = ""
	g.Game = ""
	g.Workshop = ""

	g.Money = 0

	g.Maps = map[string][]string{}
	g.Mods = []data.Mod{}
	g.Nations = []string{}

	g.Squads = make(map[string]map[string]data.Squad)
	g.Infantry = make(map[string]*data.Infantry)
	g.Vehicles = make(map[string]data.Vehicle)
	g.Items = make(map[string][]data.Item)

	g.Shop = Store{
		Squads:   []StoreSquad{},
		Infantry: []StoreInfantry{},
		Vehicles: []StoreVehicle{},
		Items:    []StoreItem{},
	}

	g.State = &data.Save{
		Campaign: &data.Campaign{
			Infantry:    []*data.Infantry{},
			Vehicles:    []*data.Vehicle{},
			Inventories: []*data.Inventory{},
			Squads:      []*data.Squad{},
		},
		Status: &data.Status{
			Mods:        []string{},
			Timestamp:   time.Now().Unix(),
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

	// g.paused = false
	// g.preview = false
	// g.tab = 0
	// g.ownscroll = rl.NewVector2(0, 0)
	// g.ownview = rl.NewRectangle(0, 0, 0, 0)
	// g.strscroll = rl.NewVector2(0, 0)
	// g.strview = rl.NewRectangle(0, 0, 0, 0)
	// g.examine = []examinable{}

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

	// rl.ClearBackground(rl.Black)
	//
	// width := float32(app.CurApp.GetWidth())
	// height := float32(app.CurApp.GetHeight())
	//
	// gui.SetState(gui.STATE_NORMAL)
	// gui.Panel(rl.NewRectangle(0, 0, width, 20), "")
	// menu := gui.Button(rl.NewRectangle(width-50, 0, 50, 20), "Menu")
	//
	//	if menu {
	//		g.paused = true
	//	}
	//
	// preview := gui.Button(rl.NewRectangle(width-150, 0, 100, 20), "Preview Battle")
	//
	//	if preview {
	//		g.preview = true
	//	}
	//
	//	if finished, err := g.Finished(); err != nil || !finished {
	//		if err != nil {
	//			log.Printf("%+v\n", err)
	//		}
	//		gui.SetState(gui.STATE_DISABLED)
	//	}
	//
	// read := gui.Button(rl.NewRectangle(width/2, 0, 80, 20), "Read Result")
	//
	//	if read {
	//		err := g.State.Read(g.Profile + "/" + g.State.Status.Name + ".sav")
	//		if err != nil {
	//			log.Print(err)
	//		}
	//		g.ReRoll()
	//	}
	//
	// gui.SetState(gui.STATE_NORMAL)
	// start := gui.Button(rl.NewRectangle(width/2-80, 0, 80, 20), "Start Battle")
	//
	//	if start {
	//		g.State.Write(g.Profile)
	//	}
	//
	// gui.Label(rl.NewRectangle(8, 0, 80, 20), fmt.Sprintf("Money: %d", g.Money))
	//
	// gui.SetState(gui.STATE_NORMAL)
	// gui.Panel(rl.NewRectangle(0, 20, width/2, height-20), "Owned:")
	//
	//	if g.tab == 0 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// ownsqd := gui.Button(rl.NewRectangle(50, 22, 50, 20), "Squads")
	//
	//	if ownsqd {
	//		g.tab = 0
	//	}
	//
	//	if g.tab == 1 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// ownveh := gui.Button(rl.NewRectangle(100, 22, 50, 20), "Vehicles")
	//
	//	if ownveh {
	//		g.tab = 1
	//	}
	//
	//	if g.tab == 2 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// ownsdr := gui.Button(rl.NewRectangle(150, 22, 50, 20), "Soldiers")
	//
	//	if ownsdr {
	//		g.tab = 2
	//	}
	//
	//	if g.tab == 3 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// ownitm := gui.Button(rl.NewRectangle(200, 22, 50, 20), "Items")
	//
	//	if ownitm {
	//		g.tab = 3
	//	}
	//
	// gui.SetState(gui.STATE_NORMAL)
	// gui.Panel(rl.NewRectangle(width/2, 20, width/2, height-20), "Store:")
	//
	//	if g.tab == 0 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// strsqd := gui.Button(rl.NewRectangle(width/2+50, 22, 50, 20), "Squads")
	//
	//	if strsqd {
	//		g.tab = 0
	//	}
	//
	//	if g.tab == 1 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// strveh := gui.Button(rl.NewRectangle(width/2+100, 22, 50, 20), "Vehicles")
	//
	//	if strveh {
	//		g.tab = 1
	//	}
	//
	//	if g.tab == 2 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// strsdr := gui.Button(rl.NewRectangle(width/2+150, 22, 50, 20), "Soldiers")
	//
	//	if strsdr {
	//		g.tab = 2
	//	}
	//
	//	if g.tab == 3 {
	//		gui.SetState(gui.STATE_DISABLED)
	//	} else {
	//
	//		gui.SetState(gui.STATE_NORMAL)
	//	}
	//
	// stritm := gui.Button(rl.NewRectangle(width/2+200, 22, 50, 20), "Items")
	//
	//	if stritm {
	//		g.tab = 3
	//	}
	//
	// gui.SetState(gui.STATE_NORMAL)
	//
	// switch g.tab {
	// case 0:
	//
	//	gui.ScrollPanel(rl.NewRectangle(0, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.State.Campaign.Squads)*20+10)), &g.ownscroll, &g.ownview)
	//	rl.BeginScissorMode(int32(g.ownview.X), int32(g.ownview.Y), int32(g.ownview.Width), int32(g.ownview.Height))
	//	for i, ent := range g.State.Campaign.Squads {
	//		preview := gui.Button(rl.NewRectangle(11, float32(48+20*i)+g.strscroll.Y, 80, 20), "View")
	//		if preview {
	//			g.examine = append(g.examine, ent)
	//		}
	//		gui.Label(rl.NewRectangle(101, float32(48+20*i)+g.ownscroll.Y, width/2-124, 20), ent.Name)
	//	}
	//	rl.EndScissorMode()
	//
	// case 1:
	//
	//	gui.ScrollPanel(rl.NewRectangle(0, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.State.Campaign.Vehicles)*20+10)), &g.ownscroll, &g.ownview)
	//	rl.BeginScissorMode(int32(g.ownview.X), int32(g.ownview.Y), int32(g.ownview.Width), int32(g.ownview.Height))
	//	for i, ent := range g.State.Campaign.Vehicles {
	//		preview := gui.Button(rl.NewRectangle(11, float32(48+20*i)+g.strscroll.Y, 80, 20), "View")
	//		if preview {
	//			g.examine = append(g.examine, ent)
	//		}
	//		gui.Label(rl.NewRectangle(101, float32(48+20*i)+g.ownscroll.Y, width/2-124, 20), ent.Name)
	//	}
	//	rl.EndScissorMode()
	//
	// case 2:
	//
	//	gui.ScrollPanel(rl.NewRectangle(0, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.State.Campaign.Infantry)*20+10)), &g.ownscroll, &g.ownview)
	//	rl.BeginScissorMode(int32(g.ownview.X), int32(g.ownview.Y), int32(g.ownview.Width), int32(g.ownview.Height))
	//	for i, ent := range g.State.Campaign.Infantry {
	//		preview := gui.Button(rl.NewRectangle(11, float32(48+20*i)+g.strscroll.Y, 80, 20), "View")
	//		if preview {
	//			g.examine = append(g.examine, ent)
	//		}
	//		gui.Label(rl.NewRectangle(101, float32(48+20*i)+g.ownscroll.Y, width/2-124, 20), ent.Name)
	//	}
	//	rl.EndScissorMode()
	//
	// case 3:
	//
	//		gui.ScrollPanel(rl.NewRectangle(0, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.State.Campaign.Items)*20+10)), &g.ownscroll, &g.ownview)
	//		rl.BeginScissorMode(int32(g.ownview.X), int32(g.ownview.Y), int32(g.ownview.Width), int32(g.ownview.Height))
	//		for i, ent := range g.State.Campaign.Items {
	//			preview := gui.Button(rl.NewRectangle(11, float32(48+20*i)+g.strscroll.Y, 80, 20), "View")
	//			if preview {
	//				g.examine = append(g.examine, ent)
	//			}
	//			gui.Label(rl.NewRectangle(101, float32(48+20*i)+g.ownscroll.Y, width/2-124, 20), ent.Name)
	//		}
	//		rl.EndScissorMode()
	//	}
	//
	// switch g.tab {
	// case 0:
	//
	//	gui.ScrollPanel(rl.NewRectangle(width/2, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.Shop.Squads)*20+10)), &g.strscroll, &g.strview)
	//	rl.BeginScissorMode(int32(g.strview.X), int32(g.strview.Y), int32(g.strview.Width), int32(g.strview.Height))
	//	for i, ent := range g.Shop.Squads {
	//		if g.Money < int64(ent.Squad.Cost) {
	//			gui.SetState(gui.STATE_DISABLED)
	//		}
	//		strbuy := gui.Button(rl.NewRectangle(width/2+11, float32(48+20*i)+g.strscroll.Y, 80, 20), fmt.Sprintf("Buy %d", ent.Squad.Cost))
	//		if strbuy {
	//			sqd := ent.Squad
	//			for i, _ := range sqd.Soldiers {
	//				inf := &data.Infantry{}
	//				*inf = *sqd.Soldiers[i]
	//				inf.Nid = fmt.Sprintf("%02d %02d", rand.Int32N(10000), rand.Int32N(10000))
	//				items := inf.Inv.Items
	//				inf.Inv = &data.Inventory{
	//					Name:  inf.Inv.Name,
	//					Items: []data.Item{},
	//				}
	//				for _, item := range items {
	//					if item.Amount == 0 {
	//						inf.Inv.Items = append(inf.Inv.Items, item)
	//					} else {
	//						item.Amount = float64(rng.RollOnce(item.Amount))
	//						if item.Amount > 0 {
	//							inf.Inv.Items = append(inf.Inv.Items, item)
	//						}
	//					}
	//				}
	//				sqd.Soldiers[i] = inf
	//			}
	//			g.State.Campaign.Squads = append(g.State.Campaign.Squads, &sqd)
	//			g.Money -= int64(ent.Squad.Cost)
	//			g.Shop.Squads[i].Amount -= 1
	//			if g.Shop.Squads[i].Amount < 1 {
	//				g.Shop.Squads = slices.Delete(g.Shop.Squads, i, i+1)
	//			}
	//		}
	//		gui.SetState(gui.STATE_NORMAL)
	//		preview = gui.Button(rl.NewRectangle(width/2+91, float32(48+20*i)+g.strscroll.Y, 80, 20), "Preview")
	//		if preview {
	//			g.examine = append(g.examine, ent)
	//		}
	//		gui.Label(rl.NewRectangle(width/2+181, float32(48+20*i)+g.strscroll.Y, width/2-204, 20), fmt.Sprintf("%s (%d)", ent.Squad.Name, ent.Amount))
	//	}
	//	rl.EndScissorMode()
	//
	// case 1:
	//
	//	gui.ScrollPanel(rl.NewRectangle(width/2, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.Shop.Vehicles)*20+10)), &g.strscroll, &g.strview)
	//	rl.BeginScissorMode(int32(g.strview.X), int32(g.strview.Y), int32(g.strview.Width), int32(g.strview.Height))
	//	for i, ent := range g.Shop.Vehicles {
	//		if g.Money < int64(ent.Vehicle.Cost) {
	//			gui.SetState(gui.STATE_DISABLED)
	//		}
	//		strbuy := gui.Button(rl.NewRectangle(width/2+11, float32(48+20*i)+g.strscroll.Y, 80, 20), fmt.Sprintf("Buy %d", ent.Vehicle.Cost))
	//		if strbuy {
	//			veh := ent.Vehicle
	//			for i, _ := range veh.Crew {
	//				inf := &data.Infantry{}
	//				*inf = *veh.Crew[i]
	//				inf.Nid = fmt.Sprintf("%02d %02d", rand.Int32N(100), rand.Int32N(100))
	//				items := inf.Inv.Items
	//				inf.Inv = &data.Inventory{
	//					Name:  inf.Inv.Name,
	//					Items: []data.Item{},
	//				}
	//				for _, item := range items {
	//					if item.Amount == 0 {
	//						inf.Inv.Items = append(inf.Inv.Items, item)
	//					} else {
	//						item.Amount = float64(rng.RollOnce(item.Amount))
	//						if item.Amount > 0 {
	//							inf.Inv.Items = append(inf.Inv.Items, item)
	//						}
	//					}
	//				}
	//				veh.Crew[i] = inf
	//			}
	//			g.State.Campaign.Vehicles = append(g.State.Campaign.Vehicles, &veh)
	//			g.Money -= int64(ent.Vehicle.Cost)
	//			g.Shop.Vehicles[i].Amount -= 1
	//			if g.Shop.Vehicles[i].Amount < 1 {
	//				g.Shop.Vehicles = slices.Delete(g.Shop.Vehicles, i, i+1)
	//			}
	//		}
	//		gui.SetState(gui.STATE_NORMAL)
	//		preview = gui.Button(rl.NewRectangle(width/2+91, float32(48+20*i)+g.strscroll.Y, 80, 20), "Preview")
	//		if preview {
	//			g.examine = append(g.examine, ent)
	//		}
	//		gui.Label(rl.NewRectangle(width/2+181, float32(48+20*i)+g.strscroll.Y, width/2-204, 20), fmt.Sprintf("%s (%d)", ent.Vehicle.Name, ent.Amount))
	//	}
	//	rl.EndScissorMode()
	//
	// case 2:
	//
	//	gui.ScrollPanel(rl.NewRectangle(width/2, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.Shop.Infantry)*20+10)), &g.strscroll, &g.strview)
	//	rl.BeginScissorMode(int32(g.strview.X), int32(g.strview.Y), int32(g.strview.Width), int32(g.strview.Height))
	//	for i, ent := range g.Shop.Infantry {
	//		if g.Money < int64(ent.Infantry.Cost) {
	//			gui.SetState(gui.STATE_DISABLED)
	//		}
	//		strbuy := gui.Button(rl.NewRectangle(width/2+11, float32(48+20*i)+g.strscroll.Y, 80, 20), fmt.Sprintf("Buy %d", ent.Infantry.Cost))
	//		if strbuy {
	//			inf := &data.Infantry{}
	//			*inf = *ent.Infantry
	//			inf.Nid = fmt.Sprintf("%02d %02d", rand.Int32N(100), rand.Int32N(100))
	//			items := inf.Inv.Items
	//			inf.Inv = &data.Inventory{
	//				Name:  inf.Inv.Name,
	//				Items: []data.Item{},
	//			}
	//			for _, item := range items {
	//				if item.Amount == 0 {
	//					inf.Inv.Items = append(inf.Inv.Items, item)
	//				} else {
	//					item.Amount = float64(rng.RollOnce(item.Amount))
	//					if item.Amount > 0 {
	//						inf.Inv.Items = append(inf.Inv.Items, item)
	//					}
	//				}
	//			}
	//			g.State.Campaign.Infantry = append(g.State.Campaign.Infantry, inf)
	//			g.Money -= int64(ent.Infantry.Cost)
	//			g.Shop.Infantry[i].Amount -= 1
	//			if g.Shop.Infantry[i].Amount < 1 {
	//				g.Shop.Infantry = slices.Delete(g.Shop.Infantry, i, i+1)
	//			}
	//		}
	//		gui.SetState(gui.STATE_NORMAL)
	//		preview = gui.Button(rl.NewRectangle(width/2+91, float32(48+20*i)+g.strscroll.Y, 80, 20), "Preview")
	//		if preview {
	//			g.examine = append(g.examine, ent)
	//		}
	//		gui.Label(rl.NewRectangle(width/2+181, float32(48+20*i)+g.strscroll.Y, width/2-204, 20), fmt.Sprintf("%s (%d)", ent.Infantry.Name, ent.Amount))
	//	}
	//	rl.EndScissorMode()
	//
	// case 3:
	//
	//		gui.ScrollPanel(rl.NewRectangle(width/2, 43, width/2, height-43), "", rl.NewRectangle(0, 0, width/2-14, float32(len(g.Shop.Items)*20+10)), &g.strscroll, &g.strview)
	//		rl.BeginScissorMode(int32(g.strview.X), int32(g.strview.Y), int32(g.strview.Width), int32(g.strview.Height))
	//		for i, ent := range g.Shop.Items {
	//			if g.Money < int64(ent.Item.Cost) {
	//				gui.SetState(gui.STATE_DISABLED)
	//			}
	//			strbuy := gui.Button(rl.NewRectangle(width/2+11, float32(48+20*i)+g.strscroll.Y, 80, 20), fmt.Sprintf("Buy %d", ent.Item.Cost))
	//			if strbuy {
	//				g.State.Campaign.Items = append(g.State.Campaign.Items, &ent.Item)
	//				g.Money -= int64(ent.Item.Cost)
	//			}
	//			gui.SetState(gui.STATE_NORMAL)
	//			preview = gui.Button(rl.NewRectangle(width/2+91, float32(48+20*i)+g.strscroll.Y, 80, 20), "Preview")
	//			if preview {
	//				g.examine = append(g.examine, ent)
	//			}
	//			gui.Label(rl.NewRectangle(width/2+181, float32(48+20*i)+g.strscroll.Y, width/2-204, 20), ent.Item.Name)
	//		}
	//		rl.EndScissorMode()
	//	}
	//
	//	if len(g.examine) > 0 {
	//		g.examine[len(g.examine)-1].Display(true, func() {
	//			g.examine = g.examine[:len(g.examine)-1]
	//		})
	//	}
	//
	//	if g.paused {
	//		gui.Panel(rl.NewRectangle((width-300)/2, (height-200)/2, 300, 200), "Menu")
	//		cnbtn := gui.Button(rl.NewRectangle((width-300)/2+278, (height-200)/2+2, 20, 20), "X")
	//		if cnbtn {
	//			g.paused = false
	//		}
	//		svbtn := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+35, 250, 30), "Save")
	//		if svbtn {
	//			g.SaveGame()
	//			g.paused = false
	//		}
	//		exbtn1 := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+75, 250, 30), "Exit to Main Menu")
	//		if exbtn1 {
	//			menu := MainMenu{}
	//			err := menu.Init()
	//			if err != nil {
	//				log.Printf("%+v", err)
	//				return
	//			}
	//			app.CurApp.SetStage(&menu)
	//		}
	//		exbtn2 := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+115, 250, 30), "Exit to Desktop")
	//		if exbtn2 {
	//			app.CurApp.Exit()
	//		}
	//		clbtn := gui.Button(rl.NewRectangle((width-300)/2+25, (height-200)/2+155, 250, 30), "Cancel")
	//		if clbtn {
	//			g.paused = false
	//		}
	//	}
	//
	//	if g.preview {
	//		gui.Panel(rl.NewRectangle((width-400)/2, (height-200)/2, 400, 200), "Next Battle")
	//		cnbtn := gui.Button(rl.NewRectangle((width-400)/2+378, (height-200)/2+2, 20, 20), "X")
	//		if cnbtn {
	//			g.preview = false
	//		}
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+30, 350, 20), fmt.Sprintf("Enemy: %s", g.State.Status.EnemyArmy))
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+50, 350, 20), fmt.Sprintf("Landscape: %s", g.State.Status.Landscape))
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+70, 350, 20), fmt.Sprintf("Region: %s", g.State.Status.Region))
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+90, 350, 20), fmt.Sprintf("Map: %s", g.State.Status.Map))
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+110, 350, 20), fmt.Sprintf("Gamemode: %s", g.State.Status.Gamemode))
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+130, 350, 20), fmt.Sprintf("Risk: %s", g.State.Status.Risk))
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+150, 350, 20), fmt.Sprintf("Level: %d", g.State.Status.PlayedGames))
	//		gui.Label(rl.NewRectangle((width-400)/2+25, (height-200)/2+170, 350, 20), fmt.Sprintf("Won: %d", g.State.Status.WonGames))
	//	}
	//
	// return
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

func (g *Game) Finished() (bool, error) {
	info, err := os.Stat(g.Profile + "/" + strcase.ToLowerCamel(g.State.Status.Name) + ".sav")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return info.ModTime().Unix() > g.State.Status.Timestamp, nil
}

func (g *Game) SaveGame() error {
	if g == nil {
		return fmt.Errorf("Game is nil")
	}
	// Marshall the settings
	data, err := json.Marshal(*g)
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
	return os.WriteFile(exPath+"/"+strcase.ToLowerCamel(g.State.Status.Name)+".json", data, 0644)
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
