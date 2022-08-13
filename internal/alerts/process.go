package alerts

import (
	"time"
	"weather-alerts/internal/discord"
	"weather-alerts/internal/state"
	"weather-alerts/internal/types"
)

func FilterWarnings(warnings []*types.Warning) []*types.Warning {
	filteredWarnings := []*types.Warning{}

	for _, warning := range warnings {
		if state.GetWarning(warning.ID) != nil {
			continue
		}

		if warning.Expires.Sub(time.Now()) <= 0 {
			continue
		}

		if warning.Notification == types.Update || warning.Notification == types.Cancel {
			tracked := false

			for _, ref := range warning.References {
				oldWarning := state.GetWarning(ref.ID0)
				if oldWarning != nil {
					warning.OldWarning = oldWarning
					tracked = true
					break
				}
			}

			if !tracked {
				continue
			}
		}

		filteredWarnings = append(filteredWarnings, warning)
	}

	return filteredWarnings
}

func ProcessWarnings(warnings []*types.Warning) {
	for _, warning := range warnings {
		switch warning.Notification {
		case types.New:
			state.AddWarning(warning)
			id := warning.ID
			warning.Timer = time.AfterFunc(warning.Expires.Sub(time.Now()), func() {
				warning := state.GetWarning(id)
				state.RemoveWarning(warning.ID)
				discord.DeleteAlert(warning, 0)
			})
		case types.Update:
			warning.OldWarning.Timer.Stop()
			id := warning.ID
			warning.Timer = time.AfterFunc(warning.Expires.Sub(time.Now()), func() {
				warning := state.GetWarning(id)
				state.RemoveWarning(warning.ID)
				discord.DeleteAlert(warning, 0)
			})
			state.RemoveWarning(warning.OldWarning.ID)
			state.AddWarning(warning)
		case types.Cancel:
			warning.OldWarning.Timer.Stop()
			state.RemoveWarning(warning.OldWarning.ID)
		}
	}
}
