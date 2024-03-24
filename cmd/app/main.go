package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"

	"github.com/stijn577/gotth/internal/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("./internal/static"))
	http.Handle("/", fs)

	http.Handle("/clicked", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create the type that we expect to decode from the JSON input
		type Counter struct {
			Value int `json:"value"`
		}

		// close the body at the end
		defer r.Body.Close()

		// create decoder for the JSON body
		decoder := json.NewDecoder(r.Body)

		counter := Counter{}

		// decode the JSON into the struct
		err := decoder.Decode(&counter)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

		fmt.Printf("INFO: \"/clicked\" called, counter value: %d\n", counter.Value)

		// input JSON to templ generated go function to get HTML output
		html, err := templ.ToGoHTML(context.Background(), handlers.Clicked(counter.Value))

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

		// write HTTP 200 OK with ssr html output
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))

	http.Handle("/card", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		err := r.ParseMultipartForm(2_000_000)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		image, handler, err := r.FormFile("image")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer image.Close()

		ext := filepath.Ext(handler.Filename)
		filename := fmt.Sprintf("%s%s", strings.ReplaceAll(handler.Filename, ext, ""), ext)

		dst, err := os.Create(filepath.Join("internal/static/images", filename))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		defer dst.Close()

		if _, err := io.Copy(dst, image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		image_path := fmt.Sprintf("/images/%s", filename)

		alt := r.FormValue("alt")
		title := r.FormValue("title")
		description := r.FormValue("description")
		hashtags := r.FormValue("hashtags")

		fmt.Printf("INFO: \"/card\" called, title: \"%s\"\n", title)
		fmt.Println(image_path)

		html, err := templ.ToGoHTML(context.Background(), handlers.Card(image_path, alt, title, description, hashtags))

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))

	port := ":8080"

	fmt.Printf("Listening on http://localhost%s/\n\n", port)

	http.ListenAndServe(port, nil)
}
