package service

import (
	"devopscontrol-api/internal/model"
	"devopscontrol-api/internal/util"
	"log"
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
	s.apps[app.ID] = app

	return app
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
