package service

import (
	"api/internal/model"
	"api/internal/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	tmpPath, err := utils.CloneGitRepo(app.GitRepo, app.Branch)
	if err != nil {
		log.Printf("❌ Git clone failed: %v", err)
		app.Status = "clone_failed"
	} else {
		log.Printf("✅ Repo cloned to: %s", tmpPath)
		app.Status = "cloned"
	}
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

func GenerateManifest(input model.App) error {
	chartPath := filepath.Join("generated", input.Name)
	valuesFile := filepath.Join(chartPath, "values.yaml")

	cmd := exec.Command("helm", "template", input.Name, chartPath, "--values", valuesFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("helm error: %s", string(out))
	}

	manifestDir := filepath.Join(chartPath, "manifests")
	if err := os.MkdirAll(manifestDir, 0755); err != nil {
		return fmt.Errorf("failed to create manifests directory: %w", err)
	}

	manifestPath := filepath.Join(manifestDir, input.Name+".yaml")
	if err := os.WriteFile(manifestPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write manifest to file: %w", err)
	}

	fmt.Printf("✅ Manifest generated and saved to %s\n", manifestPath)
	return nil
}

func InstallHelmRelease(app model.App) error {
	chartPath := filepath.Join("generated", app.Name)
	valuesFile := filepath.Join(chartPath, "values.yaml")
	releaseName := app.Name
	namespace := "default"

	cmd := exec.Command("helm", "upgrade", "--install", releaseName, chartPath,
		"--namespace", namespace,
		"--values", valuesFile)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("helm install failed: %v\nOutput: %s", err, string(out))
	}

	fmt.Println("✅Helm Output:\n", string(out))
	return nil
}

func ListDeployedApps() ([]model.AppStatus, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-A", "-l", "app.kubernetes.io/instance", "-o", "custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name,STATUS:.status.phase,READY:.status.containerStatuses[0].ready,RESTARTS:.status.containerStatuses[0].restartCount,AGE:.metadata.creationTimestamp", "--no-headers")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("kubectl get pods failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	appMap := make(map[string]*model.AppStatus)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}

		ns := parts[0]
		podName := parts[1]
		status := parts[2]
		ready := parts[3]
		restarts := parts[4]
		age := parts[5]

		appName := strings.Split(podName, "-")[0]

		app, ok := appMap[appName]
		if !ok {
			app = &model.AppStatus{
				Name:      appName,
				Namespace: ns,
				Status:    status,
				Age:       age,
			}
			appMap[appName] = app
		}

		app.Pods = append(app.Pods, model.PodStatus{
			Name:     podName,
			Status:   status,
			Ready:    ready,
			Restarts: restarts,
			Age:      age,
		})
	}

	// convert to slice
	apps := []model.AppStatus{}
	for _, a := range appMap {
		apps = append(apps, *a)
	}

	return apps, nil
}
