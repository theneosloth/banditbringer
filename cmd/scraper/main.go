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
	}

	for _, character := range characters {
		scrapeCharacter(character)
	}
}

func scrapeCharacter(name string) {
	name = strings.Replace(name, " ", "_", -1)
	moves := make([]move.Move, 0)

	url := fmt.Sprintf("https://www.dustloop.com/wiki/index.php?title=GGST/%s/Frame_Data", name)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(".cargoDynamicTable ", func(e *colly.HTMLElement) {

		// Store the order of table headings, since it's inconsistent between tables
		structure := make([]string, 0)

		e.ForEach("thead tr th", func(_ int, th *colly.HTMLElement) {
			if len(th.Text) == 0 {
				return
			}
			structure = append(structure, th.Text)
		})

		e.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			m := move.Move{}

			tr.ForEach("td", func(i int, td *colly.HTMLElement) {

				// First element in the table is blank, the rest align with the table headings
				if i == 0 {
					return
				}
				i = i - 1

				key := strings.Title(structure[i])
				err := m.SetFieldByName(key, td.Text)
				if err != nil {
					fmt.Println("Failed to set the field", err)
				}

			})

			if len(m.Input) > 0 || len(m.Name) > 0 {
				moves = append(moves, m)
			}

		})
	})

	c.Visit(url)
	writeJson(strings.ToLower(name), moves)
}

func writeJson(name string, moves []move.Move) error {

	fname := fmt.Sprintf("%s.json", name)
	fpath, err := filepath.Abs(filepath.Join("json", fname))

	if err != nil {
		return fmt.Errorf("failed to generate filepath %s", fpath)
	}

	file, err := os.Create(fpath)

	if err != nil {
		return fmt.Errorf("failed to create file %s", fname)
	}

	defer file.Close()

	character := character.Character{}
	character.SetName(name)
	character.Moves = moves

	data, err := json.MarshalIndent(character, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal %s", name)
	}

	ioutil.WriteFile(fpath, data, 0777)
	return nil
}
