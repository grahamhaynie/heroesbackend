package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gorestapi/internal/hero"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	TIMEOUT = 10 * time.Second
)

// set CORS Headers
func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// handler for default route
func handleBase(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	// wrap request context in timeout
	ctx, cancel := context.WithTimeout(r.Context(), TIMEOUT)
	defer cancel()

	switch r.Method {
	// get request with a name query parameter returns a specific hero by name
	// get request without a name query parameter returns all heroes
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		var qheroes []hero.Hero
		var err error
		if len(name) == 0 {
			qheroes, err = db.GetAll(ctx)
			if err != nil {
				fmt.Printf("error querying database: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			qheroes, err = db.GetByName(ctx, name)
			if err != nil {
				fmt.Printf("error querying database: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		h, err := json.Marshal(qheroes)
		if err != nil {
			fmt.Printf("error marashalling json: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", h)

	// post request contains a hero, add hero to database
	case http.MethodPost:
		hero := hero.Hero{}
		err := json.NewDecoder(r.Body).Decode(&hero)
		if err != nil {
			fmt.Printf("error demarshalling json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = db.AddHero(ctx, hero); err != nil {
			fmt.Printf("error adding hero to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	// put request updates a hero in database
	case http.MethodPut:
		hero := hero.Hero{}
		err := json.NewDecoder(r.Body).Decode(&hero)
		if err != nil {
			fmt.Printf("error demarshalling json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = db.UpdateHero(ctx, hero); err != nil {
			fmt.Printf("hero %v does not exist and could not be updated \n", hero)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	// for options, return nothing as this is a CORS request
	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handler for the /{id} route
func handleId(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	// wrap request context in timeout
	ctx, cancel := context.WithTimeout(r.Context(), TIMEOUT)
	defer cancel()

	// parse id from url path
	vars := mux.Vars(r)
	i, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(i)
	if err != nil {
		fmt.Println("could not convert id to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {

	// get returns hero with specified id
	case http.MethodGet:
		h, err := db.GetById(ctx, id)
		if err != nil {
			fmt.Printf("error querying database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if h == nil {
			fmt.Println("hero with id " + i + " not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		hs, err := json.Marshal(*h)
		if err != nil {
			fmt.Printf("error marashalling json: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", hs)

	// delete hero with specified id
	case http.MethodDelete:
		if err := db.DeleteHero(ctx, id); err != nil {
			fmt.Println("hero with id " + i + " could not be deleted as it does not exist")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	// options for cors requests so just return
	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handler for getting a photo file
func getPhoto(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	// parse filename
	vars := mux.Vars(r)
	n, ok := vars["fname"]
	if !ok {
		fmt.Println("filename is missing in parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fileBytes, err := os.ReadFile(resourceDir + "/" + n)
		if err != nil {
			fmt.Println("file could not be found")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)

	// options for cors requests so just return
	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handler for photo upload route
func handlePhoto(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	// parse filename
	vars := mux.Vars(r)
	n, ok := vars["fname"]
	if !ok {
		fmt.Println("filename is missing in parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut, http.MethodPost:

		// save file
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("could not read body as file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		err = os.WriteFile(resourceDir+"/"+n, body, 0644)
		if err != nil {
			fmt.Println("could not write body to file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update hero with new file location
		response := map[string]string{"url": "http://localhost:8080/photo/" + n}
		w.Header().Set("Content-Type", "application/json")
		m, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("error marashalling json: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", m)

	// options for cors requests so just return
	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
