package main

type Deployment struct {
	ApiVersion string         `json:"apiVersion,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Metadata   Metadata       `json:"metadata"`
	Spec       DeploymentSpec `json:"spec"`
}

type Service struct {
	ApiVersion string      `json:"apiVersion,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	Metadata   Metadata    `json:"metadata"`
	Spec       ServiceSpec `json:"spec"`
}

type Metadata struct {
	Name            string            `json:"name"`
	GenerateName    string            `json:"generateName"`
	ResourceVersion string            `json:"resourceVersion"`
	Labels          map[string]string `json:"labels"`
	Annotations     map[string]string `json:"annotations"`
	Uid             string            `json:"uid"`
}

type DeploymentSpec struct {
	Replicas int64       `json:"replicas"`
	Template PodTemplate `json:"template"`
}

type ServiceSpec struct {
	ClusterIP string `json:"clusterIP"`
	Type      string `json:"type"`
	Ports     []Port `json:"ports"`
}

type Port struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int64  `json:"port"`
	TargetPort int64  `json:"targetPort"`
}

type PodTemplate struct {
	Metadata Metadata `json:"metadata"`
	Spec     PodSpec  `json:"spec"`
}

type PodSpec struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}
