package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const OpenWeatherMapAPIKey = "0e264e7da407afac79c77c3f2d1feeeb"

// Creating a struct
type WeatherMapApiResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func getWeather(city string) (string, error) {
	// Try to open WeatherAPI URL
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, OpenWeatherMapAPIKey)

	//Try to request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parsing the response

	var weatherData WeatherMapApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return "", err
	}

	// Forming the result
	return fmt.Sprintf("The weather in %s is %s with a temperature of %.1f°C", city, weatherData.Weather[0].Description, weatherData.Main.Temp), nil
}

func main() {

	// Insert Telegram Bot Token
	const TELEGRAM_BOT_TOKEN = "7850205681:AAG_tWa2WJg8j1eFMGrxrizp_QKTlv6Jl28"

	telegramToken := os.Getenv(TELEGRAM_BOT_TOKEN)
	if telegramToken == "" {
		log.Fatal("TelegramBotToken is not set")
	}

	//Initialize the bot

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("authorized on account %s", bot.Self.UserName)

	// Set update configuration
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Start pooling Telegram for updates
	updates := bot.GetUpdatesChan(u)

	// Handle messages

	for update := range updates {
		if update.Message != nil {
			// Handle "/weather <city> command"
			if strings.HasPrefix(update.Message.Text, "/weather") {
				// + extract the city name from the message
				parts := strings.SplitN(update.Message.Text, " ", 2)
				if len(parts) == 2 {
					city := parts[1]
					// Fetch weather data
					weather, err := getWeather(city)
					if err != nil {
						weather = fmt.Sprintf("Could not fetch weather data: %s", err.Error())
					}
					// Send weather data back to user
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, weather)
					bot.Send(msg)
				} else {
					// Send an error message if no city was provided
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide a city name like: /weather Praga")
					bot.Send(msg)
				}
			}
		}
	}
}
