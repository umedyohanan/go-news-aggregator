package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"newsaggregator/pkg/storage"
	"newsaggregator/pkg/storage/mongo"
	"testing"
)

func TestAPI_newsHandler(t *testing.T) {
	dbase, err := mongo.New("mongodb://localhost:27017/")
	if err != nil {
		t.Errorf(err.Error())
	}
	api := New(dbase)
	// Create HTTP-request.
	req := httptest.NewRequest(http.MethodGet, "/news/1", nil)
	// Create object to save response handler.
	rr := httptest.NewRecorder()
	// Getting router. Router for request path and method
	// will call handler. Handler will write response to created object.
	api.r.ServeHTTP(rr, req)
	// Checking response code.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("response code is incorrect: got %d, expectedf %d", rr.Code, http.StatusOK)
	}
	// Read response body.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("failed to decode server response: %v", err)
	}
	// Decode JSON to list of posts.
	var data []storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("failed to decode server response: %v", err)
	}
	// Check that there is only one element in list.
	const wantLen = 1
	if len(data) != wantLen {
		t.Fatalf("got %d records, expected %d", len(data), wantLen)
	}
}
