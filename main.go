package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/notifications"
	"sync"
)

type EventResp struct{
	Events []notifications.Event
}

type CatalogResp struct {
	Repositories []string `json:"repositories"`
}

// PoC keep the repository on memory
var repository [] string
var repositorySync sync.Mutex
var catalogLoaded = false

func addRepository(newRepo string){
	repositorySync.Lock()
	for _, repo := range repository{
		if newRepo == repo {
			repositorySync.Unlock()
			return
		}
	}
	repository = append(repository, newRepo)
	repositorySync.Unlock()
}

func processPush(event notifications.Event){
	logrus.Info(event)
	addRepository(event.Target.Repository)
}

func ProcessEvent(w http.ResponseWriter, r *http.Request)  {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var eventResp EventResp;
	err := decoder.Decode(&eventResp)
	if err != nil{
		logrus.Fatal(err)
	}

	for _, event := range eventResp.Events{
		switch event.Action{
		case "push":
			processPush(event)
		default:
			logrus.Errorf("Action %s ignored", event.Action)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func GetCatalog(w http.ResponseWriter, r *http.Request){
	// Return a conflict if the catalog is not fully load
	if !catalogLoaded {
		w.WriteHeader(http.StatusConflict)
		return
	}

	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(CatalogResp{Repositories: repository})
}


func LoadCatalog(){
	logrus.Info("Loading catalog from http://localhost:5000")

	var catalogResp CatalogResp
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:5000/v2/_catalog?n=5000",nil)
	if err != nil {
		logrus.Fatal(err)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&catalogResp)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, repo := range catalogResp.Repositories {
		addRepository(repo)
		logrus.Infof("Repo added to the catalog: %s", repo)
	}
	logrus.Info("Catalog loaded")
	catalogLoaded = true
}


func main() {
	r := mux.NewRouter()
	r.Methods(http.MethodPost).Path("/event").HandlerFunc(ProcessEvent)
	r.Methods(http.MethodGet).Path("/v2/_catalog").HandlerFunc(GetCatalog)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:5001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go LoadCatalog()
	logrus.Fatal(srv.ListenAndServe())

}
