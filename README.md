# Forum

## Description
This is a project about creating a web forum that allows:


- communication between users-comments
- associating categories to posts
- liking and disliking posts and comments
- filtering posts


## port difficulties

When you attempt to run the server e.g

```bash
go run .
```

then you face port difficulties, i.e the port is currently in use, you can switch to a different port easily by creating an environment variable called PORT.

```bash
export PORT=9000
```

you can now safely restart the server and it will use the port you just provided.

```bash
go run .
```


# **Forum Web Application**

## **Table of Contents**
- [Introduction](#introduction)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [System Architecture](#system-architecture)
- [File System Structure and Functionalities](#file-system-structure-and-functionalities)
- [Installation and Setup](#installation-and-setup)
  - [Prerequisites](#prerequisites)
  - [Cloning the Repository](#cloning-the-repository)
  - [Setting up the Database](#setting-up-the-database)
  - [Running the Application](#running-the-application)
- [Testing](#testing)
- [Usage Guide](#usage-guide)
  - [User Authentication](#user-authentication)
  - [Posting and Commenting](#posting-and-commenting)
  - [Liking and Disliking](#liking-and-disliking)
  - [Filtering](#filtering)
- [Docker](#docker)
- [Error Handling](#error-handling)
- [Contributors](#contributors)
- [License](#license)
- [Open Source Contributions](#open-source-contributions)

---

## **Introduction**
This project is a fully functional **web Forum** built using **Go, Javascript, MySQL, and Docker**. The forum allows users to communicate via posts and comments, categorize discussions, like/dislike content, and filter posts based on various criteria.

The forum allows:

- communication between users via posts and comments
- associating categories to posts
- liking and disliking posts and comments
- filter posts based on various criteria(liked posts, posts by categories)

## **Features**
- **User Authentication** 
  - Secure user registration and login
  - Session management
  - User profile management
- **Secure Password Storage** (bcrypt encryption)
- **Posting and Commenting** (Only for registered users)
- **Likes & Dislikes(posts and comments)** (Only for registered users)
- **Filtering Mechanism** (Categories, Created Posts, Liked Posts)
- **Public Access** (Non-registered users can view posts and comments)
- **Database Management with MySQL**
- **Containerized Deployment with Docker**
- **Error Handling and HTTP Status Codes**
- **Unit Testing**

## **Technologies Used**
- **Backend:** Go (Golang)
- **Database:** MySQL
- **Encryption:** bcrypt
- **Session Management:** Cookies
- **Containerization:** Docker
- **SQL queries** 
- **Web Design:** HTML, CSS,javascript (no frontend frameworks used)

## **System Architecture and File system**
1. **User Authentication Module** (Handles registration, login, and session management)
2. **Database Module** (Manages users, posts, media, comments, categories, likes, and dislikes)
3. **Web Server** (Handles HTTP requests and serves HTML pages)
4. **Filtering Mechanism** (Categorizes and sorts posts based on user preferences)

## **File System Structure and Functionalities**

- **main.go**: Entry point of the application, initializes the server and routes.


- **database/**: Contains all database-related files and queries.
  - `init.go`: Initializes the database connection and schema.
  - `session_*.go`: Handles user session management (login, logout, check status).
  - `posts_*.go`: Manages post creation, retrieval, and likes/dislikes.
  - `comment_*.go`: Manages comment creation and retrieval.
  - `user_*.go`: Manages user-related queries (registration, updates).

- **handlers/**: Contains HTTP request handlers for various routes.
  - `posts/`: Handles post-related requests (create, view, like, dislike).
  - `users/`: Handles user-related requests (profile, authentication).
  - `errors/`: Handles error responses and logging.

- **models/**: Defines data models used in the application.

- **web/**: Contains HTML templates and static assets (CSS, JS).

- **utils/**: Utility functions that support various functionalities throughout the application.

- **Dockerfile**: Instructions for building a Docker image for the application.
- **run-docker.sh**: Script for running the application inside a Docker container.

- **README.md**: Project documentation and guidelines.

## **Installation and Setup**

### **Prerequisites**
Ensure you have the following installed:
- [Go](https://golang.org/dl/)
- [MySQL](https://dev.mysql.com/downloads/)
- [Docker](https://www.docker.com/get-started)
- [Git](https://git-scm.com/)

### **Cloning the Repository**
```sh
git clone https://learn.zone01kisumu.ke/git/anoduor/forum.git

cd forum
```
### **Running the Application**
```sh
go mod tidy
go run main.go
```
Visit `http://localhost:8080` in your browser.

---

**port difficulties**

When you attempt to run the server e.g

```bash
go run .
```

then you face port difficulties, i.e the port is currently in use, you can switch to a different port easily by creating an environment variable called PORT.

```bash
export PORT=9000
```

you can now safely restart the server and it will use the port you just provided.

```bash
go run .
```

## **Testing**
Run unit tests using:
```bash
go test ./...
```

---

## **Forum UI Usage Guide**

### **User Authentication**
- **Register** with an email, username, and password.
- **Login** to create a session (stored in cookies).
- **Logout** to end the session.

### **Posting and Commenting**
- Only logged-in users can **create posts, like. dislike or add comments**.
- Posts can be assigned to **multiple categories**.
- Posts, Comments, likes and dislikes are **publicly visible**.

### **Liking and Disliking**
- Registered users can **like/dislike posts and comments**.
- The number of likes/dislikes is visible to **everyone**.

### **Filtering**
- **By Categories**: View posts under a specific category.
- **By Created Posts**: View only the posts you created under your profile.
- **By Liked Posts**: View posts you have liked.

---

## **Docker**
To run the project inside a Docker container:

### **Building the Docker Image**
```sh
docker build -t forum-app .
```

### **Running the Container**
```sh
docker run -p 8080:8080 forum-app
```

or

run the script to starts the container. The The forum will be available at `http://localhost:8080`.

```bash
./run-docker.sh
```

---

## **Error Handling**
- **User Authentication Errors**
  - Invalid email/password → **401 Unauthorized**
  - Email already registered → **409 Conflict**
- **Database Errors**
  - Connection failure → **500 Internal Server Error**
- **Post and Comment Errors**
  - Unauthenticated user trying to post → **403 Forbidden**
- **Invalid Routes**
  - Page not found → **404 Not Found**

---

## **Contributors**
- **Anne Okingo** ([github.com](https://github.com/Anne-Okingo))
- **David Jesse** ([github.com](https://github.com/DavJesse))
- **Cynthia Oketch** ([github.com](https://github.com/CynthiaOketch))
- **Rodney Ochieng** ([github.com](https://github.com/rodneyo1))
- **Antony Oduor**([github.com](https://github.com/oduortoni))

---

## **License**
This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.

---

## **Open Source Contributions**
We welcome contributions from the community! To contribute:
1. **Fork** the repository
2. **Create** a new branch (`feature-branch`)
3. **Commit** changes and push
4. **Submit** a pull request

For major changes, please open an issue first to discuss your proposal.


---
