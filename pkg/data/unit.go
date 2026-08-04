package data

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Infantry struct {
	Name string     `json:"name"`
	Side string     `json:"side"`
	Tier int        `json:"tier"`
	Era  string     `jons:"era"`
	Nid  string     `json:"nid"`
	Id   int        `json:"-"`
	Cost int        `json:"cost"`
	Inv  *Inventory `json:"inv,omitempty"`
}

type Squad struct {
	Name     string      `json:"name"`
	Side     string      `json:"side"`
	Era      string      `json:"era"`
	Cost     int         `json:"cost"`
	Soldiers []*Infantry `json:"soldiers,omitempty"`
}

type Vehicle struct {
	Name string      `json:"name"`
	Side string      `json:"side"`
	Era  string      `json:"era"`
	Cost int         `json:"cost"`
	Crew []*Infantry `json:"crew,omitempty"`
	Inv  *Inventory  `json:"inv,omitempty"`
}

func FindUnits(game string, mods []Mod) (infantry map[string]*Infantry, squads map[string]map[string]Squad, vehicles map[string]Vehicle) {
	inventories := FindInventories(game, mods)
	files := SearchContent(game, mods, "/resource/set/multiplayer/units/conquest/", ".set")

	infantry = make(map[string]*Infantry)
	squads = make(map[string]map[string]Squad)
	vehicles = make(map[string]Vehicle)

	for _, file := range files {
		for _, unitset := range file.Body {
			switch {
			case unitset.Name == "define":
				continue
			case strings.HasPrefix(unitset.Name, "mp/"):
				if len(unitset.Body) < 1 {
					continue
				}

				props := unitset.Body[0]
				tier := 2
				if strings.Contains(props.Name, "tier") {
					n, err := strconv.ParseInt(strings.Split(props.Name, "tier")[1], 10, 64)
					if err == nil {
						tier = int(n)
					}
				}
				cost := tier * 50
				side := ""
				for _, arg := range props.Args {
					if strings.Contains(arg, "side") {
						side = arg[5 : len(arg)-1]
					}
					if strings.Contains(arg, "cost") {
						n, err := strconv.ParseFloat(arg[5:len(arg)-1], 64)
						if err == nil {
							cost = int(n)
						}
					}
				}
				era := "late"
				spl := strings.Split(unitset.Name, "/")
				if len(spl) > 2 {
					era = spl[len(spl)-2]
				}
				var inv *Inventory
				if i, ok := inventories[unitset.Name]; ok {
					inv = &i
				}
				inf := &Infantry{
					Name: unitset.Name,
					Side: side,
					Tier: tier,
					Era:  era,
					Nid:  "",
					Cost: cost,
					Inv:  inv,
				}
				infantry[inf.Name] = inf
			case strings.HasPrefix(unitset.Name, "squad"):
				side := ""
				name := ""
				era := ""
				soldiers := []*Infantry{}

				for _, arg := range unitset.Args {
					if len(arg) < 1 {
						continue
					}
					switch {
					case strings.HasPrefix(arg, "side"):
						side = arg[5 : len(arg)-1]
					case strings.HasPrefix(arg, "name"):
						name = arg[5 : len(arg)-1]
					case strings.HasPrefix(arg, "period"):
						era = arg[7 : len(arg)-1]
					case strings.HasPrefix(arg, "c"):
						ctext := strings.Split(arg[:len(arg)-1], "(")
						if len(ctext) > 1 {
							centry := strings.Split(ctext[1], ":")
							if len(centry) > 1 {
								soldier := centry[0]
								num, err := strconv.ParseInt(centry[1], 10, 64)
								if err == nil {
									for i := int64(0); i < num; i++ {
										soldiers = append(soldiers, &Infantry{Name: soldier})
									}
								}
							}
						}
					}
				}
				squad := Squad{
					Side:     side,
					Name:     name,
					Era:      era,
					Soldiers: soldiers,
					Cost:     0,
				}
				if _, ok := squads[squad.Side]; !ok {
					squads[squad.Side] = make(map[string]Squad)
				}
				squads[squad.Side][squad.Name] = squad
			default: //vehicle (includes emplacements)
				if len(unitset.Body) < 1 {
					continue
				}
				cost := 0
				side := ""
				crew := []*Infantry{}
				era := ""
				for _, arg := range unitset.Body[0].Args {
					if len(arg) < 1 {
						continue
					}
					switch {
					case strings.HasPrefix(arg, "side"):
						side = arg[5 : len(arg)-1]
					case strings.HasPrefix(arg, "period"):
						era = arg[7 : len(arg)-1]
					case strings.HasPrefix(arg, "crew"):
						ctext := strings.Split(arg[:len(arg)-1], "(")
						if len(ctext) > 1 {
							centry := strings.Split(ctext[1], ":")
							if len(centry) > 1 {
								soldier := centry[0]
								num, err := strconv.ParseInt(centry[1], 10, 64)
								if err == nil {
									for i := int64(0); i < num; i++ {
										crew = append(crew, &Infantry{Name: soldier})
									}
								}
							}
						}
					}
				}
				if unitset.Body[1].Name == "cost" && len(unitset.Body[1].Args) > 0 {
					n, err := strconv.ParseFloat(unitset.Body[1].Args[0], 64)
					if err == nil {
						cost = int(n)
					}
				}
				if cost < 1 {
					continue
				}
				vehicle := Vehicle{
					Name: unitset.Name,
					Cost: cost,
					Side: side,
					Crew: crew,
					Era:  era,
					Inv:  &Inventory{},
				}
				vehicles[vehicle.Name] = vehicle
			}
		}
	}

	//fixing up squads with acutal costs and actual infantry
	for _, side := range squads {
		for _, squad := range side {
			cost := 0
			for i, ent := range squad.Soldiers {
				inf := infantry[fmt.Sprintf("mp/%s/%s/%s", squad.Side, squad.Era, ent.Name)]
				squad.Soldiers[i] = inf
				if inf != nil {
					cost += inf.Cost
				} else {
					log.Printf("Failed to find conquest infantry: mp/%s/%s/%s\n", squad.Side, squad.Era, ent.Name)
				}
			}
			squad.Cost = int(float64(cost) * 1.2) // extra for convenience
			squads[squad.Side][squad.Name] = squad
		}
	}

	//fixing up vehicle crews
	for _, vehicle := range vehicles {
		for i, ent := range vehicle.Crew {
			inf := infantry[fmt.Sprintf("mp/%s/%s/%s", vehicle.Side, vehicle.Era, ent.Name)]
			vehicle.Crew[i] = inf
			if inf == nil {
				log.Printf("Failed to find conquest infantry: mp/%s/%s/%s\n", vehicle.Side, vehicle.Era, ent.Name)
			}
		}
	}

	return infantry, squads, vehicles
}
