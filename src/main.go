package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	InitLogger()
	app := tview.NewApplication().EnableMouse(true)

	// branches := []string{"main", "dev", "feature", "release"}
	dropdown := tview.NewDropDown().
		SetLabel("Select branch: ").
		SetOptions(branches, func(text string, index int) {
			log.Println("selected branch ", text, index)
		})

	mainTextTitle := tview.NewTextView().
		SetText("Main content - waiting for input...").
		SetDynamicColors(true)

	mainContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainTextTitle, 0, 1, false).
		AddItem(dropdown, 3, 0, true)

	mainContent.SetBorder(true).
		SetTitle("Main Page").
		SetTitleAlign(tview.AlignCenter)

	// switchToMain is a callback passed into dialogPage
	switchToMain := func() {
		mainContent.SetTitle(fmt.Sprintf("  %s :: %s/%s  ", selectedBranch, username, reponame))
		app.SetRoot(mainContent, true)
		app.SetFocus(dropdown)
	}

	// Start with the dialog, switch to main after input
	dialog := dialogPage(app, mainTextTitle, switchToMain, dropdown)

	if err := app.SetRoot(dialog, true).Run(); err != nil {
		panic(err)
	}
}

func dialogPage(app *tview.Application, mainTextTitle *tview.TextView, switchToMain func(), dropdown *tview.DropDown) tview.Primitive {
	errorText := tview.NewTextView().SetText("").SetTextColor(tcell.ColorRed.TrueColor())
	urlField := tview.NewInputField().
		SetLabel("Repo URL: ").
		// torvalds/linux
		SetText("https://github.com/high-horse/c-programming-practice").
		SetFieldWidth(50)

	submit := func() {
		text := urlField.GetText()
		if !validateRepoLink(text) {
			log.Println("invalid link ", text)
			errorText.SetText(fmt.Sprintf("Invalid repo link URL: %s", text))
			return
		}
		errorText.SetText("")
		url = text
		mainTextTitle.SetText("You entered " + text)

		if err := fetchContents(text); err != nil {
			errorText.SetText("Failed to fetch branches: " + err.Error())
			return
		}

		// Update dropdown options dynamically
		dropdown.SetOptions(branches, func(text string, index int) {
			log.Println("selected branch ", text, index)
		})
		for i, b := range branches {
			if b == selectedBranch {
				dropdown.SetCurrentOption(i)
				break
			}
		}

		switchToMain() // <-- swap root back to main
	}
	urlField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			submit()
		}
	})

	inputForm := tview.NewForm().
		AddFormItem(urlField).
		AddButton("OK", func() {
			submit()
		}).
		AddButton("Cancel", func() {
			app.Stop()
		})

	formWithErrTxt := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputForm, 5, 1, true).
		AddItem(errorText, 0, 1, false)

	dialogBox := tview.NewFrame(formWithErrTxt).
		SetBorders(1, 1, 1, 1, 1, 1)

	dialogBox.
		SetBorder(true).
		SetTitle(" Enter Repo URL ").
		SetTitleAlign(tview.AlignCenter)

	// Horizontal centering
	centered := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(dialogBox, 70, 1, true).
		AddItem(nil, 0, 1, false)

	// Vertical centering + height control
	dialog := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(centered, 10, 1, true).
		AddItem(nil, 0, 1, false)

	return dialog

}
