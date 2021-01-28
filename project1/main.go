package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	 "helper1"
	 "models1"
	 "time"
	 
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


//Connection mongoDB with helper1 class

var collection = helper1.ConnectDB()

func main() {
	//Init Router
	r := mux.NewRouter()

	// arrange our route
	r.HandleFunc("/api/users", getUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/api/users", createUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")
	
	
	log.Fatal(http.ListenAndServe(":8000", r))
	

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var users []models1.User

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper1.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var user models1.User
		// & character returns the memory address of the following variable.
		err := cur.Decode(&user) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(users) // encode similar to serialize process.
}


func getUser(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var user models1.User
	
	// we get params with mux.
	var params = mux.Vars(r)
	
	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		// there was an error
		w.WriteHeader(404)
		w.Write([]byte("UserID not found"))
		return
	}

	//error checking
	// if id >= len(users) {
	// 	w.WriteHeader(404)
	// 	w.Write([]byte("No post found with specified ID"))
	// 	return
	// }


	if err != nil {
		helper1.GetError(err, w)
		return
	}
	//user := users[id]
    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    var httpUsers models1.HTTPUsers
	var user models1.User
	
	_ = json.NewDecoder(r.Body).Decode(&httpUsers)
	user = httpUsers.Data
    
	user.Created = time.Now()
	result, err := collection.InsertOne(context.TODO(), user)
	
	if err != nil {
		helper1.GetError(err, w)
		return
	}
	user.Created = time.Now()
	
	user.ID = result.InsertedID.(primitive.ObjectID)
	httpUsers.Data = user
	//json.NewEncoder(w).Encode(result)
	json.NewEncoder(w).Encode(httpUsers)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var httpUsers models1.HTTPUsers
	var user models1.User
    
	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&user)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"data", httpUsers.Data},
			{"$set", bson.D{
			{"firstname", user.Firstname},
			{"lastname", user.Lastname},
			{"email",user.Email},
			{"created",user.Created},
             }},
		}},
	}

				

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&httpUsers)

	if err != nil {
		helper1.GetError(err, w)
		return
	}

	user.ID = id

	json.NewEncoder(w).Encode(httpUsers)
}


func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper1.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
