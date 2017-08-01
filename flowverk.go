package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nsf/termbox-go"
	"gopkg.in/yaml.v2"
)

// SECTION: constants {{{
var fields = []string{"id", "summary", "assignee"}

const confFileName = ".flowverk.conf"

//}}}

type query struct {
	Jql    string   `json:"jql"`
	Fields []string `json:"fields"`
}

// type jUser struct {
// 	Name        string `json:"name"`
// 	DisplayName string `jsong:"displayName"`
// }
type issueFields struct {
	Summary  string      `json:"summary"`
	Assignee interface{} `json:"assignee"`
}

// Issue structure represents narrowed issue from jira response */
type Issue struct {
	Expand string      `json:"-"`
	ID     string      `json:"id"`
	Self   string      `json:"self"`
	Key    string      `json:"key"`
	Fields issueFields `json:"fields"`
}

type jiraResp struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

//Config struct
type Config struct {
	JiraURL     string `yaml:"jiraURL"`
	ProjectName string `yaml:"projectName"`
	User        string
	Pass        string
}

func main() {

	/* Read configuration */
	var config Config
	confFile, err := ioutil.ReadFile(confFileName)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(confFile, &config)

	/* initialise read input tool */
	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	/* Build jira query */
	jql := fmt.Sprintf("status=\"To Do\" AND project=\"%s\"", config.ProjectName)

	question := query{Jql: jql, Fields: fields}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(question)

	connector := &http.Client{}

	/* build request to jira */
	req, err := http.NewRequest("POST", config.JiraURL+"search", body)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(config.User, config.Pass)

	/* fetch response from jira */
	resp, err := connector.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	var jResp jiraResp
	rBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rBody, &jResp)

	/* Draw all issues */
	var assignee map[string]interface{}
	screenSize := len(jResp.Issues)
	pointerIndex := 0
	for {
		for index, issue := range jResp.Issues {
			pointer := " "
			assignee = map[string]interface{}{"displayName": "nil"}
			if issue.Fields.Assignee != nil {
				assignee = issue.Fields.Assignee.(map[string]interface{})
			}

			if pointerIndex == index {
				pointer = "*"
			}
			fmt.Printf("[%s] %s %s <%s>\n", pointer, issue.Key, issue.Fields.Summary, assignee["displayName"])
		}

		fmt.Println("-------------")
		fmt.Println(pointerIndex, screenSize)
		fmt.Println("-------------")
		action := termbox.PollEvent()
		if (action.Key == termbox.KeyArrowDown) && (pointerIndex <= screenSize) {
			pointerIndex++
		}
		if action.Key == termbox.KeyArrowUp && (0 <= pointerIndex) {
			pointerIndex--
		}

		if action.Key == termbox.KeyEnter {
			break
		}
	}

	/* close termbox */
	defer termbox.Close()
}
