package model

import (
	"encoding/json"
	"errors"
	"net/http"

	"groupie-tracker/config"
)

func GetAllAtritst(Data *[]Artist) error { // этой функцией достаем всю инфу по всем артистам из Json ссылки
	netClient := http.Client{} // создаем HTTP клиента

	resp, err := netClient.Get(config.ArtistURL) // отправляем GET запрос на ArtistURL, и получаем resp
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { // проверяем статус resp, который должен быть 2хх
		errorText := "invlid status of the api's response; url: " + config.ArtistURL + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body) // открываем Body в json формате
	defer resp.Body.Close()

	err = dec.Decode(Data) // декодированные данные из dec отправляем в структуру Data

	if err != nil {
		return err
	}

	return nil
}

func GetAtritst(id string, Data *Artist) error { // этой функцией достаем всю инфу по определенному артисту из Json ссылки
	netClient := http.Client{}

	resp, err := netClient.Get(config.ArtistURL + "/" + id)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorText := "invlid status of the api's response; url: " + config.ArtistURL + "/" + id + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(Data)

	if err != nil {
		return err
	}

	return nil
}

func GetAllRelations(Data *map[string][]DateLocation) error {
	netClient := http.Client{}

	resp, err := netClient.Get(config.RelationURL)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorText := "invlid status of the api's response; url: " + config.RelationURL + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(Data)

	if err != nil {
		return err
	}

	return nil
}

func GetRelation(id string, Data *Artist) error {
	netClient := http.Client{}

	resp, err := netClient.Get(config.RelationURL + "/" + id)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorText := "invlid status of the api's response; url: " + config.RelationURL + "/" + id + "; status: " + resp.Status
		return errors.New(errorText)
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(Data)

	if err != nil {
		return err
	}

	return nil
}
