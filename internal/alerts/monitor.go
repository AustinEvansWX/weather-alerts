package alerts

import (
	"fmt"
	"time"
	"weather-alerts/internal/discord"
)

func Monitor() {
	for {
		warnings, err := FetchAlerts()

		if err == nil {
			warnings = FilterWarnings(warnings)
			for _, warning := range warnings {
				discord.SendAlert(warning)
			}
			ProcessWarnings(warnings)
		} else {
			fmt.Printf("Bad Request: %v\n", err)
		}

		time.Sleep(time.Second * 10)
	}
}
