package main

import (
	"api/internal/util"
	"fmt"
)

func main() {
	val := util.HelmValues{
		AppName: "myapp",
		Domain:  "myapp.local",
	}
	val.Image.Repository, val.Image.Tag = util.ParseImage("nginx:alpine")
	val.Image.Port = 8080

	err := util.RenderHelmAndSave(val, "./internal/helm/app", "./tmp")
	if err != nil {
		fmt.Println("❌ Helm Render Error:", err)
	} else {
		fmt.Println("✅ Helm templates rendered successfully.")
	}
}
