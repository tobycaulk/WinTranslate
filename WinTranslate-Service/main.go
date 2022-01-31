package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/translate"
	"github.com/go-redis/redis/v8"
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

var ctx = context.Background()
var redisClient redis.Client
var googleClient translate.Client

func getHashedText(text string) string {
	return b64.StdEncoding.EncodeToString([]byte(text))
}

func storeTranslatedTextInRedis(originalText string, translatedText string) {
	hashedOriginalText := getHashedText(originalText)
	hashedTranslatedText := getHashedText(translatedText)

	err := redisClient.Set(ctx, hashedOriginalText, translatedText, 0).Err()
	if err != nil {
		fmt.Printf("Error while storing translation %s - %s \n", originalText, translatedText)
	} else {
		fmt.Printf("Stored translation %s - %s in Redis %s \n", originalText, translatedText, hashedOriginalText)
	}

	err = redisClient.Set(ctx, hashedTranslatedText, originalText, 0).Err()
	if err != nil {
		fmt.Printf("Error while storing translation %s - %s \n", translatedText, originalText)
	} else {
		fmt.Printf("Stored translation %s - %s in Redis %s \n", translatedText, originalText, hashedTranslatedText)
	}
}

func getTranslatedTextFromRedis(text string) string {
	hashedText := getHashedText(text)
	translatedText, err := redisClient.Get(ctx, hashedText).Result()
	if err != nil {
		fmt.Printf("Error while looking up %s in redis %e \n", text, err)
		return ""
	}

	return translatedText
}

func getTranslatedTextFromGoogle(text string, languageCode string) string {
	lang, err := language.Parse(languageCode)
	if err != nil {
		fmt.Printf("Error while parsing language code %s %e \n", languageCode, err)
		return ""
	}

	res, err := googleClient.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		fmt.Printf("Error while retrieving translation from google cloud %e \n", err)
		return ""
	}

	if len(res) <= 0 {
		return ""
	}

	return res[0].Text
}

func getParsedLanguageFromCode(languageCode string) (language.Tag, error) {
	lang, err := language.Parse(languageCode)
	if err != nil {
		return language.Tag{}, err
	}

	return lang, nil
}

func getTranslatedText(text string, languageCode string) string {
	redisTranslation := getTranslatedTextFromRedis(text)
	if len(redisTranslation) != 0 && redisTranslation != "" {
		fmt.Printf("Redis translation found - %s \n", text)
		return redisTranslation
	}

	googleCloudTranslation := getTranslatedTextFromGoogle(text, languageCode)
	if len(googleCloudTranslation) != 0 && googleCloudTranslation != "" {
		fmt.Printf("Google Cloud translation found - %s \n", googleCloudTranslation)
		storeTranslatedTextInRedis(text, googleCloudTranslation)
		return googleCloudTranslation
	}

	fmt.Printf("No translation found for %s to language %s \n", text, languageCode)

	return ""
}

func handleTranslate(res http.ResponseWriter, req *http.Request) {
	var translationRequest TranslationRequest
	err := json.NewDecoder(req.Body).Decode(&translationRequest)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(res, err.Error(), 400)
		return
	}

	fmt.Printf("Handling translation request for %s to language %s \n", translationRequest.Text, translationRequest.LanguageCode)

	translationResponse := TranslationResponse{Text: getTranslatedText(translationRequest.Text, translationRequest.LanguageCode)}
	json.NewEncoder(res).Encode(translationResponse)
}

func main() {
	setup()

	http.HandleFunc("/translate", handleTranslate)

	fmt.Println("Server running on port 3000")
	http.ListenAndServe("127.0.0.1:3000", nil)
}

func getGoogleClient() (translate.Client, error) {
	client, err := translate.NewClient(ctx, option.WithCredentialsFile("C:\\Users\\tobyc\\Documents\\Git\\Environment\\wintranslate-api-key.json"))
	return *client, err
}

func getRedisClient() redis.Client {
	return *redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})
}

func setup() {
	gc, err := getGoogleClient()
	if err != nil {
		fmt.Println("Error while retrieving Google Cloud client", err)
		return
	}
	googleClient = gc

	redisClient = getRedisClient()
	fmt.Println("Setup finished")
}
