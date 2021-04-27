package handler

import (
	"errors"
	"fmt"
	"text/template"

	"github.com/satyamsoni2211/app/internal"
)

func HandleTemplate(name string) (*template.Template, error) {
	data, hasTemplate := internal.Data[name]
	if !hasTemplate {
		return nil, errors.New(fmt.Sprintf("Missing template in encoded map %s", name))
	}
	tpl, err := template.New("index.html").Parse(string(data))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error occurred while rendering %s", err))
	}
	return tpl, nil
}
