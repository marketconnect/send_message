package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RwBotToken string `env:"RW_BOT_TOKEN" env-required:"true"`
}

var cfg Config

func sendMessage(botToken string, chatID int64, message string) error {
	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	data := url.Values{}
	data.Set("chat_id", fmt.Sprintf("%d", chatID))
	data.Set("text", message)

	resp, err := http.PostForm(telegramAPI, data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		chatIDStr := r.FormValue("chat_id")
		message := r.FormValue("message")

		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Chat ID", http.StatusBadRequest)
			return
		}

		err = sendMessage(cfg.RwBotToken, chatID, message)
		if err != nil {
			http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Message sent successfully to Chat ID: %d", chatID)
		return
	}

	tmpl, err := template.ParseFiles("form.html")
	if err != nil {
		http.Error(w, "Failed to load form", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	http.HandleFunc("/", formHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
