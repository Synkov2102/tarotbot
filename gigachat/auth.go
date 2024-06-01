package gigachat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"taro-bot/certs"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Value   string
	Created time.Time
}

func GetToken() (Token, error) {

	url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	data := make(map[string]string)
	data["scope"] = "GIGACHAT_API_PERS"

	// Преобразование данных в формат x-www-form-urlencoded
	formData := formatData(data)

	method := "POST"

	client := certs.GetClientWithCerts()

	req, err := http.NewRequest(method, url, strings.NewReader(formData))
	id := uuid.New()

	if err != nil {
		fmt.Println(err)
		return Token{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("RqUID", id.String())
	req.Header.Add("Authorization", "Basic "+os.Getenv("GIGACHAT_TOKEN"))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Token{}, err
	}
	defer res.Body.Close()

	// Предполагаем, что res - это объект *http.Response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return Token{}, err
	}

	// Преобразуем байты в строку
	bodyString := string(bodyBytes)

	// Декодируем строку в структуру или карту, содержащую access_token
	var responseStruct struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	}

	err = json.Unmarshal([]byte(bodyString), &responseStruct)

	if err != nil {
		fmt.Println(err)
		return Token{}, err
	}
	// Получаем access_token
	return Token{Value: responseStruct.AccessToken, Created: time.Unix(responseStruct.ExpiresAt, 0)}, nil
}

func formatData(data map[string]string) string {
	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, value)
	}
	return formData.Encode()
}
