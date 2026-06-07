package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	base := "http://localhost:8080/api/v1"
	client := &http.Client{}

	// Register
	regBody := `{"student_id":"integ01","password":"test123456","nickname":"integ","school_id":1}`
	resp, _ := client.Post(base+"/auth/register", "application/json", bytes.NewBufferString(regBody))
	buf, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println("REGISTER:", string(buf))

	var regResp struct {
		Code int `json:"code"`
		Data struct {
			Token  string `json:"token"`
			UserID uint   `json:"user_id"`
		} `json:"data"`
	}
	json.Unmarshal(buf, &regResp)
	token := regResp.Data.Token

	// Create Class
	classBody := `{"school_id":1,"grade":"2024","name":"integ_test_class"}`
	req, _ := http.NewRequest("POST", base+"/classes", bytes.NewBufferString(classBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println("CREATE CLASS:", string(buf))

	var classResp struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(buf, &classResp)
	classID := classResp.Data.ID
	fmt.Printf("classID=%d\n", classID)

	// Create Poll (no auto_recommend)
	pollBody := fmt.Sprintf(`{"title":"integ_test_poll","scope_type":"class","scope_id":%d,"auto_recommend":false}`, classID)
	req, _ = http.NewRequest("POST", base+"/polls", bytes.NewBufferString(pollBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("CREATE POLL (status=%d): %s\n", resp.StatusCode, string(buf))

	var pollResp struct {
		Code int `json:"code"`
		Data struct {
			PollID         uint   `json:"poll_id"`
			Status         string `json:"status"`
			OptionsCreated int    `json:"options_created"`
		} `json:"data"`
	}
	json.Unmarshal(buf, &pollResp)
	fmt.Printf("poll_id=%d status=%s options=%d\n", pollResp.Data.PollID, pollResp.Data.Status, pollResp.Data.OptionsCreated)
}
