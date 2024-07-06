package mobile_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vetner360-backend/database/mongodb"
	"vetner360-backend/model"
	"vetner360-backend/utils/helping"
	data_type "vetner360-backend/utils/type"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLatestDoctorClinic(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.Doctor]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}

	var filter = bson.M{"roleId": "665cec7fc6206b06eddaacca"}
	page := requestBody.Page
	limit := requestBody.Limit
	opts := options.FindOptions{}
	opts.Sort = bson.D{{"created_at", -1}}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	records, err := mongodb.GetAll[model.Doctor](&filter, &opts, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	total, err := mongodb.TotalDocs[model.User](&filter, "users")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	var requestResponse = data_type.Response[model.Doctor]{Status: true, Message: "Successfully Completed Request", Records: &records, Count: &total}
	jsonData, err := json.Marshal(requestResponse)

	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}

func GetLatestPetByType(response http.ResponseWriter, request *http.Request) {
	var requestBody data_type.PaginationType[model.PetSell]
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	validate := helping.GetValidator()
	err = helping.ValidatingData(requestBody, response, validate)
	if err != nil {
		return
	}

	page := requestBody.Page
	limit := requestBody.Limit
	skip := int64((page - 1) * limit)
	opts := options.AggregateOptions{}

	pipeline := mongo.Pipeline{}
	var petType = chi.URLParam(request, "type")
	var petTypeBson = bson.D{{"$match", bson.M{"petDetail.type": bson.M{"$ne": petType}}}}
	searchPetBson := bson.D{}

	if requestBody.Search != nil {
		searchPetBson = bson.D{{"$text", bson.M{"$search": *requestBody.Search}}}
	}

	if petType != "all" {
		petTypeBson = bson.D{{"$match", bson.M{"petDetail.type": petType}}}
	}

	fmt.Println(requestBody)

	if requestBody.Search != nil {
		fmt.Println("DF")
		pipeline = mongo.Pipeline{
			bson.D{{"$lookup", bson.M{
				"from":         "pets",
				"localField":   "petId",
				"foreignField": "token",
				"as":           "petDetail",
			}}},
			bson.D{{"$sort", bson.D{{"created_at", -1}}}},
			bson.D{{"$unwind", "$petDetail"}},
			bson.D{{"$limit", limit}},
			bson.D{{"$skip", skip}},
			petTypeBson,
			searchPetBson,
			bson.D{{
				"$replaceRoot", bson.M{
					"newRoot": bson.M{
						"$mergeObjects": bson.A{"$$ROOT", "$petDetail"},
					},
				},
			}},
		}
	} else {
		pipeline = mongo.Pipeline{
			bson.D{{"$lookup", bson.M{
				"from":         "pets",
				"localField":   "petId",
				"foreignField": "token",
				"as":           "petDetail",
			}}},
			bson.D{{"$sort", bson.D{{"created_at", -1}}}},
			bson.D{{"$unwind", "$petDetail"}},
			bson.D{{"$limit", limit}},
			bson.D{{"$skip", skip}},
			petTypeBson,
			bson.D{{
				"$replaceRoot", bson.M{
					"newRoot": bson.M{
						"$mergeObjects": bson.A{"$$ROOT", "$petDetail"},
					},
				},
			}},
		}
	}

	records, err := mongodb.GetAllUsingPipeline[model.PetSell](pipeline, &opts, "petsells")
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	if records == nil {
		records = []model.PetSell{}
	}

	var requestResponse = data_type.Response[model.PetSell]{Status: true, Message: "Successfully Completed Request", Records: &records}
	jsonData, err := json.Marshal(requestResponse)
	if err != nil {
		helping.InternalServerError(response, err, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "application/json")
	response.Write(jsonData)
}
