package model

import (
	"regexp"
	"strconv"
	"strings"
)

func (artist *Artist) ChangeKey() {
	newMap := make(map[string][]string)
	for key, value := range artist.DatesLocations {
		newKey := strings.ToUpper(key)
		reg := regexp.MustCompile(`-`)
		newKey = reg.ReplaceAllString(newKey, "$1 - $2")
		reg = regexp.MustCompile(`_`)
		newKey = reg.ReplaceAllString(newKey, "$1-$2")
		newMap[newKey] = value
	}
	artist.DatesLocations = newMap
}

func (artistsList *ArtistsList) ChangeKeys() {
	for i := 0; i < len(artistsList.List); i++ {
		artistsList.List[i].ChangeKey()
	}
}

func (artistsList *ArtistsList) AddDatesLocations(inputList []DateLocation) { // этой функцией из relations копируем даты и места в
	for i := 0; i < len(artistsList.List); i++ {
		for j := 0; j < len(inputList); j++ {
			if artistsList.List[i].Id == inputList[j].Id { // если ID из artists и relations совпадают, то
				artistsList.List[i].DatesLocations = inputList[j].DatesLocations // мы заполняем пустые данные ArtistsList.List (DateLocations) из RelationList.List (DatesLocations)
				break
			}
		}
	}
}

func (artist *Artist) ContainsInput(input string) bool {
	if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(input)) {
		return true
	} else if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(input)) {
		return true
	} else if strings.Contains(strconv.Itoa(artist.CreationDate), input) {
		return true
	}

	for i := 0; i < len(artist.Members); i++ {
		if strings.Contains(strings.ToLower(artist.Members[i]), strings.ToLower(input)) {
			return true
		}
	}

	for key, value := range artist.DatesLocations {
		if strings.Contains(strings.ToLower(key), strings.ToLower(input)) {
			return true
		}
		for i := 0; i < len(value); i++ {
			if strings.Contains(strings.ToLower(value[i]), strings.ToLower(input)) {
				return true
			}
		}
	}

	return false
}

func (searchResultArtists *ArtistsList) SearchInputInArtistsList(artists []Artist, input string) {
	for i := 0; i < len(artists); i++ {
		if artists[i].ContainsInput(input) {
			searchResultArtists.List = append(searchResultArtists.List, artists[i])
		}
	}
}
