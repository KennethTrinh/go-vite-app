package initializers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/KennethTrinh/go-vite-app/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const telegramMaxLength = 4096

type TelegramHook struct{}

func (t *TelegramHook) Run(
	e *zerolog.Event,
	level zerolog.Level,
	message string,
) {
	if level > zerolog.WarnLevel {
		go func() {
			timestamp := time.Now().In(time.FixedZone("EST", -5*3600)).Format("2006-01-02 15:04:05")
			fullMessage := fmt.Sprintf("[%s] %s", timestamp, message)
			SendToTelegramBot(fullMessage)
		}()
	}
}

func splitMessage(message string, maxLength int) []string {
	var messages []string
	runes := []rune(message)
	for len(runes) > maxLength {
		splitIndex := maxLength
		for splitIndex > 0 && runes[splitIndex] != ' ' {
			splitIndex--
		}
		if splitIndex == 0 {
			splitIndex = maxLength
		}
		messages = append(messages, string(runes[:splitIndex]))
		runes = runes[splitIndex:]
	}
	messages = append(messages, string(runes))
	return messages
}

// https://github.com/yumusb/email_router/blob/14ef6c289972ec7c00b6bc5bc2f5e267bd03511b/func.go#L31
func SendToTelegramBot(message string, optBotToken ...string) {
	botToken := config.Env.TelegramBotToken
	if len(optBotToken) > 0 {
		botToken = optBotToken[0]
	}

	chatID := config.Env.TelegramChatID
	if len(optBotToken) > 1 {
		chatID = optBotToken[1]
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	messages := splitMessage(message, telegramMaxLength)

	for _, msgPart := range messages {
		payload := map[string]interface{}{
			"chat_id": chatID,
			"text":    msgPart,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// if resp.StatusCode != 200 {
		// 	bodyBytes, _ := io.ReadAll(resp.Body)
		// 	bodyString := string(bodyBytes)
		// 	log.Error().Msgf("Failed to send message to Telegram bot - Status: %s - Body: %s", resp.Status, bodyString)
		// }
	}
}

func InitLogger() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	logger := zerolog.New(os.Stderr).With().Caller().Logger() // timestamp added in fluent-bit
	logger = logger.Hook(&TelegramHook{})
	log.Logger = logger
}
