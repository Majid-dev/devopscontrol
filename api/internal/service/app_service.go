package service

import (
	"api/internal/model"
	"api/internal/util"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type AppService struct {
	apps map[string]model.App
	mu   sync.RWMutex
}

func NewAppService() *AppService {
	return &AppService{
		apps: make(map[string]model.App),
	}
}

func (s *AppService) CreateApp(app model.App) model.App {
	s.mu.Lock()
	defer s.mu.Unlock()

	app.ID = uuid.New().String()
	app.Status = "created"
	// Try to clone the repo
	tmpPath, err := util.CloneGitRepo(app.GitRepo, app.Branch)
	if err != nil {
		log.Printf("❌ Git clone failed: %v", err)
		app.Status = "clone_failed"
	} else {
		log.Printf("✅ Repo cloned to: %s", tmpPath)
		app.Status = "cloned"
	}

	val := util.HelmValues{
		AppName: app.Name,
		Domain:  app.Domain,
	}
	val.Image.Repository, val.Image.Tag = parseImage(app.Image)
	val.Image.Port = app.Port

	err = util.RenderHelmAndSave(val, "./chart", "./tmp/"+app.Name)
	s.apps[app.ID] = app

	return app
}

func parseImage(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return image, "latest"
}

func (s *AppService) GetApp(id string) (model.App, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	app, exists := s.apps[id]
	return app, exists
}

func (s *AppService) ListApps() []model.App {
	s.mu.RLock()
	defer s.mu.RUnlock()

	apps := make([]model.App, 0, len(s.apps))
	for _, app := range s.apps {
		apps = append(apps, app)
	}
	return apps
}
