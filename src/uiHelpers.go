package main

import (
	"fmt"
	
	"github.com/rivo/tview"
)

func getSelectionState(c *Content) (selected, total int) {
	if !c.IsDir {
		if c.Selected {
			return 1, 1
		}
		return 0, 1
	}

	for i := range c.Children {
		s, t := getSelectionState(&c.Children[i])
		selected += s
		total += t
	}

	return
}

func _formatLabel_(c *Content) string {
	check := "[ ]"

	if c.IsDir {
		selected, total := getSelectionState(c)

		switch {
		case selected == 0:
			check = "[ ]"
		case selected == total:
			check = "[ ✓ ]"
		default:
			check = "[ - ]"
		}

		check = tview.Escape(check)
		return fmt.Sprintf("%s 📁 %s (%d/%d)", check, c.Name, selected, total)
	}

	if c.Selected {
		check = "[ ✓ ]"
	}

	check = tview.Escape(check)
	return fmt.Sprintf("%s %s", check, c.Name)
}


// func updateParents(node *tview.TreeNode) {
// 	parent := node.GetParent()
// 	if parent == nil {
// 		return
// 	}

// 	ref := parent.GetReference()
// 	if ref != nil {
// 		c := ref.(*Content)
// 		parent.SetText(formatLabel(c))
// 	}

// 	updateParents(parent)
// }

func formatLabel(c *Content) string {
	check := "[ ]"
	if c.IsDir {
		selectedCount := 0
		for _, child := range c.Children {
			if child.Selected {
				selectedCount++
			}
		}
		if selectedCount == len(c.Children) && len(c.Children) > 0 {
			check = "[ ✓ ]"
		} else if selectedCount > 0 {
			check = "[ - ]" // partial
		}
	} else if c.Selected {
		check = "[ ✓ ]"
	}

	check = tview.Escape(check)

	if c.IsDir {
		return fmt.Sprintf("%s 📁 %s", check, c.Name)
	}
	return fmt.Sprintf("%s %s", check, c.Name)
}

func formatLabel_(c *Content) string {
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