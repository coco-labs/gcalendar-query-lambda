package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Query string `json:"query"`
}

type Response struct {
	Count int `json:"count"`
}

type CalendarEvent struct {
	Id      string `json:"id"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
}

type CalendarResponse struct {
	Items []CalendarEvent `json:"items"`
}

// Filter the events because cancelled events are not valid
func GetValidEvents(calendarResponse CalendarResponse, query string) []CalendarEvent {
	validEvents := make([]CalendarEvent, 0)
	for _, event := range calendarResponse.Items {
		if event.Status != "cancelled" && strings.Contains(event.Summary, query) {
			validEvents = append(validEvents, event)
		}
	}
	return validEvents
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarID := url.QueryEscape(os.Getenv("CALENDAR_ID"))
	apiKey := os.Getenv("GCALENDAR_API_KEY")
	query := request.QueryStringParameters["query"]
	safeQuery := url.QueryEscape(query)
	timeMax := time.Now().Format(time.RFC3339)

	endpoint := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events?q=%s&timeMax=%s", calendarID, safeQuery, timeMax)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println("Request error: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	client := &http.Client{}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request error: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	if resp.StatusCode >= 300 {
		log.Println("Request error: ", resp.Status, resp.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: resp.StatusCode,
			Body:       "Request was not successful, check the logs",
		}, nil
	}
	defer resp.Body.Close()

	var calendarResponse CalendarResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&calendarResponse); err != nil {
		log.Println("Request error: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	// Filter the events to avoid cancelled events
	validEvents := GetValidEvents(calendarResponse, query)

	response := Response{Count: len(validEvents)}
	json, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(json),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
