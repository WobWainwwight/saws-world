package templates

import (
	"embed"
	"html/template"
)

//go:embed index.html south-america.html south-america-image.html
var fs embed.FS

var Index = "index"
var SouthAmerica = "south-america"

func GetTemplates() (*template.Template, error) {
	tmps, err := template.ParseFS(fs, "index.html", "south-america.html", "south-america-image.html")
	if err != nil {
		return nil, err
	}
	return tmps, nil
}
