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
	imageUrl := fmt.Sprintf("https://www.dustloop.com/wiki/index.php?title=File:GGST_%s_Portrait.png", name)

	c := colly.NewCollector()

	char := &character.Character{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// System Data table
	c.OnHTML(".cargoTable", func(table *colly.HTMLElement) {
		structure := make([]string, 0)
		table.ForEach("thead tr th", func(_ int, th *colly.HTMLElement) {
			if len(th.Text) == 0 {
				return
			}
			structure = append(structure, th.Text)
		})

		table.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			// TODO: Handle Giovannas Shrek install
			if i > 1 {
				return
			}
			tr.ForEach("td", func(j int, td *colly.HTMLElement) {
				key := strings.ReplaceAll(strings.Title(structure[j]), " ", "")
				err := char.SetFieldByName(key, td.Text)
				if err != nil {
					fmt.Println("Failed to set the field", err)
				}

			})
		})

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

	c.OnHTML("#mw-imagepage-section-filehistory img", func(img *colly.HTMLElement) {
		char.ImageUrl = "https://www.dustloop.com" + img.Attr("src")
	})

	c.OnScraped(func(r *colly.Response) {
		char.Name = strings.ToLower(name)
		char.Moves = moves
		char.DustloopUrl = url
		writeJson(*char)
	})

	c.Visit(url)
	c.Visit(imageUrl)
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
