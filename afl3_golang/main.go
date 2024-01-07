// Import Package
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux" // Untuk routing yang lebih fleksibel
    "encoding/json"
    "html/template"
	_ "github.com/go-sql-driver/mysql"
)

// Struct untuk data pengguna
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

// ErrorResponse untuk information ketika terjadi Pesan Kesalahan
type ErrorResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

var users []User

func main() {
    router := mux.NewRouter()
    
    // Endpoint 1: Menampilkan pesan sambutan
    router.HandleFunc("/", rootHandler).Methods("GET")

    // Endpoint 2: Menampilkan pesan sapaan sederhana
    router.HandleFunc("/hi", hiHandler).Methods("GET")

    // Endpoint 3: Menampilkan data pengguna berdasarkan ID
    router.HandleFunc("/users/{id}", getUserHandler).Methods("GET")

    // Endpoint 4: Membuat pengguna baru
    router.HandleFunc("/users/add", createUserHandler).Methods("POST")

    // Endpoint 5: Menampilkan daftar pengguna
    router.HandleFunc("/users", listUsersHandler).Methods("GET")

	// Endpoint 6: Mendelete semua users
    router.HandleFunc("/users/delete", deleteAllUsersHandler).Methods("DELETE")

    // Endpoint 7: Mencari pengguna berdasarkan query
    router.HandleFunc("/users/search", searchUsersHandler).Methods("GET")

    log.Fatal(http.ListenAndServe(":8081", router))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    /*
    Endpoint: /
    Method: GET
    Description: Menampilkan pesan sambutan.
    Response:
    {
        "message": "Hello, world!"
    }
    */
    // fmt.Fprintf(w, `{"message": "Hello, world!"}`)
    http.ServeFile(w, r, "views/index.html")
}

func hiHandler(w http.ResponseWriter, r *http.Request) {
    /*
    Endpoint: /hi
    Method: GET
    Description: Menampilkan pesan sapaan sederhana.
    Response:
    {
        "message": "Hi"
    }
    */
    fmt.Fprintf(w, `{"message": "Hi"}`)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    /*
    Endpoint: /user/{id}
    Method: GET
    Description: Menampilkan data pengguna berdasarkan ID.
    Response:
    {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com"
    }
    */
    vars := mux.Vars(r)
    userID := vars["id"]

    for _, user := range users {
        if fmt.Sprintf("%d", user.ID) == userID {
            // Found the user, return their details
            json.NewEncoder(w).Encode(user)
            return
        }
    }

    // User not found, return 404 status code
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, `{"error": "User not found"}`)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    // Parse form data from the request
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    // Get the values of "name" and "email" fields from the form
    name := r.FormValue("name")
    email := r.FormValue("email")

    // Create a new User
    newUser := User{
        // You may need to generate a unique ID for the user.
        // ID:   <generate a unique ID>,
        Name:  name,
        Email: email,
    }

    // Add the new user to your users slice or database
    users = append(users, newUser)

    // Respond with a success message
    response := map[string]string{"message": "User created successfully"}
    jsonResponse, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
    /*
    Endpoint: /users
    Method: GET
    Description: Menampilkan daftar pengguna.
    Response:
    [
        {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com"
        },
        {
            "id": 2,
            "name": "Jane Smith",
            "email": "jane@example.com"
        }
    ]
    */
    // json.NewEncoder(w).Encode(users)

    tmpl, err := template.ParseFiles("views/UserList.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Assuming you have a 'users' variable that contains the user data
    // users := []User{
    //     {ID: 1, Name: "John Doe", Email: "john@example.com"},
    //     {ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
    // }

    err = tmpl.Execute(w, users)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func deleteAllUsersHandler(w http.ResponseWriter, r *http.Request) {
    /*
    Endpoint: /users/delete
    Method: DELETE
    Description: Menghapus semua data pengguna.
    Response:
    {
        "message": "All users deleted"
    }
    */
    
    // Hapus semua data pengguna (implementasi terserah Anda).
    // Misalnya, Anda dapat mengosongkan slice "users".
    if len(users) == 0 {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, `{"error": "No users to delete"}`)
        return
    }

    // Clear the users slice
    users = []User{}

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"message": "All users deleted"}`)
}

func searchUsersHandler(w http.ResponseWriter, r *http.Request) {
    /*
    Endpoint: /users/search?query={query}
    Method: GET
    Description: Mencari pengguna berdasarkan query.
    Response:
    [
        {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com"
        },
        {
            "id": 2,
            "name": "Jane Smith",
            "email": "jane@example.com"
        }
    ]
    */

    query := r.URL.Query().Get("query")
    if query == "" {
        // Jika query kosong, kembalikan daftar pengguna lengkap
        json.NewEncoder(w).Encode(users)
        return
    }

    var searchResults []User

    // Lakukan pencarian berdasarkan query
    // for _, user := range users {
    //     if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) ||
    //         strings.Contains(strings.ToLower(user.Email), strings.ToLower(query)) {
    //         searchResults = append(searchResults, user)
    //     }
    // }

    json.NewEncoder(w).Encode(searchResults)
}
