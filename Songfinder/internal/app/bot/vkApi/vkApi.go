package vkApi

import (
	"Songfinder/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SearchMusic() (string, error) {
	var URL = fmt.Sprintf("https://api.vk.com/method/users.get?access_token=%s&user_id=%s&fields=status&v=5.131", config.AccessToken, config.VkID)

	resp, err := http.Get(URL)
	if err != nil {
		return "", fmt.Errorf("can't do request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read data: %w", err)
	}
	var data Result
	if err = json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("can't parse json:%w", err)
	}
	if data.Resp[0].Audio.Artist != "" {
		str := fmt.Sprintf("Сейчас играет: %s - %s", data.Resp[0].Audio.Artist, data.Resp[0].Audio.Title)
		strings.TrimSpace(str)
		return str, nil
	}
	return "", nil
}
