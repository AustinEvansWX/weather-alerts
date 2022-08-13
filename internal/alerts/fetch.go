package alerts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"weather-alerts/internal/types"
)

func FetchAlerts() ([]*types.Warning, error) {
	rawBody, err := request("active", map[string]string{
		"status": "actual",
		"event":  "Special%20Weather%20Statement,Flash%20Flood%20Warning,Severe%20Thunderstorm%20Warning,Tornado%20Warning",
	})

	if err != nil {
		return nil, err
	}

	body := types.ActiveResponse{}
	err = json.Unmarshal(rawBody, &body)

	if err != nil {
		return nil, err
	}

	warnings := []*types.Warning{}

	for _, alert := range body.Features {
		properties := alert.Properties

		warning := &types.Warning{
			Type:          types.WarningType(properties.Event),
			Notification:  types.New,
			References:    properties.References,
			ID:            properties.ID0,
			NWSCenter:     properties.SenderName,
			AffectedAreas: strings.Split(properties.AreaDesc, "; "),
			Expires:       properties.Expires,
		}

		if properties.MessageType == "Cancel" {
			warning.Notification = types.Cancel
		} else if len(properties.References) > 0 {
			warning.Notification = types.Update
		}

		if warning.Notification != types.Cancel {
			if len(alert.Geometry.Coordinates) == 0 {
				continue
			}

			warning.Path = alert.Geometry.Coordinates[0]

			windThreat := ""
			hailThreat := ""
			tornadoSource := ""

			switch types.WarningType(properties.Event) {
			case types.SpecialWeatherStatement:
				if len(properties.Parameters.MaxWindGust) == 0 {
					continue
				}
				windThreat = properties.Parameters.MaxWindGust[0]
				hailSize := properties.Parameters.MaxHailSize[0]
				hailThreat = fmt.Sprintf("%s\" %s", hailSize, GetHailSizeEquivalent(hailSize))
			case types.SevereThunderstorm:
				if len(properties.Parameters.WindThreat) == 0 {
					continue
				}
				windThreat = fmt.Sprintf("%s\n%s", properties.Parameters.WindThreat[0], properties.Parameters.MaxWindGust[0])
				hailSize := properties.Parameters.MaxHailSize[0]
				hailThreat = fmt.Sprintf("%s\n%s\" %s", properties.Parameters.HailThreat[0], hailSize, GetHailSizeEquivalent(hailSize))
			case types.Tornado:
				if len(properties.Parameters.TornadoDetection) == 0 {
					continue
				}
				tornadoSource = properties.Parameters.TornadoDetection[0]
				hailSize := properties.Parameters.MaxHailSize[0]
				hailThreat = fmt.Sprintf("%s\" %s", hailSize, GetHailSizeEquivalent(hailSize))
			}

			warning.WindThreat = windThreat
			warning.HailThreat = hailThreat
			warning.TornadoSource = tornadoSource
		}

		warnings = append(warnings, warning)
	}

	return warnings, nil
}

func request(path string, params map[string]string) ([]byte, error) {
	queryString := ""

	for name, value := range params {
		queryString += fmt.Sprintf("%s=%s&", name, value)
	}

	resp, err := http.Get(fmt.Sprintf("https://api.weather.gov/alerts/%s?%s", path, queryString))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
