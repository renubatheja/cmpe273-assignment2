package controllers

import (  
    "encoding/json"
    "fmt"
    "net/http"
	"gopkg.in/mgo.v2"
    "github.com/julienschmidt/httprouter"
    "gopkg.in/mgo.v2/bson"
    "model"
)

/*
* LocationController represents the controller for operating on the location resource
*/
type (  
    LocationController struct {  
    	session *mgo.Session
	}
)

func NewLocationController(s *mgo.Session) *LocationController {  
    return &LocationController{s}
}

/*
* GetAllLocations - retrieves all location resources stored in mongolab
*/
func (lc LocationController) GetAllLocations(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // Stub locations
    var locations []model.Location

    // Fetch locations
    if err := lc.session.DB("cmpe273-assignment2").C("locations").Find(bson.M{}).All(&locations); err != nil {
        w.WriteHeader(404)
        return
    }

	results := ""
	for index := 0; index < len(locations); index++ {
	    // Marshal provided interface into JSON structure
	    locationJson, _ := json.Marshal(locations[index])
		results = results + string(locationJson) + "\n "
	}
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", results)
	fmt.Println("All Location resources retrieved successfully from MongoLab!")
	fmt.Println("------------------------------------------------------------")	
}


/*
* GetLocation - retrieves an individual location resource (based on id) stored in mongolab
*/
func (lc LocationController) GetLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // Grab id
    id := p.ByName("id")

    // Verify id is ObjectId, otherwise bail
    if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    // Grab id
    oid := bson.ObjectIdHex(id)
    // Stub location
    location := model.Location{}

    // Fetch location
    if err := lc.session.DB("cmpe273-assignment2").C("locations").FindId(oid).One(&location); err != nil {
        w.WriteHeader(404)
        return
    }

    // Marshal provided interface into JSON structure
    locationJson, _ := json.Marshal(location)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", locationJson)
	fmt.Println("Location resource retrieved successfully from MongoLab!")
	fmt.Println("------------------------------------------------------------")
}

/*
* CreateLocation - creates an individual location resource and stores it in mongolab
*/
func (lc LocationController) CreateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // Stub a location to be populated from the body
    request := model.Location{}

    // Populate the location data
    json.NewDecoder(r.Body).Decode(&request)

    // Add an Id
    request.Id = bson.NewObjectId()

	//Call Google Maps API
	addressStr, fieldsMissing := FormatAddressString(request)
	if(fieldsMissing == true) {
		    newErrorResponse := model.Error{}
		    newErrorResponse.Code = "MISSING_REQUIRED_FIELDS"
		    newErrorResponse.Description = "Address fields are empty. Please enter a valid address."
		    newErrorResponse.Fieldname="Address,City,State,ZipCode"
		    newErrorJson, _ := json.Marshal(newErrorResponse)
		    	
		    w.Header().Set("Content-Type", "application/json")
		    w.WriteHeader(400)
		    fmt.Fprintf(w, "%s", newErrorJson)	
	} else {
		latitude, longitude, invalidAddress, errorCode, errorMsg := CallGoogleAPI(addressStr)
		if(invalidAddress) {	
		    newErrorResponse := model.Error{}
		    newErrorResponse.Code = errorCode
		    newErrorResponse.Description = errorMsg
		    newErrorResponse.Fieldname=""

		    newErrorJson, _ := json.Marshal(newErrorResponse)
		    	
		    w.Header().Set("Content-Type", "application/json")
		    w.WriteHeader(400)
		    fmt.Fprintf(w, "%s", newErrorJson)
			
		} else {			
			request.Coordinate.Latitude = latitude
			request.Coordinate.Longitude = longitude
			
		    //Write the location to mongo
		    lc.session.DB("cmpe273-assignment2").C("locations").Insert(request)
		    fmt.Println("Location data written successfully to MongoLab!")
			fmt.Println("------------------------------------------------------------")		
		
		    // Marshal provided interface into JSON structure
		    requestJson, _ := json.Marshal(request)
		
		    // Write content-type, statuscode, payload
		    w.Header().Set("Content-Type", "application/json")
		    w.WriteHeader(201)
		    fmt.Fprintf(w, "%s", requestJson)
	    }
    }
}


/*
* UpdateLocation - updates an individual location resource (based on id) stored in mongolab
*/
func (lc LocationController) UpdateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Grab id
    id := p.ByName("id")

    // Verify id is ObjectId, otherwise bail
    if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    // Grab id
    oid := bson.ObjectIdHex(id)

    // Stub a location to be populated from the body
    request := model.Location{}

    // Populate the location data
    json.NewDecoder(r.Body).Decode(&request)

	//Call Google Maps API
	addressStr, fieldsMissing := FormatAddressString(request)
	if(fieldsMissing == true) {
		    newErrorResponse := model.Error{}
		    newErrorResponse.Code = "MISSING_REQUIRED_FIELDS"
		    newErrorResponse.Description = "Address fields are empty. Please enter a valid address."
		    newErrorResponse.Fieldname="Address,City,State,ZipCode"
		    newErrorJson, _ := json.Marshal(newErrorResponse)
		    	
		    w.Header().Set("Content-Type", "application/json")
		    w.WriteHeader(400)
		    fmt.Fprintf(w, "%s", newErrorJson)	
	} else {
	
		latitude, longitude, invalidAddress, errorCode, errorMsg := CallGoogleAPI(addressStr)
		if(invalidAddress) {
		    newErrorResponse := model.Error{}
		    newErrorResponse.Code = errorCode
		    newErrorResponse.Description = errorMsg
		    newErrorResponse.Fieldname=""
		    newErrorJson, _ := json.Marshal(newErrorResponse)
		    	
		    w.Header().Set("Content-Type", "application/json")
		    w.WriteHeader(400)
		    fmt.Fprintf(w, "%s", newErrorJson)
			
		} else {	
			//fmt.Printf("Returned Values:%f,%f\n", (latitude),(longitude));
			
			request.Coordinate.Latitude = latitude
			request.Coordinate.Longitude = longitude
		    //Use the existing Id
		    request.Id = oid
	
			//To Update the location to mongodb
			collection := lc.session.DB("cmpe273-assignment2").C("locations")
		
		    // Stub location For now, just to fetch the "name"
		    location := model.Location{}	
		    if err := collection.FindId(oid).One(&location); err != nil {
		        fmt.Println("Error : Error received while fetching existing location resource: ",err)
				fmt.Println("------------------------------------------------------------")		        
		    } else {
			    var name string	    
			    if(request.Name == "") {
			    	name = location.Name
			    	//fmt.Println("No name was provided for update, so using the existing one : ",name)
			    } else {
			    	name = request.Name
			    	//fmt.Println("Name was provided for update, so using the new one : ",name)
			    }
						
			    // Fetch location
				oldLocation := bson.M{"id": oid}
			    if err := collection.FindId(oid).One(&oldLocation); err != nil {
			        w.WriteHeader(404)
			        return
			    }
				
				//fmt.Println("OldCollection : ", oldLocation)
				newLocation := bson.M{"$set": bson.M{"name": name , "address": request.Address, "city": request.City, "state": request.State, "zipcode": request.Zipcode,
										"coordinate" : bson.M{"lat": request.Coordinate.Latitude, "lng": request.Coordinate.Longitude} }}
										//bson.M{"$coordinate": bson.M{"lat": request.Coordinate.Latitude, "lng": request.Coordinate.Longitude } }
				//fmt.Println("NewCollection : ", newLocation)								
				err := collection.Update(oldLocation, newLocation)
				if err != nil {
					fmt.Println("Error : Error received while updating lcoation : ",err)
					fmt.Println("------------------------------------------------------------")					
				} else {   
				    //Fetch new location
				    newLocationResponse := model.Location{}	
				    if err := collection.FindId(oid).One(&newLocationResponse); err != nil {
				        fmt.Println("Error : Error received while fetching updated location!: ",err)
				    }
				
				    // Marshal provided interface into JSON structure
				    newLocationJson, _ := json.Marshal(newLocationResponse)
				
				    // Write content-type, statuscode, payload
				    w.Header().Set("Content-Type", "application/json")
				    w.WriteHeader(201)
				    fmt.Fprintf(w, "%s", newLocationJson)
					fmt.Println("Location resource updated successfully in MongoLab!")
					fmt.Println("------------------------------------------------------------")				    
			    }
		    }
		}
	}    
}

/*
* RemoveLocation - removes an individual location resource (based on id) stored in mongolab
*/
func (lc LocationController) RemoveLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // Grab id
    id := p.ByName("id")

    // Verify id is ObjectId, otherwise bail
    if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    // Grab id
    oid := bson.ObjectIdHex(id)

    // Remove user
    if err := lc.session.DB("cmpe273-assignment2").C("locations").RemoveId(oid); err != nil {
        w.WriteHeader(404)
        return
    }

    // Write status
    w.WriteHeader(200)
    fmt.Println("Location resource removed successfully from MongoLab!")
	fmt.Println("------------------------------------------------------------")    
}


