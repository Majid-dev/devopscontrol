package util

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type HelmValues struct {
	AppName string
	Image   struct {
		Repository string
		Tag        string
		Port       int
	}
	Domain string
}

const helmValuesTemplate = `
appName: "{{ .AppName }}"
image:
  repository: "{{ .Image.Repository }}"
  tag: "{{ .Image.Tag }}"
  port: {{ .Image.Port }}

domain: "{{ .Domain }}"
`

func RenderHelmAndSave(values HelmValues, chartDir, outputDir string) error {
	// 1. Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// 2. Create and write values.yaml
	valuesYAML, err := yaml.Marshal(values)
	if err != nil {
		return err
	}
	valuesPath := filepath.Join(outputDir, "values.yaml")
	if err := os.WriteFile(valuesPath, valuesYAML, 0644); err != nil {
		return err
	}

	// 3. Run helm template
	cmd := exec.Command("helm", "template", values.AppName, chartDir, "-f", valuesPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	// 4. Write output manifest.yaml
	manifestPath := filepath.Join(outputDir, "manifest.yaml")
	if err := os.WriteFile(manifestPath, out.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

// ParseImage splits an image string like "nginx:alpine" into repository and tag.
func ParseImage(image string) (repository string, tag string) {
	parts := strings.SplitN(image, ":", 2)
	repository = parts[0]
	if len(parts) > 1 {
		tag = parts[1]
	} else {
		tag = "latest"
	}
	return
}
