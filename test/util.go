package test

import (
	"bytes"
	"encoding/json"
	"kayak-backend/global"
	"kayak-backend/model"
	"net/http"
	"net/http/httptest"
)

func Get(url string, token string, query map[string][]string, dest interface{}) int {
	var buf bytes.Buffer
	req, _ := http.NewRequest("GET", url, &buf)
	if token != "" {
		req.Header.Add(global.TokenHeader, token)
	}
	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			for _, c := range v {
				q.Add(k, c)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	w := httptest.NewRecorder()
	global.Router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), dest)
	return w.Code
}

func Post(url string, token string, src interface{}, dest interface{}) int {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(src)
	req, _ := http.NewRequest("POST", url, &buf)
	if token != "" {
		req.Header.Add(global.TokenHeader, token)
	}
	w := httptest.NewRecorder()
	global.Router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), dest)
	return w.Code
}

func Delete(url string, token string, src interface{}, dest interface{}) int {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(src)
	req, _ := http.NewRequest("DELETE", url, &buf)
	if token != "" {
		req.Header.Add(global.TokenHeader, token)
	}
	w := httptest.NewRecorder()
	global.Router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), dest)
	return w.Code
}

func Put(url string, token string, src interface{}, dest interface{}) int {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(src)
	req, _ := http.NewRequest("PUT", url, &buf)
	if token != "" {
		req.Header.Add(global.TokenHeader, token)
	}
	w := httptest.NewRecorder()
	global.Router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), dest)
	return w.Code
}

func randomUser() model.User {
	return randomUsers[r.Intn(len(randomUsers))]
}

func randomGroup() model.Group {
	return randomGroups[r.Intn(len(randomGroups))]
}

func randomProblemSet() model.ProblemSet {
	return randomProblemSets[r.Intn(len(randomProblemSets))]
}

func randomProblem() model.ProblemType {
	return randomProblems[r.Intn(len(randomProblems))]
}

func randomNote() model.Note {
	return randomNotes[r.Intn(len(randomNotes))]
}

func randomDiscussion() model.Discussion {
	return randomDiscussions[r.Intn(len(randomDiscussions))]
}
