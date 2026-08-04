package stage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"roguefront/pkg/app"
	"roguefront/pkg/data"

	"github.com/AllenDang/giu"
	"github.com/iancoleman/strcase"
	"github.com/sqweek/dialog"
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

	examine []examinable `json:"-"`
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

	g.examine = []examinable{}

	return nil
}

// render hook
func (g *Game) Build() {
	if g == nil {
		return
	}

	height := float32(app.CurApp.GetHeight())
	width := float32(app.CurApp.GetWidth())

	giu.Column(
		giu.Row(
			giu.AlignManually(giu.AlignLeft, giu.Labelf("Money: %d", g.Money), 100, true),
			giu.Dummy(width/2-200, 0),
			giu.Button("Start Battle").Size(100, 20).OnClick(func() {
				err := g.State.Write(g.Profile)
				if err != nil {
					log.Print(err)
				}
			}),
			giu.Button("Read Result").Disabled(func() bool {
				fin, err := g.Finished()
				if err != nil {
					log.Print(err)
				}
				return !fin
			}()).Size(100, 20).OnClick(func() {
				err := g.State.Read(g.Profile + "/" + g.State.Status.Name + ".sav")
				if err != nil {
					log.Print(err)
				}
				g.ReRoll()
			}),
			giu.Dummy(width/2-330, 0),
			giu.Button("Preview Battle").Size(125, 20).OnClick(func() { giu.OpenPopup("Preview") }),
			giu.PopupModal("Preview").Layout(
				giu.Column(
					giu.Labelf("Enemy: %s", g.State.Status.EnemyArmy),
					giu.Labelf("Landscape: %s", g.State.Status.Landscape),
					giu.Labelf("Region: %s", g.State.Status.Region),
					giu.Labelf("Map: %s", g.State.Status.Map),
					giu.Labelf("Gamemode: %s", g.State.Status.Gamemode),
					giu.Labelf("Risk: %s", g.State.Status.Risk),
					giu.Labelf("Level: %d", g.State.Status.PlayedGames),
					giu.Labelf("Won: %d", g.State.Status.WonGames),
					giu.Button("Close").Size(100, 20).OnClick(func() { giu.CloseCurrentPopup() }),
				),
			),
			giu.Button("Menu").Size(75, 20).OnClick(func() { giu.OpenPopup("Menu") }),
			giu.PopupModal("Menu").Layout(
				giu.Column(
					giu.Button("Save").Size(200, 20).OnClick(func() {
						g.SaveGame()
						giu.CloseCurrentPopup()
					}),
					giu.Button("Load").Size(200, 20).OnClick(func() {
						file, err := dialog.File().
							Filter("Json Files", "json").
							Title("Load Previous Roguefront Campaign").Load()
						if err == nil {
							g.LoadGame(file)
						}
					}),
					giu.Button("Exit to Main Menu").Size(200, 20).OnClick(func() {
						menu := MainMenu{}
						err := menu.Init()
						if err != nil {
							log.Printf("%+v", err)
							return
						}
						app.CurApp.SetStage(&menu)
					}),
					giu.Button("Exit to Desktop").Size(200, 20).OnClick(func() { app.CurApp.Exit() }),
					giu.Button("Cancel").Size(200, 20).OnClick(func() { giu.CloseCurrentPopup() }),
				),
			),
		),
		giu.Row(
			giu.Child().Size(width/2-10, height-40).Layout(
				giu.TabBar().ID("Owned").TabItems(
					giu.TabItem("Squads").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for _, elem := range g.State.Campaign.Squads {
								rows = append(rows, giu.Row(
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
					giu.TabItem("Vehicles").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for _, elem := range g.State.Campaign.Vehicles {
								rows = append(rows, giu.Row(
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
					giu.TabItem("Infantry").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for _, elem := range g.State.Campaign.Infantry {
								rows = append(rows, giu.Row(
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
					giu.TabItem("Items").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for _, elem := range g.State.Campaign.Items {
								rows = append(rows, giu.Row(
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
				),
			),
			giu.Child().Size(width/2-10, height-40).Layout(
				giu.TabBar().ID("Store").TabItems(
					giu.TabItem("Squads").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for i, elem := range g.Shop.Squads {
								rows = append(rows, giu.Row(
									giu.Buttonf("Buy: %d", elem.Squad.Cost).Size(100, 20).OnClick(func() {
										g.BuySquad(i, elem)
									}).Disabled(g.Money < int64(elem.Squad.Cost)),
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Squad.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
					giu.TabItem("Vehicles").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for i, elem := range g.Shop.Vehicles {
								rows = append(rows, giu.Row(
									giu.Buttonf("Buy: %d", elem.Vehicle.Cost).Size(100, 20).OnClick(func() {
										g.BuyVehicle(i, elem)
									}).Disabled(g.Money < int64(elem.Vehicle.Cost)),
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Vehicle.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
					giu.TabItem("Infantry").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for i, elem := range g.Shop.Infantry {
								rows = append(rows, giu.Row(
									giu.Buttonf("Buy: %d", elem.Infantry.Cost).Size(100, 20).OnClick(func() {
										g.BuyInfantry(i, elem)
									}).Disabled(g.Money < int64(elem.Infantry.Cost)),
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Infantry.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
					giu.TabItem("Items").Layout(
						func() giu.Widget {
							rows := []giu.Widget{}
							for _, elem := range g.Shop.Items {
								rows = append(rows, giu.Row(
									giu.Buttonf("Buy: %d", elem.Item.Cost).Size(100, 20).OnClick(func() {
										g.BuyItem(elem)
									}).Disabled(g.Money < int64(elem.Item.Cost)),
									giu.Buttonf("Preview").Size(100, 20).OnClick(func() {
										g.examine = append(g.examine, elem)
									}),
									giu.Labelf("%s", elem.Item.Name),
								))
							}
							return giu.Column(rows...)
						}(),
					),
				),
			),
		),
		func() giu.Widget {
			if len(g.examine) > 0 {
				giu.OpenPopup("Examine")
				return giu.PopupModal("Examine").Layout(
					g.examine[len(g.examine)-1].Display(true, func() {
						g.examine = g.examine[:len(g.examine)-1]
						giu.CloseCurrentPopup()
					}),
				)
			}
			return nil
		}(),
	).Build()
}

// handle resize event
func (g *Game) OnResize(w int, h int) {
	if g == nil {
		return
	}
}

// handle update cycle
func (g *Game) Update(dt int64) {
	if g == nil {
		return
	}
}

// handle player input
func (g *Game) OnInput(dt int64) {
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
