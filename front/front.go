package front

import (
	"bytes"
	"html/template"
	"io"
	"path/filepath"

	"github.com/pkg/errors"
)

// Expected Errors
var (
	ErrTemplateNotFound = errors.New("template not found")
)

type HTMLRenderer struct {
	templates *template.Template
}

func NewHTMLRenderer() (*HTMLRenderer, error) {
	templatesPaths, err := filepath.Glob("./front/templates/*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the HTML template folder")
	}

	templates, err := template.New("base").ParseFiles(templatesPaths...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse one of the HTML template")
	}

	return &HTMLRenderer{templates: templates}, nil
}

func (t *HTMLRenderer) Render(w io.Writer, templateName string, params interface{}) error {
	template := t.templates.Lookup(templateName)
	if template == nil {
		return ErrTemplateNotFound
	}

	// If a key is missing return an error instead of filling with the default value.
	template = template.Option("missingkey=error")

	buf := new(bytes.Buffer)
	err := template.Execute(buf, params)
	if err != nil {
		return errors.Wrapf(err, "failed to execute the template %q", templateName)
	}

	_, err = buf.WriteTo(w)

	return err
}
