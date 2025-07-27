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
