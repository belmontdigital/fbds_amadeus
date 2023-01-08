package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"syscall"
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

	Events struct {
		FunctionRoomGroup FunctionRoomGroupsResponse
		DefiniteEvent     DefiniteEventSearchResponse
	}

	RefreshAuthTokenRequest struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}

	FunctionRoomGroupRequest struct {
		LocationIDs  []string `json:"LocationIds"`
		RecordStatus string   `json:"-"` // TODO: add note
	}

	DefiniteEventSearchRequest struct {
		BookingEventDateTimeBegin string `json:"BookingEventDateTimeBegin"`
		BookingEventDateTimeEnd   string `json:"BookingEventDateTimeEnd"`
		FunctionRoomGroupId       string `json:"FunctionRoomGroupId"`
		LocationId                string `json:"LocationId"`
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
)

const (
	kAHWSBaseURL                 string        = "https://api-release.amadeus-hospitality.com"
	kAuthPath                    string        = "/release/2.0/OAuth2"
	kAccessTokenPath             string        = kAuthPath + "/AccessToken"
	kRefreshAccessTokenPath      string        = kAuthPath + "/RefreshAccessToken"
	kAPIPath                     string        = "/api/release"
	kLocationSearchPath          string        = kAPIPath + "/Location/Search"
	kLocationsByExternalID       string        = kAPIPath + "/location/ExternalLocationId"
	kLocationsByID               string        = kAPIPath + "/location/LocationId"
	kFunctionRoomGroupSearchPath string        = kAPIPath + "/functionroomgroup/Search"
	kFunctionRoomsSearchPath     string        = kAPIPath + "/functionroom/Search"
	kDefiniteEventSearchPath     string        = kAPIPath + "/bookingEvent/DefiniteEventSearch"
	kTTL                         time.Duration = time.Minute * 15
)

var (
	ocpApimSubscriptionKey string = os.Getenv("AHWS_APIM_SUBSCRIPTION_KEY")
	apiCache                      = cache.New(15*time.Minute, 20*time.Minute)
)

func main() {
	gob.Register(AuthTokenResponse{})
	gob.Register([]LocationResponse{})
	gob.Register([]FunctionRoomGroupsResponse{})
	gob.Register([]DefiniteEventSearchResponse{})

	log.Println("load cache")
	loadCacheGob()

	// AuthToken = GetAuthToken()
	// log.Println("Authenticated with token:" + AuthToken)

	// apiCache.LoadFile("cache.file")
	log.Println("loaded cache")

	// log.Println("flush cache")
	// apiCache.Flush()
	// log.Println("flushed")

	// programAborted := make(chan bool)
	cancelChan := make(chan os.Signal, 1)

	// catch SIGTERM or SIGINT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		http.HandleFunc("/view/cover", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "cover_screen.html")
		})

		http.HandleFunc("/view/schedule", func(w http.ResponseWriter, r *http.Request) {
			var locationIDs []string
			locations := GetLocations()
			for _, location := range locations {
				locationIDs = append(locationIDs, location.Id)
			}
			functionRoomGroups := GetFunctionRoomGroups(FunctionRoomGroupRequest{
				LocationIDs:  locationIDs,
				RecordStatus: "Active",
			})

			tmpl, err := template.ParseFiles("schedule_screen.html.template")
			LogError(err)
			events := ScheduleScreen{}
			for _, frg := range functionRoomGroups {
				definiteEvents := GetBookingEventDetails(DefiniteEventSearchRequest{
					LocationId:                frg.LocationId,
					FunctionRoomGroupId:       frg.Id,
					BookingEventDateTimeBegin: time.Now().Format("2006-01-02"),
					BookingEventDateTimeEnd:   time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
				})
				var e []DefiniteEventSearchResponse
				for _, event := range definiteEvents {
					re := regexp.MustCompile("T(.*:.*):.*$")
					startTime := re.FindStringSubmatch(event.StartDateTime)
					endTime := re.FindStringSubmatch(event.EndDateTime)
					event.StartDateTime = startTime[1]
					event.EndDateTime = endTime[1]
					e = append(e, event)
					log.Println(event.StartDateTime)
				}
				sort.Slice(e, func(i, j int) bool {
					return e[i].StartDateTime < e[j].StartDateTime
				})
				events[frg.Name] = e
			}
			err = tmpl.Execute(w, events)
			LogError(err)
			// http.ServeFile(w, r, "cover_screen.html")
		})

		// http.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		// 	w.Header().Add("Content-Type", "application/json")
		// 	w.Write(locations.Bytes())
		// })

		http.ListenAndServe(":8080", nil)
	}()

	sig := <-cancelChan
	log.Printf("Caught signal %v, waiting 3 seconds for graceful shutdown.", sig)

	log.Println("save cache")
	saveCacheGob()
	log.Println("saved cache")

	// programAborted <- true
	// close(programAborted)
	time.Sleep(time.Second * 3)
	log.Println("Goodbye.")
}

func GetAuthToken() string {
	if cached, expiresAt, found := apiCache.GetWithExpiration("AccessToken"); found {
		log.Println("AccessToken.string:" + cached.(string))
		log.Println("AccessToken.expiresAt:" + expiresAt.String())
		// AuthToken = cached.(string)
		return cached.(string)
	} else {
		if cached, expiresAt, found := apiCache.GetWithExpiration("RefreshAccessToken"); found {
			log.Println("RefreshAccessToken.string:" + cached.(string))
			log.Println("RefreshAccessToken.expiresAt:" + expiresAt.String())
			return refreshAccessToken(cached.(string))
		} else {
			return authenticate()
		}
	}
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

	var responseData AuthTokenResponse

	// Encode the dictionary as JSON.
	jsonRequestBody, err := json.Marshal(authRequest)
	LogError(err)

	log.Println(authRequest)

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kAccessTokenPath, jsonRequestBody))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

	apiCache.Set("AuthTokenResponse", responseData, cache.DefaultExpiration)
	apiCache.Set("AccessToken", responseData.AuthToken, cache.DefaultExpiration)
	apiCache.Set("RefreshAccessToken", responseData.RefreshToken, time.Hour*71)

	return responseData.AuthToken
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

func GetLocations() []LocationResponse {
	cached, expiresAt, found := apiCache.GetWithExpiration("GetLocations")
	if found {
		log.Println("hit cache(GetLocations) - expires in: " + expiresAt.String())
		LogPrettyPrintJSON(cached)
		return cached.([]LocationResponse)
	}

	var responseData []LocationResponse
	body := doHTTPRequest(newHTTPPostRequest(kAHWSBaseURL + kLocationSearchPath))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
	apiCache.Set("GetLocations", responseData, cache.DefaultExpiration)
	return responseData
}

func GetLocationsbyID() []LocationResponse {
	cached, found := apiCache.Get("GetLocationsByID")
	if found {
		log.Println("hit cache: ")
		LogPrettyPrintJSON(cached)
		return cached.([]LocationResponse)
	}

	var responseData []LocationResponse
	body := doHTTPRequest(newHTTPGetRequest(kAHWSBaseURL + kLocationsByID))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
	apiCache.Set("GetLocationsByID", responseData, cache.DefaultExpiration)
	return responseData
}

func GetLocationsByExternalID() []LocationResponse {
	cached, found := apiCache.Get("GetLocationsByExternalID")
	if found {
		log.Println("hit cache: ")
		LogPrettyPrintJSON(cached)
		return cached.([]LocationResponse)
	}

	var responseData []LocationResponse
	body := doHTTPRequest(newHTTPGetRequest(kAHWSBaseURL + kLocationsByExternalID))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)
	apiCache.Set("GetLocationsByExternalID", responseData, cache.DefaultExpiration)
	return responseData
}

func GetFunctionRoomGroups(functionRoomGroupsRequest FunctionRoomGroupRequest) []FunctionRoomGroupsResponse {
	var cachedFunctionRoomGroups []FunctionRoomGroupsResponse

	for idx, key := range functionRoomGroupsRequest.LocationIDs {
		cached, found := apiCache.Get("functionRoomGroups:AtLocation:" + key)
		if found {
			functionRoomGroupsRequest.LocationIDs = append(
				functionRoomGroupsRequest.LocationIDs[:idx],
				functionRoomGroupsRequest.LocationIDs[idx+1:]...)
			cachedFunctionRoomGroups = append(cachedFunctionRoomGroups, cached.([]FunctionRoomGroupsResponse)...)
		}
	}

	jsonRequestBody, err := json.Marshal(functionRoomGroupsRequest)
	LogError(err)

	var responseData []FunctionRoomGroupsResponse

	body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kFunctionRoomGroupSearchPath, jsonRequestBody))

	LogError(json.Unmarshal(body, &responseData))
	LogPrettyPrintJSON(responseData)

	cacheStg := map[string][]FunctionRoomGroupsResponse{}
	for _, key := range responseData {
		cacheStg[key.LocationId] = append(cacheStg[key.LocationId], key)
	}

	for k, v := range cacheStg {
		apiCache.Add("functionRoomGroups:AtLocation:"+k, v, cache.DefaultExpiration)
	}

	responseData = append(responseData, cachedFunctionRoomGroups...)

	return responseData
}

func GetBookingEventDetails(definiteEventSearchRequest DefiniteEventSearchRequest) []DefiniteEventSearchResponse {
	var cachedEvents []DefiniteEventSearchResponse

	cached, found := apiCache.Get("GetBookingEventsDetails:OnDate:" + definiteEventSearchRequest.BookingEventDateTimeBegin + ":AtLocation:" + definiteEventSearchRequest.LocationId)
	if found {
		cachedEvents = cached.([]DefiniteEventSearchResponse)
		log.Println("hit cache: ")
		LogPrettyPrintJSON(cachedEvents)
		return cachedEvents
	} else {
		jsonRequestBody, err := json.Marshal(definiteEventSearchRequest)
		LogError(err)

		var responseData []DefiniteEventSearchResponse

		body := doHTTPRequest(newHTTPPostJSONRequest(kAHWSBaseURL+kDefiniteEventSearchPath, jsonRequestBody))

		LogError(json.Unmarshal(body, &responseData))
		LogPrettyPrintJSON(responseData)

		apiCache.Set(
			"GetBookingEventsDetails:OnDate:"+
				definiteEventSearchRequest.BookingEventDateTimeEnd+
				":AtLocation:"+definiteEventSearchRequest.LocationId,
			responseData, cache.DefaultExpiration)

		return responseData
	}
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

	// query := req.URL.Query()
	// query.Add("subscription-key", ocpApimSubscriptionKey)
	// req.URL.RawQuery = query.Encode()

	if json != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if strings.Contains(URI, kAPIPath) {
		req.Header.Set("Authorization", "OAuth "+GetAuthToken())
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

	pc, _, _, _ := runtime.Caller(1)
	callerMethod := runtime.FuncForPC(pc).Name()
	log.Println("HTTP code:" + resp.Status + " - in function:" + callerMethod)

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		log.Println(err)
		return
	}
	return (body)
}

func saveCacheGob() {
	log.Println("save cache gob:")
	cachedItems := apiCache.Items()

	// Create a file to store the gob
	file, err := os.Create("cache.gob")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a gob encoder
	encoder := gob.NewEncoder(file)

	// Encode the map
	err = encoder.Encode(cachedItems)
	if err != nil {
		panic(err)
	}
}

func loadCacheGob() {
	log.Println("load cache gob:")
	file, err := os.Open("cache.gob")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a gob decoder
	decoder := gob.NewDecoder(file)

	// Decode the gob into a map
	var items map[string]cache.Item
	err = decoder.Decode(&items)
	if err == nil {
		apiCache = cache.NewFrom(15*time.Minute, 20*time.Minute, items)
	}
}

func expirationTimeFromUnixTime(utime int64) time.Duration {
	t := time.Unix(utime, 0)
	// Calculate the duration from the Unix time to now
	return time.Since(t)
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
