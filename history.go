package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type ChatMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChannelHistory struct {
	Messages []ChatMessage `json:"messages"`
}

type HistoryManager struct {
	Histories  map[string]*ChannelHistory `json:"histories"`
	Filename   string
	Mutex      sync.Mutex
	RotateDays int // 保存期間
	ticker     *time.Ticker
	quit       chan struct{}
}

func NewHistoryManager(filename string, days int) *HistoryManager {
	hm := &HistoryManager{
		Histories:  make(map[string]*ChannelHistory),
		Filename:   filename,
		RotateDays: days,
		// 1分間隔で保存
		ticker: time.NewTicker(1 * time.Minute),
		quit:   make(chan struct{}),
	}
	_ = hm.Load() // ファイルがなければ無視
	go hm.periodicSave()
	return hm
}

func (hm *HistoryManager) AddMessage(channelID, role, content string) {
	hm.Mutex.Lock()
	defer hm.Mutex.Unlock()
	msg := ChatMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
	hist, ok := hm.Histories[channelID]
	if !ok {
		hist = &ChannelHistory{Messages: []ChatMessage{}}
		hm.Histories[channelID] = hist
	}
	hist.Messages = append(hist.Messages, msg)
	hm.cleanup(channelID)
}

func (hm *HistoryManager) GetMessages(channelID string) []ChatMessage {
	hm.Mutex.Lock()
	defer hm.Mutex.Unlock()
	if hist, ok := hm.Histories[channelID]; ok {
		return hist.Messages
	}
	return nil
}

func (hm *HistoryManager) cleanup(channelID string) {
	if hm.RotateDays == 0 {
		return
	}
	threshold := time.Now().AddDate(0, 0, -hm.RotateDays)
	hist := hm.Histories[channelID]
	idx := 0
	for i, m := range hist.Messages {
		if m.Timestamp.After(threshold) {
			idx = i
			break
		}
	}
	hist.Messages = hist.Messages[idx:]
}

func (hm *HistoryManager) Save() error {
	hm.Mutex.Lock()
	defer hm.Mutex.Unlock()
	f, err := os.Create(hm.Filename)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	return e.Encode(hm)
}

func (hm *HistoryManager) Load() error {
	f, err := os.Open(hm.Filename)
	if err != nil {
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	return d.Decode(hm)
}

func (hm *HistoryManager) periodicSave() {
	for {
		select {
		case <-hm.ticker.C:
			_ = hm.Save()
		case <-hm.quit:
			return
		}
	}
}

func (hm *HistoryManager) Stop() {
	close(hm.quit)
	hm.Save()
}

func (hm *HistoryManager) Forget(channelID string) error {
	hm.Mutex.Lock()
	// 現在の history.json をバックアップする
	t := time.Now().Format("20060102-150405")
	_ = os.Rename(hm.Filename, "history-"+t+".json")
	// リセット
	hm.Histories[channelID] = &ChannelHistory{Messages: []ChatMessage{}}
	hm.Mutex.Unlock()
	err := hm.Save()
	if err != nil {
		log.Printf("Error: An error happend when removing history: %s", err)
		return err
	} else {
		return nil
	}
}
