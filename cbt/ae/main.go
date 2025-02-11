package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jonreiter/govader"
)

// Struktur untuk request dan response
type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

// Variabel global untuk database dan NLP model
var (
	db       *sql.DB
	analyzer *govader.SentimentIntensityAnalyzer
)

func init() {
	var err error

	// Koneksi ke MySQL
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/mental_health")
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	// Cek koneksi
	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak bisa diakses:", err)
	}

	fmt.Println("Koneksi database berhasil!")

	// Inisialisasi Sentiment Analyzer
	analyzer = govader.NewSentimentIntensityAnalyzer()
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Analisis Sentimen
	sentimentScore := analyzer.PolarityScores(req.Message).Compound

	// Tentukan respons chatbot
	reply := analyzeMessage(req.Message, sentimentScore)

	// Simpan chat ke database
	saveChatToDB(req.Message, reply, sentimentScore)

	// Kirim response ke frontend
	res := ChatResponse{Reply: reply}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func analyzeMessage(msg string, score float64) string {
	if score < -0.2 {
		return "Sepertinya kamu sedang overthinking. Aku di sini untuk mendengarkanmu. Jika butuh bantuan lebih lanjut, coba hubungi konselor: 081234567890."
	} else if score > 0.2 {
		return "Wah, senang mendengar hal itu! Terus semangat ya! 💙"
	}
	return "Aku di sini untuk mendengarkanmu. Ceritakan apa yang sedang kamu pikirkan. 💙"
}

func saveChatToDB(userMessage, botResponse string, sentimentScore float64) {
	_, err := db.Exec("INSERT INTO chat_logs (user_message, bot_response, sentiment_score) VALUES (?, ?, ?)", userMessage, botResponse, sentimentScore)
	if err != nil {
		log.Println("Gagal menyimpan chat ke database:", err)
	}
}

func main() {
	http.HandleFunc("/chat", chatHandler)

	fmt.Println("Server berjalan di port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
