package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"

	"github.com/wizzldev/mailer/types"
)

//go:embed templates/*.html
var templates embed.FS

// Template represents an interface for loading and rendering templates.
type Template interface {
	// ParseComponent parses a component template and returns its content as a string.
	ParseComponent(name string, data types.TemplateProps) (string, error)
}

// template implements the Template interface.
type template struct {
	root    string
	section string
	content []byte
}

// NewTemplate creates a new template instance.
func NewTemplate(rootLayout string, mainSection string, data types.TemplateProps) (Template, error) {
	root, err := templates.Open(fmt.Sprintf("templates/%s.html", rootLayout))
	if err != nil {
		return nil, err
	}
	defer root.Close()

	content, err := io.ReadAll(root)
	if err != nil {
		return nil, err
	}

	for key, value := range data {
		content = bytes.ReplaceAll(content, []byte(fmt.Sprintf("{%s}", key)), []byte(fmt.Sprintf("%v", value)))
	}

	return &template{
		root:    rootLayout,
		section: mainSection,
		content: content,
	}, nil
}

// ParseComponent parses a component template and returns its content as a string.
func (t *template) ParseComponent(name string, data types.TemplateProps) (string, error) {
	if name == t.root {
		return "", errors.New("root template cannot be loaded as a component")
	}

	component, err := templates.Open(fmt.Sprintf("templates/%s.html", name))
	if err != nil {
		return "", err
	}
	defer component.Close()

	content, err := io.ReadAll(component)
	if err != nil {
		return "", err
	}

	for key, value := range data {
		content = bytes.ReplaceAll(content, []byte(fmt.Sprintf("{%s}", key)), []byte(fmt.Sprintf("%v", value)))
	}

	return string(bytes.Replace(t.content, []byte("@"+t.section), content, 1)), nil
}
