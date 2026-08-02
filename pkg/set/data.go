package set

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"roguefront/res"
)

var ()

type Infantry struct {
	Name string     `json:"name"`
	Side string     `json:"side"`
	Tier int        `json:"tier"`
	Era  string     `jons:"era"`
	Nid  string     `json:"name"`
	Cost int        `json:"cost"`
	Inv  *Inventory `json:"inv,omitempty"`
}

type Squad struct {
	Name     string     `json:"name"`
	Side     string     `json:"side"`
	Era      string     `json:"era"`
	Cost     int        `json:"cost"`
	Soldiers []Infantry `json:"soldiers"`
}

type Vehicle struct {
	Name string     `json:"name"`
	Side string     `json:"side"`
	Era  string     `json:"era"`
	Cost int        `json:"cost"`
	Crew []Infantry `json:"crew"`
}

type Inventory struct {
	Name  string     `json:"name"`
	Items []res.Item `json:"items"`
}

func SearchPak(path string, prefix string, postfix string) (sets []*Set) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Printf("Error on %s: %+v\n", path, err)
		return nil
	}
	defer r.Close()

	for _, f := range r.File {
		split := strings.Split(f.Name, "/")
		if strings.HasPrefix(f.Name, prefix) && strings.HasSuffix(f.Name, postfix) && !strings.HasPrefix(split[len(split)-1], ".") {
			fr, err := f.Open()
			if err != nil {
				log.Printf("Error on %s: %+v\n", path+"/"+f.Name, err)
			}

			data, err := io.ReadAll(fr)
			fr.Close()
			if err != nil {
				log.Printf("Error on %s: %+v\n", path+"/"+f.Name, err)
			}
			s, err := Parse(path+"/"+f.Name, data)
			if err != nil {
				log.Printf("Error on %s: %+v\n", path+"/"+f.Name, err)
			} else {
				sets = append(sets, s)
			}
		}
	}
	return sets
}

func FindPaks(path string, prefix string, postfix string) (sets []*Set) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Error on %s: %+v\n", path, err)
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pak") {
			sets = append(sets, SearchPak(path+"/"+entry.Name(), prefix, postfix)...)
		}
	}

	return sets
}

func FindLoose(path string, prefix string, postfix string) (sets []*Set) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), postfix) && !strings.HasPrefix(entry.Name(), ".") {
			data, err := os.ReadFile(path + "/" + entry.Name())
			if err != nil {
				log.Printf("Error on %s: %+v\n", path+"/"+entry.Name(), err)
			}
			s, err := Parse(path+"/"+entry.Name(), data)
			if err != nil {
				log.Printf("Error on %s: %+v\n", path+"/"+entry.Name(), err)
			} else {
				sets = append(sets, s)
			}
		}
		if entry.IsDir() && prefix == "" && !strings.HasPrefix(entry.Name(), ".") {
			sets = append(sets, FindLoose(path+"/"+entry.Name(), prefix, postfix)...)
		}
	}

	return sets
}

func SearchContent(game string, mods []res.Mod, search string, ft string) (sets []*Set) {
	split := strings.Split(search, "/")
	end := split[len(split)-1]
	root := strings.Join(split[0:2], "/")
	prefix := strings.Join(split[0:len(split)-1], "/")
	postfix := strings.Join(split[2:], "/")

	sets = append(FindPaks(game+root, postfix, end+ft), FindLoose(game+prefix, end, ft)...)
	for _, mod := range mods {
		if mod.Enabled {
			sets = append(sets, FindPaks(mod.Path+root, postfix, end+ft)...)
			sets = append(sets, FindLoose(mod.Path+prefix, end, ft)...)
		}
	}

	return sets
}

func FindNations(game string, mods []res.Mod) (nations []string) {
	files := SearchContent(game, mods, "/resource/set/dynamic_campaign/values", ".set")

	for _, sets := range files {
		for _, values := range sets.Body {
			if values.Name != "Regions" {
				continue
			}
			for _, region := range values.Body {
				matches := region.FindAll("AvailableMatchups")
				if matches == nil {
					continue
				}

				for _, match := range matches {
					for _, matchup := range match.Args {
						nats := strings.Split(matchup, " ")
						if len(nats) == 2 {
							if !slices.Contains(nations, nats[0]) {
								nations = append(nations, nats[0])
							}
							if !slices.Contains(nations, nats[1]) {
								nations = append(nations, nats[1])
							}
						}
					}
				}
			}
		}
	}
	return nations
}

func FindLevels(game string, mods []res.Mod) (levels []string) {
	files := SearchContent(game, mods, "/resource/set/dynamic_campaign", ".set")

	for _, sets := range files {
		for _, lvl := range sets.Body {
			if lvl.Name != "ResupplyCost" && lvl.Name != "GameModes" && !slices.Contains(levels, lvl.Name) {
				levels = append(levels, lvl.Name)
			}
		}
	}
	return levels
}

func FindResources(game string, mods []res.Mod) (resources []string) {
	files := SearchContent(game, mods, "/resource/set/dynamic_campaign", ".set")

	for _, sets := range files {
		for _, ressets := range sets.Body {
			resset := ressets.Find("Resources")
			if resset == nil {
				continue
			}

			for _, res := range resset.Body {
				if !slices.Contains(resources, res.Name) {
					resources = append(resources, res.Name)
				}
			}
		}
	}

	return resources
}

func FindMaps(game string, mods []res.Mod) (maps map[string][]string) {
	files := SearchContent(game, mods, "/resource/set/dynamic_campaign/map_points", ".set")
	maps = make(map[string][]string)

	for _, file := range files {
		for _, regsets := range file.Body {
			if regsets.Name != "Regions" {
				continue
			}
			for _, regset := range regsets.Body {
				region := regset.Find("Maps")
				if region == nil {
					continue
				}
				for _, mapset := range region.Body {
					if !slices.Contains(maps[regset.Name], mapset.Name) {
						maps[regset.Name] = append(maps[regset.Name], mapset.Name)
					}
				}
			}
		}
	}

	return maps
}

func FindUnits(game string, mods []res.Mod) (infantry map[string]Infantry, squads map[string]map[string]Squad, vehicles map[string]Vehicle) {
	inventories := FindInventories(game, mods)
	files := SearchContent(game, mods, "/resource/set/multiplayer/units/conquest/", ".set")

	infantry = make(map[string]Infantry)
	squads = make(map[string]map[string]Squad)
	vehicles = make(map[string]Vehicle)

	for _, unitsets := range files {
		for _, unitset := range unitsets.Body {
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
				inf := Infantry{
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
				soldiers := []Infantry{}

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
										soldiers = append(soldiers, Infantry{Name: soldier})
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
				crew := []Infantry{}
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
										crew = append(crew, Infantry{Name: soldier})
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
				soldier := infantry[fmt.Sprintf("mp/%s/%s/%s", squad.Side, squad.Era, ent.Name)]
				squad.Soldiers[i] = soldier
				cost += soldier.Cost
			}
			squad.Cost = int(float64(cost) * 1.2) // extra for convenience
			squads[squad.Side][squad.Name] = squad
		}
	}

	//fixing up vehicle crews
	for _, vehicle := range vehicles {
		for i, ent := range vehicle.Crew {
			soldier := infantry[fmt.Sprintf("mp/%s/%s/%s", vehicle.Side, vehicle.Era, ent.Name)]
			vehicle.Crew[i] = soldier
		}
	}

	return infantry, squads, vehicles
}

func FindInventories(game string, mods []res.Mod) (inventories map[string]Inventory) {
	sets := SearchContent(game, mods, "/resource/set/breed/mp/", ".set")
	inventories = make(map[string]Inventory)

	for _, invset := range sets {
		name := strings.Join(strings.Split(invset.Name, "breed/")[1:], "breed/")
		name = name[:len(name)-4]
		inv := Inventory{
			Name:  name,
			Items: []res.Item{},
		}
		if len(invset.Body) > 0 {
			if len(invset.Body[0].Body) > 0 {
				props := invset.Body[0].Body
				for _, prop := range props {
					switch prop.Name {
					case "armors":
						for _, def := range prop.Body {
							item := res.Item{
								Name:   def.Args[0],
								Amount: 0,
								Equip:  fmt.Sprintf("{user \"%s\"}", def.Name),
							}
							inv.Items = append(inv.Items, item)
						}
					case "inventory":
						for _, def := range prop.Body {
							if def.Name == "item" && len(def.Args) > 0 {
								n := float64(0)
								extra := ""
								if len(def.Args) > 1 {
									if def.Args[1] != "filling" && def.Args[1] != "filled" {
										if len(def.Args) > 2 && def.Args[2] != "filled" {
											min, err := strconv.ParseFloat(def.Args[1], 64)
											if err != nil {
												log.Printf("Error parsing %s: %+v\n", invset.Name, err)
												continue
											}
											max, err := strconv.ParseFloat(def.Args[2], 64)
											if err != nil {
												log.Printf("Error parsing %s: %+v\n", invset.Name, err)
												continue
											}
											// not perfect, but uses average to represent range
											n = (max-min)/2 + min
										} else {
											num, err := strconv.ParseFloat(def.Args[1], 64)
											if err != nil {
												log.Printf("Error parsing %s: %+v\n", invset.Name, err)
												continue
											}
											n = num
										}
									} else {
										extra = "filled"
									}
								}
								item := res.Item{
									Name:   def.Args[0],
									Amount: n,
									Equip:  extra,
								}
								inv.Items = append(inv.Items, item)
							}
						}
					}
				}
			}
		}
		inventories[inv.Name] = inv
	}
	return inventories
}

func FindItems(game string, mods []res.Mod) (items map[string][]res.Item) {
	sets := SearchContent(game, mods, "/resource/set/stuff/", "")
	items = make(map[string][]res.Item)

	costs := map[string]int{
		"bazooka":   20,
		"rifle":     10,
		"explosive": 20,
		"grenade":   8,
		"med":       5,
		"melee":     2,
		"pistol":    7,
		"smg":       12,
		"special":   18,
		"flame":     30,
		"mgun":      50,
		"shell":     3,
	}

	for _, itemset := range sets {
		name := strings.Join(strings.Split(itemset.Name, "stuff/")[1:], "stuff/")
		split := strings.Split(name, "/")
		switch split[0] {
		// only processing items which would be in an inventory
		case "bazooka", "rifle", "explosive", "grenade", "med", "melee", "pistol", "smg", "special", "flame", "mgun", "shell":
			name = split[len(split)-1]
			cost := 10
			if c, ok := costs[split[0]]; ok {
				if strings.HasSuffix(name, ".ammo") {
					c /= 4
				}
				cost = c
			}
			item := res.Item{
				Name: name,
				Cost: cost,
			}
			items[split[0]] = append(items[split[0]], item)
		}
	}
	return items
}
