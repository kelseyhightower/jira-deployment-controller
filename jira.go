package main

import (
	"fmt"
	"log"
	"strconv"

	jira "github.com/andygrunwald/go-jira"
)

type customFields struct {
	image    string
	replicas int
	expose   bool
	name     string
}

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

		cf, err := getCustomFields(jiraClient, issue.ID)

		// Do the deployment.
		err = syncDeployment(cf.name, cf.image, cf.replicas)
		if err != nil {
			log.Println(err)

			// Mark the issue failed.
			_, err := jiraClient.Issue.DoTransition(issue.ID, failTransitionId)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		message := fmt.Sprintf("Deployed image: %s replicas: %d exposed: %v successfully.", cf.image, cf.replicas, cf.expose)
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

func getCustomFields(client *jira.Client, issueID string) (customFields, error) {
	var cf customFields

	req, _ := client.NewRequest("GET", "rest/api/2/issue/"+issueID, nil)
	data := new(map[string]interface{})
	_, err := client.Do(req, data)
	if err != nil {
		return cf, err
	}

	d := *data
	f := d["fields"]
	if fields, ok := f.(map[string]interface{}); ok {
		for field, value := range fields {
			switch field {
			case nameFieldId:
				cf.name = value.(string)
			case imageFieldId:
				cf.image = value.(string)
			case replicasFieldId:
				for k, v := range value.(map[string]interface{}) {
					if k == "value" {
						cf.replicas, err = strconv.Atoi(v.(string))
						if err != nil {
							return cf, err
						}
					}
				}
			case exposeFieldId:
				if value != nil {
					cf.expose = true
				}
			}
		}
	}
	return cf, nil
}
