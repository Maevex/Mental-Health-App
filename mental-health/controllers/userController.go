package controllers

import (
	"encoding/json"
	"mental-health/config"
	middlewares "mental-health/middleware"
	"mental-health/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// struct request login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// struct response login
type LoginResponse struct {
	Token string `json:"token"`
}

// struct untuk view user
type UserResponse struct {
	Nama  string `json:"nama"`
	Email string `json:"email"`
}

// func login
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi input kosong
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		writeJSONError(w, "Email dan password tidak boleh kosong", http.StatusBadRequest)
		return
	}

	var user models.User
	err = config.DB.QueryRow("SELECT user_id, nama, email, password, role FROM user WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Nama, &user.Email, &user.Password, &user.Role)
	if err != nil {
		writeJSONError(w, "Email atau password salah", http.StatusUnauthorized)
		return
	}

	// Bandingkan password terenkripsi
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if user.Password != req.Password {
			writeJSONError(w, "Email atau password salah", http.StatusUnauthorized)
			return
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		writeJSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}

// Ambil semua user
func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT user_id, nama, email FROM user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nama, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// ambil data user dengan login
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	var user models.User
	err := config.DB.QueryRow("SELECT user_id, nama, email FROM user WHERE user_id = ?", userID).
		Scan(&user.ID, &user.Nama, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := UserResponse{
		Nama:  user.Nama,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// regist
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi input kosong
	if strings.TrimSpace(user.Nama) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		http.Error(w, "Nama, email, dan password tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// Enkripsi password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal mengenkripsi password", http.StatusInternalServerError)
		return
	}

	_, err = config.DB.Exec("INSERT INTO user (nama, email, password) VALUES (?, ?, ?)", user.Nama, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

//update user

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = config.DB.Exec("UPDATE user SET nama = ?, email = ?, password = ? WHERE user_id = ?",
		user.Nama, user.Email, user.Password, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	_, err := config.DB.Exec("DELETE FROM user WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func PublicCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Server is up and running!"}`))
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
