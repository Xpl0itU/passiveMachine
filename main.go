package main

import (
	"github.com/atotto/clipboard"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication().
		EnablePaste(true).
		EnableMouse(true)

	list := tview.NewList()

	form := tview.NewForm()

	frame := tview.NewFrame(list).
		AddText("Passive Machine", true, tview.AlignCenter, tview.Styles.PrimaryTextColor).
		AddText("Tip: Use the Register button to sign up for the service", true, tview.AlignCenter, tview.Styles.ContrastSecondaryTextColor)

	menuItems := []MenuItem{
		{"Honeygain", "Earn passive income by sharing your internet connection", &HoneygainConfig{}},
		{"EarnApp", "Earn passive income by sharing your internet connection", &EarnAppConfig{}},
		{"PawnsApp", "Earn passive income by sharing your internet connection", &PawnsAppConfig{}},
		{"PacketStream", "Earn passive income by sharing your internet connection", &PacketStreamConfig{}},
		{"Grass", "Earn passive income by sharing your internet connection", &GrassConfig{}},
		{"Mysterium", "Earn passive income by sharing your internet connection", &MystConfig{}},
		{"Peer2Profit", "Earn passive income by sharing your internet connection", &Peer2ProfitConfig{}},
		{"Watchtower (Automatic Updates) (Recommended)", "Automatically update your docker containers", &WatchtowerConfig{Configured: true}},
	}

	for _, item := range menuItems {
		list.AddItem(item.GetName(), item.GetDescription(), 0, func() {
			form.Clear(true)
			item.Config.ConfigureForm(form, frame, app)
			app.SetRoot(form, true)
		})
	}

	list.AddItem("[*] Build Docker Compose File", "Build a docker-compose file from the selected items (Advanced)", 0, func() {
		dockerComposeFile := buildDockerComposeFile(menuItems)
		form.Clear(true)
		form.AddTextView("Docker Compose File", dockerComposeFile, 0, 0, true, true)
		form.AddButton("Copy to Clipboard", func() {
			if err := clipboard.WriteAll(dockerComposeFile); err != nil {
				form.AddTextView("Error", "Failed to copy to clipboard", 0, 1, true, false)
			} else {
				form.AddTextView("Success", "Copied to clipboard", 0, 1, true, false)
			}
		})
		form.AddButton("Return", func() {
			returnToMenu(frame, app)
		})
		app.SetRoot(form, true)
	})

	list.AddItem("[*] Create Docker Containers", "Create docker containers from the selected items (Recommended)", 0, func() {
		logView := tview.NewTextView()
		form.Clear(true)
		form.AddFormItem(logView)
		errors := batchCreateDockerContainers(menuItems, logView)
		if len(errors) == 0 {
			form.AddTextView("Success", "All containers created successfully", 0, 0, true, false)
		} else {
			form.AddTextView("Errors", "Some containers failed to create", 0, 0, true, false)
			for _, err := range errors {
				form.AddTextView("Error", err.Error(), 0, 1, true, false)
			}
		}
		form.AddButton("Return", func() {
			returnToMenu(frame, app)
		})
		app.SetRoot(form, true)
	})
	app.SetRoot(frame, true).Run()
}
