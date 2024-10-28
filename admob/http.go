package admob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GCP_API_TOKEN")))

	client := &http.Client{}
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
