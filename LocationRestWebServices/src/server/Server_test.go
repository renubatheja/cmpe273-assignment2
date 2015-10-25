package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"gopkg.in/mgo.v2"
	"github.com/julienschmidt/httprouter"
  	"model"
  	"controllers"
  	"testing"
  	"strings"
  	"encoding/json"
)

/*
* A Test file - For Testing REST API calls (GET, POST, PUT, DELETE)
*/

var Router httprouter.Router

func TestAllRESTAPIs(t *testing.T) {
	  //Setup httprouter and handlers for REST APIs for location service
	  Router := httprouter.New()
	  lc := controllers.NewLocationController(GetRemoteMGOSession())
	
	  Router.GET("/locations", lc.GetAllLocations)
	  Router.GET("/locations/:id", lc.GetLocation)
	  Router.POST("/locations", lc.CreateLocation)
	  Router.PUT("/locations/:id", lc.UpdateLocation)
	  Router.DELETE("/locations/:id", lc.RemoveLocation)
	  
	  //Test - POST a new location 
	  locationJson := `{"name": "Marshal Eriksen", "address": "1024", "city": "Sunnyvale", "state": "CA","zipcode": "94086"}`
	  reader := strings.NewReader(locationJson) //Convert string to reader
	  
	  req, _ := http.NewRequest("POST", "/locations", reader)
	
	  w := httptest.NewRecorder()
	  Router.ServeHTTP(w, req)
	
	  if !(w.Code == 201) {
	      t.Errorf("Failed to get response from server %s\n", w.Body.String())
	      t.Errorf("Error Code %s\n", w.Code)
	      t.FailNow()
	  }
  
      responseFromPost := model.Location{}
      json.NewDecoder(w.Body).Decode(&responseFromPost)
      oid := fmt.Sprintf("%x", string(responseFromPost.Id))
  
	  s := []string{"/locations/", string(oid)};
	  uri := strings.Join(s, "");
	  req, _ = http.NewRequest("GET", uri, nil)
	  w = httptest.NewRecorder()
	  Router.ServeHTTP(w, req)
	
	  if !(w.Code == http.StatusOK) {
	      t.Errorf("Failed to get response from server %s\n", w.Body.String())
	      t.Errorf("Error Code %s\n", w.Code)
	      t.FailNow()
	  }
	  
	  //Test - PUT the above location 
	  locationJson = `{"address": "978", "city": "New York", "state": "New York","zipcode": ""}`
	  reader = strings.NewReader(locationJson) //Convert string to reader
	  
	  req, _ = http.NewRequest("PUT", uri, reader)
	
	  w = httptest.NewRecorder()
	  Router.ServeHTTP(w, req)
	
	  if !(w.Code == 201) {
	      t.Errorf("Failed to get response from server %s\n", w.Body.String())
	      t.Errorf("Error Code %s\n", w.Code)
	      t.FailNow()
	  }
	  
	  //Test - DELETE the above location
	  req, _ = http.NewRequest("DELETE", uri, nil)
	  w = httptest.NewRecorder()
	  Router.ServeHTTP(w, req)
	
	  if !(w.Code == http.StatusOK) {
	      t.Errorf("Failed to get response from server %s\n", w.Body.String())
	      t.Errorf("Error Code %s\n", w.Code)
	      t.FailNow()
	  }

	  //Test - GET and verify that the above location resource does not exist
	  req, _ = http.NewRequest("GET", uri, nil)
	  w = httptest.NewRecorder()
	  Router.ServeHTTP(w, req)
	
	  if !(w.Code == 404) {
	      t.Errorf("Failed to get response from server %s\n", w.Body.String())
	      t.Errorf("Error Code %s\n", w.Code)
	      t.FailNow()
	  }
}


func TestGetAll(t *testing.T) {
      //Setup httprouter and handlers for REST APIs for location service
	  Router := httprouter.New()
	  lc := controllers.NewLocationController(GetRemoteMGOSession())
	
	  Router.GET("/locations", lc.GetAllLocations)
	  Router.GET("/locations/:id", lc.GetLocation)
	  Router.POST("/locations", lc.CreateLocation)
	  Router.PUT("/locations/:id", lc.UpdateLocation)
	  Router.DELETE("/locations/:id", lc.RemoveLocation)
	
	  
	  //Test - GET all locations
	  req, _ := http.NewRequest("GET", "/locations", nil)
	  w := httptest.NewRecorder()
	  Router.ServeHTTP(w, req)
	
	  if !(w.Code == http.StatusOK) {
	      t.Errorf("Failed to get response from server %s\n", w.Body.String())
	      t.FailNow()
	  }
}


func GetRemoteMGOSession() *mgo.Session {
	  mongolab_uri :="mongodb://renubatheja:renubatheja@ds035844.mongolab.com:35844/cmpe273-assignment2"
	  session, err := mgo.Dial(mongolab_uri)
	  if err != nil {
	    fmt.Printf("Can't connect to mongo, go error %v\n", err)
	  }	
	  session.SetSafe(&mgo.Safe{})
	  return session
}
