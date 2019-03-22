package yin

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type ClientConfig struct {
	Directory             string
	BaseHref              string
	SinglePageApplication bool
}

func ServeClient(c ClientConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if c.BaseHref != "" {
			url = strings.TrimPrefix(r.URL.Path, "/"+c.BaseHref)
		}
		serveThis := path.Join(c.Directory, url)

		file, _ := os.Stat(serveThis)
		if file != nil {
			Res(w).File(r, serveThis)
			return
		}

		if c.SinglePageApplication == false {
			file, _ = os.Stat(serveThis + ".html")
			if file != nil {
				Res(w).File(r, serveThis+".html")
				return
			}

			file, _ = os.Stat(serveThis + "/index.html")
			if file != nil {
				Res(w).File(r, serveThis+"/index.html")
				return
			}
		}

		if c.SinglePageApplication == true {
			Res(w).File(r, c.Directory+"/index.html")
		} else {
			Res(w).SendStatus(http.StatusNotFound)
		}
	}
}
