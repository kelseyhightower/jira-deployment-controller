package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	apiHost             = "127.0.0.1:8001"
	deploymentsEndpoint = "/apis/extensions/v1beta1/namespaces/default/deployments"
	servicesEndpoint    = "/api/v1/namespaces/default/services"
)

func syncDeployment(name, image string, replicas int) error {
	deployment, err := getDeployment(name)
	if err == ErrNotExist {
		return createDeployment(name, image, replicas)
	}
	if err != nil {
		return err
	}

	deployment.Spec.Replicas = int64(replicas)
	deployment.Spec.Template.Spec.Containers[0].Image = image

	b := make([]byte, 0)
	body := bytes.NewBuffer(b)
	err = json.NewEncoder(body).Encode(deployment)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s%s/%s", apiHost, deploymentsEndpoint, name)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Updating deployment failed:" + resp.Status)
	}
	return nil
}

var ErrNotExist = errors.New("does not exist")

func getDeployment(name string) (*Deployment, error) {
	var deployment Deployment

	path := fmt.Sprintf("%s/%s", deploymentsEndpoint, name)

	request := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:   apiHost,
			Path:   path,
			Scheme: "http",
		},
	}

	request.Header.Set("Accept", "application/json, */*")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrNotExist
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Get deployment error non 200 reponse: " + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&deployment)
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

func createDeployment(name, image string, replicas int) error {
	labels := make(map[string]string)
	labels["run"] = name

	containers := []Container{
		Container{Image: image, Name: name},
	}

	deployment := Deployment{
		ApiVersion: "extensions/v1beta1",
		Kind:       "Deployment",
		Metadata:   Metadata{Name: name},
		Spec: DeploymentSpec{
			Replicas: int64(replicas),
			Template: PodTemplate{
				Metadata: Metadata{
					Labels: labels,
				},
				Spec: PodSpec{
					Containers: containers,
				},
			},
		},
	}

	var b []byte
	body := bytes.NewBuffer(b)
	err := json.NewEncoder(body).Encode(deployment)
	if err != nil {
		return err
	}

	request := &http.Request{
		Body:          ioutil.NopCloser(body),
		ContentLength: int64(body.Len()),
		Header:        make(http.Header),
		Method:        http.MethodPost,
		URL: &url.URL{
			Host:   apiHost,
			Path:   deploymentsEndpoint,
			Scheme: "http",
		},
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		log.Println(string(data))
		return errors.New("Deployment: Unexpected HTTP status code" + resp.Status)
	}
	return nil
}
