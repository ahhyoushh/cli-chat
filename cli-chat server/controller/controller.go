package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TanishkBansode/cli-chat/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database_string = os.Getenv("DATABASE_STRING")
var connectionString = database_string // ADD THE DATABASE STRING FROM MONGO DB

const dbName = "chatdata"
const chatcol = "chat"
const usercol = "user"

var chatcollection *mongo.Collection
var usercollection *mongo.Collection

// connect with MongoDB

func init() {
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected successfully!")

	chatcollection = client.Database(dbName).Collection(chatcol)
	usercollection = client.Database(dbName).Collection(usercol)

	fmt.Println("Collection reference is ready!")
}

// Helpers

func insertMessage(msg model.Conversations) string {
	val, err := chatcollection.InsertOne(context.Background(), msg)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted")
	insertedID := val.InsertedID.(primitive.ObjectID).Hex()

	return insertedID
}

func updateConversations(username string, newConvID string) {
	filter := bson.M{"username": username}
	update := bson.M{"$push": bson.M{"conversation": newConvID}}

	result, err := usercollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated %v\n", result.ModifiedCount)
}

func getAllMessages(username string) []model.Conversations {
	var user model.User
	filter := bson.M{"username": username}
	err := usercollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	option := options.FindOne().SetSort(bson.D{{Key: "time", Value: 1}})

	var messages []model.Conversations
	for _, msgID := range user.Conversation {
		fmt.Printf("msgID: %s\n", msgID)
		var message model.Conversations
		id, _ := primitive.ObjectIDFromHex(msgID)
		filter := bson.M{"_id": id}
		err := chatcollection.FindOne(context.Background(), filter, option).Decode(&message)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}

	return messages
}

func getUnreadMessages(username string) []model.Conversations {
	filter := bson.M{"receiver": username, "read": false}

	option := options.Find().SetSort(bson.D{{Key: "time", Value: 1}})

	// Find unread messages matching the filter
	cursor, err := chatcollection.Find(context.Background(), filter, option)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Decode messages from cursor
	var messages []model.Conversations
	for cursor.Next(context.Background()) {
		var message model.Conversations
		if err := cursor.Decode(&message); err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return messages
}

func SetAllMessagesRead(username string) {
	filter := bson.M{"receiver": username}
	update := bson.M{"$set": bson.M{"read": true}}

	result, err := chatcollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated %v documents to mark messages as read\n", result.ModifiedCount)
}

// now real stuff
type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetAllMsgs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds userCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Check if the provided username exists in the database
	var user model.User
	filter := bson.M{"username": creds.Username}
	err := usercollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify if the provided password matches the one associated with the username
	if user.Password != creds.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// If username and password are valid, proceed to get messages
	allMsgs := getAllMessages(user.Username)
	SetAllMessagesRead(user.Username)
	json.NewEncoder(w).Encode(allMsgs)
}

func GetUnreadMsgs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds userCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Check if the provided username exists in the database
	var user model.User
	filter := bson.M{"username": creds.Username}
	err := usercollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify if the provided password matches the one associated with the username
	if user.Password != creds.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// If username and password are valid, proceed to get unread messages
	unreadMsgs := getUnreadMessages(user.Username)
	SetAllMessagesRead(user.Username)
	json.NewEncoder(w).Encode(unreadMsgs)
}

func CreateMsg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var msg model.Conversations
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	msg.Time = time.Now()
	convid := insertMessage(msg)
	fmt.Println(convid)
	id, _ := primitive.ObjectIDFromHex(convid)
	msg.Conversation_ID = id

	// Check if receiver exists
	var user model.User
	filter := bson.M{"username": msg.Receiver}
	err := usercollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "Receiver not found", http.StatusNotFound)
		return
	}

	// Update conversations for both sender and receiver using their usernames
	if msg.Sender == msg.Receiver {
		updateConversations(msg.Sender, convid)
	} else {
		updateConversations(msg.Sender, convid)
		updateConversations(msg.Receiver, convid)
	}

	json.NewEncoder(w).Encode(msg)
}

// Handlers for user signup and login

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	filter := bson.M{"username": user.Username}
	var existingUser model.User
	err := usercollection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Insert new user into the database
	user.ID = primitive.NewObjectID()
	user.Conversation = []string{}
	_, err = usercollection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds userCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	filter := bson.M{"username": creds.Username}
	var user model.User
	err := usercollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify password
	if user.Password != creds.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Login successful")
}
