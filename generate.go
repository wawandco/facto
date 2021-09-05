package facto

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

var (
	//go:embed templates/factory.go.tmpl
	factoryTmpl string
)

// Generate a factory on a given root directory and
// a name for it (passed in args).
func Generate(root string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("generate: please provide a name for the factory")
	}

	name := args[1]
	tmpl, err := template.New("factory").Parse(factoryTmpl)
	if err != nil {
		return err
	}

	folder := filepath.Join(root, "factories")
	err = os.MkdirAll(folder, 0777)
	if err != nil {
		return fmt.Errorf("error creating factories folder: %w", err)
	}

	file, err := os.Create(filepath.Join(folder, snakeCase(name)+".go"))
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	data := struct {
		Name string
	}{
		Name: camelCase(name),
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error creating executing template: %w", err)
	}

	return nil
}

func camelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	p := strings.Fields(s)

	var g []string
	for _, value := range p {
		g = append(g, strings.Title(value))
	}

	return strings.Join(g, "")
}

func snakeCase(s string) string {
	var res = make([]rune, 0, len(s))
	var p = '_'

	for i, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			res = append(res, '_')
		} else if unicode.IsUpper(r) && i > 0 {
			if unicode.IsLetter(p) && !unicode.IsUpper(p) || unicode.IsDigit(p) {
				res = append(res, '_', unicode.ToLower(r))
			} else {
				res = append(res, unicode.ToLower(r))
			}
		} else {
			res = append(res, unicode.ToLower(r))
		}

		p = r
	}

	return string(res)
}
