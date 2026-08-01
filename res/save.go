package res

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type Save struct {
	Campaign *Campaign
	Status   *Status
}

type Status struct {
	Mods        []string
	Timestamp   int64
	Seed        int64
	Name        string
	Army        string
	EnemyArmy   string
	Difficulty  string
	Resources   int
	FogOfWar    string
	Region      string
	PlayedGames int
	WonGames    int
	Landscape   string
	Gamemode    string
	Map         string
	Risk        string
	Texmod      string
}

type Campaign struct {
	Soldiers    []*Soldier
	Inventories []*Inventory
	Squads      []*Squad
}

type Soldier struct {
	Id   int
	Path string
	Name string
	Mid  string
}

type Inventory struct {
	Id    int
	Items []Item
}

type Item struct {
	Name   string
	Amount int
	Pos    *struct {
		X int
		Y int
	}
	Equip string
}

type Squad struct {
	Name     string
	Soldiers []*Soldier
}

func ReadSave(filename string) (*Save, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	s := &Save{}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(rc)
		rc.Close()

		if err != nil {
			return nil, err
		}

		switch f.Name {
		case "campaign.scn":
			s.Campaign, err = ParseCampaign(data)
			if err != nil {
				return nil, err
			}

		case "status":
			s.Status, err = ParseStatus(data)
			if err != nil {
				return nil, err
			}
		}
	}

	if s.Campaign == nil {
		return nil, fmt.Errorf("%s\n", "campaign.scn missing")
	}

	if s.Status == nil {
		return nil, fmt.Errorf("%s\n", "status missing")
	}

	return s, nil
}

func (s *Save) Write(path string) error {
	file, err := os.Create(path + "/" + strcase.ToLowerCamel(s.Status.Name) + ".sav")
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	write := func(name string, data []byte) error {
		f, err := w.Create(name)
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		return err
	}

	if err := write("campaign.scn", s.Campaign.Bytes()); err != nil {
		return err
	}

	if err := write("status", s.Status.Bytes()); err != nil {
		return err
	}

	return nil
}

func (s *Status) Bytes() []byte {
	var w bytes.Buffer

	fmt.Fprintf(&w, "{saveinfo\n")

	fmt.Fprintf(&w, "\t{version 9}\n")
	fmt.Fprintf(&w, "\t{gameVersion \"1.065.0\"}\n")

	fmt.Fprintf(&w, "\t{mods\n")
	for _, mod := range s.Mods {
		fmt.Fprintf(&w, "\t\t\"%s\"\n", mod)
	}
	fmt.Fprintf(&w, "\t}\n")
	fmt.Fprintf(&w, "\t{timestamp %d}\n", s.Timestamp)
	fmt.Fprintf(&w, "\t{seed %d}\n", s.Seed)
	fmt.Fprintf(&w, "\t{name \"%s\"}\n", s.Name)
	fmt.Fprintf(&w, "\t{army \"%s\"}\n", s.Army)
	fmt.Fprintf(&w, "\t{enemyArmy \"%s\"}\n", s.EnemyArmy)
	fmt.Fprintf(&w, "\t{difficulty \"%s\"}\n", s.Difficulty)
	fmt.Fprintf(&w, "\t{resources %d}\n", s.Resources)
	fmt.Fprintf(&w, "\t{fogofwar %s}\n", s.FogOfWar)
	fmt.Fprintf(&w, "\t{region \"%s\"}\n", s.Region)
	fmt.Fprintf(&w, "\t{playedGames %d}\n", s.PlayedGames)
	fmt.Fprintf(&w, "\t{wonGames %d}\n", s.WonGames)
	fmt.Fprintf(&w, "\t{landscape \"%s\"}\n", s.Landscape)

	fmt.Fprintf(&w, "\t{mp 0}\n")
	fmt.Fprintf(&w, "\t{sp 0}\n")
	fmt.Fprintf(&w, "\t{ap 0}\n")
	fmt.Fprintf(&w, "\t{rp 0}\n")
	fmt.Fprintf(&w, "\t{duration 1}\n")
	fmt.Fprintf(&w, "\t{unlockedResearch {\"reinforcement_stage_1\"}\n\t\t{\"reinforcement_stage_2\"}\n\t\t{\"reinforcement_stage_3\"}\n\t\t{\"reinforcement_stage_4\"}\n\t\t{\"reinforcement_stage_5\"}\n\t}\n")
	fmt.Fprintf(&w, "\t{manualControlMode 3}\n")
	fmt.Fprintf(&w, "\t{selectedMapPoint hq_b}\n")
	fmt.Fprintf(&w, "\t{attacking}\n")

	fmt.Fprintf(&w, "\t{mapPoints\n")
	fmt.Fprintf(&w, "\t\t{\n")
	fmt.Fprintf(&w, "\t\t\t{name hq_a}\n")
	fmt.Fprintf(&w, "\t\t\t{landscape wood}\n")
	fmt.Fprintf(&w, "\t\t\t{gamemode campaign_capture_the_flag}\n")
	fmt.Fprintf(&w, "\t\t\t{ownerTeam a}\n")
	fmt.Fprintf(&w, "\t\t\t{adjacentMaps {\"hq_b\"}}\n")
	fmt.Fprintf(&w, "\t\t\t{risk %s}\n", s.Risk)
	fmt.Fprintf(&w, "\t\t\t{reward none}\n")
	fmt.Fprintf(&w, "\t\t\t{map \"%s\"}\n", s.Map)
	fmt.Fprintf(&w, "\t\t\t{texmod %s}\n", s.Texmod)
	fmt.Fprintf(&w, "\t\t}\n")
	fmt.Fprintf(&w, "\t\t{\n")
	fmt.Fprintf(&w, "\t\t\t{name hq_b}\n")
	fmt.Fprintf(&w, "\t\t\t{landscape wood}\n")
	fmt.Fprintf(&w, "\t\t\t{gamemode %s}\n", s.Gamemode)
	fmt.Fprintf(&w, "\t\t\t{ownerTeam b}\n")
	fmt.Fprintf(&w, "\t\t\t{adjacentMaps {\"hq_a\"}}\n")
	fmt.Fprintf(&w, "\t\t\t{risk %s}\n", s.Risk)
	fmt.Fprintf(&w, "\t\t\t{reward none}\n")
	fmt.Fprintf(&w, "\t\t\t{map \"%s\"}\n", s.Map)
	fmt.Fprintf(&w, "\t\t\t{texmod %s}\n", s.Texmod)
	fmt.Fprintf(&w, "\t\t}\n")
	fmt.Fprintf(&w, "\t}\n")
	fmt.Fprintf(&w, "\t{roundsHistory}\n")
	fmt.Fprintf(&w, "}\n")

	return w.Bytes()
}

func (c *Campaign) Bytes() []byte {
	var w bytes.Buffer

	fmt.Fprintf(&w, "{campaign\n")
	for _, soldier := range c.Soldiers {
		fmt.Fprintf(&w, "%s", soldier.Bytes())
	}
	for _, inventory := range c.Inventories {
		fmt.Fprintf(&w, "%s", inventory.Bytes())
	}
	fmt.Fprintf(&w, "\t{CampaignSquads\n")
	for _, squad := range c.Squads {
		fmt.Fprintf(&w, "%s", squad.Bytes())
	}
	fmt.Fprintf(&w, "\t}\n")
	fmt.Fprintf(&w, "}\n")

	return w.Bytes()
}

func (s *Soldier) Bytes() []byte {
	var w bytes.Buffer

	fmt.Fprintf(&w, "\t{human \"%s\" 0xc%03d\n", s.Path, s.Id)
	fmt.Fprintf(&w, "\t\t{Position 0 0}\n")
	fmt.Fprintf(&w, "\t\t{TexMod \"auto\"}\n")
	fmt.Fprintf(&w, "\t\t{SpawnedInFog}\n")
	fmt.Fprintf(&w, "\t\t{Volume \"ram\"\n")
	fmt.Fprintf(&w, "\t\t\t{able {visible 0}{bullet 0}{throwing 0}{obstacle 0}{contact 0}{contact_ground 0}{blast 0}{select 0}{touch 0}{blockcamera 0}}\n")
	fmt.Fprintf(&w, "\t\t\t{disabled}\n")
	fmt.Fprintf(&w, "\t\t}\n")
	fmt.Fprintf(&w, "\t\t{Player 0}\n")
	fmt.Fprintf(&w, "\t\t{MID %s}\n", s.Mid)
	fmt.Fprintf(&w, "\t\t{NameId %s}\n", s.Name)
	fmt.Fprintf(&w, "\t\t{FsmState \"stand_noaim\"}\n")
	fmt.Fprintf(&w, "\t}\n")

	return w.Bytes()
}

func (i *Inventory) Bytes() []byte {
	var w bytes.Buffer

	fmt.Fprintf(&w, "\t{Inventory\n")
	fmt.Fprintf(&w, "\t\t{box\n")
	fmt.Fprintf(&w, "\t\t\t{clear}\n")
	for _, item := range i.Items {
		fmt.Fprintf(&w, "%s", item.Bytes())
	}
	fmt.Fprintf(&w, "\t\t}\n")
	fmt.Fprintf(&w, "\t}\n")

	return w.Bytes()
}

func (i *Item) Bytes() []byte {
	var w bytes.Buffer

	fmt.Fprintf(&w, "\t\t\t{item %s ", i.Name)
	if i.Amount > 0 {
		fmt.Fprintf(&w, "%d ", i.Amount)
	}
	if i.Pos != nil {
		fmt.Fprintf(&w, "{cell %d %d}", i.Pos.X, i.Pos.Y)
	}
	fmt.Fprintf(&w, "%s}\n", i.Equip)

	return w.Bytes()
}

func (s *Squad) Bytes() []byte {
	var w bytes.Buffer

	fmt.Fprintf(&w, "\t\t{\"%s\" \"stage_1\"", s.Name)
	for _, soldier := range s.Soldiers {
		fmt.Fprintf(&w, "0xc%03d", soldier.Id)
	}
	fmt.Fprintf(&w, "}\n")

	return w.Bytes()
}

func ParseStatus(data []byte) (*Status, error) {
	s := &Status{}

	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "{mods"):
			if err := ParseMods(scanner, s); err != nil {
				return nil, err
			}

		case strings.HasPrefix(line, "{timestamp"):
			s.Timestamp = ParseInt64(line)

		case strings.HasPrefix(line, "{seed"):
			s.Seed = ParseInt64(line)

		case strings.HasPrefix(line, "{name"):
			s.Name = ParseString(line)

		case strings.HasPrefix(line, "{army"):
			s.Army = ParseString(line)

		case strings.HasPrefix(line, "{enemyArmy"):
			s.EnemyArmy = ParseString(line)

		case strings.HasPrefix(line, "{difficulty"):
			s.Difficulty = ParseString(line)

		case strings.HasPrefix(line, "{resources"):
			s.Resources = int(ParseInt64(line))

		case strings.HasPrefix(line, "{fogofwar"):
			s.FogOfWar = ParseString(line)

		case strings.HasPrefix(line, "{region"):
			s.Region = ParseString(line)

		case strings.HasPrefix(line, "{playedGames"):
			s.PlayedGames = int(ParseInt64(line))

		case strings.HasPrefix(line, "{wonGames"):
			s.WonGames = int(ParseInt64(line))

		case strings.HasPrefix(line, "{landscape"):
			s.Landscape = ParseString(line)

		case strings.HasPrefix(line, "{map"):
			s.Map = ParseString(line)

		case strings.HasPrefix(line, "{risk"):
			s.Risk = ParseString(line)
		}
	}

	return s, scanner.Err()
}

func ParseCampaign(data []byte) (*Campaign, error) {
	c := &Campaign{}

	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "{Human"):
			soldier, err := ParseSoldier(scanner, line)
			if err != nil {
				return nil, err
			}
			c.Soldiers = append(c.Soldiers, soldier)
		case strings.HasPrefix(line, "{Inventory"):
			inventory, err := ParseInventory(scanner, line)
			if err != nil {
				return nil, err
			}
			c.Inventories = append(c.Inventories, inventory)
		case strings.HasPrefix(line, "{CampaignSquads"):
		}
	}

	return c, scanner.Err()
}

func ParseSoldier(scanner *bufio.Scanner, line string) (*Soldier, error) {
	s := &Soldier{}

	fields := strings.Fields(line)
	if len(fields) >= 2 {
		s.Path = strings.Trim(fields[1], "\"")
	}
	if len(fields) >= 3 {
		s.Id = int(ParseInt64(strings.Trim(fields[2], "0xc")))
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "{MID"):
			s.Mid = ParseString(line)
		case strings.HasPrefix(line, "{NameId"):
			s.Name = ParseString(line)
		case strings.HasPrefix(line, "}"):
			break
		}
	}

	return s, scanner.Err()
}

func ParseInventory(scanner *bufio.Scanner, line string) (*Inventory, error) {
	i := &Inventory{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "{item"):
		}
	}

	return i, scanner.Err()
}

func FieldValue(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "{")

	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	return strings.TrimSpace(strings.TrimSuffix(parts[1], "}"))
}

func ParseInt64(line string) int64 {
	n, _ := strconv.ParseInt(FieldValue(line), 10, 64)
	return n
}

func ParseFloat(line string) float64 {
	f, _ := strconv.ParseFloat(FieldValue(line), 64)
	return f
}

func ParseString(line string) string {
	return strings.Trim(FieldValue(line), "\"")
}

func ParseMods(scanner *bufio.Scanner, s *Status) error {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "}" {
			return nil
		}

		mod := ParseString(line)
		if mod != "" {
			s.Mods = append(s.Mods, mod)
		}
	}

	return scanner.Err()
}
