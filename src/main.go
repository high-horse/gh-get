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
		// SetText("Main content - waiting for input...").
		SetDynamicColors(true)

	tree := tview.NewTreeView()
	tree.SetBorder(false).SetBorderPadding(0, 0, 5 ,0)
	tree.SetTitle("Repository Contents")

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}

		c := ref.(*Content)

		if !c.IsDir {
			log.Println("Download:", c.DownloadUrl)
			return
		}

		if c.Fetched {
			node.SetExpanded(!node.IsExpanded())
			return
		}

		children, err := fetchContentAtPath(owner, reponame, selectedBranch, c.Path)
		if err != nil {
			log.Println("error:", err)
			return
		}

		c.Children = children
		c.Fetched = true

		node.ClearChildren()
		buildTree(node, c.Children)
		node.SetExpanded(true)
	})
	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == ' ' {

			node := tree.GetCurrentNode()
			if node == nil {
				return nil
			}

			ref := node.GetReference()
			if ref == nil {
				return nil
			}

			c := ref.(*Content)

			if c.IsDir {
				toggleRecursive(c, !c.Selected)
			} else {
				c.Selected = !c.Selected
			}
			// update label
			//
			node.SetText(formatLabel(c))

			return nil // consume event
		}
		return event
	})

	mainContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		// AddItem(mainTextTitle, 0, 1, false).
		AddItem(dropdown, 3, 0, false).
		AddItem(tree, 0, 5, true)

	mainContent.SetBorder(true).
		SetTitle("Main Page").
		SetTitleAlign(tview.AlignCenter)

	// switchToMain is a callback passed into dialogPage
	switchToMain := func() {
		mainContent.SetTitle(fmt.Sprintf("  %s :: %s/%s  ", selectedBranch, owner, reponame))

		rootContents, err := fetchContentAtPath(owner, reponame, selectedBranch, "")
		if err != nil {
			log.Println("error fetching root:", err)
			return
		}

		rootNode := tview.NewTreeNode(fmt.Sprintf("%s/%s", owner, reponame)).
			SetColor(tcell.ColorGreen)

		buildTree(rootNode, rootContents)

		tree.SetRoot(rootNode).SetCurrentNode(rootNode)

		app.SetRoot(mainContent, true)
		// app.QueueUpdateDraw(func() {
		// 	app.SetFocus(tree)
		// })
		go func() {
			app.QueueUpdateDraw(func() {
				app.SetFocus(tree)
			})
		}()
		// app.SetFocus(tree)
	}

	// Start with the dialog, switch to main after input
	dialog := dialogPage(app, mainTextTitle, switchToMain, dropdown)

	if err := app.SetRoot(dialog, true).Run(); err != nil {
		panic(err)
	}
}

func buildTree(parent *tview.TreeNode, contents []Content) {
	for i := range contents {
		c := &contents[i]

		label := formatLabel(c)

		node := tview.NewTreeNode(label).
			SetReference(c)

		if c.IsDir {
			node.SetColor(tcell.ColorBlue)
			// node.AddChild(tview.NewTreeNode("loading..."))
		} else {
			node.SetColor(tcell.ColorWhite)
		}

		parent.AddChild(node)
	}
}

func formatLabel(c *Content) string {
	check := "[ ]"
	if c.Selected {
		check = "[ ✓ ]"
	}

	// escape brackets for tview
	check = tview.Escape(check)

	if c.IsDir {
		return fmt.Sprintf("%s 📁 %s", check, c.Name)
	}
	return fmt.Sprintf("%s %s", check, c.Name)
}

func collectSelected(contents []Content, result *[]Content) {
	for _, c := range contents {
		if c.IsDir {
			collectSelected(c.Children, result)
		} else if c.Selected {
			*result = append(*result, c)
		}
	}
}

func toggleRecursive(c *Content, selected bool) {
	c.Selected = selected
	for i := range c.Children {
		toggleRecursive(&c.Children[i], selected)
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
		// mainTextTitle.SetText("You entered " + text)

		if err := fetchContents(text); err != nil {
			log.Println()
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
