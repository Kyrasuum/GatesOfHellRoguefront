package data

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Inventory struct {
	Id    int    `json:"-"`
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

func FindInventories(game string, mods []Mod) (inventories map[string]Inventory) {
	sets := SearchContent(game, mods, "/resource/set/breed/mp/", ".set")
	inventories = make(map[string]Inventory)

	for _, invset := range sets {
		name := strings.Join(strings.Split(invset.Name, "breed/")[1:], "breed/")
		name = name[:len(name)-4]
		inv := Inventory{
			Name:  name,
			Items: []Item{},
		}
		if len(invset.Body) > 0 {
			if len(invset.Body[0].Body) > 0 {
				props := invset.Body[0].Body
				for _, prop := range props {
					switch prop.Name {
					case "armors":
						for _, def := range prop.Body {
							item := Item{
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
								item := Item{
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
