package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprintf(w, r.URL.Path+"\n")
	_, _ = fmt.Fprintf(w, "status: Online \n")
	_, _ = fmt.Fprintf(w, strconv.Itoa(app.config.port))

}

func (app *application) genericCRUD(w http.ResponseWriter, r *http.Request) {

	tokens := strings.Split(r.URL.Path, "/")
	path := "/" + tokens[1] + "/" + tokens[2]

	switch path {
	case "/v1/customers":
		objCreateOrUpdate[Customers](app, w, r)
	case "/v1/systems":
		objCreateOrUpdate[Systems](app, w, r)
	case "/v1/users":
		objCreateOrUpdate[Users](app, w, r)
	case "/v1/locations":
		objCreateOrUpdate[Locations](app, w, r)
	case "/v1/assets":
		objCreateOrUpdate[Assets](app, w, r)
	case "/v1/assetdata":
		objCreateOrUpdate[AssetData](app, w, r)
	case "/v1/assetdataresources":
		objCreateOrUpdate[AssetDataResources](app, w, r)
	case "/v1/reporttemplates":
		objCreateOrUpdate[ReportTemplates](app, w, r)
	case "/v1/reports":
		objCreateOrUpdate[Reports](app, w, r)
	case "/v1/serviceasset":
		objCreateOrUpdate[ServiceAsset](app, w, r)
	case "/v1/servicebillingpolicy":
		objCreateOrUpdate[ServiceBillingPolicy](app, w, r)
	case "/v1/serviceproviderattributes":
		objCreateOrUpdate[ServiceProviderAttributes](app, w, r)
	case "/v1/serviceprovider":
		objCreateOrUpdate[ServiceProvider](app, w, r)
	case "/v1/servicetypes":
		objCreateOrUpdate[ServiceTypes](app, w, r)
	case "/v1/services":
		objCreateOrUpdate[Services](app, w, r)
	case "/v1/sessions":
		objCreateOrUpdate[Sessions](app, w, r)
	case "/v1/systemdata":
		objCreateOrUpdate[SystemData](app, w, r)
	case "/v1/systemdataresources":
		objCreateOrUpdate[SystemDataResources](app, w, r)
	case "/v1/task":
		objCreateOrUpdate[Task](app, w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}

}

func objCreateOrUpdate[T any](app *application, w http.ResponseWriter, r *http.Request) {
	//Get json data from post request.
	var obj T
	var objs []T
	w.Header().Set("Content-Type", "application/json")

	//If the request is a get request, then query the database for the data and return the data json encoded.
	if r.Method == "GET" {
		//check if there is an ID in the URL
		tokens := strings.Split(r.URL.Path, "/")
		if len(tokens) > 3 {
			id, err := strconv.Atoi(tokens[3])
			if err != nil {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			query := app.db().Where("id = ?", id)
			query.Find(&obj)
			jsonEncoder := json.NewEncoder(w)
			err = jsonEncoder.Encode(&obj)
			if err != nil {
				return
			}
			return
		}
		query := app.db()

		//query the database for the data. and return the data json encoded.
		if len(r.URL.Query()) > 0 {
			for key, value := range r.URL.Query() {
				query = query.Where(key+" = ?", value[0])
			}
			query.Find(&objs)
		} else {
			query.Find(&objs)
		}

		//return array of objs
		jsonEncoder := json.NewEncoder(w)
		err := jsonEncoder.Encode(&objs)
		if err != nil {
			return
		}
		//curl command for testing
		//curl -X GET http://localhost:4000/v1/customers
		//curl -X GET http://localhost:4000/v1/customers?name=test

	}

	if r.Method == "POST" {
		//Get the json data from the post request.
		jsonDecoder := json.NewDecoder(r.Body)
		err := jsonDecoder.Decode(&obj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//Insert the data into the database.
		pk := app.db().Create(&obj)
		//query the database for the data. and return the data json encoded.
		app.db().First(&obj, pk)
		jsonEncoder := json.NewEncoder(w)
		err = jsonEncoder.Encode(&obj)
		if err != nil {
			return
		}
		//curl command for testing
		//curl -X POST -H "Content-Type: application/json" -d '{"name":"test"}' http://localhost:4000/v1/customers

	}
	if r.Method == "PATCH" {
		//Get the json data from the post request.
		jsonDecoder := json.NewDecoder(r.Body)
		customer := jsonDecoder.Decode(&obj)
		//Insert the data into the database.
		app.db().Save(&customer)

		jsonEncoder := json.NewEncoder(w)
		err := jsonEncoder.Encode(&customer)
		if err != nil {
			return
		}

	}

}
