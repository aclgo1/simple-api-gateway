package service

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

func SpaHandler(embedFS embed.FS, root string) http.HandlerFunc{
	subFs,err := fs.Sub(embedFS,root)
	if err != nil {
		log.Fatalf("error loading static files: %v", err)
	}

	fileServer := http.FileServer(http.FS(subFs))

	return func (w http.ResponseWriter, r *http.Request)  {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == ""{
			path = "index.html"
		}

		_,err := fs.Stat(subFs,path)
		if err != nil {
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w,r)
	}
}