package sniper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/fidesy/me-sniper/internal/models"
)

const storageTime = time.Duration(time.Second * 30)

var (
	cache = make(map[string]*models.Floor)
	mutex sync.Mutex
	cli   = &http.Client{}
)

func GetFloor(symbol string) float64 {
	if _, ok := cache[symbol]; ok && time.Since(cache[symbol].Time) < storageTime {

	} else {
		request, _ := http.NewRequest("GET", fmt.Sprintf("https://api-mainnet.magiceden.dev/v2/collections/%s/stats", symbol), nil)
		resp, err := cli.Do(request)
		if err != nil {
			return cache[symbol].Value

		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return cache[symbol].Value
		}

		body, _ := io.ReadAll(resp.Body)
		var floorResp models.FloorResponse
		json.Unmarshal(body, &floorResp)

		mutex.Lock()
		defer mutex.Unlock()
		cache[symbol] = &models.Floor{Value: floorResp.FloorPrice / 1e9, Time: time.Now()}
	}

	return cache[symbol].Value
}
