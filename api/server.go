package api

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	"net/http"

	"html/template"
	"strings"

	// "github.com/google/uuid"
	"github.com/gorilla/mux"
)

// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####
//	Usefull Structs
// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####

type Server struct {
	*mux.Router
}

type Data struct {
	Error error
	Name  string
	Theme string
	From  string
}

// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####
//	Main Functions for Creating ServerRoutes
// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####

func NewServer() *Server {
	server := &Server{
		Router: mux.NewRouter(),
	}
	server.routes()
	return server
}

func (server *Server) routes() {
	server.HandleFunc("/", server.sendIndex()).Methods("GET")
	server.NotFoundHandler = server.sendHTML()
	server.HandleFunc("/static/{type}/{filename}", server.sendFile()).Methods("GET")
}

// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####
//	Route Functions
// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####

func (server *Server) sendFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := fmt.Sprintf("static/%s/%s", mux.Vars(r)["type"], mux.Vars(r)["filename"])

		data, err := ioutil.ReadFile(string(path))

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - File Not Found!"))
		}

		var contentType string
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".svg") {
			contentType = "image/svg+xml"
		} else {
			contentType = "text/plain"
		}

		w.Header().Add("Content-Type", contentType)
		w.Write(data)
	}
}

func (server *Server) sendIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{nil, "Jo", "dark", ""}

		data.From = "/"

		allFiles := getAllTemplateFiles(fmt.Sprintf("%vindex.html", data.From))

		templates, _ := template.ParseFiles(allFiles...)

		executeTemplates(templates, w, data)
	}
}

func (server *Server) sendHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{nil, "", "", ""}

		data.From = getPathFromURL(r.URL.Path)

		allFiles := getAllTemplateFiles(data.From)

		templates, err := template.ParseFiles(allFiles...)
		if err != nil {
			allFiles[len(allFiles)-1] = "./views/error/index.html"
			templates, _ = template.ParseFiles(allFiles...)
		}

		executeTemplates(templates, w, data)
	}
}

// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####
//
// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####

func getPathFromURL(urlPath string) string {
	filepath := strings.Split(urlPath, "/")
	path := ""
	for index, pathpart := range filepath {
		if index == 0 {
			continue
		}
		if index == len(filepath)-1 {
			if !strings.Contains(pathpart, ".") {
				pathpart += "/index.html"
			}
		}
		path = fmt.Sprintf("%s/%s", path, pathpart)
	}
	return path
}

func getAllTemplateFiles(path string) []string {
	var allFiles []string
	files, err := ioutil.ReadDir("./views/templates")
	if err != nil {
		fmt.Println("Err on reading templates")
	}

	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./views/templates/"+filename)
		}
	}
	allFiles = append(allFiles, fmt.Sprintf("./views%s", path))
	return allFiles
}

func executeTemplates(templates *template.Template, w http.ResponseWriter, data Data) {
	s1 := templates.Lookup("header.html")
	s1.ExecuteTemplate(w, "header", data)
	s2 := templates.Lookup("navigation.html")
	s2.ExecuteTemplate(w, "navigation", data)
	s3 := templates.Lookup("index.html")
	s3.ExecuteTemplate(w, "index", data)
	s4 := templates.Lookup("footer.html")
	s4.ExecuteTemplate(w, "footer", data)
	s4.Execute(w, data)
}
