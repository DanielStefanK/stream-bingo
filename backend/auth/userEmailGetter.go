package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"

	"github.com/DanielStefanK/stream-bingo/config"
)

type UserInfo struct {
	ID        string
	Name      string
	Email     string
	AvatarURL string
}

func GetUserInfoFromProvider(provider string, token string) (UserInfo, error) {
	if provider == "github" {
		return getUserInfoFromGithub(token)
	}

	defaultData, err := fetchDefault(provider, config.GetConfig().OAuth.Providers[provider].UserURL, token)

	if err != nil {
		log.Printf("Error fetching user info from %s: %v", provider, err)
		return UserInfo{}, err
	}

	return UserInfo{
		ID:        defaultData["id"],
		Name:      defaultData["name"],
		Email:     defaultData["email"],
		AvatarURL: defaultData["avatar_url"],
	}, nil
}

func getUserInfoFromGithub(token string) (UserInfo, error) {

	userInfo := UserInfo{}

	client := &http.Client{}

	infoReq, _ := http.NewRequest("GET", config.GetConfig().OAuth.Providers["github"].UserURL, nil)
	infoReq.Header.Set("Authorization", "Bearer "+token)
	rspInfo, err := client.Do(infoReq)
	if err != nil {
		log.Printf("Error making request for github: %v", err)
		return userInfo, err
	}
	defer rspInfo.Body.Close()

	// get only id from response (the response container strings and numbers)
	var idRsp map[string]interface{}
	err = json.NewDecoder(rspInfo.Body).Decode(&idRsp)
	if err != nil {
		log.Printf("Error decoding user info JSON from github: %v", err)
		return userInfo, err
	}

	userInfo.ID = fmt.Sprint(uint(idRsp["id"].(float64)))
	userInfo.Name = idRsp["login"].(string)
	userInfo.AvatarURL = idRsp["avatar_url"].(string)

	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request for github: %v", err)
		return userInfo, err
	}
	defer resp.Body.Close()

	var emailRsp []map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&emailRsp)
	if err != nil {
		log.Printf("Error decoding user info JSON from github: %v", err)

		return userInfo, err
	}

	userInfo.Email = emailRsp[0]["email"].(string)

	log.Printf("Successfully fetched user info from github")
	return userInfo, nil
}

func fetchDefault(providerName string, userURL, accessToken string) (map[string]string, error) {
	log.Printf("Fetching user info from URL: %s", userURL)

	req, _ := http.NewRequest("GET", userURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to %s: %v", userURL, err)
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]string

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Printf("Error decoding user info JSON from %s: %v", userURL, err)

		return nil, err
	}

	log.Printf("Successfully fetched user info from %s", userURL)
	return userInfo, nil
}
