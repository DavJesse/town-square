package messages

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/database"
	"forum/models"
)

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database.IsLoggedIn(r)
	senderID := session.UserID
	receiverIDStr := r.URL.Query().Get("receiver_id") // receiver's ID passed in query
	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.Response{Message: "Error: wrong recipient id"})
		return
	}
	message := r.FormValue("message")

	// Insert message into DB
	err = database.SendMessage(senderID, receiverID, message)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{Message: "Message sent successfully"})
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := database.IsLoggedIn(r)
	userID := session.UserID
	receiverIDStr := r.URL.Query().Get("receiver_id") // receiver's ID passed in query
	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.Response{Message: "Error: wrong recipient id"})
		return
	}
	limit := 10

	messages, err := database.GetMessages(userID, receiverID, limit, 0)
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
