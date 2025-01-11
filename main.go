package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type APP struct {
	hotDir     string
	histDir    string
	title      string
	listenAddr string
	staticDir  string
}

func genHomePage(title string, hotResultMD string, histResultMD string) string {
	homePageTemplate := `
# {{ .Title }}

{{ .HotResultMD }}
{{ .HistResultMD }}
`
	// Create template data structure
	data := struct {
		Title        string
		HotResultMD  string
		HistResultMD string
	}{
		Title:        title,
		HotResultMD:  hotResultMD,
		HistResultMD: histResultMD,
	}

	// Create new template and parse the template string
	tmpl, err := template.New("homepage").Parse(homePageTemplate)
	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	// Create buffer to store the result
	var buf bytes.Buffer

	// Execute template with data
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return buf.String()
}

func (app *APP) findFileInHist(filename string) string {
	filesHist, err := listFiles(app.histDir, 3, "*.html,*.md")
	if err != nil {
		return ""
	}

	//find filename in files
	for _, file := range filesHist {
		log.Println(file)
		if filepath.Base(file) == filename {
			return filepath.ToSlash(strings.TrimPrefix(file, app.histDir))
		}
	}
	return ""
}

func (app *APP) handlePath(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request for %s\n", r.URL.Path)

	path := r.URL.Path

	// Serve home page for root path
	if path == "/" {
		hotResultMD := renderMarkdown("最近更新", app.hotDir, 1, "*.html,*.md")
		histResultMD := renderMarkdown("历史存档", app.histDir, 3, "*.html,*.md")
		homePage := genHomePage(app.title, hotResultMD, histResultMD)
		log.Println(homePage)
		html, err := markdown2html(homePage, app.title)
		if err != nil {
			http.Error(w, "Error generating home page", http.StatusInternalServerError)
			return
		}
		log.Println(html)
		// fmt.Fprintf(w, html)
		w.Write([]byte(html))
		return
	}

	// serve static files
	if strings.HasPrefix(path, "/static/") {
		app.serveFile(w, r, filepath.Join(app.staticDir, path[1:]))
		return
	}

	// Remove leading slash
	path = path[1:]

	fullPathInHot := filepath.Join(app.hotDir, path)

	// Check if path contains a directory separator and not exists in hotDir
	if strings.Contains(path, "/") && !fileExists(fullPathInHot) {
		// Handle files from histDir (e.g., /2024-12/filename)
		fullPath := filepath.Join(app.histDir, path)
		app.serveFile(w, r, fullPath)
		return
	} else {
		// Handle files from hotDir (e.g., /filename.md)
		fullPath := filepath.Join(app.hotDir, path)

		//check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if file := app.findFileInHist(path); file != "" {
				http.Redirect(w, r, file, http.StatusTemporaryRedirect)
				return
			}
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		app.serveFile(w, r, fullPath)
		return
	}
}

// Helper function to serve files
func (app *APP) serveFile(w http.ResponseWriter, r *http.Request, fullPath string) {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filename := filepath.Base(fullPath)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	title := filename

	ext := filepath.Ext(fullPath)
	if ext == ".md" {
		html, err := markdown2html(string(content), title)
		if err != nil {
			http.Error(w, "Error generating home page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}

	// 使用标准库检测 MIME 类型
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "text/plain" // 默认 MIME 类型
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}

func main() {
	// Add flag definitions
	hotDirFlag := flag.String("hot", `C:\Users\line\OneDrive\data\Released\JY-AI`, "Path to hot directory")
	histDirFlag := flag.String("hist", `C:\Users\line\OneDrive\data\Released\JY-AI\history`, "Path to history directory")
	titleFlag := flag.String("title", "学习英语", "Title of the page")
	listenAddrFlag := flag.String("listen", "127.0.0.1:8080", "Server listen address")
	flag.Parse()

	app := APP{
		hotDir:     *hotDirFlag,
		histDir:    *histDirFlag,
		title:      *titleFlag,
		listenAddr: *listenAddrFlag,
	}

	// Setup HTTP routes
	http.HandleFunc("/", app.handlePath)

	// Start the server
	log.Printf("Starting server on %s\n", app.listenAddr)
	if err := http.ListenAndServe(app.listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
