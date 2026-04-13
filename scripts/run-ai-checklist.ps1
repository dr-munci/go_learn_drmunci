param(
    [string]$BaseUrl = "http://localhost:8090"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$results = New-Object System.Collections.Generic.List[object]

function Add-Result {
    param(
        [string]$Id,
        [string]$Name,
        [string]$Status,
        [string]$Detail
    )
    $results.Add([pscustomobject]@{
            Test   = $Id
            Name   = $Name
            Status = $Status
            Detail = $Detail
        })
}

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Url,
        [hashtable]$Headers = @{},
        $Body = $null
    )
    $maxAttempts = 3
    $attempt = 0
    while ($attempt -lt $maxAttempts) {
        $attempt++
        try {
            if ($null -ne $Body) {
                $json = $Body | ConvertTo-Json -Depth 10 -Compress
                $resp = Invoke-WebRequest -Method $Method -Uri $Url -Headers $Headers -ContentType "application/json" -Body $json -UseBasicParsing
            }
            else {
                $resp = Invoke-WebRequest -Method $Method -Uri $Url -Headers $Headers -UseBasicParsing
            }
            $parsed = $null
            try { $parsed = $resp.Content | ConvertFrom-Json } catch { $parsed = $resp.Content }
            Start-Sleep -Milliseconds 220
            return @{
                Status = [int]$resp.StatusCode
                Body   = $parsed
                Raw    = $resp.Content
            }
        }
        catch {
            $status = 0
            $content = ""
            if ($_.Exception.Response) {
                $status = [int]$_.Exception.Response.StatusCode
                $content = (New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())).ReadToEnd()
            }

            # Rate limit'e takılınca kısa bekleyip tekrar dene.
            if ($status -eq 429 -and $attempt -lt $maxAttempts) {
                Start-Sleep -Milliseconds 600
                continue
            }

            $parsed = $null
            try { $parsed = $content | ConvertFrom-Json } catch { $parsed = $content }
            Start-Sleep -Milliseconds 220
            return @{
                Status = $status
                Body   = $parsed
                Raw    = $content
            }
        }
    }

    # Güvenlik fallback'i
    return @{
        Status = 0
        Body   = $null
        Raw    = ""
    }
}

function Get-NestedValue {
    param(
        $Object,
        [string[]]$Path
    )

    $current = $Object
    foreach ($segment in $Path) {
        if ($null -eq $current) { return $null }
        $props = $current.PSObject.Properties.Name
        if (-not ($props -contains $segment)) { return $null }
        $current = $current.$segment
    }
    return $current
}

$suffix = [Guid]::NewGuid().ToString("N").Substring(0, 8)
$teacherEmail = "teacher_$suffix@test.com"
$studentEmail = "student_$suffix@test.com"

# T1
$t1 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/auth/register" -Body @{
    name     = "Teacher $suffix"
    email    = $teacherEmail
    password = "123456"
    role     = "teacher"
}
if ($t1.Status -eq 201 -and $t1.Body.token) { Add-Result "T1" "POST /api/auth/register (teacher)" "PASS" "201 + token" } else { Add-Result "T1" "POST /api/auth/register (teacher)" "FAIL" "Status=$($t1.Status)" }

# T2
$t2 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/auth/register" -Body @{
    name     = "Student $suffix"
    email    = $studentEmail
    password = "123456"
    role     = "student"
}
if ($t2.Status -eq 201 -and $t2.Body.token) { Add-Result "T2" "POST /api/auth/register (student)" "PASS" "201 + token" } else { Add-Result "T2" "POST /api/auth/register (student)" "FAIL" "Status=$($t2.Status)" }

# T3
$t3 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/auth/login" -Body @{ email = $teacherEmail; password = "123456" }
$teacherToken = Get-NestedValue -Object $t3.Body -Path @("token")
if ($t3.Status -eq 200 -and $teacherToken) { Add-Result "T3" "POST /api/auth/login" "PASS" "200 + token" } else { Add-Result "T3" "POST /api/auth/login" "FAIL" "Status=$($t3.Status)" }

$studentLogin = Invoke-Api -Method "POST" -Url "$BaseUrl/api/auth/login" -Body @{ email = $studentEmail; password = "123456" }
$studentToken = Get-NestedValue -Object $studentLogin.Body -Path @("token")

# T4
$t4 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/auth/login" -Body @{ email = $teacherEmail; password = "yanlis" }
if ($t4.Status -eq 401) { Add-Result "T4" "POST /api/auth/login (wrong password)" "PASS" "401" } else { Add-Result "T4" "POST /api/auth/login (wrong password)" "FAIL" "Status=$($t4.Status)" }

$teacherHeaders = @{ Authorization = "Bearer $teacherToken" }
$studentHeaders = @{ Authorization = "Bearer $studentToken" }

# T5
$t5 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/courses" -Headers $teacherHeaders -Body @{
    title       = "Go Course $suffix"
    description = "desc"
    category    = "go"
}
$courseId = Get-NestedValue -Object $t5.Body -Path @("data", "ID")
if ($t5.Status -eq 201 -and $courseId) { Add-Result "T5" "POST /api/courses (teacher)" "PASS" "201" } else { Add-Result "T5" "POST /api/courses (teacher)" "FAIL" "Status=$($t5.Status)" }

# T6
$t6 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/courses" -Headers $studentHeaders -Body @{ title = "Nope" }
if ($t6.Status -eq 403) { Add-Result "T6" "POST /api/courses (student)" "PASS" "403" } else { Add-Result "T6" "POST /api/courses (student)" "FAIL" "Status=$($t6.Status)" }

# T7
$t7 = Invoke-Api -Method "GET" -Url "$BaseUrl/api/courses?page=1&limit=5" -Headers $teacherHeaders
$keys7 = @("data", "page", "limit", "total")
$ok7 = $t7.Status -eq 200
foreach ($k in $keys7) { $ok7 = $ok7 -and ($t7.Body.PSObject.Properties.Name -contains $k) }
if ($ok7) { Add-Result "T7" "GET /api/courses pagination" "PASS" "200 + fields" } else { Add-Result "T7" "GET /api/courses pagination" "FAIL" "Status=$($t7.Status)" }

# T8
$t8 = Invoke-Api -Method "GET" -Url "$BaseUrl/api/courses/$courseId" -Headers $teacherHeaders
if ($t8.Status -eq 200) { Add-Result "T8" "GET /api/courses/:id" "PASS" "200" } else { Add-Result "T8" "GET /api/courses/:id" "FAIL" "Status=$($t8.Status)" }

# T9
$t9 = Invoke-Api -Method "PUT" -Url "$BaseUrl/api/courses/$courseId" -Headers $teacherHeaders -Body @{
    title       = "Go Course Updated $suffix"
    description = "updated"
    category    = "go"
}
if ($t9.Status -eq 200) { Add-Result "T9" "PUT /api/courses/:id (owner)" "PASS" "200" } else { Add-Result "T9" "PUT /api/courses/:id (owner)" "FAIL" "Status=$($t9.Status)" }

# T10 (create another course then delete)
$tmp = Invoke-Api -Method "POST" -Url "$BaseUrl/api/courses" -Headers $teacherHeaders -Body @{
    title       = "Delete Course $suffix"
    description = "temp"
    category    = "go"
}
$tmpId = Get-NestedValue -Object $tmp.Body -Path @("data", "ID")
$t10 = Invoke-Api -Method "DELETE" -Url "$BaseUrl/api/courses/$tmpId" -Headers $teacherHeaders
if ($t10.Status -eq 200) { Add-Result "T10" "DELETE /api/courses/:id (owner)" "PASS" "200" } else { Add-Result "T10" "DELETE /api/courses/:id (owner)" "FAIL" "Status=$($t10.Status)" }

# T11
$t11 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/courses/$courseId/lessons" -Headers $teacherHeaders -Body @{
    title   = "Lesson $suffix"
    content = "content"
    order   = 1
}
$lessonId = Get-NestedValue -Object $t11.Body -Path @("data", "ID")
if ($t11.Status -eq 201 -and $lessonId) { Add-Result "T11" "POST /api/courses/:id/lessons" "PASS" "201" } else { Add-Result "T11" "POST /api/courses/:id/lessons" "FAIL" "Status=$($t11.Status)" }

# T12
$t12 = Invoke-Api -Method "GET" -Url "$BaseUrl/api/courses/$courseId/lessons" -Headers $teacherHeaders
if ($t12.Status -eq 200) { Add-Result "T12" "GET /api/courses/:id/lessons" "PASS" "200" } else { Add-Result "T12" "GET /api/courses/:id/lessons" "FAIL" "Status=$($t12.Status)" }

# T13
$t13 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/lessons/$lessonId/quiz" -Headers $teacherHeaders -Body @{
    title     = "Quiz $suffix"
    questions = @(
        @{ text = "Q1"; option_a = "A"; option_b = "B"; option_c = "C"; option_d = "D"; correct = "b" },
        @{ text = "Q2"; option_a = "A"; option_b = "B"; option_c = "C"; option_d = "D"; correct = "b" }
    )
}
$quizId = Get-NestedValue -Object $t13.Body -Path @("data", "ID")
if ($t13.Status -eq 201 -and $quizId) { Add-Result "T13" "POST /api/lessons/:id/quiz" "PASS" "201" } else { Add-Result "T13" "POST /api/lessons/:id/quiz" "FAIL" "Status=$($t13.Status)" }

# T14
$quizGet = Invoke-Api -Method "GET" -Url "$BaseUrl/api/lessons/$lessonId/quiz" -Headers $studentHeaders
$answers = @{}
$quizQuestions = Get-NestedValue -Object $quizGet.Body -Path @("data", "questions")
if ($null -ne $quizQuestions) {
    foreach ($q in $quizQuestions) { $answers[[string]$q.ID] = "b" }
}
$t14 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/quiz/$quizId/submit" -Headers $studentHeaders -Body @{ answers = $answers }
$score14 = Get-NestedValue -Object $t14.Body -Path @("data", "score")
if ($t14.Status -eq 200 -and $null -ne $score14 -and $score14 -ge 1) { Add-Result "T14" "POST /api/quiz/:id/submit" "PASS" "200 + score" } else { Add-Result "T14" "POST /api/quiz/:id/submit" "FAIL" "Status=$($t14.Status)" }

# T15
$t15 = Invoke-Api -Method "POST" -Url "$BaseUrl/api/lessons/$lessonId/complete" -Headers $studentHeaders
if ($t15.Status -eq 200) { Add-Result "T15" "POST /api/lessons/:id/complete" "PASS" "200" } else { Add-Result "T15" "POST /api/lessons/:id/complete" "FAIL" "Status=$($t15.Status)" }

# T16
$t16 = Invoke-Api -Method "GET" -Url "$BaseUrl/api/my/progress" -Headers $studentHeaders
if ($t16.Status -eq 200) { Add-Result "T16" "GET /api/my/progress" "PASS" "200" } else { Add-Result "T16" "GET /api/my/progress" "FAIL" "Status=$($t16.Status)" }

# T17
$t17 = Invoke-Api -Method "GET" -Url "$BaseUrl/api/courses?page=1&limit=5" -Headers $teacherHeaders
$keys17 = @("data", "page", "limit", "total")
$ok17 = $t17.Status -eq 200
foreach ($k in $keys17) { $ok17 = $ok17 -and ($t17.Body.PSObject.Properties.Name -contains $k) }
if ($ok17) { Add-Result "T17" "Pagination fields validation" "PASS" "data/page/limit/total" } else { Add-Result "T17" "Pagination fields validation" "FAIL" "Missing fields" }

# T18
$codes = @()
for ($i = 1; $i -le 20; $i++) {
    try {
        $rawBody = @{ email = $teacherEmail; password = "123456" } | ConvertTo-Json -Compress
        $r = Invoke-WebRequest -Method "POST" -Uri "$BaseUrl/api/auth/login" -ContentType "application/json" -Body $rawBody -UseBasicParsing
        $codes += [int]$r.StatusCode
    }
    catch {
        if ($_.Exception.Response) {
            $codes += [int]$_.Exception.Response.StatusCode
        }
        else {
            $codes += 0
        }
    }
}
if ($codes -contains 429) { Add-Result "T18" "20 hızlı istek" "PASS" "429 observed" } else { Add-Result "T18" "20 hızlı istek" "FAIL" "No 429 observed" }

# T19 (WebSocket connect + echo)
try {
    # T18'de yoğun login trafikten sonra limiter etkisini azaltmak için kısa bekleme.
    Start-Sleep -Seconds 2

    $wsPassed = $false
    $wsDetail = ""
    for ($attempt = 1; $attempt -le 3 -and -not $wsPassed; $attempt++) {
        $ws1 = $null
        $ws2 = $null
        try {
            $ws1 = [System.Net.WebSockets.ClientWebSocket]::new()
            $ws1.Options.SetRequestHeader("Authorization", "Bearer $teacherToken")
            $ws2 = [System.Net.WebSockets.ClientWebSocket]::new()
            $ws2.Options.SetRequestHeader("Authorization", "Bearer $studentToken")

            $wsUri = [Uri]"ws://localhost:8090/ws/classroom/$courseId"
            $ws1.ConnectAsync($wsUri, [Threading.CancellationToken]::None).Wait()
            $ws2.ConnectAsync($wsUri, [Threading.CancellationToken]::None).Wait()

            $buf = New-Object byte[] 4096
            $seg = [ArraySegment[byte]]::new($buf)
            $null = $ws2.ReceiveAsync($seg, [Threading.CancellationToken]::None).Result

            $payload = '{"text":"echo-test"}'
            $bytes = [Text.Encoding]::UTF8.GetBytes($payload)
            $ws1.SendAsync([ArraySegment[byte]]::new($bytes), [System.Net.WebSockets.WebSocketMessageType]::Text, $true, [Threading.CancellationToken]::None).Wait()

            $res = $ws2.ReceiveAsync($seg, [Threading.CancellationToken]::None).Result
            $msg = [Text.Encoding]::UTF8.GetString($buf, 0, $res.Count)
            if ($msg -like "*echo-test*") {
                $wsPassed = $true
                $wsDetail = "Echo OK"
            } else {
                $wsDetail = "Echo not received"
            }
        }
        catch {
            $wsDetail = $_.Exception.Message
            Start-Sleep -Milliseconds 700
        }
        finally {
            if ($null -ne $ws1) { $ws1.Dispose() }
            if ($null -ne $ws2) { $ws2.Dispose() }
        }
    }

    if ($wsPassed) {
        Add-Result "T19" "WebSocket connect + mesaj" "PASS" $wsDetail
    } else {
        Add-Result "T19" "WebSocket connect + mesaj" "FAIL" $wsDetail
    }
}
catch {
    Add-Result "T19" "WebSocket connect + mesaj" "FAIL" $_.Exception.Message
}

# T20 (Docker)
try {
    $dockerCmd = Get-Command docker -ErrorAction Stop
    if ($null -ne $dockerCmd) {
        docker compose up --build -d | Out-Null
        Start-Sleep -Seconds 4
        $dockerPs = docker compose ps --format json
        if ($dockerPs) {
            Add-Result "T20" "docker compose up --build" "PASS" "Container started"
        }
        else {
            Add-Result "T20" "docker compose up --build" "FAIL" "No container output"
        }
        docker compose down | Out-Null
    }
}
catch {
    Add-Result "T20" "docker compose up --build" "SKIP" "Docker not available"
}

# T21
$t21 = Invoke-Api -Method "GET" -Url "$BaseUrl/swagger/index.html"
if ($t21.Status -eq 200) { Add-Result "T21" "GET /swagger/index.html" "PASS" "200" } else { Add-Result "T21" "GET /swagger/index.html" "FAIL" "Status=$($t21.Status)" }

# T22 (basic structure + naming sanity)
$requiredPaths = @(
    "config", "database", "handlers", "middleware", "models", "docs", "Dockerfile", "docker-compose.yml", "main.go"
)
$missing = @()
foreach ($p in $requiredPaths) {
    if (-not (Test-Path $p)) { $missing += $p }
}
if ($missing.Count -eq 0) {
    Add-Result "T22" "Klasör yapısı + naming" "PASS" "Required structure exists"
}
else {
    Add-Result "T22" "Klasör yapısı + naming" "FAIL" ("Missing: " + ($missing -join ", "))
}

Write-Host ""
Write-Host "=== 22 AI TEST SENARYOSU RAPORU ===" -ForegroundColor Cyan
$results | Sort-Object Test | Format-Table -AutoSize

$pass = ($results | Where-Object { $_.Status -eq "PASS" } | Measure-Object).Count
$fail = ($results | Where-Object { $_.Status -eq "FAIL" } | Measure-Object).Count
$skip = ($results | Where-Object { $_.Status -eq "SKIP" } | Measure-Object).Count

Write-Host ""
Write-Host "PASS: $pass | FAIL: $fail | SKIP: $skip" -ForegroundColor Yellow
if ($fail -gt 0) { exit 1 } else { exit 0 }

