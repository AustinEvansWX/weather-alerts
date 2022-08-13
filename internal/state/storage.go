package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"weather-alerts/internal/types"
)

func Read() (types.WarningMap, error) {
	data, err := ioutil.ReadFile("./storage.json")

	if err != nil {
		return nil, err
	}

	var warnings types.WarningMap
	err = json.Unmarshal(data, &warnings)
	if err != nil {
		return nil, err
	}

	return warnings, nil
}

func Store(warnings types.WarningMap) error {
	data, err := json.Marshal(warnings)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./storage.json", data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func Save() {
	err := Store(warningState.Warnings)
	if err != nil {
		fmt.Printf("Error writing to storage file: %v\n", err)
	}
}
