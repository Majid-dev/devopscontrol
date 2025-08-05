package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"api/internal/model"
)

func GenerateHelmFiles(app model.App) error {
	// export tag of image name
	imageTag := "latest"
	if parts := strings.Split(app.Image, ":"); len(parts) == 2 {
		imageTag = parts[1]
	}

	// file pathes
	basePath := filepath.Join("generated", app.Name)
	if err := os.MkdirAll(filepath.Join(basePath, "templates"), 0755); err != nil {
		return err
	}

	// data for template
	data := map[string]interface{}{
		"Name":       app.Name,
		"Image":      app.Image,
		"Port":       app.Port,
		"Domain":     app.Domain,
		"AppVersion": imageTag,
	}

	// input and output template files
	files := []struct {
		TemplatePath string
		OutputPath   string
	}{
		{"internal/helm/values.yaml", filepath.Join(basePath, "values.yaml")},
		{"internal/helm/Chart.yaml", filepath.Join(basePath, "Chart.yaml")},
	}

	for _, file := range files {
		tpl, err := template.ParseFiles(file.TemplatePath)
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		if err := os.WriteFile(file.OutputPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

func CopyTemplates(app model.App) error {
	src := "internal/helm/templates"
	dest := filepath.Join("generated", app.Name, "templates")

	return CopyDir(src, dest)
}

func CopyDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		target := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(target)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
