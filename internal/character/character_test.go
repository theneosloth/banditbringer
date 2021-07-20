package character

import (
	"testing"
)

func TestIsValidCharName(t *testing.T) {
	exactMatch := "ramlethal_valentine"
	roughMatch := "Chipp Zanuff"
	validAlias := "Kyle"
	invalidName := "Devil Jin"

	_, isExactMatch := IsValidCharName(exactMatch)

	if !isExactMatch {
		t.Errorf("%s should be an exact match", exactMatch)
	}

	roughMatchNormalized, _ := IsValidCharName(roughMatch)

	if roughMatchNormalized != "chipp_zanuff" {
		t.Errorf("%s should be normalized to", "chipp_zanuff")
	}

	validAliasNormalized, _ := IsValidCharName(validAlias)
	if validAliasNormalized != "ky_kiske" {
		t.Errorf("alias %s should be normalized to", "ky_kiske")
	}

	_, foundInvalidName := IsValidCharName(invalidName)

	if foundInvalidName {
		t.Errorf("%s is not a valid name", invalidName)
	}

}

func TestSetName(t *testing.T) {
	validCharacter := Character{}
	err := validCharacter.SetName("Millia Rage")

	if err != nil {
		t.Errorf("%s should be a valid name", "Millia Rage")
	}

	invalidCharacter := Character{}
	err = invalidCharacter.SetName("Marshall Law")

	if err == nil {
		t.Errorf("%s should bt an invalid name", "Marshall Law")
	}

}
