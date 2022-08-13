package state

import (
	"fmt"
	"sync"
	"weather-alerts/internal/types"
)

type WarningState struct {
	Mu       sync.RWMutex
	Warnings types.WarningMap
}

var (
	warningState = WarningState{Warnings: types.WarningMap{}}
)

func Initialize() {
	warnings, err := Read()

	if err != nil {
		fmt.Printf("Error reading from storage file: %v\n", err)
		return
	}

	warningState.Warnings = warnings
}

func GetWarning(id string) *types.Warning {
	return warningState.Warnings[id]
}

func GetWarnings(ids []string) types.WarningMap {
	fetchedWarnings := map[string]*types.Warning{}
	for _, id := range ids {
		fetchedWarnings[id] = warningState.Warnings[id]
	}
	return fetchedWarnings
}

func GetAllWarnings() types.WarningMap {
	return warningState.Warnings
}

func AddWarning(warning *types.Warning) {
	warningState.Mu.Lock()
	defer warningState.Mu.Unlock()
	warningState.Warnings[warning.ID] = warning
	Save()
}

func AddWarnings(warnings []*types.Warning) {
	warningState.Mu.Lock()
	defer warningState.Mu.Unlock()
	for _, warning := range warnings {
		warningState.Warnings[warning.ID] = warning
	}
	Save()
}

func RemoveWarning(id string) {
	warningState.Mu.Lock()
	defer warningState.Mu.Unlock()
	delete(warningState.Warnings, id)
	Save()
}

func RemoveWarnings(ids []string) {
	warningState.Mu.Lock()
	defer warningState.Mu.Unlock()
	for _, id := range ids {
		delete(warningState.Warnings, id)
	}
	Save()
}
