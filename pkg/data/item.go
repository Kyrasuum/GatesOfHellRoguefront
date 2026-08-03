package data

import (
	"strings"
)

type Item struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Pos    *struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"pos"`
	Equip string `json:"equip"`
	Cost  int    `json:"cost"`
}

func FindItems(game string, mods []Mod) (items map[string][]Item) {
	sets := SearchContent(game, mods, "/resource/set/stuff/", "")
	items = make(map[string][]Item)

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
			item := Item{
				Name: name,
				Cost: cost,
			}
			items[split[0]] = append(items[split[0]], item)
		}
	}
	return items
}
