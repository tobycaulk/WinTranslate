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
	Text string
	LanguageCode string
}

type TranslationResponse struct {
	Text string
}

func getGoogleClient() translate.Client {
	ctx := context.Background()
	client, err := translate.NewClient(ctx, option.WithCredentialsFile("C:\\Users\\tobyc\\Documents\\Git\\Environment\\wintranslate-api-key.json"))
	if err != nil {
		fmt.Println("Error while retrieving Google Cloud client", err)
	}

	return *client
}

func getTranslatedText(text string, languageCode string) string {
	ctx := context.Background()
	var client = getGoogleClient()
	defer client.Close()

	lang, _ := language.Parse(languageCode)

	res, _ := client.Translate(ctx, []string{text}, lang, nil)
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

	translationResponse := TranslationResponse{Text: getTranslatedText(translationRequest.Text, translationRequest.LanguageCode)}
	json.NewEncoder(res).Encode(translationResponse)
}

func main() {
	http.HandleFunc("/translate", handleTranslate)
	http.ListenAndServe(":3000", nil)
}
