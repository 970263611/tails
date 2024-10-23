package net

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var hs map[string]func(req map[string]any) (resp []byte)

func Web(port int, handlers map[string]func(req map[string]any) (resp []byte)) {
	hs = handlers
	for k, _ := range hs {
		http.HandleFunc(k, handler)
	}
	portStr := strconv.Itoa(port)
	fmt.Println("Web will start at " + portStr)
	http.ListenAndServe(":"+portStr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	f := hs[uri]
	if f == nil {
		return
	}
	params := r.URL.Query()
	r.ParseForm()
	maps := mergeMaps(mapChange(r.Form, false), mapChange(params, true))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body", err)
	}
	var p3 map[string]any
	err = json.Unmarshal(body, &p3)
	if err != nil {
		fmt.Println(w, "Error parsing JSON", err)
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
