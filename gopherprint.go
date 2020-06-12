package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/mcuadros/go-octoprint"
)

var serverURL = "http://localhost:8000"
var apiKey = ""

func getMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(
		items,
		menuet.MenuItem{
			Text: "Open OctoPrint",
			Clicked: func() {
				exec.Command("open", serverURL).Start()
			},
		},
	)
	return items
}

func updateMenubarTitle(client *octoprint.Client) {
	progressStr := ""

	stateReq := octoprint.StateRequest{}
	state, stateErr := stateReq.Do(client)
	if stateErr != nil {
		log.Fatalf("error requesting octoprint state: %s", stateErr)
	}

	if state.State.Flags.Printing {
		jobReq := octoprint.JobRequest{}
		job, jobErr := jobReq.Do(client)
		if jobErr != nil {
			log.Fatalf("error requesting octoprint job: %s", jobErr)
		}

		log.Print("Octoprint Job:", job.Progress, job.Job)
		progressStr = fmt.Sprintf(` - %.2f %%`, job.Progress.Completion)
	}

	menuet.App().SetMenuState(&menuet.MenuState{
		Title: fmt.Sprintf("🐙 %s%s", state.State.Text, progressStr),
	})
}

func updateApp(client *octoprint.Client) {
	for {
		updateMenubarTitle(client)
		time.Sleep(time.Second * 1)
	}
}

func main() {
	println(serverURL)
	client := octoprint.NewClient(serverURL, apiKey)

	go updateApp(client)

	menuet.App().Label = "com.github.thibmaek.gopherprint"
	menuet.App().Children = getMenuItems
	menuet.App().RunApplication()
}
