package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	exposeFieldId          string
	failTransitionId       string
	filterId               string
	host                   string
	imageFieldId           string
	inProgressTransitionId string
	password               string
	replicasFieldId        string
	successTransitionId    string
	syncInterval           int
	username               string
)

func main() {
	flag.StringVar(&exposeFieldId, "expose-field-id", "", "The expose custom field ID.")
	flag.StringVar(&failTransitionId, "fail-transition-id", "", "The failed transition ID.")
	flag.StringVar(&filterId, "filter-id", "", "The Jira search filter ID")
	flag.StringVar(&host, "host", "http://127.0.0.1:8080", "The Jira host address.")
	flag.StringVar(&imageFieldId, "image-field-id", "", "The image custom field ID.")
	flag.StringVar(&inProgressTransitionId, "in-progress-transition-id", "", "The in progress transition ID.")
	flag.StringVar(&replicasFieldId, "replicas-field-id", "", "The replicas custom field ID.")
	flag.StringVar(&successTransitionId, "success-transition-id", "", "The success transition ID.")
	flag.IntVar(&syncInterval, "sync-interval", 30, "The sync interval in seconds.")

	username = os.Getenv("JIRA_USERNAME")
	password = os.Getenv("JIRA_PASSWORD")

	flag.Parse()
	log.Println("Starting Jira Deployment Controller...")
	processIssues()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-time.After(time.Duration(syncInterval) * time.Second):
			processIssues()
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			os.Exit(0)
		}
	}
}
