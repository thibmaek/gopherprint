package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/haklop/gnotifier"
	"github.com/mcuadros/go-octoprint"
)

var serverURL = "http://localhost:8000"
var apiKey = ""

const (
	displayName = "Gopherprint"
	bundleID    = "com.github.thibmaek.gopherprint"
)

const (
	stateKey         = "DISPLAY_VALUE"
	temperatureValue = "TEMPERATURE"
	progressValue    = "PROGRESS"
)

var client octoprint.Client

func sendNotification(title string, message string) {
	notification := gnotifier.Notification(title, message)
	notification.GetConfig().Expiration = 2000
	notification.GetConfig().ApplicationName = displayName
	notification.Push()
}

func getDisplayString(value string) string {
	return fmt.Sprintf("🐙 %s", value)
}

func getBaseMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(
		items,
		menuet.MenuItem{
			Text: fmt.Sprintf("Open OctoPrint (%s)", serverURL),
			Clicked: func() {
				exec.Command("open", serverURL).Start()
			},
		},
		menuet.MenuItem{Type: "separator"},
		menuet.MenuItem{
			Text: "Display",
			Children: func() []menuet.MenuItem {
				return []menuet.MenuItem{
					menuet.MenuItem{
						Text:  "Job progress",
						State: menuet.Defaults().String(stateKey) == progressValue,
						Clicked: func() {
							menuet.Defaults().SetString(stateKey, progressValue)
						},
					},
					menuet.MenuItem{
						Text:  "Temperature",
						State: menuet.Defaults().String(stateKey) == temperatureValue,
						Clicked: func() {
							menuet.Defaults().SetString(stateKey, temperatureValue)
						},
					},
				}
			},
		},
	)

	return items
}

func handleUpdatePrinterState(client *octoprint.Client) {
	req := octoprint.StateRequest{}
	state, err := req.Do(client)
	if err != nil {
		if err.Error() == "Printer is not operational" {
			menuet.App().SetMenuState(&menuet.MenuState{
				Title: getDisplayString("Printer not available"),
			})
			return
		}

		log.Fatalf("error requesting octoprint state: %s", err)
	}

	progressStr := ""

	menuItems := getBaseMenuItems()
	menubarValue := menuet.Defaults().String(stateKey)

	if state.State.Flags.Printing {
		req := octoprint.JobRequest{}
		job, err := req.Do(client)
		if err != nil {
			log.Fatalf("error requesting octoprint job: %s", err)
		}

		jobFinished := job.Progress.Completion >= 100
		if jobFinished {
			sendNotification(
				"Finished printing!",
				fmt.Sprintf("%s finished in %.2f minutes", job.Job.File.Name, job.Progress.PrintTime/30),
			)
		}

		progressStr = fmt.Sprintf(` - %.2f %%`, job.Progress.Completion)
		menuItems = append(
			menuItems,
			menuet.MenuItem{
				Text: fmt.Sprintf("File: %s", job.Job.File.Name),
			},
		)
	}

	menuet.App().Children = func() []menuet.MenuItem { return menuItems }

	switch menubarValue {
	case temperatureValue:
		menubarValue = getDisplayString(
			fmt.Sprintf(
				"Bed: %d / Tool: %d",
				int64(state.Temperature.Current["bed"].Actual),
				int64(state.Temperature.Current["tool0"].Actual),
			),
		)
	case progressValue:
		menubarValue = getDisplayString(fmt.Sprintf("%s %s", state.State.Text, progressStr))
	default:
		menubarValue = getDisplayString(fmt.Sprintf("%s", state.State.Text))
	}

	menuet.App().SetMenuState(&menuet.MenuState{
		Title: menubarValue,
	})
}

func updateApp(client *octoprint.Client) {
	for {
		handleUpdatePrinterState(client)
		time.Sleep(time.Second * 1)
	}
}

func main() {
	octoprintClient := octoprint.NewClient(serverURL, apiKey)

	go updateApp(octoprintClient)

	menuet.App().Label = bundleID
	menuet.App().RunApplication()
}
