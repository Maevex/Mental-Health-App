package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mental-health/config"
	middlewares "mental-health/middleware"
	"mental-health/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jonreiter/govader"
)

// Struct for user request

type UserRequest struct {
	Message string `json:"message"`
	SesiID  *int   `json:"sesi_id,omitempty"`
}

// Struct for gemini response parsing
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// function for chat handling
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

			// session validation: check if sesiID belongs to user
			if !isSesiMilikUser(db, userID, sesiID) {
				http.Error(w, "Forbidden: Sesi ini bukan milikmu", http.StatusForbidden)
				return
			}
		} else {

			sesiID = createNewSesi(db, userID)
		}

		// Analyze sentiment of user input
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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": chatbotResponse,
			"sesi_id":  sesiID,
		})

	}
}

// create new chat session
func createNewSesi(db *sql.DB, userID int) int {
	result, _ := db.Exec("INSERT INTO sesi_konsultasi (user_id) VALUES (?)", userID)
	sesiID64, _ := result.LastInsertId()
	return int(sesiID64)
}

// govader sentiment analysis
func analyzeSentiment(text string) float64 {
	analyzer := govader.NewSentimentIntensityAnalyzer()
	return analyzer.PolarityScores(text).Compound
}

// Save user message to database
func saveUserMessage(db *sql.DB, sesiID int, message string, sentiment float64) int {
	result, _ := db.Exec("INSERT INTO pesan_konsultasi (sesi_id, isi_pesan, sentiment_score) VALUES (?, ?, ?)", sesiID, message, sentiment)
	pesanID64, _ := result.LastInsertId()
	return int(pesanID64)
}

// Save chatbot response to database
func saveChatbotResponse(db *sql.DB, pesanID int, response string) {
	db.Exec("INSERT INTO respons_chatbot (pesan_id, isi_respons) VALUES (?, ?)", pesanID, response)
}

// Get response from Gemini API
func getGeminiResponse(userInput string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-2.0-flash:generateContent?key=" + apiKey

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

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes)) // log

	// Parse the response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}
	return "Tidak ada respons dari chatbot.", nil
}

// function to check session ownership
func isSesiMilikUser(db *sql.DB, userID, sesiID int) bool {
	var ownerID int

	err := db.QueryRow("SELECT user_id FROM sesi_konsultasi WHERE sesi_id = ?", sesiID).Scan(&ownerID)
	if err != nil {

		fmt.Println("Error query sesi:", err)
		return false
	}

	fmt.Printf("SesiID: %d, OwnerID: %d, UserID: %d\n", sesiID, ownerID, userID)

	return ownerID == userID
}

// GetSessions retrieves all sessions for the authenticated user
func GetSessions(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized: Missing user_id", http.StatusUnauthorized)
		return
	}
	fmt.Printf("UserID from context: %d\n", userID)

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

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(sessions)
}

// get chat detail for a specific session
func GetDetailSesiHandler(w http.ResponseWriter, r *http.Request) {
	db := config.DB

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized: Missing user_id", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	sesiIDStr := vars["id"]
	if sesiIDStr == "" {
		http.Error(w, "Missing sesi_id", http.StatusBadRequest)
		return
	}

	// session conversion to int
	sesiID, err := strconv.Atoi(sesiIDStr)
	if err != nil {
		http.Error(w, "Invalid sesi_id", http.StatusBadRequest)
		return
	}

	// session validation: check if sesiID belongs to user
	if !isSesiMilikUser(db, userID, sesiID) {
		http.Error(w, "Forbidden: Sesi ini bukan milikmu", http.StatusForbidden)
		return
	}

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

	// chat response struct
	type ChatData struct {
		Message         string `json:"message"`
		ChatbotResponse string `json:"chatbot_response"`
	}

	// slice to hold chat history
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatHistory)
}

// AnalisaKeluhanHandler handles the analysis of early user complaints
func AnalisaKeluhanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
		fmt.Printf("UserID from context: %d\n", userID)
		if !ok {
			http.Error(w, "Unauthorized: Missing user_id", http.StatusUnauthorized)
			return
		}

		var req struct {
			Kategori    string `json:"kategori"`
			Keluhan     string `json:"keluhan"`
			SubKategori string `json:"sub_kategori"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		sentimentScore := analyzeSentiment(req.Keluhan)
		prompt := fmt.Sprintf(`Sebagai seorang konsultan profesional di bidang kesehatan mental, analisis dan buatkan ringkasan masalah berdasarkan keluhan berikut.

		Kategori: %s
		Sub-kategori: %s
		Keluhan: %s

		Tuliskan ringkasan dari sudut pandang konsultan, lalu berikan saran atau langkah awal yang dapat dilakukan oleh pengguna untuk mengatasi masalah ini. Gunakan bahasa yang empatik dan mudah dimengerti.`, req.Kategori, req.SubKategori, req.Keluhan)

		kesimpulan, err := getGeminiResponse(prompt)

		if err != nil {
			http.Error(w, "Gagal mengambil kesimpulan dari Gemini", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`INSERT INTO keluhan_awal (user_id, kategori, keluhan, sentimen_score, kesimpulan_chatbot) 
	VALUES (?, ?, ?, ?, ?)`, userID, req.Kategori, req.Keluhan, sentimentScore, kesimpulan)
		if err != nil {
			fmt.Printf("DB error: %v\n", err)
			http.Error(w, "Gagal menyimpan keluhan", http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"kesimpulan": kesimpulan,
			"saran":      "Apakah kamu mau lanjut chat dengan chatbot kami?",
			"konsultan":  []map[string]string{},
		}

		if sentimentScore < -0.5 {
			rows, err := config.DB.Query(`SELECT nama, spesialisasi, pengalaman, no_telepon, email FROM konsultan_kontak`)
			if err != nil {
				http.Error(w, "Gagal mengambil data konsultan", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var konsultanList []map[string]string
			for rows.Next() {
				var nama, spesialisasi, pengalaman, noTelp, email string
				if err := rows.Scan(&nama, &spesialisasi, &pengalaman, &noTelp, &email); err != nil {
					http.Error(w, "Gagal membaca data konsultan", http.StatusInternalServerError)
					return
				}
				konsultanList = append(konsultanList, map[string]string{
					"nama":         nama,
					"spesialisasi": spesialisasi,
					"pengalaman":   pengalaman,
					"no_telepon":   noTelp,
					"email":        email,
				})
			}
			resp["rekomendasi"] = "Kami menyarankan kamu untuk menghubungi konsultan berikut:"
			resp["konsultan"] = konsultanList
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
