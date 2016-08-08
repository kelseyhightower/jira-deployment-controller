package main

import (
	"fmt"
	"log"
	"strconv"

	jira "github.com/andygrunwald/go-jira"
)

func processIssues() {
	jiraClient, err := jira.NewClient(nil, host)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = jiraClient.Authentication.AcquireSessionCookie(username, password)
	if err != nil {
		log.Println(err)
		return
	}

	// Search for deployment issues using a Jira filter.
	query := fmt.Sprintf("filter=%s", filterId)
	issues, resp, err := jiraClient.Issue.Search(query, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.Total == 0 {
		return
	}

	log.Printf("Processing %d issues.", resp.Total)
	for _, issue := range issues {
		log.Println("Processing issue", issue.ID)

		// Mark the issue in progress.
		_, err = jiraClient.Issue.DoTransition(issue.ID, inProgressTransitionId)
		if err != nil {
			log.Println(err)
			continue
		}

		// Extract deployment info from the Jira custom fields.
		image := ""
		replicas := 0
		expose := false

		req, _ := jiraClient.NewRequest("GET", "rest/api/2/issue/"+issue.ID, nil)
		jiraIssue := new(map[string]interface{})
		_, err = jiraClient.Do(req, issue)
		if err != nil {
			log.Fatal(err)
		}

		m := *jiraIssue
		f := m["fields"]

		if rec, ok := f.(map[string]interface{}); ok {
			for key, val := range rec {
				switch key {
				case imageFieldId:
					image = val.(string)
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
		err := syncDeployment("adhoc", image, replicas)
		if err != nil {
			log.Println(err)

			// Mark the issue failed.
			_, err := jiraClient.Issue.DoTransition(issue.ID, failTransitionId)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		message := fmt.Sprintf("Deployed image: %s replicas: %d exposed: %v successfully.", image, replicas, expose)
		log.Println(message)

		_, _, err = jiraClient.Issue.AddComment(issue.ID, &jira.Comment{Body: message})
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = jiraClient.Issue.DoTransition(issue.ID, successTransitionId)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
