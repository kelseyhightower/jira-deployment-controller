FROM scratch
ADD jira-deployment-controller /jira-deployment-controller
ENTRYPOINT ["/jira-deployment-controller"]
