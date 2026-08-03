package data

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"roguefront/pkg/set"
)

var ()

func SearchPak(path string, prefix string, postfix string) (sets []*set.Set) {
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
			s, err := set.Parse(path+"/"+f.Name, data)
			if err != nil {
				log.Printf("Error on %s: %+v\n", path+"/"+f.Name, err)
			} else {
				sets = append(sets, s)
			}
		}
	}
	return sets
}

func FindPaks(path string, prefix string, postfix string) (sets []*set.Set) {
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

func FindLoose(path string, prefix string, postfix string) (sets []*set.Set) {
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
			s, err := set.Parse(path+"/"+entry.Name(), data)
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

func SearchContent(game string, mods []Mod, search string, ft string) (sets []*set.Set) {
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

func FindNations(game string, mods []Mod) (nations []string) {
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

func FindLevels(game string, mods []Mod) (levels []string) {
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

func FindResources(game string, mods []Mod) (resources []string) {
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

func FindMaps(game string, mods []Mod) (maps map[string][]string) {
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
