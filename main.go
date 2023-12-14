package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Hero struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Power    string `json:"power"`
	AlterEgo string `json:"alterEgo"`
	PhotoURL string `json:"photoURL"`
}

const (
	basepath = "/api/heroes"
)

var (
	heroes []Hero = []Hero{
		{Id: 12, Name: "Dr. Nice", Power: "bein nice", AlterEgo: "nobody", PhotoURL: "http://localhost:8080/photo/minion.jpg"},
		{Id: 13, Name: "Bombasto", Power: "throwing stuf"},
		{Id: 14, Name: "Celeritas", Power: "celebrity", AlterEgo: "tom cruise"},
		{Id: 15, Name: "Magneta", Power: "not sure tbh"},
		{Id: 16, Name: "RubberMan", Power: "elastic arms", AlterEgo: "steve"},
		{Id: 17, Name: "Dynama", Power: "dynamite"},
		{Id: 18, Name: "Dr. IQ", Power: "talking", AlterEgo: "michael"},
		{Id: 19, Name: "Magma", Power: "making rocks"},
		{Id: 20, Name: "Tornado", Power: "spinning in an office chair", AlterEgo: "nick"},
	}
)

func getHero(id int) *Hero {
	for _, he := range heroes {
		if he.Id == id {
			return &he
		}
	}
	return nil
}

func matchName(name string) []Hero {
	hs := make([]Hero, 0)
	for _, he := range heroes {
		if strings.Contains(strings.ToLower(he.Name), strings.ToLower(name)) {
			hs = append(hs, he)
		}
	}
	return hs
}

func updateHero(h Hero) bool {
	for i, he := range heroes {
		if he.Id == h.Id {
			heroes[i] = h
			return true
		}
	}
	return false
}

func deleteHero(id int) bool {
	del := -1
	for i, he := range heroes {
		if he.Id == id {
			del = i
			break
		}
	}
	if del == -1 {
		return false
	}
	heroes = append(heroes[:del], heroes[del+1:]...)
	sortHeroes()
	return true
}

func addHero(h Hero) {
	// make sure id is not duplicate
	// works because heroes is sorted
	for _, he := range heroes {
		if he.Id == h.Id {
			h.Id++
		}
	}
	heroes = append(heroes, h)
	sortHeroes()
}

func sortHeroes() {
	// sort heroes
	sort.Slice(heroes, func(i int, j int) bool {
		return heroes[i].Id < heroes[j].Id
	})
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
}

func handleBase(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		qheroes := heroes
		if len(name) > 0 {
			qheroes = matchName(name)
		}

		h, err := json.Marshal(qheroes)
		if err != nil {
			fmt.Printf("error marashalling json: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", h)
	case http.MethodPost:
		hero := Hero{}
		err := json.NewDecoder(r.Body).Decode(&hero)
		if err != nil {
			fmt.Printf("error demarshalling json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		addHero(hero)
	case http.MethodPut:
		hero := Hero{}
		err := json.NewDecoder(r.Body).Decode(&hero)
		if err != nil {
			fmt.Printf("error demarshalling json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !updateHero(hero) {
			fmt.Printf("hero %v does not exist and could not be updated \n", hero)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleId(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
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
	case http.MethodGet:
		h := getHero(id)
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
	case http.MethodDelete:
		if !deleteHero(id) {
			fmt.Println("hero with id " + i + " could not be deleted as it does not exist")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPhoto(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	vars := mux.Vars(r)
	n, ok := vars["fname"]
	if !ok {
		fmt.Println("filename is missing in parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fileBytes, err := os.ReadFile("resources/" + n)
		if err != nil {
			fmt.Println("file could not be found")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePhoto(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	vars := mux.Vars(r)
	n, ok := vars["fname"]
	if !ok {
		fmt.Println("filename is missing in parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
	case http.MethodPut, http.MethodPost:
		// save file
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("could not read body as file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		err = os.WriteFile("resources/"+n, body, 0644)
		if err != nil {
			fmt.Println("could not write body to file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// update hero with new file location
		h := getHero(id)
		if h != nil {
			// TODO - fix this
			h.PhotoURL = "http://localhost:8080/photo/" + n
			if !updateHero(*h) {
				fmt.Printf("could not update photo for hero %d\n", id)
			}
		}

		fmt.Println(heroes)

	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	sortHeroes()
	r := mux.NewRouter()
	r.HandleFunc(basepath, handleBase).Methods(http.MethodGet, http.MethodPut, http.MethodPost, http.MethodOptions)
	r.HandleFunc(basepath+"/", handleBase).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(basepath+"/{id}", handleId).Methods(http.MethodGet, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/photo/"+"{fname}", getPhoto).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(basepath+"/photo/{id}/"+"{fname}", handlePhoto).Methods(http.MethodPut, http.MethodPost, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":8080", r)
}
