package gigachat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"taro-bot/certs"
)

func Generate(cards string, question string, token string) (string, error) {

	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	method := "POST"

	task := fmt.Sprintf(`{
		"model": "GigaChat",
		"messages": [
			  {
				  "role": "system",
				  "content": "Ты — опытный таролог. Ты делаешь расклад таро, объясни выпавшие карты. Описание необходимо делать в контексте вопроса."
			  },
			  {
				  "role": "user",
				  "content": "Опиши выпавшие карты %s. Был задан вопрос %s  "
			  }
		  ]
	  }`, cards, question)
	payload := strings.NewReader(task)

	client := certs.GetClientWithCerts()
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	// Предполагаем, что res - это объект *http.Response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Преобразуем байты в строку
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	// Декодируем строку в структуру или карту, содержащую access_token
	type Choice struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	}

	var ChatCompletion struct {
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Object  string `json:"object"`
		Usage   struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Choices []Choice `json:"choices"`
	}

	err = json.Unmarshal([]byte(bodyString), &ChatCompletion)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Получаем Сообщение
	if len(ChatCompletion.Choices) > 0 {
		message := ChatCompletion.Choices[0].Message.Content
		fmt.Println(message)
		return message, nil
	} else {
		return "err", err
	}

}
