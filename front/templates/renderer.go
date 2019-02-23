package templates

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/pkg/errors"
)

// Expected Errors
var (
	ErrTemplateNotFound = errors.New("template not found")
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer() (*Renderer, error) {
	templatesPaths, err := filepath.Glob("./front/templates/*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the  template folder")
	}

	templates, err := template.New("base").ParseFiles(templatesPaths...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse one of the  template")
	}

	return &Renderer{templates: templates}, nil
}

func (t *Renderer) Render(w io.Writer, templateName string, params interface{}) error {
	template := t.templates.Lookup(templateName)
	if template == nil {
		return ErrTemplateNotFound
	}

	// If a key is missing return an error instead of filling with the default value.
	template = template.Option("missingkey=error")

	err := template.Execute(w, params)
	if err != nil {
		return errors.Wrapf(err, "failed to execute the template %q", templateName)
	}

	return nil
}
