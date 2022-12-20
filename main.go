package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const ()

type (
	Fixturer interface {
		FixturePath() string
		Bytes() []byte
	}

	CreateAuthToken struct {
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

	RefreshAuthToken struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}

	ConsumerCache struct {
		Data      []byte    `json:"-"`
		UpdatedAt time.Time `json:"-"`
	}

	GetLocations struct {
		ConsumerCache
	}
	GetLocationsByID struct {
		ConsumerCache
	}
	GetLocationsByExternalID struct {
		ConsumerCache
	}

	GetFunctionRoomGroups struct {
		LocationIds  []string `json:"LocationIds"`
		RecordStatus string   `json:"-"`
		ConsumerCache
	}

	GetDefiniteEvents struct {
		BookingEventDateTimeBegin string `json:"BookingEventDateTimeBegin"`
		BookingEventDateTimeEnd   string `json:"BookingEventDateTimeEnd"`
		LocationId                string `json:"LocationId"`
		ConsumerCache
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
		ExternalId                         string      `json:"ExternalId"`
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

func (authRequest CreateAuthToken) Do() (responseData AuthTokenResponse) {
	body := authRequest.Bytes()
	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

	return responseData
}

func (authRequest CreateAuthToken) Bytes() []byte {
	jsonRequestBody, err := json.Marshal(authRequest)
	LogError(err)

	return doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kAccessTokenPath, jsonRequestBody))
}

func (CreateAuthToken) FixturePath() string {
	return "fixtures/authtoken.json"
}

func (refreshTokenRequest RefreshAuthToken) Do() (responseData AuthTokenResponse) {
	jsonRequestBody, err := json.Marshal(refreshTokenRequest)
	LogError(err)

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kRefreshAccessTokenPath, jsonRequestBody))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

	return responseData
}

func (r GetLocations) Do() (responseData []LocationResponse) {
	body := r.Bytes()
	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

	return responseData
}

func (r GetLocations) Bytes() []byte {
	return r.bytes()
}

func (r *GetLocations) bytes() []byte {
	if len(r.Data) != 0 {
		return r.Data
	}

	r.Data = doHTTPRequest(newHTTPPostRequest(kAHWSBaseURL + kLocationSearchPath))
	return r.Data
}

func (GetLocations) FixturePath() string {
	return "fixtures/locations.json"
}

func (r GetLocationsByID) Do() {
	var responseData []LocationResponse
	body := r.Bytes()

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
}

func (r GetLocationsByID) Bytes() []byte {
	return r.bytes()
}

func (r *GetLocationsByID) bytes() []byte {
	// TODO: Implement proper caching here
	if len(r.Data) != 0 {
		return r.Data
	}

	r.Data = doHTTPRequest(newHTTPPostRequest(kAHWSBaseURL + kLocationsByID))
	return r.Data
}

func (GetLocationsByID) FixturePath() string {
	return "fixtures/locationsbyid.json"
}

func (r GetLocationsByExternalID) Do() {
	var responseData []LocationResponse
	body := r.Bytes()

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
}

func (r GetLocationsByExternalID) Bytes() []byte {
	return r.bytes()
}

func (r *GetLocationsByExternalID) bytes() []byte {
	// TODO: Implement proper caching here
	if len(r.Data) != 0 {
		return r.Data
	}

	r.Data = doHTTPRequest(newHTTPPostRequest(kAHWSBaseURL + kLocationsByExternalID))
	return r.Data
}

func (GetLocationsByExternalID) FixturePath() string {
	return "fixtures/locationsbyexternalid.json"
}

func (f GetFunctionRoomGroups) Do() {
	var responseData FunctionRoomGroupsResponse
	LogError(json.Unmarshal(f.Bytes(), &responseData))
	LogPrettyPrintJSON(responseData)
}

func (r GetFunctionRoomGroups) Bytes() []byte {
	return r.bytes()
}

func (r *GetFunctionRoomGroups) bytes() []byte {
	// TODO: Implement proper caching here
	jsonRequestBody, err := json.Marshal(r)
	LogError(err)

	if len(r.Data) != 0 {
		return r.Data
	}

	r.Data = doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kFunctionRoomGroupSearchPath, jsonRequestBody))
	return r.Data
}

func (f GetFunctionRoomGroups) FixturePath() string {
	return "fixtures/functionroomgroups.json"
}

func (r GetDefiniteEvents) Do() {
	var responseData DefiniteEventSearchResponse
	LogError(json.Unmarshal(r.Bytes(), &responseData))
	LogPrettyPrintJSON(responseData)
}

func (r GetDefiniteEvents) Bytes() []byte {
	return r.bytes()
}

func (r *GetDefiniteEvents) bytes() []byte {
	jsonRequestBody, err := json.Marshal(r)
	LogError(err)

	// TODO: Implement proper caching here
	if len(r.Data) != 0 {
		return r.Data
	}

	return doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kDefiniteEventSearchPath, jsonRequestBody))
}

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

func genFixtures() {
	// populateFixture(authTokenRequest)

	// Populate location fixture
	locationClient := GetLocations{}
	populateFixture(locationClient)

	// Decode location data for parsing
	var locations []LocationResponse
	locationData := locationClient.Bytes()

	json.Unmarshal(locationData, &locations)

	// Populate function room group fixture for all locations
	var locationIds []string
	for _, location := range locations {
		locationIds = append(locationIds, location.Id)
	}

	functionRoomGroups := GetFunctionRoomGroups{
		LocationIds: locationIds,
	}
	populateFixture(functionRoomGroups)
}

func populateFixture(fixturer Fixturer) {
	resp := fixturer.Bytes()

	var out bytes.Buffer
	err := json.Indent(&out, resp, "", "    ")
	LogError(err)

	log.Println("Populating fixture:", fixturer.FixturePath(), "of size:", len(resp))

	f, err := os.Create(fixturer.FixturePath())
	LogError(err)

	f.Write(out.Bytes())
	log.Println("Successfully wrote fixture:", fixturer.FixturePath())
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
