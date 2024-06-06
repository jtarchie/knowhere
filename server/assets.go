package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed assets
var embeddedFiles embed.FS

func assetHandler() http.Handler {
	fsys, err := fs.Sub(embeddedFiles, "assets")
	if err != nil {
		panic(fmt.Sprintf("could embed files: %s", err))
	}

	return  http.FileServer(http.FS(fsys))
}