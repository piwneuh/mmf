package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Glicko struct {
	Rating    float64 `json:"rating"`
	Deviation float64 `json:"deviation"`
}

type Perf struct {
	Glicko   Glicko `json:"glicko"`
	Nb       int    `json:"nb"`
	Progress int    `json:"progress"`
}

type User struct {
	Name string `json:"name"`
}

type LichessResponse struct {
	User       User    `json:"user"`
	Perf       Perf    `json:"perf"`
	Rank       *int    `json:"rank"`
	Percentile float64 `json:"percentile"`
}

func getGlicko(username, perf string) (LichessResponse, error) {

	apiKey := os.Getenv("LICHESS_API_KEY")
	if apiKey == "" {
		return LichessResponse{}, fmt.Errorf("LICHESS_API_KEY not found in environment variables")
	}

	url := fmt.Sprintf("https://lichess.org/api/user/%s/perf/%s", username, perf)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LichessResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return LichessResponse{}, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LichessResponse{}, fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LichessResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	var lr LichessResponse
	err = json.Unmarshal(body, &lr)
	if err != nil {
		return LichessResponse{}, fmt.Errorf("error parsing JSON response: %v", err)
	}

	log.Println("Elo from lichess for", username, "is", lr.Perf.Glicko.Rating, "with deviation", lr.Perf.Glicko.Deviation)

	return lr, nil
}