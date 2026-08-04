package stage

import (
	"fmt"
	"log"
	"maps"
	rand "math/rand/v2"
	"slices"

	"roguefront/pkg/data"
	"roguefront/pkg/rng"

	"github.com/AllenDang/giu"
)

type Store struct {
	Squads   []StoreSquad    `json:"sqd"`
	Infantry []StoreInfantry `json:"inf"`
	Vehicles []StoreVehicle  `json:"veh"`
	Items    []StoreItem     `json:"itm"`
}

type StoreSquad struct {
	Squad  data.Squad `json:"sqd"`
	Amount int        `json:"num"`
}

type StoreInfantry struct {
	Infantry *data.Infantry `json:"inf"`
	Amount   int            `json:"num"`
}

type StoreVehicle struct {
	Vehicle data.Vehicle `json:"veh"`
	Amount  int          `json:"num"`
}

type StoreItem struct {
	Item   data.Item `json:"itm"`
	Amount int       `json:"num"`
}

func (ent StoreInfantry) Display(editable bool, exit func()) giu.Widget {
	return ent.Infantry.Display(false, exit)
}
func (ent StoreVehicle) Display(editable bool, exit func()) giu.Widget {
	return ent.Vehicle.Display(false, exit)
}
func (ent StoreSquad) Display(editable bool, exit func()) giu.Widget {
	return ent.Squad.Display(false, exit)
}
func (ent StoreItem) Display(editable bool, exit func()) giu.Widget {
	return ent.Item.Display(false, exit)
}

type examinable interface {
	Display(editable bool, exit func()) giu.Widget
}

func (g *Game) Populate() error {
	g.PopulateMaps()
	g.PopulateUnits()
	g.PopulateItems()

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
	g.Maps = data.FindMaps(g.Game, g.Mods)
}

func (g *Game) PopulateUnits() {
	g.Infantry, g.Squads, g.Vehicles = data.FindUnits(g.Game, g.Mods)
}

func (g *Game) PopulateItems() {
	g.Items = data.FindItems(g.Game, g.Mods)
}

func (g *Game) ReRoll() {
	g.RollLandscape()
	g.RollRegion()
	g.RollMap()
	g.RollEnemy()
	g.RollRisk()
	g.RollStore()
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

func (g *Game) RollRisk() {
	risks := []string{
		"low",
		"standard",
		"high",
	}
	g.State.Status.Risk = risks[rand.IntN(len(risks))]
}

func (g *Game) RollStore() {
	g.Shop.Squads = []StoreSquad{}
	g.Shop.Infantry = []StoreInfantry{}
	g.Shop.Vehicles = []StoreVehicle{}
	g.Shop.Items = []StoreItem{}

	progress := g.State.Status.PlayedGames / 10
	era := 1 + min(progress, 2)
	eras := map[string]int{
		"early": 0,
		"mid":   1,
		"late":  2,
	}

	rolls := int(rand.Int32N(16)) + 15 + progress*5
	for i := 0; i < rolls; i++ {
		category := rand.Int32N(100)
		switch {
		case category < 20: // squad roll
			squads := slices.Collect(maps.Keys(g.Squads[g.State.Status.Army]))
			for {
				squad := g.Squads[g.State.Status.Army][squads[rand.Int32N(int32(len(squads)))]]
				if eras[squad.Era] <= era {
					g.Shop.Squads = append(g.Shop.Squads, StoreSquad{
						Squad:  squad,
						Amount: int(rand.Int32N(3) + 1),
					})
					break
				}
			}
		case category >= 20 && category < 50: // infantry roll
			infantry := slices.Collect(maps.Keys(g.Infantry))
			for {
				inf := g.Infantry[infantry[rand.Int32N(int32(len(infantry)))]]
				if inf.Side == g.State.Status.Army && eras[inf.Era] <= era {
					g.Shop.Infantry = append(g.Shop.Infantry, StoreInfantry{
						Infantry: inf,
						Amount:   int(rand.Int32N(3) + 1),
					})
					break
				}
			}
		case category >= 50 && category < 60: // vehicle roll
			vehicles := slices.Collect(maps.Keys(g.Vehicles))
			for {
				veh := g.Vehicles[vehicles[rand.Int32N(int32(len(vehicles)))]]
				if veh.Side == g.State.Status.Army && eras[veh.Era] <= era {
					g.Shop.Vehicles = append(g.Shop.Vehicles, StoreVehicle{
						Vehicle: veh,
						Amount:  1,
					})
					break
				}
			}
		case category >= 60 && category < 100: // item roll
			itemcats := slices.Collect(maps.Keys(g.Items))
			items := g.Items[itemcats[rand.Int32N(int32(len(itemcats)))]]
			item := items[rand.Int32N(int32(len(items)))]
			g.Shop.Items = append(g.Shop.Items, StoreItem{
				Item:   item,
				Amount: 0,
			})
		}
	}
}

func (g *Game) BuyVehicle(i int, ent StoreVehicle) {
	if g.Money < int64(ent.Vehicle.Cost) {
		return
	}
	veh := ent.Vehicle
	for i, _ := range veh.Crew {
		inf := &data.Infantry{}
		*inf = *veh.Crew[i]
		inf.Nid = fmt.Sprintf("%02d %02d", rand.Int32N(100), rand.Int32N(100))
		items := inf.Inv.Items
		inf.Inv = &data.Inventory{
			Name:  inf.Inv.Name,
			Items: []data.Item{},
		}
		for _, item := range items {
			if item.Amount == 0 {
				inf.Inv.Items = append(inf.Inv.Items, item)
			} else {
				item.Amount = float64(rng.RollOnce(item.Amount))
				if item.Amount > 0 {
					inf.Inv.Items = append(inf.Inv.Items, item)
				}
			}
		}
		veh.Crew[i] = inf
	}
	g.State.Campaign.Vehicles = append(g.State.Campaign.Vehicles, &veh)
	g.Money -= int64(ent.Vehicle.Cost)
	g.Shop.Vehicles[i].Amount -= 1
	if g.Shop.Vehicles[i].Amount < 1 {
		g.Shop.Vehicles = slices.Delete(g.Shop.Vehicles, i, i+1)
	}
}

func (g *Game) BuySquad(i int, ent StoreSquad) {
	if g.Money < int64(ent.Squad.Cost) {
		return
	}
	sqd := ent.Squad
	for i, _ := range sqd.Soldiers {
		inf := &data.Infantry{}
		*inf = *sqd.Soldiers[i]
		inf.Nid = fmt.Sprintf("%02d %02d", rand.Int32N(10000), rand.Int32N(10000))
		items := inf.Inv.Items
		inf.Inv = &data.Inventory{
			Name:  inf.Inv.Name,
			Items: []data.Item{},
		}
		for _, item := range items {
			if item.Amount == 0 {
				inf.Inv.Items = append(inf.Inv.Items, item)
			} else {
				item.Amount = float64(rng.RollOnce(item.Amount))
				if item.Amount > 0 {
					inf.Inv.Items = append(inf.Inv.Items, item)
				}
			}
		}
		sqd.Soldiers[i] = inf
	}
	g.State.Campaign.Squads = append(g.State.Campaign.Squads, &sqd)
	g.Money -= int64(ent.Squad.Cost)
	g.Shop.Squads[i].Amount -= 1
	if g.Shop.Squads[i].Amount < 1 {
		g.Shop.Squads = slices.Delete(g.Shop.Squads, i, i+1)
	}
}

func (g *Game) BuyInfantry(i int, ent StoreInfantry) {
	if g.Money < int64(ent.Infantry.Cost) {
		return
	}
	inf := &data.Infantry{}
	*inf = *ent.Infantry
	inf.Nid = fmt.Sprintf("%02d %02d", rand.Int32N(100), rand.Int32N(100))
	items := inf.Inv.Items
	inf.Inv = &data.Inventory{
		Name:  inf.Inv.Name,
		Items: []data.Item{},
	}
	for _, item := range items {
		if item.Amount == 0 {
			inf.Inv.Items = append(inf.Inv.Items, item)
		} else {
			item.Amount = float64(rng.RollOnce(item.Amount))
			if item.Amount > 0 {
				inf.Inv.Items = append(inf.Inv.Items, item)
			}
		}
	}
	g.State.Campaign.Infantry = append(g.State.Campaign.Infantry, inf)
	g.Money -= int64(ent.Infantry.Cost)
	g.Shop.Infantry[i].Amount -= 1
	if g.Shop.Infantry[i].Amount < 1 {
		g.Shop.Infantry = slices.Delete(g.Shop.Infantry, i, i+1)
	}
}

func (g *Game) BuyItem(ent StoreItem) {
	if g.Money < int64(ent.Item.Cost) {
		return
	}
	g.State.Campaign.Items = append(g.State.Campaign.Items, &ent.Item)
	g.Money -= int64(ent.Item.Cost)
}
