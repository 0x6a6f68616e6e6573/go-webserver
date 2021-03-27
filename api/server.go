package api

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	name  string
	theme string
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
		parsedTemplate, _ := template.ParseFiles(fmt.Sprintf("static/%s/%s", mux.Vars(r)["type"], mux.Vars(r)["filename"]))

		if err := parsedTemplate.Execute(w, nil); err != nil {
			log.Println("Error executing template or file not found")
		}
	}
}

func (server *Server) sendIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{nil, "we", "ar", "e"}
		data.From = "/"

		templates, err := template.ParseFiles(getAllTemplateFiles("/index.html")...)
		if err != nil {
			w.Write([]byte("????????"))
		}
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
		fmt.Println(err)
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
	s1.ExecuteTemplate(w, "header", nil)
	s2 := templates.Lookup("navigation.html")
	s2.ExecuteTemplate(w, "navigation", nil)
	s3 := templates.Lookup("index.html")
	s3.ExecuteTemplate(w, "index", nil)
	s4 := templates.Lookup("footer.html")
	s4.ExecuteTemplate(w, "footer", nil)
	s4.Execute(w, data)
}
