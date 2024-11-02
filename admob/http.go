package admob

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type DateRange struct {
	StartDate Date `json:"startDate"`
	EndDate   Date `json:"endDate"`
}

type LocalizationSettings struct {
	CurrencyCode string `json:"currencyCode"`
	LanguageCode string `json:"languageCode"`
}

type Header struct {
	DateRange            DateRange            `json:"dateRange"`
	LocalizationSettings LocalizationSettings `json:"localizationSettings"`
}

type DimensionValue struct {
	Value        string `json:"value"`
	DisplayLabel string `json:"displayLabel"`
}

type MetricValue struct {
	IntegerValue string `json:"integerValue,omitempty"`
	MicrosValue  string `json:"microsValue,omitempty"`
}

type DimensionValues struct {
	App DimensionValue `json:"APP"`
}

type MetricValues struct {
	Clicks            MetricValue `json:"CLICKS"`
	AdRequests        MetricValue `json:"AD_REQUESTS"`
	Impressions       MetricValue `json:"IMPRESSIONS"`
	EstimatedEarnings MetricValue `json:"ESTIMATED_EARNINGS"`
}

type Row struct {
	DimensionValues DimensionValues `json:"dimensionValues"`
	MetricValues    MetricValues    `json:"metricValues"`
}

type Footer struct {
	MatchingRowCount string `json:"matchingRowCount"`
}

type ReportItem struct {
	Header *Header `json:"header,omitempty"`
	Row    *Row    `json:"row,omitempty"`
	Footer *Footer `json:"footer,omitempty"`
}

func RequestAdmobApi(apiUrl string, requestBody map[string]interface{}) ([]ReportItem, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("JSONのエンコードに失敗しちゃった！: %v", err)
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("HTTPリクエストの作成に失敗しちゃった！: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client, err := oauth2Authentication()
	if err != nil {
		log.Fatalf("OAuth2認証に失敗しちゃった！: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTPリクエストに失敗しちゃった！: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("HTTPレスポンスの読み取りに失敗しちゃった！: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		var result []ReportItem
		err := json.Unmarshal([]byte(body), &result)
		if err != nil {
			log.Fatalf("JSONのパースに失敗しちゃった！: %v", err)
		}
		return result, nil
	} else {
		fmt.Printf("レスポンスが異常だったみたい！: %d\n", resp.StatusCode)
		fmt.Println("れすぽんす:", string(body))
		return nil, err
	}
}

func MakeRequestBody(cmd *NoticeSummaryCmd) (map[string]interface{}, time.Time, time.Time) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)
	switch {
	case cmd.w:
		startDate = startDate.AddDate(0, 0, -int(startDate.Weekday()))
	case cmd.m:
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.Local)
		endDate = startDate.AddDate(0, 1, 0)
	case cmd.y:
		startDate = time.Date(startDate.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		endDate = startDate.AddDate(1, 0, 0)
	}

	// TODO: パラメータを外部から設定できるようにする
	return map[string]interface{}{
		"report_spec": map[string]interface{}{
			"date_range": map[string]interface{}{
				"start_date": map[string]int{"year": startDate.Year(), "month": int(startDate.Month()), "day": startDate.Day()},
				"end_date":   map[string]int{"year": endDate.Year(), "month": int(endDate.Month()), "day": endDate.Day()},
			},
			"dimensions": []string{"APP"},
			"metrics":    []string{"CLICKS", "AD_REQUESTS", "IMPRESSIONS", "ESTIMATED_EARNINGS"},
			"sort_conditions": []map[string]string{
				{"metric": "CLICKS", "order": "DESCENDING"},
			},
			"localization_settings": map[string]string{
				"currency_code": "JPY",
				"language_code": "ja-JP",
			},
		},
	}, startDate, endDate
}

func oauth2Authentication() (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_OAUTH2_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/admob.readonly"},
		Endpoint:     google.Endpoint,
	}

	ctx := context.Background()
	token, err := getTokenFromFile("config/token.json")

	if err != nil {
		u := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
		decoded, err := url.QueryUnescape(u)
		if err != nil {
			log.Fatal(err)
		}
		exec.Command("open", decoded).Run()

		fmt.Printf("ログイン後のリダイレクト先にあるcodeを入力してEnterを押下してね！")

		// ブラウザでの認証後、リダイレクトURLに返されたコードを使ってトークンを取得
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}

		token, err = config.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		saveToken("config/token.json", token)
	}

	client := config.Client(ctx, token)
	resp, err := client.Get("https://accounts.google.com/o/oauth2/auth")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("OAuth2認証が完了したよ！")

	return client, nil
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("トークンの保存に失敗しちゃった！: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	return &token, err
}
