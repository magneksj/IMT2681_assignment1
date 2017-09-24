package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
Test_malformedURL test for wrong number of URL segments
*/
func Test_malformedURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/projectinfo/v1/github.com/golang/go/f",
		ts.URL + "/projectinfo/v1/github.com/",
		ts.URL + "/projectinfo/v1/",
	}

	for _, test := range testCases {
		resp, err := http.Get(test)
		if err != nil {
			t.Errorf("Get request error: %s", err.Error())
			return
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected StatusCode %d from %s, recieved %d", http.StatusBadRequest, test, resp.StatusCode)
			return
		}
	}
}

func Test_wrongURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	testCases := []string{
		ts.URL + "/info/v1/github.com/golang/go/",
		ts.URL + "/projectinfo/v2/github.com/golang/go/",
		ts.URL + "/projectinfo/v1/bitbucket.org/golang/go/",
		ts.URL + "/pi/v3/ntnu.blackboard.com/golang/go/",
	}

	for _, test := range testCases {
		resp, err := http.Get(test)
		if err != nil {
			t.Errorf("Get request error: %s", err.Error())
			return
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected StatusCode %d from %s, recieved %d", http.StatusBadRequest, test, resp.StatusCode)
			return
		}
	}
}

func Test_getGithubInfo_pass(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	testCases := []string{
		ts.URL + "/projectinfo/v1/github.com/golang/go",
		ts.URL + "/projectinfo/v1/github.com/magneksj/TestGithubAPI/",
		ts.URL + "/projectinfo/v1/github.com/freeCodeCamp/freeCodeCamp/",
	}

	for _, test := range testCases {
		resp, err := http.Get(test)
		if err != nil {
			t.Errorf("Get request error: %s", err.Error())
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected StatusCode %d from %s, recieved %d", http.StatusOK, test, resp.StatusCode)
			return
		}

		var data Response
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			t.Errorf("Error parsing json body: %s", err.Error())
		}

		if data.Owner == "" {
			t.Error("Expected a repo owner")
		}
	}
}

func Test_getGithubInfo_fail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	URL := ts.URL + "/projectinfo/v1/github.com/magneksj/notreal/"

	resp, err := http.Get(URL)
	if err != nil {
		t.Errorf("Get request error: %s", err.Error())
		return
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusCode %d from %s, recieved %d", http.StatusNotFound, URL, resp.StatusCode)
		return
	}
}
