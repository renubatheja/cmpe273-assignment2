package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "model"
    "io/ioutil"
    "strings"
)

//-------------------------------Structs to Process Google Maps API----------------------------
type Response struct {
	Results []ResultStruct `json:"results"`
}

type ResultStruct struct {
	Geometry GeometryStruct `json:"geometry"`
}

type GeometryStruct struct {
	Location LocationStruct `json:"location"`
}

type LocationStruct struct {
	Latitude float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}


/*
* Function to make Google Maps API call for fetching latitude and longitude for an address
*/
func CallGoogleAPI(FullAddressStr string) (float64, float64, bool, string, string){
	var query string
	var invalidAddressEntered bool
	invalidAddressEntered = false
	var errorDescription string
	var errorCode string
	
	errorDescription = ""
	errorCode = ""
	API_KEY := "AIzaSyDK-CekuqhDm14oObk81E2VEoL4uPHAHl4"

	query = "https://maps.google.com/maps/api/geocode/json?address=" + FullAddressStr + "&sensor=false&key=" + API_KEY
	resp, err := http.Get(query)
	
	if err != nil {
		invalidAddressEntered = true
		errorCode = "GOOGLE_MAPS_API_ERROR"
		errorDescription = "Error received while executing Google Maps API call. Please check the address."
		fmt.Println("Error : Error received while executing Google Maps API call! : ", err);
		fmt.Println("------------------------------------------------------------")
		return 0, 0, invalidAddressEntered, errorCode, errorDescription
	}
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		errorCode = "GOOGLE_MAPS_API_ERROR"
		errorDescription = "Unable to read response from Google Maps API. Please try again later."
		fmt.Println("Error : Error received while reading response from Google Maps API call! : ", err);
		fmt.Println("------------------------------------------------------------")	
		invalidAddressEntered = true
		return 0, 0, invalidAddressEntered, errorCode, errorDescription
	}
		
	var response Response
	errUnmarshal := json.Unmarshal(body, &response)
	if errUnmarshal != nil {
		errorCode = "GOOGLE_MAPS_API_ERROR"
		errorDescription = "Unable to unmarshall response from Google Maps API. Please try again later."
		fmt.Println("Error : Error received while unmashalling response from Google Maps API call!: ", errUnmarshal);
		fmt.Println("------------------------------------------------------------")	
		invalidAddressEntered = true
		return 0, 0, invalidAddressEntered, errorCode, errorDescription
	}
		
	var latitude float64
	var longitude float64
	if(len(response.Results) > 0){
		latitude = response.Results[0].Geometry.Location.Latitude
		longitude = response.Results[0].Geometry.Location.Longitude
		errorDescription = ""
	} else {
		errorCode = "INVALID_ADDRESS"
		errorDescription = "No results returned! Please enter a valid address"
		fmt.Println("Error : No results returned! Please enter a valid address")
		fmt.Println("------------------------------------------------------------")
		invalidAddressEntered = true
		return 0, 0, invalidAddressEntered, errorCode, errorDescription
	}
	
	return latitude, longitude, invalidAddressEntered, errorCode, errorDescription
}


/*
* Function for formatting address fields to be passed to Google Maps API
*/
func FormatAddressString(request model.Location) (string,bool) {
	//Grab and store request params in strings
	address := request.Address
	city := request.City
	zipcode := request.Zipcode
	state := request.State

	if(address == "" && city == "" && zipcode == "" && state == "") {
		fmt.Println("Error : Address fields are empty!")
		fmt.Println("------------------------------------------------------------")
		return "", true
	}
	
	//Make one string with ' ' delimiters
	s := []string{address, city, state, zipcode};
	fullAddressStr := strings.Join(s, "+");

	//add + sign in between every two words
	var addressWords []string
	addressWords = strings.Split(fullAddressStr, " ")
	fullAddress := ""
	for index := 0; index < len(addressWords); index++ {
		fullAddress = fullAddress + addressWords[index] + "+" 
	}
	return fullAddress, false
}