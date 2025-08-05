package model

type App struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	GitRepo string `json:"git_repo"`
	Branch  string `json:"branch"`
	Domain  string `json:"domain"`
	Port    int    `json:"port"`
	Image   string `json:"image"`
	Status  string `json:"status"`
}

type PodStatus struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Ready    string `json:"ready"`
	Restarts string `json:"restarts"`
	Age      string `json:"age"`
}

type AppStatus struct {
	Name      string      `json:"name"`
	Namespace string      `json:"namespace"`
	Status    string      `json:"status"`
	Age       string      `json:"age"`
	Pods      []PodStatus `json:"pods"`
}
