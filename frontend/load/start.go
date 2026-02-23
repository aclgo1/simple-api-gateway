package load

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type load struct {
	Tmpl        *template.Template
	HtmlDirName string
	CssDirName  string
	PathCss     string
	Root        string
}

func NewLoad(htmlDirName, cssDirName, root string) *load {
	return &load{
		HtmlDirName: htmlDirName,
		CssDirName:  cssDirName,
		Root:        root,
	}

}

func (l *load) Start() (*Pages, error) {
	if err := l.searchAndParseHtmlPages(); err != nil {
		return nil, err
	}

	if err := l.searchPathCss(); err != nil {
		return nil, err
	}

	if err := l.validate(); err != nil {
		return nil, err
	}

	pages := &Pages{tmpl: l.Tmpl, pathCss: l.PathCss}

	return pages, nil
}

func (l *load) validate() error {
	if l.Tmpl == nil {
		return errors.New("l.Tmpl is nil")
	}

	// if l.PathCss == "" {
	// 	return errors.New("l.PathCss is empty")
	// }

	return nil
}

func (l *load) searchPathCss() error {
	wk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() && info.Name() == l.CssDirName {
			l.PathCss = path
		}

		return nil
	}

	return filepath.Walk(l.Root, wk)
}

func (l *load) searchAndParseHtmlPages() error {

	wk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() && info.Name() == l.HtmlDirName {
			pattern := fmt.Sprintf("%s/*.html", path)

			glob, err := template.ParseGlob(pattern)
			if err != nil {
				return err
			}

			l.Tmpl = template.Must(glob, nil)

			return nil
		}

		return nil
	}

	return filepath.Walk(l.Root, wk)
}
