package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/patrickmn/go-cache"
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
		ExpiresAt    int64       `json:"-"`
	}

	ScheduleScreen map[string][]DefiniteEventSearchResponse
	CoverScreen    struct {
		EventName string
		StartTime string
		EndTime   string
	}

	Events struct {
		FunctionRoomGroup FunctionRoomGroupsResponse
		DefiniteEvent     DefiniteEventSearchResponse
	}

	RefreshAuthTokenRequest struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}

	FunctionRoomGroupRequest struct {
		LocationIDs []string `json:"LocationIds"`
		// RecordStatus []string `json:"RecordStatus"`
	}

	FunctionRoomRequest struct {
		LocationIDs  []string `json:"LocationIds"`
		RecordStatus string   `json:"-"`
	}

	DefiniteEventSearchRequest struct {
		BookingEventDateTimeBegin string `json:"BookingEventDateTimeBegin"`
		BookingEventDateTimeEnd   string `json:"BookingEventDateTimeEnd"`
		FunctionRoomGroupId       string `json:"FunctionRoomGroupId,omitempty"`
		LocationId                string `json:"LocationId"`
		MaxResultCount            int    `json:"MaxResultCount"`
	}

	DefiniteEventCacheSearchResults struct {
		BookingEventDateTime string `-`
		LocationId           string `-`
		Found                bool
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

	RoomGroups struct {
		RoomGroup string   `json:"RoomGroup"`
		Rooms     []string `json:"Rooms"`
	}
)

const (
	// kAHWSBaseURL string = "https://api-release.amadeus-hospitality.com"
	// kAuthPath    string = "/release/2.0/OAuth2"
	// kAPIPath     string = "/api/release"
	kAHWSBaseURL                 string        = "https://api.newmarketinc.com"
	kAuthPath                    string        = "/2.0/OAuth2"
	kAPIPath                     string        = "/api"
	kAccessTokenPath             string        = kAuthPath + "/AccessToken"
	kRefreshAccessTokenPath      string        = kAuthPath + "/RefreshAccessToken"
	kLocationSearchPath          string        = kAPIPath + "/Location/Search"
	kLocationsByExternalID       string        = kAPIPath + "/location/ExternalLocationId"
	kLocationsByID               string        = kAPIPath + "/location/LocationId"
	kFunctionRoomGroupSearchPath string        = kAPIPath + "/functionroomgroup/Search"
	kFunctionRoomsSearchPath     string        = kAPIPath + "/functionroom/Search"
	kDefiniteEventSearchPath     string        = kAPIPath + "/bookingEvent/DefiniteEventSearch"
	kTTL                         time.Duration = time.Minute * 15
)

const (
	DebugLevelNone = iota
	DebugLevelErrors
	DebugLevelVerbose
	DebugLevelTrace
)

const (
	CacheLevelNone = iota
	CacheLevelAuth
	CacheLevelAll
)

var (
	ocpApimSubscriptionKey string = os.Getenv("AHWS_APIM_SUBSCRIPTION_KEY")
	apiCache                      = cache.New(15*time.Minute, 20*time.Minute)
	debugLevel                    = DebugLevelTrace
	cacheLevel                    = CacheLevelAll
	hasTokenChan                  = make(chan bool)
	HasAuthToken                  = false
	roomGroups             []RoomGroups

	authRequest = AuthTokenRequest{
		ClientID:     os.Getenv("AHWS_CLIENT_ID"),
		ClientSecret: os.Getenv("AHWS_CLIENT_SECRET"),
		Username:     os.Getenv("AHWS_USERNAME"),
		Password:     os.Getenv("AHWS_PASSWORD"),
		GrantType:    "password",
	}
)

func main() {

	if authRequest.ClientID == "" ||
		authRequest.Username == "" ||
		authRequest.Password == "" ||
		authRequest.GrantType == "" ||
		ocpApimSubscriptionKey == "" {
		log.Panicln("FATAL: Environment Vars for authentication not set")
	}

	gob.Register(AuthTokenResponse{})
	gob.Register([]LocationResponse{})
	gob.Register([]FunctionRoomGroupsResponse{})
	gob.Register(FunctionRoomGroupsResponse{})
	gob.Register([]DefiniteEventSearchResponse{})
	gob.Register([]RoomGroups{})
	loadCacheGob()
	roomGroups, _ = loadJSONMapping()
	//amadeusIntegrationTesting()

	cancelChan := make(chan os.Signal, 1)

	// catch SIGTERM or SIGINT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go httpServer(cancelChan)
	sig := <-cancelChan
	log.Printf("Caught signal %v, waiting 3 seconds for graceful shutdown.", sig)

	saveCacheGob()

	time.Sleep(time.Second * 3)
	log.Println("Goodbye.")
}

func scheduleView(w http.ResponseWriter, definiteEvents []DefiniteEventSearchResponse) {
	events := ScheduleScreen{}
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.ParseFiles("schedule_screen.html.template")
	LogError(err)

	for _, event := range definiteEvents {
		if !event.IsPosted {
			continue
		}

		re := regexp.MustCompile("T(.*:.*):.*$")
		startTime := re.FindStringSubmatch(event.StartDateTime)
		endTime := re.FindStringSubmatch(event.EndDateTime)
		event.StartDateTime = startTime[1]
		event.EndDateTime = endTime[1]
		events[event.BookingPostAs] = append(events[event.BookingPostAs], event)
	}
	for _, v := range events {
		sort.Slice(v, func(i, j int) bool {
			return v[i].StartDateTime < v[j].StartDateTime
		})
	}
	LogError(tmpl.Execute(w, events))

}

func FindRoomInRoomGroups(roomName string, eventRoom string) bool {
	cached, found := apiCache.Get("RoomGroupByGroupName:" + eventRoom)
	if found {
		for _, room := range cached.(RoomGroups).Rooms {
			if room == roomName {
				return true
			}
		}
	}
	return false
}

func coverView(w http.ResponseWriter, roomId string, definiteEvents []DefiniteEventSearchResponse) {
	cs := CoverScreen{}
	w.Header().Add("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("cover_screen.html.template")
	LogError(err)

	for _, event := range definiteEvents {
		if !event.IsPosted || (roomId != event.FunctionRoomName && !FindRoomInRoomGroups(roomId, event.FunctionRoomName)) {
			continue
		}

		// NOTE: We don't know what tz... Defaults to EST, can accept the tz from the client if required.
		loc, err := time.LoadLocation("EST")
		LogError(err)

		now := time.Now().In(loc)

		re := regexp.MustCompile("T(.*:.*):.*$")
		startTime := re.FindStringSubmatch(event.StartDateTime)
		endTime := re.FindStringSubmatch(event.EndDateTime)

		startTimeSlice := strings.Split(startTime[1], ":")
		startTimeHour, _ := strconv.Atoi(startTimeSlice[0])
		startTimeMinute, _ := strconv.Atoi(startTimeSlice[1])

		endTimeSlice := strings.Split(endTime[1], ":")
		endTimeHour, _ := strconv.Atoi(endTimeSlice[0])
		endTimeMinute, _ := strconv.Atoi(endTimeSlice[1])

		startDate := time.Date(
			now.Year(),
			now.Month(),
			now.Day(), startTimeHour, startTimeMinute, 0, 0, loc)

		endDate := time.Date(
			now.Year(),
			now.Month(),
			now.Day(), endTimeHour, endTimeMinute, 0, 0, loc)

		currentTime := now

		// NOTE: For manipulating time for testing.
		// currentTime = time.Date(
		// 	time.Now().Year(),
		// 	time.Now().Month(),
		// 	time.Now().Add(time.Hour*24).Day(), 00, 13, 0, 0, loc)

		if startDate.Unix() <= currentTime.Unix() && currentTime.Unix() < endDate.Unix() {
			cs.StartTime = startTime[1]
			cs.EndTime = endTime[1]
			cs.EventName = event.Name
		}
	}

	if cs.EventName == "" {
		cs.EventName = "No Current Event"
	}

	err = tmpl.Execute(w, cs)
	LogError(err)
}

func httpServer(cancelChan chan<- os.Signal) {
	http.HandleFunc("/view/cover", func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("location-id") || !r.URL.Query().Has("room-id") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("location-id, and room-id must be provided"))
			return
		}

		definiteEvents, _ := GetBookingEventDetails(DefiniteEventSearchRequest{
			LocationId:                r.URL.Query().Get("location-id"),
			BookingEventDateTimeBegin: time.Now().Format("2006-01-02"),
			BookingEventDateTimeEnd:   time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		})
		coverView(w, r.URL.Query().Get("room-id"), definiteEvents)
	})

	http.HandleFunc("/view/schedule", func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("location-id") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("location-id must be provided"))
			return
		}

		w.Header().Add("Content-Type", "text/html")
		definiteEvents, _ := GetBookingEventDetails(DefiniteEventSearchRequest{
			LocationId:                r.URL.Query().Get("location-id"),
			FunctionRoomGroupId:       r.URL.Query().Get("group-id"),
			BookingEventDateTimeBegin: time.Now().Format("2006-01-02"),
			BookingEventDateTimeEnd:   time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		})
		scheduleView(w, definiteEvents)
	})

	// http.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Add("Content-Type", "application/json")
	// 	w.Write(locations.Bytes())
	// })
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

func GetAuthToken() string {

	authToken := ""
	if cacheLevel == CacheLevelNone {
		authToken = authenticate()
	} else {
		if cached, expiresAt, found := apiCache.GetWithExpiration("AccessToken"); found {
			DebugPrint("AccessToken.string:"+cached.(string), DebugLevelVerbose)
			DebugPrint("AccessToken.expiresAt:"+expiresAt.String(), DebugLevelVerbose)
			authToken = cached.(string)
		} else {
			if cached, expiresAt, found := apiCache.GetWithExpiration("RefreshAccessToken"); found {
				DebugPrint("RefreshAccessToken.string:"+cached.(string), DebugLevelVerbose)
				DebugPrint("RefreshAccessToken.expiresAt:"+expiresAt.String(), DebugLevelVerbose)
				authToken = refreshAccessToken(cached.(string))
			} else {
				authToken = authenticate()
			}
		}
	}

	if authToken == "" {
		// Replace with alert and retry
		DebugPrint("Failed to obtain auth token!", DebugLevelErrors)
		// return "", errors.New("Failed to obtain auth token!")
	}
	return authToken
}

// func authenticate() (authTokenResponse AuthTokenResponse) {

func authenticate() string {
	authRequest := AuthTokenRequest{
		ClientID:     os.Getenv("AHWS_CLIENT_ID"),
		ClientSecret: os.Getenv("AHWS_CLIENT_SECRET"),
		Username:     os.Getenv("AHWS_USERNAME"),
		Password:     os.Getenv("AHWS_PASSWORD"),
		GrantType:    "password",
	}

	var authTokenResponse AuthTokenResponse

	// Encode the dictionary as JSON.
	jsonRequestBody, err := MarshalAndLogWithErrorOutput(authRequest)
	if err != nil {
		return ""
	}

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kAccessTokenPath, jsonRequestBody))

	err = unMarshalAndLogWithErrorOutput(body, &authTokenResponse)
	if err != nil {
		return ""
	} else {
		apiCache.Set("AuthTokenResponse", authTokenResponse, cache.DefaultExpiration)
		apiCache.Set("AccessToken", authTokenResponse.AuthToken, cache.DefaultExpiration)
		apiCache.Set("RefreshAccessToken", authTokenResponse.RefreshToken, time.Hour*71)
	}
	return authTokenResponse.AuthToken
}

// func refreshAccessToken(refreshAccessToken string) (authTokenResponse AuthTokenResponse) {

func refreshAccessToken(refreshAccessToken string) string {
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

	apiCache.Set("AuthTokenResponse", responseData, cache.DefaultExpiration)
	apiCache.Set("AccessToken", responseData.AuthToken, cache.DefaultExpiration)
	apiCache.Set("RefreshAccessToken", responseData.RefreshToken, time.Hour*71)

	return responseData.AuthToken
}

func GetLocationsbyID() ([]LocationResponse, error) {

	if cacheLevel == CacheLevelAll {
		cached, expiresAt, found := apiCache.GetWithExpiration("LocationsByID")
		if found {
			DebugPrint("hit cache(LocationsByID) - expires at: "+expiresAt.String(), DebugLevelVerbose)
			LogPrettyPrintJSON(cached)
			return cached.([]LocationResponse), nil
		}
	}
	var locationResponse []LocationResponse
	body := doHTTPRequest(newHTTPGetRequest(kAHWSBaseURL + kLocationsByID))

	err := unMarshalAndLogWithErrorOutput(body, &locationResponse)
	if len(locationResponse) > 0 && cacheLevel == CacheLevelAll {
		apiCache.Set("LocationsByID", locationResponse, cache.DefaultExpiration)
	}
	return locationResponse, err
}

func GetLocationsByExternalID() ([]LocationResponse, error) {

	if cacheLevel == CacheLevelAll {
		cached, expiresAt, found := apiCache.GetWithExpiration("LocationsByExternalID")
		if found {
			DebugPrint("hit cache(LocationsByExternalID) - expires at: "+expiresAt.String(), DebugLevelVerbose)
			LogPrettyPrintJSON(cached)
			return cached.([]LocationResponse), nil
		}
	}

	var locationResponse []LocationResponse
	body := doHTTPRequest(newHTTPGetRequest(kAHWSBaseURL + kLocationsByExternalID))
	err := unMarshalAndLogWithErrorOutput(body, &locationResponse)

	if len(locationResponse) > 0 && cacheLevel == CacheLevelAll {
		apiCache.Set("LocationsByExternalID", locationResponse, cache.DefaultExpiration)
	}
	return locationResponse, err
}

func GetFunctionRoomGroup(IDs []string) ([]FunctionRoomGroupsResponse, error) {
	if len(IDs) == 0 {
		return nil, errors.New("empty array")
	}

	var IDsToQuery []string
	var cachedFunctionRoomGroups []FunctionRoomGroupsResponse

	if cacheLevel == CacheLevelAll {

		for _, key := range IDs {
			cached, expiresAt, found := apiCache.GetWithExpiration("FunctionRoomGroup:" + key)
			if found {
				DebugPrint("hit cache(FunctionRoomGroup) - expires at: "+expiresAt.String(), DebugLevelVerbose)
				cachedFunctionRoomGroups = append(cachedFunctionRoomGroups, cached.(FunctionRoomGroupsResponse))
			} else {
				IDsToQuery = append(IDsToQuery, key)
			}
		}
	} else {
		IDsToQuery = IDs
	}

	functionRoomGroupsRequest := FunctionRoomGroupRequest{
		IDsToQuery,
	}

	jsonRequestBody, err := MarshalAndLogWithErrorOutput(functionRoomGroupsRequest)
	if err != nil {
		return nil, err
	}

	var functionRoomGroupsResponse []FunctionRoomGroupsResponse

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kFunctionRoomGroupSearchPath, jsonRequestBody))

	err = unMarshalAndLogWithErrorOutput(body, &functionRoomGroupsResponse)

	if len(functionRoomGroupsResponse) > 0 && cacheLevel == CacheLevelAll {
		functionRoomGroupsResponse = append(functionRoomGroupsResponse, cachedFunctionRoomGroups...)
		for _, key := range functionRoomGroupsResponse {
			apiCache.Add("FunctionRoomGroup:"+key.Id, key, cache.DefaultExpiration)
		}
	}
	return functionRoomGroupsResponse, err
}

// func GetFunctionRoomGroups(functionRoomGroupsRequest FunctionRoomGroupRequest) []FunctionRoomGroupsResponse {
// 	var cachedFunctionRoomGroups []FunctionRoomGroupsResponse

// 	for idx, key := range functionRoomGroupsRequest.LocationIDs {
// 		cached, found := apiCache.Get("functionRoomGroups:AtLocation:" + key)
// 		if found {
// 			functionRoomGroupsRequest.LocationIDs = append(
// 				functionRoomGroupsRequest.LocationIDs[:idx],
// 				functionRoomGroupsRequest.LocationIDs[idx+1:]...)
// 			cachedFunctionRoomGroups = append(cachedFunctionRoomGroups, cached.([]FunctionRoomGroupsResponse)...)
// 		}
// 	}

// 	jsonRequestBody, err := json.Marshal(functionRoomGroupsRequest)
// 	LogError(err)

// 	var responseData []FunctionRoomGroupsResponse

// 	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kFunctionRoomGroupSearchPath, jsonRequestBody))

// 	LogError(json.Unmarshal(body, &responseData))
// 	LogPrettyPrintJSON(responseData)

// 	cacheStg := map[string][]FunctionRoomGroupsResponse{}
// 	for _, key := range responseData {
// 		cacheStg[key.LocationId] = append(cacheStg[key.LocationId], key)
// 	}

// 	for k, v := range cacheStg {
// 		apiCache.Add("functionRoomGroups:AtLocation:"+k, v, cache.DefaultExpiration)
// 	}

// 	responseData = append(responseData, cachedFunctionRoomGroups...)

// 	return responseData
// }

// kFunctionRoomsSearchPath
func GetFunctionRooms(req FunctionRoomRequest) []any {
	// var cachedFunctionRooms []any

	// for idx, key := range req.LocationIDs {
	// 	cached, found := apiCache.Get("functionRooms:AtLocation:" + key)
	// 	if found {
	// 		req.LocationIDs = append(
	// 			req.LocationIDs[:idx],
	// 			req.LocationIDs[idx+1:]...)
	// 		cachedFunctionRooms = append(cachedFunctionRooms, cached.([]FunctionRoomsResponse)...)
	// 	}
	// }

	jsonRequestBody, err := json.Marshal(req)
	LogError(err)

	//var responseData []FunctionRoomGroupsResponse

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+"/api/V2/FunctionRoom/Export", jsonRequestBody))
	fmtData := &bytes.Buffer{}
	json.Indent(fmtData, body, "", " ")
	fmt.Println(fmtData.String())

	//LogError(json.Unmarshal(body, &responseData))
	//LogPrettyPrintJSON(responseData)

	// cacheStg := map[string][]FunctionRoomsResponse{}
	// for _, key := range responseData {
	// 	cacheStg[key.LocationId] = append(cacheStg[key.LocationId], key)
	// }

	// for k, v := range cacheStg {
	// 	apiCache.Add("functionRoomGroups:AtLocation:"+k, v, cache.DefaultExpiration)
	// }

	// responseData = append(responseData, cachedFunctionRoomGroups...)

	//return responseData
	return nil
}

func DebugPrint(s any, logLevel int) {
	if debugLevel >= logLevel {
		log.Println("DEBUG:", s)
	}
}

func GetBookingEventDetails(definiteEventSearchRequest DefiniteEventSearchRequest) ([]DefiniteEventSearchResponse, error) {

	var cachedEvents []DefiniteEventSearchResponse
	var definiteEventSearchResponse []DefiniteEventSearchResponse
	var shouldQuery bool
	var err error

	eventRange := definiteEventSearchRequest.BookingEventDateTimeBegin + ":to:" + definiteEventSearchRequest.BookingEventDateTimeEnd

	if cacheLevel == CacheLevelAll {
		cached, expiresAt, found := apiCache.GetWithExpiration("BookingEventsDetailsWithInDateRange:" + eventRange + ":AtLocation:" + definiteEventSearchRequest.LocationId)
		if found {
			shouldQuery = false
			cachedEvents = cached.([]DefiniteEventSearchResponse)
			DebugPrint("hit cache(BookingEventsDetailsWithInDateRange) - expires at: "+expiresAt.String(), DebugLevelVerbose)
			LogPrettyPrintJSON(cachedEvents)
			return cachedEvents, nil
		} else {
			shouldQuery = true
		}
	}
	if shouldQuery {
		jsonRequestBody, err := MarshalAndLogWithErrorOutput(definiteEventSearchRequest)
		if err != nil {
			return nil, err
		}

		body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kDefiniteEventSearchPath, jsonRequestBody))
		err = unMarshalAndLogWithErrorOutput(body, &definiteEventSearchResponse)
		if len(definiteEventSearchResponse) > 0 && cacheLevel == CacheLevelAll {
			apiCache.Set("BookingEventsDetailsWithInDateRange:"+eventRange+":AtLocation:"+definiteEventSearchRequest.LocationId, definiteEventSearchResponse, cache.DefaultExpiration)
		}
		return definiteEventSearchResponse, err
	}
	return definiteEventSearchResponse, err
}

// func GetBookingEventDetails(definiteEventSearchRequest DefiniteEventSearchRequest) []DefiniteEventSearchResponse {
// 	var cachedEvents []DefiniteEventSearchResponse

// 	cached, found := apiCache.Get(
// 		"GetBookingEventsDetails:" +
// 			definiteEventSearchRequest.BookingEventDateTimeBegin +
// 			":to:" +
// 			definiteEventSearchRequest.BookingEventDateTimeEnd +
// 			":AtLocation:" +
// 			definiteEventSearchRequest.LocationId)
// 	if found {
// 		cachedEvents = cached.([]DefiniteEventSearchResponse)
// 		log.Println("hit cache: ")
// 		LogPrettyPrintJSON(cachedEvents)
// 		return cachedEvents
// 	} else {
// 		jsonRequestBody, err := json.Marshal(definiteEventSearchRequest)
// 		LogError(err)

// 		var responseData []DefiniteEventSearchResponse

// 		body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kDefiniteEventSearchPath, jsonRequestBody))

// 		LogError(json.Unmarshal(body, &responseData))
// 		LogPrettyPrintJSON(responseData)

// 		apiCache.Set(
// 			"GetBookingEventsDetails:"+
// 				definiteEventSearchRequest.BookingEventDateTimeBegin+
// 				":to:"+
// 				definiteEventSearchRequest.BookingEventDateTimeEnd+
// 				":AtLocation:"+
// 				definiteEventSearchRequest.LocationId,
// 			responseData, cache.DefaultExpiration)

// 		return responseData
// 	}
// }

func MarshalAndLogWithErrorOutput(request interface{}) (jsonRequestBody []byte, err error) {
	pc, _, _, _ := runtime.Caller(1)
	callerMethod := runtime.FuncForPC(pc).Name()

	DebugPrint("Marshing type:"+TypeName(request)+" - from function:"+callerMethod, DebugLevelVerbose)

	jsonRequestBody, err = json.Marshal(request)

	if err != nil {
		DebugPrint(err, DebugLevelErrors)
	} else {
		LogPrettyPrintJSON(request)
	}
	return jsonRequestBody, err
}
func unMarshalAndLogWithErrorOutput(body []byte, request interface{}) (err error) {
	pc, _, _, _ := runtime.Caller(1)
	callerMethod := runtime.FuncForPC(pc).Name()
	DebugPrint("UnMarshing type:"+TypeName(request)+" - from function:"+callerMethod, DebugLevelVerbose)

	json.Unmarshal(body, &request)

	if err != nil {
		DebugPrint(err, DebugLevelErrors)
	} else {
		LogPrettyPrintJSON(request)
	}
	return err
}

func LogPrettyPrintJSON(data interface{}) {
	log.Println(PrettyPrintJSON(data))
}

func PrettyPrintJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")

	// Marshal the map back into a JSON byte slice with indentation
	if err != nil {
		log.Println(err)
		return ""
	}
	return (string(b))
}

func LogError(err error) {
	if err != nil && debugLevel >= DebugLevelErrors {
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
		DebugPrint(err, DebugLevelErrors)
		return
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", ocpApimSubscriptionKey)

	// query := req.URL.Query()
	// query.Add("subscription-key", ocpApimSubscriptionKey)
	// req.URL.RawQuery = query.Encode()

	if json != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if strings.Contains(req.URL.Path, kAPIPath) {
		DebugPrint(URI, DebugLevelVerbose)
		req.Header.Set("Authorization", "OAuth "+GetAuthToken())
	}

	return req
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func doHTTPRequest(req *http.Request) (body []byte) {
	body, err := DoHTTPRequest(req)
	DebugPrint(err, DebugLevelErrors)
	return body
}

func DoHTTPRequest(req *http.Request) (body []byte, err error) {

	// var err error
	DebugPrint(req.URL, DebugLevelVerbose)
	req = req.WithContext(context.Background())

	err = httpDo(context.Background(), req, func(resp *http.Response, err error) error {
		if err != nil {
			DebugPrint(err, DebugLevelErrors)
			return err
		} else {
			pc, _, _, _ := runtime.Caller(1)
			callerMethod := runtime.FuncForPC(pc).Name()
			DebugPrint("HTTP code:"+resp.Status+" - in function:"+callerMethod, DebugLevelVerbose)
			if resp.StatusCode != 200 {
				DebugPrint("HTTP code:"+resp.Status+" - in function:"+callerMethod, DebugLevelErrors)
			}
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				DebugPrint(err, DebugLevelErrors)
				return err
			}
			if debugLevel == DebugLevelTrace {
				log.Print("request header")
				LogPrettyPrintJSON(req.Header)

				log.Print("request body:")
				LogPrettyPrintJSON(req.Body)

				log.Print("response header")
				LogPrettyPrintJSON(resp.Header)

				log.Print("response body:")
				LogPrettyPrintJSON(resp.Body)
			}
		}
		return err
	})
	return body, err
}

// func doHTTPRequest(req *http.Request) (body []byte) {
// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	pc, _, _, _ := runtime.Caller(1)
// 	callerMethod := runtime.FuncForPC(pc).Name()
// 	log.Println("HTTP code:" + resp.Status + " - in function:" + callerMethod)

// 	defer resp.Body.Close()
// 	body, err = io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	return (body)
// }

func saveCacheGob() {
	if cacheLevel != CacheLevelNone {
		DebugPrint("saving cache", DebugLevelVerbose)
		cachedItems := apiCache.Items()

		// Create a file to store the gob
		file, err := os.Create("cache.gob")
		if err != nil {
			DebugPrint(err, DebugLevelErrors)
		}
		defer file.Close()

		// Create a gob encoder
		encoder := gob.NewEncoder(file)

		// Encode the map
		err = encoder.Encode(cachedItems)
		if err != nil {
			DebugPrint(err, DebugLevelErrors)
		}
		DebugPrint("saved cache", DebugLevelVerbose)
	}
}

func loadCacheGob() {
	if cacheLevel != CacheLevelNone {
		DebugPrint("loading cache", DebugLevelVerbose)
		file, err := os.Open("cache.gob")
		if err != nil {
			DebugPrint(err, DebugLevelErrors)
		}
		defer file.Close()

		// Create a gob decoder
		decoder := gob.NewDecoder(file)

		// Decode the gob into a map
		var items map[string]cache.Item
		err = decoder.Decode(&items)
		if err == nil {
			apiCache = cache.NewFrom(15*time.Minute, 20*time.Minute, items)
		} else {
			DebugPrint(err, DebugLevelErrors)
		}
		DebugPrint("loaded cache", DebugLevelVerbose)
	}
}

func loadJSONMapping() ([]RoomGroups, error) {
	// Open the JSON file
	file, err := os.Open("mapping.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data into a struct
	var roomGroups []RoomGroups
	err = json.NewDecoder(file).Decode(&roomGroups)
	// log.Print(roomGroups)

	for _, roomGroup := range roomGroups {
		apiCache.Set("RoomGroupByGroupName:"+roomGroup.RoomGroup, roomGroup, cache.NoExpiration)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return roomGroups, nil
	}

}

func expirationTimeFromUnixTime(utime int64) time.Duration {
	t := time.Unix(utime, 0)
	// Calculate the duration from the Unix time to now
	return time.Since(t)
}

func TypeName(t any) (x string) {

	n := reflect.TypeOf(t).Name()

	if n == "" {
		return reflect.TypeOf(t).Elem().String()
	} else {
		return n
	}
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
