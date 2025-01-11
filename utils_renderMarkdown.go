package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

func renderMarkdown(title string, dir string, depth int, exts string) string {
	hotMarkdownTemplate := `## {{ .Title }}
{{ range .Files }}
- [{{ .Name }}]({{ .Path }})
{{- end }}
`
	type FileInfo struct {
		Name string
		Path string
	}

	type TemplateData struct {
		Title string
		Files []FileInfo
	}

	files, err := listFiles(dir, depth, exts)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	var fileInfos []FileInfo
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		relPath, _ := filepath.Rel(dir, file)
		relPath = strings.ReplaceAll(relPath, "\\", "/")
		fileInfos = append(fileInfos, FileInfo{
			Name: name,
			Path: "./" + relPath,
		})
	}

	data := TemplateData{
		Title: title,
		Files: fileInfos,
	}

	var builder strings.Builder
	tmpl := template.Must(template.New("hot").Parse(hotMarkdownTemplate))
	err = tmpl.Execute(&builder, data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return builder.String()
}
