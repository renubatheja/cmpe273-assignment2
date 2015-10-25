package model

import (
	"gopkg.in/mgo.v2/bson"
)

/*
* Location resource - structure of location resource (Used in Request and Response formats)
*/
type (  
    Location struct {
        Id     bson.ObjectId `json:"id" bson:"_id"`
        Name   string        `json:"name" bson:"name"`
        Address string       `json:"address" bson:"address"`
        City    string       `json:"city" bson:"city"`
        State    string      `json:"state" bson:"state"`
        Zipcode string		 `json:"zipcode" bson:"zipcode"`
        Coordinate Coordinate `json:"coordinate" bson:"coordinate"`
    }    

	Coordinate struct {   
        	Latitude float64  `json:"lat" bson:"lat"`
        	Longitude float64 `json:"lng" bson:"lng"`
    }
	
)
/*
* Structure of Error - to be sent if the location service encounters any errors    
*/
type Error struct {
    	Code string `json:"code" bson:"code"`
    	Description string `json:"description" bson:"description"`
    	Fieldname string `json:"fieldname" bson:"fieldname"`
}
