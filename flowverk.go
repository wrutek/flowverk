package main

import (
	"flowverk/configuration"
	"flowverk/jira"
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

// SECTION: constants {{{

const confFileName = ".flowverk.yaml"

//}}}

func main() {

	/* Read configuration */
	config := configuration.GetConfig()

	/* initialise read input tool */
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	connector := jira.NewConnection(config)
	issues, err := connector.GetIssues("To Do")
	if err != nil {
		closeTermbox(0)

		fmt.Println(err)
		os.Exit(1)
	}

	/* Draw all issues */
	var assignee map[string]interface{}
	var parent map[string]interface{}
	screenSize := len(issues)
	termWidth, _ := termbox.Size()
	pointerIndex := 0

	for {
		var pointer string
		var pointedIssue jira.Issue

		/* draw a screen */
		for index, issue := range issues {
			assignee = map[string]interface{}{"displayName": "nil"}
			if issue.Fields.Assignee != nil {
				assignee = issue.Fields.Assignee.(map[string]interface{})
			}

			parentInfo := ""
			if issue.Fields.Parent != nil {
				parent = issue.Fields.Parent.(map[string]interface{})
				parentInfo = fmt.Sprintf("<%s> %s --> ", parent["key"], parent["fields"].(map[string]interface{})["summary"])
			}

			if pointerIndex == index {
				pointer = "*"
				pointedIssue = issue
			} else {
				pointer = " "
			}
			fmt.Printf("[%s] %s<%s> %s <%s>\n", pointer, parentInfo, issue.Key, issue.Fields.Summary, assignee["displayName"])
		}

		fmt.Printf("\n\n%s\n", strings.Repeat("-", termWidth))
		fmt.Printf("press Enter to assign: %s to yourself and put this ticket to InProgress\n", pointedIssue.Key)
		fmt.Println("q/ESC - exit")
		fmt.Println("Enter - select issue")
		fmt.Printf("%s\n", strings.Repeat("-", termWidth))
		termbox.HideCursor()
		termbox.Flush()

		/* wait for action key */
		action := termbox.PollEvent()
		if (action.Key == termbox.KeyArrowDown) && (pointerIndex < screenSize-1) {
			pointerIndex++
		}
		if action.Key == termbox.KeyArrowUp && (0 < pointerIndex) {
			pointerIndex--
		}

		if action.Key == termbox.KeyEnter {
			break
		}

		if action.Key == termbox.KeyEsc || action.Ch == 'q' {
			break
		}

		termbox.SetCursor(0, 0)
		termbox.Flush()
	}

	closeTermbox(screenSize)
}

func closeTermbox(screenSize int) {
	/* close termbox */
	termbox.SetCursor(0, screenSize+7)
	defer termbox.Close()
	termbox.Flush()
}

//
// func assignTicket(issue Issue) {
//
// }
