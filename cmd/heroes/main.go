package main

import (
	"context"
	"flag"
	"fmt"
	"gorestapi/internal/database"
	"gorestapi/internal/database/memorydb"
	"gorestapi/internal/database/mongodb"
	setFlag "gorestapi/internal/flag"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var (
	db           database.Herodb
	resourceFlag setFlag.FlagVar
	resourceDir  string
	uri          setFlag.FlagVar
	params       database.Params
)

const (
	basepath = "/api/heroes"
)

// initialize flags - different than the default flag package usage to enable checking flag set
func init() {
	flag.Var(&uri, "u", "URI of mongodb. If not specified, will use in memory database.")
	flag.Var(&resourceFlag, "r", "Resource directory for pictures. Required")
}

func main() {
	flag.Parse()
	fmt.Println("Running with flags: ")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %v\n", f.Name, f.Value)
	})

	// check if resource flag is present, if so check that the value points to a valid directory
	var err error
	if _, err = os.Stat(resourceFlag.Value); os.IsNotExist(err) {
		fmt.Println(resourceFlag.Value + " resource directory does not point to a real directory")
		os.Exit(1)
	}
	resourceDir, err = filepath.Abs(resourceFlag.Value)
	if err != nil {
		fmt.Printf("Unable to resolve absoulate path of resource directory: %v\n", resourceFlag.Value)
		os.Exit(1)
	}

	// connect to databse depending on what database is configured
	if uri.IsSet {
		fmt.Println("Using mongodb")
		db = &mongodb.Mongodb{}
		params = mongodb.MongodbParams{URI: uri.Value}
	} else {
		fmt.Println("Using memorydb")
		db = &memorydb.Memorydb{}
	}

	// use timeout context or mongo will hang forever
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Connect(ctx, params); err != nil {
		fmt.Println("Error connecting to db: " + err.Error())
		os.Exit(1)
	}
	fmt.Println("Database connected")

	defer func() {
		if err := db.Disconnect(ctx); err != nil {
			fmt.Println("Error disconnecting from db: " + err.Error())
			os.Exit(1)
		}
	}()

	// set up routes
	r := mux.NewRouter()
	r.HandleFunc(basepath, handleBase).Methods(http.MethodGet, http.MethodPut, http.MethodPost, http.MethodOptions)
	r.HandleFunc(basepath+"/", handleBase).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(basepath+"/{id}", handleId).Methods(http.MethodGet, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/photo/"+"{fname}", getPhoto).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(basepath+"/photo/{id}", handlePhoto).Methods(http.MethodPost, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":8080", r)
}
