package net

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

var hs map[string]func(req map[string]any) (resp []byte)

func Web(port int, handlers map[string]func(req map[string]any) (resp []byte)) error {
	hs = handlers
	for k, _ := range hs {
		http.HandleFunc(k, handler)
	}
	portStr := strconv.Itoa(port)
	log.Info("Web will start at " + portStr)
	err := http.ListenAndServe(":"+portStr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return err
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	f := hs[uri]
	params := r.URL.Query()
	r.ParseForm()
	maps := mergeMaps(mapChange(r.Form, false), mapChange(params, true))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Error reading request body", err)
	}
	var p3 map[string]any
	if len(body) > 0 {
		err = json.Unmarshal(body, &p3)
		if err != nil {
			log.Error(w, "Error parsing JSON", err)
		}
	}
	result := mergeMaps(maps, p3)
	resp := f(result)
	if resp != nil {
		w.Write(resp)
	}
}

func mapChange(m1 map[string][]string, query bool) map[string]any {
	changed := make(map[string]any)
	for k, v := range m1 {
		if query {
			if len(v) == 1 {
				changed[k] = v[0]
				continue
			}
		}
		changed[k] = v
	}
	return changed
}

func mergeMaps(m1 map[string]any, m2 map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}
