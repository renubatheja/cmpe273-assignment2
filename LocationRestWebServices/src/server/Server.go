package main

import (
    "fmt"
    "net/http"
    "gopkg.in/mgo.v2"
    "github.com/julienschmidt/httprouter"
    "controllers"
)

/*
* Function main - Setting up httprouter and REST API handlers (GET, POST, PUT, DELETE)
*/
func main() {    
    // Instantiate a new router
    r := httprouter.New()
    
    // Controller is going to need a mongo session to use in the CRUD methods. 
    // It would be connected to MongoLab.
    lc := controllers.NewLocationController(GetRemoteMGOSession())

	// Add handlers for REST webservices on 'location' resource
    // Get all location resources
    r.GET("/locations", lc.GetAllLocations)
    
    // Get a location resource identified by id
    r.GET("/locations/:id", lc.GetLocation)
	
	// Create a location resource
    r.POST("/locations", lc.CreateLocation)
	
	// Update a location resource identified by id
	r.PUT("/locations/:id", lc.UpdateLocation)
    
    // Delete a location resource identified by id
    r.DELETE("/locations/:id", lc.RemoveLocation)    
    
    // Fire up the server
    http.ListenAndServe("localhost:4000", r)    
}


//-----------------------Local MongoDb and MongoLab setup----------------------------
/*
* GetMGOSession - used in the server to connect to Local mongodb
*/
func GetMGOSession() *mgo.Session {  
    // Connect to our local mongo
    s, err := mgo.Dial("mongodb://localhost")

    // Check if connection error, is mongo running?
    if err != nil {
    	fmt.Printf("Error : Can't connect to mongo, go error %v\n", err)
        panic(err)
    }
    return s
}

/*
* GetRemoteMGOSession - used in the server to connect to Remote mongolab
*/
func GetRemoteMGOSession() *mgo.Session {
	mongolab_uri := "mongodb://renubatheja:renubatheja@ds035844.mongolab.com:35844/cmpe273-assignment2"
	session, err := mgo.Dial(mongolab_uri)
  	if err != nil {
    	fmt.Printf("Error : Can't connect to mongo, go error %v\n", err)
  	}
	
	session.SetSafe(&mgo.Safe{})
	return session
}

