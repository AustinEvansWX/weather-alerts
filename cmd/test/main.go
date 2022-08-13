package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weather-alerts/internal/alerts"
	"weather-alerts/internal/discord"
	"weather-alerts/internal/state"
)

func main() {
	session, err := discord.Login("MTAwNjAxMjQwMjgyOTEwNzI3MA.GmpRgi.x858SjNwtF1uba8iTEhC3h3GPQp063SItC5-Ds")
	defer session.Close()

	state.SetSession(session)

	if err != nil {
		fmt.Printf("Error launching Discord session: %v\n", err)
		return
	}

	state.Initialize()

	for _, warning := range state.GetAllWarnings() {
		id := warning.ID
		warning.Timer = time.AfterFunc(warning.Expires.Sub(time.Now()), func() {
			warning := state.GetWarning(id)
			state.RemoveWarning(warning.ID)
			discord.DeleteAlert(warning, 0)
		})
	}

	alerts.Monitor()

	waitForExit()
}

func waitForExit() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-exit
}
