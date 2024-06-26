package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
		type Card struct {
			Image_path  string `json:"image_path"`
			Alt         string `json:"alt"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Hashtags    string `json:"hashtags"`
		}

		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		card := Card{}

		err := decoder.Decode(&card)

		fmt.Printf("INFO: \"/card\" called, title: \"%s\"\n", card.Title)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

		html, err := templ.ToGoHTML(context.Background(), handlers.Card(card.Image_path, card.Alt, card.Title, card.Description, card.Hashtags))

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
