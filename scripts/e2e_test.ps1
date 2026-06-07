# E2E Integration Test - Campus Collab
$base = "http://localhost:8080/api/v1"
$pass = 0
$fail = 0

function test($name, $result, $condition) {
    if ($condition) { 
        Write-Host "  [PASS] $name" -ForegroundColor Green
        $global:pass++
    } else {
        Write-Host "  [FAIL] $name -> $result" -ForegroundColor Red
        $global:fail++
    }
}

# 1. Register
Write-Host "`n=== Auth ===" -ForegroundColor Cyan
$body = @{student_id="e2e$(Get-Date -Format 'HHmmss')"; password="test123456"; nickname="E2E"; school_id=1} | ConvertTo-Json
try {
    $r = Invoke-RestMethod -Uri "$base/auth/register" -Method POST -Body $body -ContentType "application/json"
    $token = $r.data.token
    $userId = $r.data.user_id
    test "Register" $r ($r.data.token.Length -gt 0)
} catch { test "Register" $_.Exception.Message $false }

if ($token) {
    # 2. Login
    try {
        $r2 = Invoke-RestMethod -Uri "$base/auth/login" -Method POST -Body $body -ContentType "application/json"
        test "Login" $r2 ($r2.data.token.Length -gt 0)
    } catch { test "Login" $_.Exception.Message $false }

    $headers = @{Authorization="Bearer $token"}

    # 3. Get Me
    try {
        $r3 = Invoke-RestMethod -Uri "$base/auth/me" -Headers $headers
        test "GetCurrentUser" $r3 ($r3.data.id -gt 0)
    } catch { test "GetCurrentUser" $_.Exception.Message $false }

    # 4. Refresh Token
    try {
        $r4 = Invoke-RestMethod -Uri "$base/auth/refresh-token" -Method POST -Headers $headers
        $token = $r4.data.token
        $headers = @{Authorization="Bearer $token"}
        test "RefreshToken" $r4 ($r4.data.token.Length -gt 0)
    } catch { test "RefreshToken" $_.Exception.Message $false }

    # 5. Create Class
    Write-Host "`n=== Class ===" -ForegroundColor Cyan
    try {
        $cb = @{school_id=1; grade="2024"; name="E2E测试班级"} | ConvertTo-Json
        $r5 = Invoke-RestMethod -Uri "$base/classes" -Method POST -Body $cb -ContentType "application/json" -Headers $headers
        $classId = $r5.data.id
        test "CreateClass" $r5 ($classId -gt 0)
    } catch { test "CreateClass" $_.Exception.Message $false }

    # 6. List Classes
    try {
        $r6 = Invoke-RestMethod -Uri "$base/classes" -Headers $headers
        test "ListClasses" $r6 ($r6.data.Count -gt 0)
    } catch { test "ListClasses" $_.Exception.Message $false }

    # 7. Get Class Detail
    if ($classId) {
        try {
            $r7 = Invoke-RestMethod -Uri "$base/classes/$classId" -Headers $headers
            test "GetClassDetail" $r7 ($r7.data.id -eq $classId)
        } catch { test "GetClassDetail" $_.Exception.Message $false }

        # 8. Create Timetable
        Write-Host "`n=== Timetable ===" -ForegroundColor Cyan
        $tBody = @{
            entries=@(
                @{day_of_week=1; period_start=1; period_end=2; course_name="高等数学"; teacher="张老师"; room="A101"},
                @{day_of_week=3; period_start=3; period_end=4; course_name="大学英语"; teacher="李老师"; room="B202"}
            )
        } | ConvertTo-Json -Depth 4
        try {
            $r8 = Invoke-RestMethod -Uri "$base/timetables/class/$classId" -Method POST -Body $tBody -ContentType "application/json" -Headers $headers
            test "CreateTimetable" $r8 ($r8.data.created_count -gt 0)
        } catch { test "CreateTimetable" $_.Exception.Message $false }

        # 9. Get Personal Timetable
        try {
            $r9 = Invoke-RestMethod -Uri "$base/timetables/personal?class_id=$classId" -Headers $headers
            test "GetPersonalTimetable" $r9 ($r9.data.entries.Count -gt 0)
        } catch { test "GetPersonalTimetable" $_.Exception.Message $false }

        # 10. Create Poll (with auto_recommend)
        Write-Host "`n=== Poll ===" -ForegroundColor Cyan
        $pBody = @{
            title="E2E投票测试"
            scope_type="class"
            scope_id=$classId
            deadline=(Get-Date).AddDays(3).ToString("yyyy-MM-ddTHH:mm:ssZ")
            auto_recommend=$true
            time_preference=@{day_start_hour=8; day_end_hour=22; min_duration_min=60; max_recommendations=5}
        } | ConvertTo-Json -Depth 4
        try {
            $r10 = Invoke-RestMethod -Uri "$base/polls" -Method POST -Body $pBody -ContentType "application/json" -Headers $headers
            $pollId = $r10.data.poll_id
            test "CreatePoll+Recommend" $r10 ($pollId -gt 0 -and $r10.data.options_created -gt 0)
        } catch { test "CreatePoll+Recommend" $_.Exception.Message $false }

        if ($pollId) {
            # 11. Get Poll Detail
            try {
                $r11 = Invoke-RestMethod -Uri "$base/polls/$pollId" -Headers $headers
                test "GetPollDetail" $r11 ($r11.data.id -eq $pollId)
            } catch { test "GetPollDetail" $_.Exception.Message $false }

            # 12. Get Options
            try {
                $r12 = Invoke-RestMethod -Uri "$base/polls/$pollId/options" -Headers $headers
                $optId = $r12.data.options[0].id
                test "GetOptions" $r12 ($r12.data.options.Count -gt 0)
            } catch { test "GetOptions" $_.Exception.Message $false }

            # 13. Open Poll
            try {
                $r13 = Invoke-RestMethod -Uri "$base/polls/$pollId/open" -Method POST -Headers $headers
                test "OpenPoll" $r13 ($r13.data.message.Length -gt 0)
            } catch { test "OpenPoll" $_.Exception.Message $false }

            # 14. Vote
            if ($optId) {
                try {
                    $vBody = @{votes=@(@{option_id=$optId; choice="yes"})} | ConvertTo-Json -Depth 3
                    $r14 = Invoke-RestMethod -Uri "$base/polls/$pollId/vote" -Method POST -Body $vBody -ContentType "application/json" -Headers $headers
                    test "SubmitVote" $r14 ($r14.data.voted_count -gt 0)
                } catch { test "SubmitVote" $_.Exception.Message $false }
            }

            # 15. Get Results
            try {
                $r15 = Invoke-RestMethod -Uri "$base/polls/$pollId/results" -Headers $headers
                test "GetResults" $r15 ($r15.data.items.Count -gt 0)
            } catch { test "GetResults" $_.Exception.Message $false }

            # 16. Close Poll
            try {
                $r16 = Invoke-RestMethod -Uri "$base/polls/$pollId/close" -Method POST -Headers $headers
                test "ClosePoll" $r16 ($r16.data.message.Length -gt 0)
            } catch { test "ClosePoll" $_.Exception.Message $false }

            # 17. Finalize
            if ($optId) {
                try {
                    $fBody = @{final_option_id=$optId} | ConvertTo-Json
                    $r17 = Invoke-RestMethod -Uri "$base/polls/$pollId/finalize" -Method POST -Body $fBody -ContentType "application/json" -Headers $headers
                    test "FinalizePoll" $r17 ($r17.data.message.Length -gt 0)
                } catch { test "FinalizePoll" $_.Exception.Message $false }
            }
        }
    }
}

Write-Host "`n====================================" -ForegroundColor Cyan
Write-Host "  PASS: $pass | FAIL: $fail" -ForegroundColor $(if($fail -eq 0){'Green'}else{'Red'})
Write-Host "====================================`n"
