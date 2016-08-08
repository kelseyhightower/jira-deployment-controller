# jira-deployment-controller

WIP: Jira to Kubernetes bridge. Stay tunned!

## Usage

```
jira-deployment-controller -h
```
```
Usage of jira-deployment-controller:
  -expose-field-id string
    	The expose custom field ID.
  -fail-transition-id string
    	The transition ID to use when the process fails.
  -filter-id string
    	The Jira filter id to search for deployment issues
  -host string
    	The Jira host address. (default "http://127.0.0.1:8080")
  -image-field-id string
    	The image custom field ID.
  -in-progress-transition-id string
    	The transition ID that marks an issue in progress.
  -replicas-field-id string
    	The replicas custom field ID.
  -success-transition-id string
    	The transition ID to use when the process succeeds.
```

### Example

```
export JIRA_USERNAME="ninja"
export JIRA_PASSWORD=""
```

```
jira-deployment-controller \
  -host http://127.0.0.1:8080 \
  -in-progress-transition-id 21 \
  -success-transition-id 31 \
  -fail-transition-id 41 \
  -image-field-id customfield_10103 \
  -replicas-field-id customfield_10102 \
  -expose-field-id customfield_10101 \
  -filter-id deployments
```

```
2016/08/07 11:54:39 Processing 1 issues.
2016/08/07 11:54:39 Processing issue 10000
2016/08/07 11:54:39 Deployed container: nginx:1.10 replicas: 2 exposed: true successfully
```

## Kubernetes Deployment Example

### Create Jira Secrets

```
kubectl create secret generic jira --from-literal=password=XXXXXXXX
```

### Create Jira Deployment Controller Deployment

```
kubectl create -f deployments/jira-deployment-controller.yaml
```
