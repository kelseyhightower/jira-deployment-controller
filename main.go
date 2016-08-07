package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/andygrunwald/go-jira"
)

var (
	approvedStatusId   string
	doneStatusId       string
	imageFieldId       string
	inProgressStatusId string
	projectId          string
	username           string
	password           string
	host               string
)

var jiraClient *jira.Client

func main() {
	flag.StringVar(&approvedStatusId, "approved-status-id", "", "The status ID that marks an issue approved.")
	flag.StringVar(&doneStatusId, "done-status-id", "", "The status ID that marks an issue done.")
	flag.StringVar(&imageFieldId, "image-field-id", "", "The container image custom field ID.")
	flag.StringVar(&inProgressStatusId, "in-progress-status-id", "", "The status ID that marks an issue in progress.")
	flag.StringVar(&projectId, "project-id", "", "The Jira project ID used for Kubernetes deployments.")
	flag.StringVar(&username, "username", "", "The Jira login username.")
	flag.StringVar(&password, "password", "", "The Jira login password.")
	flag.StringVar(&host, "host", "http://127.0.0.1:8080", "The Jira host address.")

	flag.Parse()

	var err error
	jiraClient, err = jira.NewClient(nil, host)
	if err != nil {
		log.Fatal(err)
	}
	_, err = jiraClient.Authentication.AcquireSessionCookie(username, password)
	if err != nil {
		log.Fatal(err)
	}

	processIssues()

	for {
		select {
		case <-time.After(30 * time.Second):
			processIssues()
		}
	}
}

func processIssues() {
	query := fmt.Sprintf("project=%s AND status=\"%s\"", projectId, "To Do")
	issues, resp, err := jiraClient.Issue.Search(query, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Processing %d issues.", resp.Total)

	for _, i := range issues {
		log.Println("Processing issue", i.ID)

		// Mark the issue in progress.
		_, err = jiraClient.Issue.DoTransition(i.ID, inProgressStatusId)
		if err != nil {
			log.Println(err)
			continue
		}

		// Extract the container image name from the container image custom field.
		containerImageName := ""

		req, _ := jiraClient.NewRequest("GET", "rest/api/2/issue/"+i.ID, nil)
		issue := new(map[string]interface{})
		_, err = jiraClient.Do(req, issue)
		if err != nil {
			log.Fatal(err)
		}

		m := *issue
		f := m["fields"]

		if rec, ok := f.(map[string]interface{}); ok {
			for key, val := range rec {
				if key == imageFieldId {
					if vmap, ok := val.(map[string]interface{}); ok {
						containerImageName = vmap["value"].(string)
						break
					}
				}
			}
		}

		// Do the deployment.
		message := fmt.Sprintf("Deployed container %s", containerImageName)
		_, _, err = jiraClient.Issue.AddComment(i.ID, &jira.Comment{Body: message})
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = jiraClient.Issue.DoTransition(i.ID, doneStatusId)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Deployed container %s successfully", containerImageName)
	}
}
