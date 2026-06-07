package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	base := "http://localhost:8080/api/v1"
	fmt.Println("=== 并发负载测试 ===")
	fmt.Println()

	// Phase 1: 准备测试数据 — 注册30个用户 + 创建1个班级 + 课表 + 投票
	fmt.Println("--- Phase 1: 准备测试数据 ---")

	// 注册管理员用户 + 创建班级
	adminID := fmt.Sprintf("perf%05d", time.Now().UnixNano()%100000)
	regBody := fmt.Sprintf(`{"student_id":"%s","password":"test123456","nickname":"Admin","school_id":1}`, adminID)
	resp, _ := http.Post(base+"/auth/register", "application/json", bytes.NewBufferString(regBody))
	buf, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var regResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(buf, &regResp)
	adminToken := regResp.Data.Token
	fmt.Printf("Admin registered: %s\n", adminID)

	// 创建班级
	classBody := `{"school_id":1,"grade":"2024","name":"perf_test_class"}`
	req, _ := http.NewRequest("POST", base+"/classes", bytes.NewBufferString(classBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	var classResp struct {
		Data struct {
			ID         uint   `json:"id"`
			InviteCode string `json:"invite_code"`
		} `json:"data"`
	}
	json.Unmarshal(buf, &classResp)
	classID := classResp.Data.ID
	inviteCode := classResp.Data.InviteCode
	fmt.Printf("Class created: ID=%d, invite=%s\n", classID, inviteCode)

	// 创建课表 (触发后台继承)
	ttBody := `{"entries":[
		{"day_of_week":1,"period_start":1,"period_end":2,"course_name":"高等数学","teacher":"张老师","room":"A101"},
		{"day_of_week":3,"period_start":3,"period_end":4,"course_name":"大学英语","teacher":"李老师","room":"B202"}
	]}`
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/timetables/class/%d", base, classID), bytes.NewBufferString(ttBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("Timetable created for class %d\n", classID)

	// 创建投票 (auto_recommend)
	pollBody := fmt.Sprintf(`{
		"title":"perf_test_poll",
		"scope_type":"class",
		"scope_id":%d,
		"deadline":"2026-06-20T00:00:00Z",
		"auto_recommend":true,
		"time_preference":{"day_start_hour":8,"day_end_hour":22,"min_duration_min":60,"max_recommendations":5}
	}`, classID)
	req, _ = http.NewRequest("POST", base+"/polls", bytes.NewBufferString(pollBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	var pollResp struct {
		Data struct {
			PollID         uint `json:"poll_id"`
			OptionsCreated int  `json:"options_created"`
		} `json:"data"`
	}
	json.Unmarshal(buf, &pollResp)
	pollID := pollResp.Data.PollID
	fmt.Printf("Poll created: ID=%d, options=%d\n", pollID, pollResp.Data.OptionsCreated)

	// 开启投票
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/polls/%d/open", base, pollID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, _ = http.DefaultClient.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println("Poll opened")

	// 获取第一个 option ID
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/polls/%d/options", base, pollID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, _ = http.DefaultClient.Do(req)
	buf, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	var optResp struct {
		Data struct {
			Options []struct{ ID uint `json:"id"` } `json:"options"`
		} `json:"data"`
	}
	json.Unmarshal(buf, &optResp)
	if len(optResp.Data.Options) == 0 {
		fmt.Println("No options, aborting")
		return
	}
	optID := optResp.Data.Options[0].ID

	// Phase 2: 并发投票压测
	fmt.Println()
	fmt.Println("--- Phase 2: 并发投票压测 (20 goroutines, 50 votes each) ---")

	concurrency := 20
	votesPerGoroutine := 50
	var wg sync.WaitGroup
	startTime := time.Now()
	successCount := int64(0)
	failCount := int64(0)
	var mu sync.Mutex

	// 预先注册并获取 tokens
	type userToken struct {
		studentID string
		token     string
	}
	users := make([]userToken, concurrency)
	for i := 0; i < concurrency; i++ {
		sid := fmt.Sprintf("perfv%05d%d", time.Now().UnixNano()%100000, i)
		ub := fmt.Sprintf(`{"student_id":"%s","password":"test123456","nickname":"Voter%d","school_id":1}`, sid, i)
		r, _ := http.Post(base+"/auth/register", "application/json", bytes.NewBufferString(ub))
		b, _ := io.ReadAll(r.Body)
		r.Body.Close()
		var ur struct {
			Data struct{ Token string `json:"token"` } `json:"data"`
		}
		json.Unmarshal(b, &ur)

		// 加入班级
		joinBody := fmt.Sprintf(`{"invite_code":"%s"}`, inviteCode)
		joinReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/classes/%d/join", base, classID), bytes.NewBufferString(joinBody))
		joinReq.Header.Set("Authorization", "Bearer "+ur.Data.Token)
		joinReq.Header.Set("Content-Type", "application/json")
		jr, _ := http.DefaultClient.Do(joinReq)
		io.ReadAll(jr.Body)
		jr.Body.Close()

		users[i] = userToken{studentID: sid, token: ur.Data.Token}
	}
	fmt.Printf("Registered %d voters and joined class\n", concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			client := &http.Client{Timeout: 10 * time.Second}
			token := users[idx].token
			voteBody := fmt.Sprintf(`{"votes":[{"option_id":%d,"choice":"yes"}]}`, optID)

			for j := 0; j < votesPerGoroutine; j++ {
				req, _ := http.NewRequest("POST", fmt.Sprintf("%s/polls/%d/vote", base, pollID),
					bytes.NewBufferString(voteBody))
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					mu.Lock()
					failCount++
					mu.Unlock()
					continue
				}
				io.ReadAll(resp.Body)
				resp.Body.Close()
				if resp.StatusCode == 200 || resp.StatusCode == 201 {
					mu.Lock()
					successCount++
					mu.Unlock()
				} else {
					mu.Lock()
					failCount++
					mu.Unlock()
				}
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	totalRequests := int64(concurrency * votesPerGoroutine)
	fmt.Printf("Total requests: %d, Success: %d, Fail: %d\n", totalRequests, successCount, failCount)
	fmt.Printf("Duration: %v, Throughput: %.0f req/s\n", elapsed, float64(totalRequests)/elapsed.Seconds())

	// Phase 3: GetResults 缓存命中测试
	fmt.Println()
	fmt.Println("--- Phase 3: GetResults 缓存命中测试 ---")
	startTime = time.Now()
	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/polls/%d/results", base, pollID), nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		resp, _ := http.DefaultClient.Do(req)
		io.ReadAll(resp.Body)
		resp.Body.Close()
	}
	elapsed = time.Since(startTime)
	fmt.Printf("100x GetResults: %v (avg: %v)\n", elapsed, elapsed/100)
}
