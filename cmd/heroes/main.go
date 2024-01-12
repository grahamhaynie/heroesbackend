package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gorestapi/internal/database"
	"gorestapi/internal/database/memorydb"
	"gorestapi/internal/database/mongodb"
	setFlag "gorestapi/internal/flag"
	"gorestapi/internal/hero"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	db          database.Herodb
	resourceDir string
	uri         setFlag.FlagVar
	params      database.Params
)

const (
	basepath    = "/api/heroes"
	resourceEnv = "RESOURCE_DIR"
)

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleBase(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		var qheroes []hero.Hero
		var err error
		if len(name) == 0 {
			qheroes, err = db.GetAll()
			if err != nil {
				fmt.Printf("error querying database: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			qheroes, err = db.GetByName(name)
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
	case http.MethodPost:
		hero := hero.Hero{}
		err := json.NewDecoder(r.Body).Decode(&hero)
		if err != nil {
			fmt.Printf("error demarshalling json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = db.AddHero(hero); err != nil {
			fmt.Printf("error adding hero to database: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		hero := hero.Hero{}
		err := json.NewDecoder(r.Body).Decode(&hero)
		if err != nil {
			fmt.Printf("error demarshalling json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = db.UpdateHero(hero); err != nil {
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
		h, err := db.GetById(id)
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
	case http.MethodDelete:
		if err := db.DeleteHero(id); err != nil {
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
		fileBytes, err := os.ReadFile(resourceDir + "/" + n)
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

	case http.MethodOptions:
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// initialize flag
func init() {
	flag.Var(&uri, "u", "URI of mongodb. If not specified, will use in memory database.")
}

func main() {
	flag.Parse()
	fmt.Println("Running with flags: ")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %v\n", f.Name, f.Value)
	})

	var err error
	e := os.Getenv(resourceEnv)
	if _, err = os.Stat(e); os.IsNotExist(err) {
		fmt.Println(resourceEnv + " env variable does not point to a real directory, please set")
		os.Exit(1)
	}
	resourceDir, err = filepath.Abs(e)
	if err != nil {
		fmt.Printf("Unable to resolve absolute path  %v\n", e)
		os.Exit(1)
	}
	fmt.Println("Running with " + resourceEnv + " set to " + resourceDir)

	if uri.IsSet {
		fmt.Println("Using mongodb")
		db = &mongodb.Mongodb{}
		// mongodb://localhost:27017
		params = mongodb.MongodbParmas{URI: uri.Value}
	} else {
		fmt.Println("Using memorydb")
		db = &memorydb.Memorydb{}
	}
	if err := db.Connect(params); err != nil {
		fmt.Println("Error connecting to db: " + err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := db.Disconnect(); err != nil {
			fmt.Println("Error disconnecting from db: " + err.Error())
			os.Exit(1)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc(basepath, handleBase).Methods(http.MethodGet, http.MethodPut, http.MethodPost, http.MethodOptions)
	r.HandleFunc(basepath+"/", handleBase).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(basepath+"/{id}", handleId).Methods(http.MethodGet, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/photo/"+"{fname}", getPhoto).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(basepath+"/photo/{id}/"+"{fname}", handlePhoto).Methods(http.MethodPut, http.MethodPost, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":8080", r)
}
