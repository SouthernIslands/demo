package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const baseURL string = "https://datastore.googleapis.com/v1/projects/"

type cacheHandler struct {
	*Server
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paras := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(paras) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	projectid := strings.Split(paras, ":")[0]
	method := strings.Split(paras, ":")[1]
	token := r.URL.RawQuery[strings.LastIndex(r.URL.RawQuery, "=")+1:]
	//token := r.Header.Get("Authorization")[0]
	log.Println(method)
	log.Println(token)

	if len(projectid) == 0 || len(method) == 0 || len(token) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		//???
		//w.Write([]byte("Missing argument."))
		return
	}

	m := r.Method
	if m == http.MethodPost {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		log.Println(req)

		tmpk := req["keys"].(map[string]interface{})
		tmpp := tmpk["path"].(map[string]interface{})
		kind := tmpp["kind"].(string)
		id := fmt.Sprintln(tmpp["id"])

		if len(req) != 0 {
			key := projectid + kind + id
			if method == "lookup" {
				b, e := h.Get(key)
				if e != nil {
					log.Println(e)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(e.Error()))
				}
				w.Header().Set("Content-Type", "application/json")
				if b == nil {
					//cache miss
					//add to map and list
					data := h.DoFetch(projectid, key, token, req)
					//res := http.ResponseWriter(resp)
					//past the response to client
					//outdated？

					//map -> json -> encode json to w
					json.NewEncoder(w).Encode(data)
				} else {
					w.Write(b)
				}
			} else if method == "commit" {
				//commit

				//set to cache
				h.Set(key, []byte("map"))
			} else {
				//filter
				//perform as a http agent
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing payload."))
		}
		return
	}
	//if m == http.MethodDelete {
	//	e := h.Del(key)
	//	if e != nil {
	//		log.Println(e)
	//		w.WriteHeader(http.StatusInternalServerError)
	//	}
	//	return
	//}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *cacheHandler) DoFetch(projectid, key, token string, message map[string]interface{}) map[string]interface{} {
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(baseURL+projectid+":lookup"+"?access_token="+token,
		"application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	//json.NewDecoder(bodyBytes).Decode(&result)
	//resp.body ([]byte)-> map

	json.Unmarshal(bodyBytes, &result)

	log.Println("Retrieved from Datastore :", resp)

	if resp.StatusCode == http.StatusOK {
		if result["found"] != nil {
			log.Println("Entity Found: ", result["found"])
		}
		if result["missing"] != nil {
			log.Println("Entity Not Found: ", result["missing"])
		}
		log.Println(result)

		h.Set(key, bodyBytes)
	} else {
		log.Println("WARNING: ", resp.StatusCode, resp.Body)
	}
	return result
}

/*
func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := r.Method
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			e := h.Set(key, b)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}
	if m == http.MethodGet {
		b, e := h.Get(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(b) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(b)
		return
	}
	if m == http.MethodDelete {
		e := h.Del(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
*/
func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}
