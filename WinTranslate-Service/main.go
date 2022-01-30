package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type TranslationRequest struct {
	Text         string
	LanguageCode string
}

type TranslationResponse struct {
	Text string
}

func getGoogleClient() translate.Client {
	ctx := context.Background()
	client, err := translate.NewClient(ctx, option.WithCredentialsFile("/Users/tobycaulk/src/env/wintranslate-api-key.json"))
	if err != nil {
		fmt.Println("Error while retrieving Google Cloud client", err)
	}

	return *client
}

func getParsedLanguageFromCode(languageCode string) (language.Tag, error) {
	lang, err := language.Parse(languageCode)
	if err != nil {
		return language.Tag{}, err
	}

	return lang, nil
}

func getTranslatedText(text string, languageCode string) string {
	ctx := context.Background()
	client := getGoogleClient()
	defer client.Close()

	language, err := getParsedLanguageFromCode(languageCode)
	if err != nil {
		fmt.Println("Error while parsing language code", err)
		return ""
	}

	res, err := client.Translate(ctx, []string{text}, language, nil)
	if err != nil {
		fmt.Println("Error while retrieving translation", err)
	}

	return res[0].Text
}

func handleTranslate(res http.ResponseWriter, req *http.Request) {
	var translationRequest TranslationRequest
	err := json.NewDecoder(req.Body).Decode(&translationRequest)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(res, err.Error(), 400)
		return
	}

	translatedText := getTranslatedText(translationRequest.Text, translationRequest.LanguageCode)
	translationResponse := TranslationResponse{Text: translatedText}
	json.NewEncoder(res).Encode(translationResponse)
}

func main() {
	http.HandleFunc("/translate", handleTranslate)
	http.ListenAndServe(":3000", nil)
}
