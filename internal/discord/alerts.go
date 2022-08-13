package discord

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strings"
	"time"
	"weather-alerts/internal/state"
	"weather-alerts/internal/types"

	"github.com/bwmarrin/discordgo"
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
)

func SendAlert(warning *types.Warning) {
	session := state.GetSession()

	if warning.Notification == types.Update || warning.Notification == types.Cancel {
		fmt.Println("Warning updated or cancelled")
		DeleteAlert(warning.OldWarning, 0)
	}

	if warning.Notification == types.Cancel {
		return
	}

	embed := CreateEmbed(warning)
	mapImage := CreateMap(warning)

	msg, err := session.ChannelMessageSendComplex("1006012166673027105", &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		File:   &discordgo.File{Name: "map.png", Reader: mapImage},
	})

	if err != nil {
		fmt.Printf("Error sending Discord alert: %v\n", err)
		return
	}

	warning.MessageID = msg.ID
}

func CreateEmbed(warning *types.Warning) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{URL: "attachment://map.png"},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Atmos Weather Alerts",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	switch warning.Notification {
	case types.New:
		embed.Author = &discordgo.MessageEmbedAuthor{Name: "[ New ]"}
	case types.Update:
		embed.Author = &discordgo.MessageEmbedAuthor{Name: "[ Update ]"}
	}

	windThreat := warning.WindThreat
	hailThreat := warning.HailThreat
	tornadoSource := warning.TornadoSource

	if warning.Notification == types.Update {
		if warning.WindThreat != warning.OldWarning.WindThreat {
			windThreat = fmt.Sprintf("OLD:\n%s\n\nNEW:\n%s", warning.OldWarning.WindThreat, warning.WindThreat)
		}
		if warning.HailThreat != warning.OldWarning.HailThreat {
			hailThreat = fmt.Sprintf("OLD:\n%s\n\nNEW:\n%s", warning.OldWarning.HailThreat, warning.HailThreat)
		}
		if warning.TornadoSource != warning.OldWarning.TornadoSource {
			tornadoSource = fmt.Sprintf("OLD:\n%s\n\nNEW:\n%s", warning.OldWarning.TornadoSource, warning.TornadoSource)
		}
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "**Affected Areas**",
			Value: fmt.Sprintf("```\n%s\n```", strings.Join(warning.AffectedAreas, "\n")),
		},
		{Name: "", Value: ""},
		{
			Name:  "**Hail**",
			Value: fmt.Sprintf("```\n%s\n```", hailThreat),
		},
		{
			Name:  "**Expires**",
			Value: fmt.Sprintf("```\n%s\n```", FormatTime(warning.Expires)),
		},
	}

	switch warning.Type {
	case types.SpecialWeatherStatement:
		embed.Title = ":cloud_rain:  **Special Weather Statement**  :cloud_rain:"
		embed.Fields[1].Name = "**Wind**"
		embed.Fields[1].Value = fmt.Sprintf("```\n%s\n```", windThreat)
		embed.Color = 14661746
	case types.FlashFloodWarning:
		embed.Title = ":ocean:  **Flash Flood Warning**  :ocean:"
		embed.Fields = append(embed.Fields[:1], embed.Fields[3:]...)
		embed.Color = 65280
	case types.SevereThunderstorm:
		embed.Title = ":cloud_lightning:  **Severe Thunderstorm Warning**  :cloud_lightning:"
		embed.Fields[1].Name = "**Wind**"
		embed.Fields[1].Value = fmt.Sprintf("```\n%s\n```", windThreat)
		embed.Color = 16753920
	case types.Tornado:
		embed.Title = ":cloud_tornado:  **Tornado Warning**  :cloud_tornado:"
		embed.Color = 13458780
		embed.Fields[1].Name = "**Tornado Source**"
		embed.Fields[1].Value = fmt.Sprintf("```\n%s\n```", tornadoSource)
	}

	return embed
}

func FormatTime(t time.Time) string {
	formatted := ""

	t = t.Local()
	hour := t.Hour()

	meridiem := "PM"

	if hour < 12 {
		meridiem = "AM"
		if hour == 0 {
			hour = 12
		}
	} else if hour > 12 {
		hour -= 12
	}

	formatted += fmt.Sprintf("%d:", hour)

	minute := t.Minute()

	if minute < 10 {
		formatted += fmt.Sprintf("0%d ", minute)
	} else {
		formatted += fmt.Sprintf("%d ", minute)
	}

	timeZome, _ := t.Zone()

	formatted += fmt.Sprintf("%s %s", meridiem, timeZome)

	return formatted
}

func CreateMap(warning *types.Warning) *bytes.Reader {
	var fillColor color.Color
	var strokeColor color.Color

	switch warning.Type {
	case types.SpecialWeatherStatement:
		fillColor = color.NRGBA{223, 184, 114, 120}
		strokeColor = color.RGBA{223, 184, 114, 255}
	case types.FlashFloodWarning:
		fillColor = color.NRGBA{0, 255, 0, 120}
		strokeColor = color.RGBA{0, 255, 0, 255}
	case types.SevereThunderstorm:
		fillColor = color.NRGBA{255, 165, 0, 120}
		strokeColor = color.RGBA{255, 165, 0, 255}
	case types.Tornado:
		fillColor = color.NRGBA{205, 93, 92, 170}
		strokeColor = color.RGBA{205, 93, 92, 255}
	}

	ctx := sm.NewContext()
	ctx.SetSize(750, 750)

	coords := []s2.LatLng{}

	for _, point := range warning.Path {
		coords = append(coords, s2.LatLngFromDegrees(point[1], point[0]))
	}

	warningArea := sm.NewArea(coords, strokeColor, fillColor, 5)

	ctx.SetBoundingBox(warningArea.Bounds())
	ctx.AddObject(warningArea)

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)

	return bytes.NewReader(buf.Bytes())
}

func DeleteAlert(warning *types.Warning, attempts int) {
	session := state.GetSession()

	err := session.ChannelMessageDelete("1006012166673027105", warning.MessageID)

	if err != nil {
		fmt.Printf("Error deleting Discord alert %s: %v\n", warning.MessageID, err)
		if attempts < 3 {
			attempts++
			DeleteAlert(warning, attempts)
		}
	}
}
