package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const ()

type (
	AuthTokenRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		GrantType    string `json:"grant_type"`
	}

	AuthTokenResponse struct {
		AuthToken    string      `json:"access_token"`
		ExpiresIn    json.Number `json:"expires_in"`
		RefreshToken string      `json:"refresh_token"`
		TokenType    string      `json:"token_type"`
	}

	RefreshAuthTokenRequest struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}

	FunctionRoomGroupRequest struct {
		LocationIDs  []string `json:"LocationIds"`
		RecordStatus string   `json:"RecordStatus"`
	}

	DefiniteEventSearchRequest struct {
		BookingEventDateTimeBegin string `json:"BookingEventDateTimeBegin"`
		BookingEventDateTimeEnd   string `json:"BookingEventDateTimeEnd"`
		LocationId                string `json:"LocationId"`
	}

	ErrorResponse struct {
		Error     string `json:"error"`
		ErrorDesc string `json:"error_description"`
		GrantType string `json:"grant_type"`
		ErrorURI  string `json:"error_uri"`
	}

	LocationResponse struct {
		Name                               string      `json:"Name"`
		Status                             string      `json:"Status"`
		AddressLine1                       string      `json:"AddressLine1"`
		AddressLine2                       string      `json:"AddressLine2"`
		AddressLine3                       string      `json:"AddressLine3"`
		City                               string      `json:"City"`
		Country                            string      `json:"Country"`
		CountryCode                        string      `json:"CountryCode"`
		DistanceToNearestAirport           json.Number `json:"DistanceTonearestAirport"`
		DistanceUnitOfMeasure              string      `json:"DistanceUnitOfMeasure"`
		DrivetimeToNearestAirportInMinutes json.Number `json:"DrivetimeToNearestAirportInMinutes"`
		Fax                                string      `json:"Fax"`
		NearestAirportCode                 string      `json:"NearestAirportCode"`
		Phone                              string      `json:"Phone"`
		PostalCode                         string      `json:"PostalCode"`
		SizeUnitofMeasure                  string      `json:"SizeUnitofMeasure"`
		StateProvince                      string      `json:"StateProvince"`
		TimeZone                           string      `json:"TimeZone"`
		WebSiteUrl                         string      `json:"WebSiteUrl"`
		Id                                 string      `json:"Id"`
	}

	FunctionRoomGroupsResponse struct {
		Id                      string   `json:"Id"`
		ExternalId              string   `json:"ExternalId"`
		RecordStatus            string   `json:"RecordStatus"`
		FunctionRoomIds         []string `json:"FunctionRoomIds"`
		ExternalFunctionRoomIds []string `json:"ExternalFunctionRoomIds"`
		LocationId              string   `json:"LocationId"`
		ExternalLocationId      string   `json:"ExternalLocationId"`
		AlternateDescription    string   `json:"AlternateDescription"`
		AlternateName           string   `json:"AlternateName"`
		Description             string   `json:"Description"`
		Name                    string   `json:"Name"`
	}

	LocationFunctionRoomsResponse struct {
		ExternalId                            string      `json:"ExternalId"`
		ExternalCreatedById                   string      `json:"ExternalCreatedById"`
		ExternalCreatedOn                     string      `json:"ExternalCreatedOn"`
		ExternalModifiedById                  string      `json:"ExternalModifiedById"`
		CreatedById                           string      `json:"CreatedById"`
		CreatedOn                             string      `json:"CreatedOn"`
		ModifiedBy                            string      `json:"ModifiedBy"`
		ModifiedOn                            string      `json:"ModifiedOn"`
		Abbreviation                          string      `json:"Abbreviation"`
		Alias                                 string      `json:"Alias"`
		AlternateFunctionRoomName             string      `json:"AlternateFunctionRoomName"`
		Area                                  string      `json:"Area"`
		Comments                              string      `json:"Comments"`
		DefaultAdministrativeChargePercentage string      `json:"DefaultAdministrativeChargePercentage"`
		DefaultGratuityPercentage             string      `json:"DefaultGratuityPercentage"`
		DefaultSetupDurationMinutes           json.Number `json:"DefaultSetupDurationMinutes"`
		DefaultTeardownDurationMinutes        json.Number `json:"DefaultTeardownDurationMinutes"`
		ExternalBuildingId                    string      `json:"ExternalBuildingId"`
		DefaultEventSetupTypeId               string      `json:"DefaultEventSetupTypeId"`
		ExternalDefaultEventSetupTypeId       string      `json:"ExternalDefaultEventSetupTypeId"`
		ExternalLevelId                       json.Number `json:"ExternalLevelId"`
		Height                                json.Number `json:"Height"`
		ImageUri                              string      `json:"ImageUri"`
		Length                                json.Number `json:"Length"`
		MaxAccessHeight                       json.Number `json:"MaxAccessHeight"`
		MaxAccessWidth                        json.Number `json:"MaxAccessWidth"`
		MinimumCapacity                       json.Number `json:"MinimumCapacity"`
		MultiRoomBlockGroup                   string      `json:"MultiRoomBlockGroup"`
		Sequence                              json.Number `json:"Sequence"`
		WebSiteUrl                            string      `json:"WebSiteUrl"`
		Width                                 json.Number `json:"Width"`
		LocationId                            string      `json:"LocationId"`
		ExternalLocationId                    string      `json:"ExternalLocationId"`
		FunctionRoomType                      string      `json:"FunctionRoomType"`
		Name                                  string      `json:"Name"`
	}

	DefiniteEventSearchResponse struct {
		ExternalId                       string      `json:"ExternalId"`
		AccountName                      string      `json:"AccountName"`
		AlternateAccountName             string      `json:"AlternateAccountName"`
		AlternateEventClassificationName string      `json:"AlternateEventClassificationName"`
		AlternateFunctionRoomName        string      `json:"AlternateFunctionRoomName"`
		BookingPostAs                    string      `json:"BookingPostAs"`
		BookingTypeName                  string      `json:"BookingTypeName"`
		EndDateTime                      string      `json:"EndDateTime"`
		EventClassificationName          string      `json:"EventClassificationName"`
		ExternalAccountId                string      `json:"ExternalAccountId"`
		ExternalFunctionRoomId           string      `json:"ExternalFunctionRoomId"`
		FunctionRoomName                 string      `json:"FunctionRoomName"`
		LocationName                     string      `json:"LocationName"`
		StartDateTime                    string      `json:"StartDateTime"`
		AgreedAttendance                 json.Number `json:"AgreedAttendance"`
		AlternateName                    string      `json:"AlternateName"`
		Description                      string      `json:"Description"`
		EstimatedAttendance              json.Number `json:"EstimatedAttendance"`
		ForecastedAttendance             json.Number `json:"ForecastedAttendance"`
		GuaranteedAttendance             json.Number `json:"GuaranteedAttendance"`
		IsPosted                         bool        `json:"IsPosted"`
		Name                             string      `json:"Name"`
		SetAttendance                    json.Number `json:"SetAttendance"`
		ExternalBookingId                string      `json:"ExternalBookingId"`
		ExternalLocationId               string      `json:"ExternalLocationId"`
		Id                               string      `json:"Id"`
	}
)

const (
	kAHWSBaseURL                 string = "https://api-release.amadeus-hospitality.com"
	kAuthPath                    string = "/release/2.0/OAuth2"
	kAccessTokenPath             string = kAuthPath + "/AccessToken"
	kRefreshAccessTokenPath      string = kAuthPath + "/RefreshAccessToken"
	kAPIPath                     string = "/api/release"
	kLocationSearchPath          string = kAPIPath + "/Location/Search"
	kLocationsByExternalID       string = kAPIPath + "/location/ExternalLocationId"
	kLocationsByID               string = kAPIPath + "/location/LocationId"
	kFunctionRoomGroupSearchPath string = kAPIPath + "/functionroomgroup/Search"
	kFunctionRoomsSearchPath     string = kAPIPath + "/functionroom/Search"
	kDefiniteEventSearchPath     string = kAPIPath + "/bookingEvent/DefiniteEventSearch"
)

var (
	ocpApimSubscriptionKey string = os.Getenv("AHWS_APIM_SUBSCRIPTION_KEY")
	vAuthTokenResponse     AuthTokenResponse
)

func main() {
	vAuthTokenResponse = authenticate()

	log.Println("Authenticated with token:" + vAuthTokenResponse.AuthToken)
	log.Println("Authenticated expires_in:" + vAuthTokenResponse.ExpiresIn)
	log.Println("Authenticated refresh token:" + vAuthTokenResponse.RefreshToken)

	GetLocations()

	ticker := time.NewTicker(time.Hour * 71)

	// Create a channel to receive the tick events from the timer
	tickChannel := ticker.C

	// Start a for loop to range over the tick events from the channel
	for range tickChannel {
		// This block of code will be executed every time the timer fires

		// You can put any code you want to run on a timer here
		log.Println("Refreshing token")

		vAuthTokenResponse = refreshAccessToken(vAuthTokenResponse.RefreshToken)

		log.Println("Refreshed auth token:" + vAuthTokenResponse.AuthToken)
	}
}

func authenticate() (authTokenResponse AuthTokenResponse) {

	authRequest := AuthTokenRequest{
		ClientID:     os.Getenv("AHWS_CLIENT_ID"),
		ClientSecret: os.Getenv("AHWS_CLIENT_SECRET"),
		Username:     os.Getenv("AHWS_USERNAME"),
		Password:     os.Getenv("AHWS_PASSWORD"),
		GrantType:    "password",
	}

	var responseData AuthTokenResponse

	// Encode the dictionary as JSON.
	jsonRequestBody, err := json.Marshal(authRequest)
	LogError(err)

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kAccessTokenPath, jsonRequestBody))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

	return responseData
}

func refreshAccessToken(refreshAccessToken string) (authTokenResponse AuthTokenResponse) {
	refreshTokenRequest := RefreshAuthTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshAccessToken,
	}

	// Encode the dictionary as JSON.
	jsonRequestBody, err := json.Marshal(refreshTokenRequest)
	LogError(err)

	var responseData AuthTokenResponse

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kRefreshAccessTokenPath, jsonRequestBody))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

	return responseData
}

func GetLocations() {
	var responseData []LocationResponse

	body := doHTTPRequest(newHTTPPostRequest(kAHWSBaseURL + kLocationSearchPath))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
}

func GetLocationsID() {

	var responseData []LocationResponse
	body := doHTTPRequest(newHTTPGetRequest(kAHWSBaseURL + kLocationsByID))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

}

func GetLocationsByExternalID() {

	var responseData []LocationResponse
	body := doHTTPRequest(newHTTPGetRequest(kAHWSBaseURL + kLocationsByExternalID))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
}

func GetFunctionRoomGroup(IDs []string, recordStatus string) {

	functionRoomGroupsRequest := FunctionRoomGroupRequest{
		IDs,
		recordStatus,
	}

	jsonRequestBody, err := json.Marshal(functionRoomGroupsRequest)
	LogError(err)

	var responseData FunctionRoomGroupsResponse

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kFunctionRoomGroupSearchPath, jsonRequestBody))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
}

func GetBookingEventDetails(definiteEventSearchRequest DefiniteEventSearchRequest) {

	jsonRequestBody, err := json.Marshal(definiteEventSearchRequest)
	LogError(err)

	var responseData DefiniteEventSearchResponse

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kDefiniteEventSearchPath, jsonRequestBody))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

}

// func SetHeadersForReq() {
// 	pc, _, _, _ := runtime.Caller(1)
// 	callerMethod := runtime.FuncForPC(pc).Name()

// 	log.Println("callerMethod = " + callerMethod)

// 	switch callerMethod {
// 	case "GetLocations":

// 	}
// }

func LogPrettyPrintJSON(data interface{}) {

	b, err := json.MarshalIndent(data, "", "    ")

	// Marshal the map back into a JSON byte slice with indentation
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
}

func LogError(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

func newHTTPGetRequest(URI string) (req *http.Request) {
	return newHTTPRequest("GET", URI, nil)
}

func newHTTPPostRequest(URI string) (req *http.Request) {
	return newHTTPRequest("POST", URI, nil)
}

func newHTTPPostJSONRequest(URI string, json []byte) (req *http.Request) {
	return newHTTPRequest("POST", URI, json)
}

func newHTTPRequest(requestType string, URI string, json []byte) (req *http.Request) {
	req, err := http.NewRequest(requestType, URI, bytes.NewReader(json))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", ocpApimSubscriptionKey)

	if json != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if strings.Contains(URI, kAPIPath) {
		req.Header.Set("Authorization", "OAuth "+vAuthTokenResponse.AuthToken)
	}

	return req
}

func doHTTPRequest(req *http.Request) (body []byte) {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	return (body)
}

// func handleHTTPCode(statusCode int) {

// 	switch statusCode {
// 	case 100:
// 		log.Println("Continue")
// 	case http.StatusOK:
// 		log.Println("OK")
// 	case 300:
// 		log.Println("Multiple Choices")
// 	case 400:
// 		log.Println("Invalid request due to bad parameter values")
// 	case 403:
// 		log.Println("Unknown Username or Incorrect Password")
// 	case 500:
// 		log.Println("Internal Server Error")
// 	default:
// 		log.Println("Unknown status code")
// 	}
// }

// q := req.URL.Query()
// q.Add("subscription-key", ocpApimSubscriptionKey)
// req.URL.RawQuery = q.Encode()
// fmt.Println(req.URL)

