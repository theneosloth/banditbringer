package main

import (
	"banditbringer/internal/character"
	"banditbringer/internal/move"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

var baseURL = "http://dustloop.com"

func main() {
	characters := [...]string{
		"Ky Kiske",
		"Anji Mito",
		"Axl Low",
		"Chipp Zanuff",
		"Faust",
		"Giovanna",
		"I-No",
		"Leo Whitefang",
		"May",
		"Millia Rage",
		"Nagoriyuki",
		"Potemkin",
		"Ramlethal Valentine",
		"Sol Badguy",
		"Zato-1",
		"Goldlewis Dickinson",
		"Jack-O",
	}

	for _, character := range characters {
		scrapeCharacter(character)
	}
}

// Fragile and ugly, probably better to construct a map	of heading names to struct names
func HeadingToStructName(name string) string {
	name = strings.TrimSpace(name)
	aliases := map[string]string{
		"R.I.S.C. Gain":          "RiscGain",
		"CH Type":                "ChType",
		"PreJump":                "Prejump",
		"R.I.S.C. Gain Modifier": "RiscModifier",
	}

	alias, hasAlias := aliases[name]

	if hasAlias {
		return alias
	}

	return strings.Title(
		strings.ReplaceAll(strings.ReplaceAll(name, "-", " "), " ", ""))
}

func scrapeCharacter(name string) {
	name = strings.Replace(name, " ", "_", -1)
	moves := make([]move.Move, 0)

	url := fmt.Sprintf("%s/wiki/index.php?title=GGST/%s/Data", baseURL, name)

	c := colly.NewCollector()

	char := &character.Character{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// System Data table
	c.OnHTML("table.sortable", func(table *colly.HTMLElement) {
		structure := make([]string, 0)

		table.ForEach("tr th", func(_ int, th *colly.HTMLElement) {
			if len(th.Text) == 0 {
				return
			}
			structure = append(structure, th.Text)
		})

		table.ForEach("tr td", func(i int, th *colly.HTMLElement) {
			field := HeadingToStructName(structure[i])

			result := ""
			// TODO: Handle multiple images properly
			if field == "Portrait" || field == "Icon" {
				result = fmt.Sprintf("%s%s", baseURL, th.ChildAttr("img", "src"))
			} else {
				result = strings.TrimSpace(th.Text)
			}

			err := char.SetFieldByName(field, result)
			if err != nil {
				fmt.Println("failed to set", err, field)
			}
		})

	})

	//individual move tables
	c.OnHTML("table.wikitable:not(.sortable)", func(table *colly.HTMLElement) {
		structure := make([]string, 0)

		m := move.Move{}
		table.ForEach("tr th", func(_ int, th *colly.HTMLElement) {
			if len(th.Text) == 0 {
				return
			}
			structure = append(structure, th.Text)
		})

		table.ForEach("tr td", func(i int, th *colly.HTMLElement) {

			// Images don't have a matching header, have to be handled first
			if i >= len(structure) {
				// This will break eventually
				imageIndex := len(structure)
				hitboxIndex := imageIndex + 1

				getUrl := func() string {
					imgUrl := th.ChildAttr("img", "src")
					if imgUrl == "" {
						return ""
					}
					return fmt.Sprintf("%s%s", baseURL, th.ChildAttr("img", "src"))
				}

				if i == imageIndex {
					m.SetFieldByName("Images", getUrl())
				}

				if i == hitboxIndex {
					m.SetFieldByName("Hitboxes", getUrl())
				}
				return
			}

			field := HeadingToStructName(structure[i])
			result := strings.TrimSpace(th.Text)

			err := m.SetFieldByName(field, result)
			if err != nil {
				fmt.Println("failed to set", err, field)
			}
		})

		if len(m.Input) > 0 || len(m.Name) > 0 {
			moves = append(moves, m)
		}

	})

	c.OnScraped(func(r *colly.Response) {
		char.Name = strings.ToLower(name)
		char.Moves = moves
		char.DustloopUrl = fmt.Sprintf("%s/wiki/index.php?title=GGST/%s", baseURL, name)
		writeJson(*char)
	})

	c.Visit(url)
}

func writeJson(char character.Character) error {

	fname := fmt.Sprintf("%s.json", char.Name)
	fpath, err := filepath.Abs(filepath.Join("json", fname))

	if err != nil {
		return fmt.Errorf("failed to generate filepath %s", fpath)
	}

	file, err := os.Create(fpath)

	if err != nil {
		return fmt.Errorf("failed to create file %s", fname)
	}

	defer file.Close()

	data, err := json.MarshalIndent(char, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal %s", char.Name)
	}

	ioutil.WriteFile(fpath, data, 0777)
	return nil
}
