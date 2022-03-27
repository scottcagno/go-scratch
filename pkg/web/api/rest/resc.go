package rest

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type ResourceHandler interface {
	ReturnAll(w http.ResponseWriter, r *http.Request)
	ReturnOne(w http.ResponseWriter, r *http.Request)
	InsertOne(w http.ResponseWriter, r *http.Request)
	UpdateOne(w http.ResponseWriter, r *http.Request)
	DeleteOne(w http.ResponseWriter, r *http.Request)
}

type Resource interface {
	GetAll() http.Handler
	Get(id string) http.Handler
	Add() http.Handler
	Set(id string) http.Handler
	Del(id string) http.Handler
}

func WriteAsJSON(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	return
}

func WriteAsXML(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(200)
	w.Header().Set("content-type", "application/xml")
	if err := xml.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	return
}
