package character

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func GetAllCharacters() []string {
	characters := make([]string, 0)
	for k := range validCharacters {
		characters = append(characters, k)
	}
	return characters
}

func LoadChar(name string) (character Character) {
	name = strings.Replace(name, " ", "_", -1)

	fname := fmt.Sprintf("%s.json", name)
	fpath, err := filepath.Abs(filepath.Join("json", fname))

	file, err := ioutil.ReadFile(fpath)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(file, &character)

	return character
}
