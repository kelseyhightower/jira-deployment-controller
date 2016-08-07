# Jira Deploy

WIP: Jira to Kubernetes bridge. Stay tunned!

## Usage

```
jira-deploy -h
```
```
Usage of jira-deploy:
  -approved-status-id string
    	The status ID that marks an issue approved.
  -done-status-id string
    	The status ID that marks an issue done.
  -host string
    	The Jira host address. (default "http://127.0.0.1:8080")
  -image-field-id string
    	The container image custom field ID.
  -in-progress-status-id string
    	The status ID that marks an issue in progress.
  -password string
    	The Jira login password.
  -project-id string
    	The Jira project ID used for Kubernetes deployments.
  -username string
    	The Jira login username.

```

### Example

```
jira-deploy \
  -host http://127.0.0.1:8080 \
  -done-status-id 31 \
  -approved-status-id 11 \
  -in-progress-status-id 21 \
  -image-field-id customfield_10100 \
  -project-id KUBE \
  -username kelseyhightower \
  -password <password>
```
