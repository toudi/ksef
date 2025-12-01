package abstract

import (
	"errors"
	"path/filepath"
	"text/template"
)

var (
	ErrNoTemplatesDefined = errors.New("nie zdefiniowano żadnych szablonów renderujących")
)

type HeaderFooterSettings struct {
	Left   string
	Center string
	Right  string
}

type TemplatesCollection struct {
	Enabled  bool
	SrcPath  string
	Header   HeaderFooterSettings
	Footer   HeaderFooterSettings
	Template *template.Template
}

type Templates struct {
	Invoice TemplatesCollection
	UPO     TemplatesCollection
}

func ReadTemplatesFromDirectory(dirname string, extension string, funcs template.FuncMap) (TemplatesCollection, error) {
	var templates = TemplatesCollection{
		Enabled: true,
		SrcPath: dirname,
	}

	var err error

	templates.Template, err = template.New("").Funcs(funcs).Delims("<<", ">>").ParseGlob(filepath.Join(dirname, "*."+extension))

	return templates, err
}
