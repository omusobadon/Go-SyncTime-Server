package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	mutex       sync.Mutex
	CurrentTime time.Time
)

// TimeResponse はWorld Time APIからのレスポンスを表す構造体です。
type TimeResponse struct {
	Dateline string `json:"dateline"`
}

// updateTime はWorld Time APIから時刻を取得し、CurrentTimeを更新
func updateTime() {
	resp, err := http.Get("http://worldtimeapi.org/api/ip")
	if err != nil {
		fmt.Println("時刻取得エラー:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		fmt.Println("レスポンス読み込みエラー:", err)
		return
	}

	var timeResp TimeResponse
	if err := json.Unmarshal(body, &timeResp); err != nil {
		fmt.Println("レスポンスのアンマーシャリングエラー:", err)
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, timeResp.Dateline)
	if err != nil {
		fmt.Println("時刻の解析エラー:", err)
		return
	}

	mutex.Lock()
	CurrentTime = parsedTime
	mutex.Unlock()

	fmt.Println("Time updated:", CurrentTime)
}

func main() {
	// World Time APIから時刻を取得し、CurrentTimeを更新
	updateTime()

	// 3時間ごとに時刻を更新
	go func() {
		for range time.Tick(10 * time.Second) {
			updateTime()
		}
	}()

	// 1秒ごとにCurrentTimeを更新する
	go func() {
		for range time.Tick(1 * time.Second) {
			mutex.Lock()
			CurrentTime = CurrentTime.Add(1 * time.Second)
			mutex.Unlock()
		}
	}()

	// 1秒待つ
	time.Sleep(time.Duration(1 * time.Second))

	// 1秒ごとに時刻を表示する
	for range time.Tick(1 * time.Second) {
		mutex.Lock()
		fmt.Println("Current Time:", CurrentTime)
		mutex.Unlock()
	}
}
