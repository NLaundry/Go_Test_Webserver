package main

import (
	//"context"
	"bytes"
	"fmt"
	"github.com/a-h/templ"
	"github.com/yuin/goldmark" // markdown parser
	"log"
	"net/http"
	"os"
)

var pages = []string{"Shower-Thoughts", "Tech", "Poetry", "Publications", "Talks"}

func main() {
	mux := http.NewServeMux()
	// Home Page
	mux.HandleFunc("/", ComponentHandler(home(pages)))

	// index pages for each section
	for _, page := range pages {
		fmt.Println(page)
        files, err := os.ReadDir("./static/" + page)
        if err != nil {
            fmt.Println(err)
        }
        fileNames := make([]string, 0)
        for _, file := range files {
            fileNames = append(fileNames, file.Name())
        }
		mux.HandleFunc("/"+page, ComponentHandler(sectionPage(page, fileNames)))
	}


	for _, page := range pages {
		mux.HandleFunc("/"+page+"/{pageName}/", func(w http.ResponseWriter, r *http.Request) {
			pageName := r.PathValue("pageName")
			filePath := "./static/" + page + "/" + pageName
			fmt.Println(filePath)
			mdBuf, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
			}
            var buf bytes.Buffer
            goldmark.Convert(mdBuf, &buf)
            articlePage(filePath, buf.String()).Render(r.Context(), w)

		})
	}

	log.Fatal(http.ListenAndServe(":8000", mux))

}

// // creates a function that matches the HandlerFunc signature and calls the component's Render method
func ComponentHandler(component templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component.Render(r.Context(), w)
	}
}
