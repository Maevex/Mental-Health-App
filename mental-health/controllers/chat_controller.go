package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"mental-health/config"
	middlewares "mental-health/middleware"
	"mental-health/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jonreiter/govader"
)

// Struct untuk menerima request dari user

type UserRequest struct {
	Message string `json:"message"`
	SesiID  *int   `json:"sesi_id,omitempty"`
}

// Struct untuk parsing response dari Gemini
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// ChatHandler menangani permintaan chat
func ChatHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
		fmt.Printf("UserID from context: %d\n", userID)
		if !ok {
			http.Error(w, "Unauthorized: Missing user_id", http.StatusUnauthorized)
			return
		}

		var userInput UserRequest
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var sesiID int
		if userInput.SesiID != nil {
			sesiID = *userInput.SesiID

			// Validasi sesi: Pastikan sesi milik user yang sedang login
			if !isSesiMilikUser(db, userID, sesiID) {
				http.Error(w, "Forbidden: Sesi ini bukan milikmu", http.StatusForbidden)
				return
			}
		} else {
			// Jika tidak ada sesiID, buat sesi baru untuk user yang sedang login
			sesiID = createNewSesi(db, userID)
		}

		sentimentScore := analyzeSentiment(userInput.Message)
		if sentimentScore < -0.5 {
			rows, err := config.DB.Query("SELECT nama, spesialisasi, pengalaman, no_telepon, email FROM konsultan_kontak")
			if err != nil {
				http.Error(w, "Gagal mengambil data konsultan", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var response string = "Sepertinya kamu sedang mengalami masalah serius. Jika kamu butuh bantuan, kamu bisa menghubungi:\n"
			for rows.Next() {
				var nama, spesialisasi, pengalaman, noTelepon, email string
				if err := rows.Scan(&nama, &spesialisasi, &pengalaman, &noTelepon, &email); err != nil {
					http.Error(w, "Gagal membaca data konsultan", http.StatusInternalServerError)
					return
				}
				response += fmt.Sprintf("*%s*\n Spesialisasi: %s\n Pengalaman: %s\n %s\n %s\n\n", nama, spesialisasi, pengalaman, noTelepon, email)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"response": response})
			return
		}

		chatbotResponse, err := getGeminiResponse(userInput.Message)
		if err != nil {
			http.Error(w, "Error getting chatbot response", http.StatusInternalServerError)
			return
		}

		pesanID := saveUserMessage(db, sesiID, userInput.Message, sentimentScore)
		saveChatbotResponse(db, pesanID, chatbotResponse)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"response": chatbotResponse})
	}
}

func createNewSesi(db *sql.DB, userID int) int {
	result, _ := db.Exec("INSERT INTO sesi_konsultasi (user_id) VALUES (?)", userID)
	sesiID64, _ := result.LastInsertId()
	return int(sesiID64)
}

func analyzeSentiment(text string) float64 {
	analyzer := govader.NewSentimentIntensityAnalyzer()
	return analyzer.PolarityScores(text).Compound
}

func saveUserMessage(db *sql.DB, sesiID int, message string, sentiment float64) int {
	result, _ := db.Exec("INSERT INTO pesan_konsultasi (sesi_id, isi_pesan, sentiment_score) VALUES (?, ?, ?)", sesiID, message, sentiment)
	pesanID64, _ := result.LastInsertId()
	return int(pesanID64)
}

func saveChatbotResponse(db *sql.DB, pesanID int, response string) {
	db.Exec("INSERT INTO respons_chatbot (pesan_id, isi_respons) VALUES (?, ?)", pesanID, response)
}

func getGeminiResponse(userInput string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-pro:generateContent?key=" + apiKey

	requestBody, _ := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": userInput}}},
		},
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var geminiResp GeminiResponse
	json.NewDecoder(resp.Body).Decode(&geminiResp)

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}
	return "Tidak ada respons dari chatbot.", nil
}

// Fungsi untuk memeriksa apakah sesi milik user
func isSesiMilikUser(db *sql.DB, userID, sesiID int) bool {
	var ownerID int
	// Perbaiki query, bandingkan sesi_id, bukan user_id
	err := db.QueryRow("SELECT user_id FROM sesi_konsultasi WHERE sesi_id = ?", sesiID).Scan(&ownerID)
	if err != nil {
		// Jika ada error (misalnya sesi tidak ditemukan), log error tersebut
		fmt.Println("Error query sesi:", err)
		return false
	}

	// Debugging: Cetak nilai ownerID untuk memeriksa apa yang didapatkan
	fmt.Printf("SesiID: %d, OwnerID: %d, UserID: %d\n", sesiID, ownerID, userID)

	return ownerID == userID
}

// GetSessionsHandler menangani permintaan untuk mendapatkan sesi milik user
func GetSessions(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari context
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized: Missing user_id", http.StatusUnauthorized)
		return
	}
	fmt.Printf("UserID from context: %d\n", userID)

	// Query untuk mengambil semua sesi yang dimiliki oleh user
	rows, err := config.DB.Query("SELECT sesi_id, user_id, status, waktu_mulai FROM sesi_konsultasi WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sessions []models.SesiKonsultasi
	for rows.Next() {
		var session models.SesiKonsultasi
		if err := rows.Scan(&session.SesiID, &session.UserID, &session.Status, &session.WaktuMulai); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessions = append(sessions, session)
	}

	// Set header untuk response
	w.Header().Set("Content-Type", "application/json")
	// Kembalikan hasil dalam format JSON
	json.NewEncoder(w).Encode(sessions)
}

func GetDetailSesiHandler(w http.ResponseWriter, r *http.Request) {
	db := config.DB

	// Ambil user_id dari context
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized: Missing user_id", http.StatusUnauthorized)
		return
	}

	// Ambil sesi_id dari parameter URL
	vars := mux.Vars(r)
	sesiIDStr := vars["id"] // Ambil sesi_id dari path parameter "id"
	if sesiIDStr == "" {
		http.Error(w, "Missing sesi_id", http.StatusBadRequest)
		return
	}

	// Convert sesi_id dari string ke int
	sesiID, err := strconv.Atoi(sesiIDStr)
	if err != nil {
		http.Error(w, "Invalid sesi_id", http.StatusBadRequest)
		return
	}

	// Validasi: Pastikan sesi ini milik user
	if !isSesiMilikUser(db, userID, sesiID) {
		http.Error(w, "Forbidden: Sesi ini bukan milikmu", http.StatusForbidden)
		return
	}

	// Ambil semua pesan dan respons di sesi ini
	rows, err := db.Query(`
		SELECT 
			p.isi_pesan, 
			r.isi_respons 
		FROM 
			pesan_konsultasi p
		LEFT JOIN 
			respons_chatbot r ON p.pesan_id = r.pesan_id
		WHERE 
			p.sesi_id = ?
		ORDER BY 
			p.pesan_id ASC`, sesiID)
	if err != nil {
		http.Error(w, "Gagal mengambil pesan", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Bentuk response
	type ChatData struct {
		Message         string `json:"message"`
		ChatbotResponse string `json:"chatbot_response"`
	}

	var chatHistory []ChatData
	for rows.Next() {
		var chat ChatData
		err := rows.Scan(&chat.Message, &chat.ChatbotResponse)
		if err != nil {
			http.Error(w, "Gagal membaca data pesan", http.StatusInternalServerError)
			return
		}
		chatHistory = append(chatHistory, chat)
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatHistory)
}
