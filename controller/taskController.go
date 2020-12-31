package controller

import (
	"encoding/json"
	"fmt"
	"github.com/DebanjanBarman/todo/db"
	"github.com/DebanjanBarman/todo/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//GetTasks returns tasks in the collection depending on the query
//if there isn't any task it returns nothing
func GetTasks(writer http.ResponseWriter, request *http.Request) {
	cursor, err := db.TaskCollection.Find(db.Context, bson.M{})

	var tasks []bson.M

	err = cursor.All(db.Context, &tasks)
	handleError(err)

	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(tasks)
	handleError(err)
}

//GetTask takes the requested taskId and returns the task
//if the id is invalid it returns "Invalid id"
//if no document is found it returns nothing
func GetTask(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	var task bson.M

	taskId, idError := primitive.ObjectIDFromHex(id)

	if idError != nil {
		_, _ = fmt.Fprintf(writer, "Invalid id")
	}

	err := db.TaskCollection.FindOne(db.Context, bson.D{{"_id", taskId}}).Decode(&task)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(task)
	if err != nil {
		_, _ = fmt.Fprintf(writer, "Cannot encode json data")
	}
}

//CreateNewTask Creates a new task and returns the InsertedID/taskID
//if there is any error it returns an error
func CreateNewTask(writer http.ResponseWriter, request *http.Request) {
	var requestBody models.Task

	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	newTask := bson.D{
		{"text", requestBody.Text},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
		{"completed", requestBody.Completed},
	}
	res, err := db.TaskCollection.InsertOne(db.Context, newTask)
	if err != nil {
		log.Fatal(err)
	}

	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(res)
}

//UpdateTask updates a task and returns 202 status code
//if error occurs it returns the error
//if no document found it returns 404 status code
func UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task bson.M
	var requestBody models.Task

	id := mux.Vars(request)["id"]
	taskId, idError := primitive.ObjectIDFromHex(id)
	if idError != nil {
		_, _ = fmt.Fprintf(writer, "Invalid id")
	}

	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.D{
		{"_id", taskId},
	}

	updatedTask := bson.D{{"$set",
		bson.D{
			{"text", requestBody.Text},
			{"updated_at", time.Now()},
			{"completed", requestBody.Completed},
		},
	}}

	err = db.TaskCollection.FindOneAndUpdate(db.Context, filter, updatedTask).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}

	writer.WriteHeader(http.StatusAccepted)
}

//DeleteTask deletes a task
//if no document found it returns 404 status code
func DeleteTask(writer http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	var deletedDocument bson.M

	taskId, idError := primitive.ObjectIDFromHex(id)
	if idError != nil {
		_, _ = fmt.Fprintf(writer, "Invalid id")
	}
	filter := bson.D{{
		"_id", taskId,
	}}
	err := db.TaskCollection.FindOneAndDelete(db.Context, filter).Decode(&deletedDocument)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}
	writer.WriteHeader(http.StatusNoContent)
}
