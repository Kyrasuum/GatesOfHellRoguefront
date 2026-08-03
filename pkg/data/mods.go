package data

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

var ()

type Mod struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

func FindMods(workshop string) (mods []Mod, err error) {
	entries, err := os.ReadDir(workshop)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// Filter out standard files, keeping only folders for mods
		if entry.IsDir() {
			data, err := os.ReadFile(workshop + "/" + entry.Name() + "/mod.info")
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return nil, err
			}

			scanner := bufio.NewScanner(bytes.NewReader(data))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())

				if line == "" {
					continue
				}

				switch {
				case strings.HasPrefix(line, "{name"):
					mods = append(mods, Mod{
						Name:    ParseString(line),
						Id:      entry.Name(),
						Enabled: false,
						Path:    workshop + "/" + entry.Name(),
					})
					break
				}
			}
		}
	}

	return mods, nil
}
