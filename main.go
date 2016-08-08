package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/andygrunwald/go-jira"
)

var (
	successTransitionId    string
	failTransitionId       string
	inProgressTransitionId string
	imageFieldId           string
	replicasFieldId        string
	exposeFieldId          string
	username               string
	password               string
	host                   string
)

var jiraClient *jira.Client
var filterId string

func main() {
	// Jira transition ID mapping.
	flag.StringVar(&successTransitionId, "success-transition-id", "", "The transition ID to use when the process succeeds.")
	flag.StringVar(&failTransitionId, "fail-transition-id", "", "The transition ID to use when the process fails.")
	flag.StringVar(&inProgressTransitionId, "in-progress-transition-id", "", "The transition ID that marks an issue in progress.")

	// Jira custom field mappings.
	flag.StringVar(&imageFieldId, "image-field-id", "", "The image custom field ID.")
	flag.StringVar(&replicasFieldId, "replicas-field-id", "", "The replicas custom field ID.")
	flag.StringVar(&exposeFieldId, "expose-field-id", "", "The expose custom field ID.")

	// Jira login info.
	flag.StringVar(&host, "host", "http://127.0.0.1:8080", "The Jira host address.")
	flag.StringVar(&filterId, "filter-id", "", "The Jira filter id to search for deployment issues")

	username = os.Getenv("JIRA_USERNAME")
	password = os.Getenv("JIRA_PASSWORD")

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
	query := fmt.Sprintf("filter=%s", filterId)
	issues, resp, err := jiraClient.Issue.Search(query, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Processing %d issues.", resp.Total)

	for _, i := range issues {
		log.Println("Processing issue", i.ID)

		// Mark the issue in progress.
		_, err = jiraClient.Issue.DoTransition(i.ID, inProgressTransitionId)
		if err != nil {
			log.Println(err)
			continue
		}

		// Extract the container image name from the container image custom field.
		containerImageName := ""
		replicas := 0
		expose := false

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
				switch key {
				case imageFieldId:
					containerImageName = val.(string)
				case replicasFieldId:
					for k, v := range val.(map[string]interface{}) {
						if k == "value" {
							i, err := strconv.Atoi(v.(string))
							if err != nil {
								log.Println(err)
								continue
							}
							replicas = i
						}
					}
				case exposeFieldId:
					if val != nil {
						expose = true
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
		_, err = jiraClient.Issue.DoTransition(i.ID, successTransitionId)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Deployed container: %s replicas: %d exposed: %v successfully", containerImageName, replicas, expose)
	}
}
