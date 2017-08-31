package jira

import (
	"bytes"
	"encoding/json"
	"flowverk/configuration"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Constants
var fields = []string{"id", "summary", "assignee"}

type issueFields struct {
	Summary  string      `json:"summary"`
	Assignee interface{} `json:"assignee"`
}

type jiraResp struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

type boardResp struct {
	Values []struct {
		ID   int `json:"id"`
		Name string
	}
}

type query struct {
	Jql    string   `json:"jql"`
	Fields []string `json:"fields"`
}

// Issue structure represents narrowed issue from jira response */
type Issue struct {
	Expand string      `json:"-"`
	ID     string      `json:"id"`
	Self   string      `json:"self"`
	Key    string      `json:"key"`
	Fields issueFields `json:"fields"`
}

/*Jira suuuper*/
type Jira struct {
	url        string
	BoardURL   string
	Connection *http.Client
	Project    string
	config     configuration.Config
}

// NewConnection creates and setup connection to jira
func NewConnection(config *configuration.Config) *Jira {

	jira := new(Jira)
	jira.url = config.JiraURL + config.IssuesURI
	jira.BoardURL = config.JiraURL + config.BoardURI + strconv.Itoa(config.Board)
	jira.Connection = &http.Client{}
	jira.Project = config.ProjectName
	jira.config = *config

	return jira
}

func (jira *Jira) getCurrentSprint() (int, error) {
	req, err := http.NewRequest("GET", jira.BoardURL+"/sprint?state=active", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(jira.config.User, jira.config.Pass)

	/* fetch response from jira */
	resp, err := jira.Connection.Do(req)
	if err != nil {
		rBody, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("Response status: %s\nFull resp: %s", resp.Status, rBody)
	}
	if 300 <= resp.StatusCode {
		rBody, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("Response status: %s\nFull resp: %s", resp.Status, rBody)
	}

	var bResp boardResp
	rBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rBody, &bResp)
	return bResp.Values[0].ID, nil
}

//GetIssues return issues of given status from project
func (jira *Jira) GetIssues(status string) ([]Issue, error) {
	sprint, err := jira.getCurrentSprint()
	if err != nil {
		return nil, err
	}

	jql := fmt.Sprintf("status=\"%s\" AND project=\"%s\" AND Sprint=%d",
		status, jira.Project, sprint)

	question := query{Jql: jql, Fields: fields}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(question)
	req, err := http.NewRequest("POST", jira.url+"search", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(jira.config.User, jira.config.Pass)

	/* fetch response from jira */
	resp, err := jira.Connection.Do(req)
	if err != nil {
		return nil, err
	}

	if 300 <= resp.StatusCode {
		rBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Response status: %s\nFull resp: %s", resp.Status, rBody)
	}

	var jResp jiraResp
	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(rBody, &jResp)

	return jResp.Issues, nil
}
