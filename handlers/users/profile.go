package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	"forum/models"
)

// ViewUserProfile handler
func ViewUserProfile(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Accept") == "application/json" {
        serveUserProfileJSON(w, r)
        return
    }

    // Serve index.html for frontend SPA routing
    http.ServeFile(w, r, "./web/templates/index.html")
}

func serveUserProfileJSON(w http.ResponseWriter, r *http.Request) {
    session, logged := database.IsLoggedIn(r)
    if !logged {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    userData, err := database.GetUserbySessionID(session.SessionID)
    if err != nil {
        log.Printf("Error getting user: %v\n", err)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(models.Response{
            Code:     http.StatusTemporaryRedirect,
            Message:  "Error getting user",
            Redirect: "/login",
        })
        return
    }

    UserPosts, err := database.PostsFilterByUser(userData.ID)
    if err != nil {
        log.Printf("Error getting posts: %v\n", err)
    }

    profileData := struct {
        User  models.User
        Posts []models.Post
    }{
        User:  userData,
        Posts: UserPosts,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.Response{
        Code:    http.StatusOK,
        Message: "Profile data retrieved",
        Data:    profileData,
    })
}

// func ViewUserProfile(w http.ResponseWriter, r *http.Request) {
// 	session, logged := database.IsLoggedIn(r)
// 	if !logged {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	userData, err := database.GetUserbySessionID(session.SessionID)
// 	// fmt.Printf("UserData retrieved: %+v\n", userData)  // Add debug logging
// 	if err != nil {
// 		log.Printf("Error getting user: %v\n", err) // Add error logging
// 		w.Header().Set("Content-Type", "application/json")
// 		r := models.Response{
// 			Code:     http.StatusTemporaryRedirect,
// 			Message:  "Error getting user",
// 			Redirect: "/login",
// 		}
// 		json.NewEncoder(w).Encode(r)
// 		// http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	UserPosts, err := database.PostsFilterByUser(userData.ID)
// 	if err != nil {
// 		log.Printf("Error getting posts: %v\n", err) // Add error logging
// 		// w.Header().Set("Content-Type", "application/json")
// 		// r := models.Response{
// 		// 	Code:     http.StatusOK,
// 		// 	Message:  "Error getting posts",
// 		// 	Redirect: "/login",
// 		// }
// 		// json.NewEncoder(w).Encode(r)
// 		// return
// 	}

// 	// Combine user data and user posts into a single struct
// 	profileData := struct {
// 		User  models.User
// 		Posts []models.Post
// 	}{
// 		User:  userData,
// 		Posts: UserPosts,
// 	}
// 	// Debug logs
// 	// fmt.Println(profileData.User)
// 	// fmt.Println("-------")
// 	// fmt.Println(profileData.Posts)

// 	w.Header().Set("Content-Type", "application/json")
// 	req := models.Response{
// 		Code:     http.StatusOK,
// 		Message:  "Error getting posts",
// 		Redirect: "",
// 		Data:     profileData,
// 	}
// 	json.NewEncoder(w).Encode(req)
// }
