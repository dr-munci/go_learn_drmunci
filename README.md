# Mahmut KOÇOĞLU
Öğrenci No: 24080410003

## GoLearn - Uzaktan Eğitim Platformu API

GoLearn, Gin + GORM + SQLite tabanlı bir uzaktan eğitim platformu backend API projesidir.  
Projede kimlik doğrulama (JWT), rol bazlı yetkilendirme (RBAC), kurs/lesson/quiz yönetimi, ilerleme takibi, rate limiting, WebSocket canlı sınıf iletişimi, Swagger dokümantasyonu ve Docker desteği bulunmaktadır.

## Teknoloji Yığını

- Go
- Gin (HTTP framework)
- GORM (ORM)
- SQLite (pure Go driver: `github.com/glebarez/sqlite`)
- JWT (`github.com/golang-jwt/jwt/v5`)
- WebSocket (`github.com/gorilla/websocket`)
- Swagger (`swaggo`)
- Docker + Docker Compose

## Proje Yapısı

```text
golearn/
├── config/
├── database/
│   └── db.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── handlers/
│   ├── auth.go
│   ├── course.go
│   ├── lesson.go
│   ├── progress.go
│   ├── quiz.go
│   └── websocket.go
├── middleware/
│   ├── auth.go
│   ├── rbac.go
│   └── ratelimit.go
├── models/
│   ├── course.go
│   ├── lesson.go
│   ├── progress.go
│   ├── quiz.go
│   └── user.go
├── scripts/
│   └── run-ai-checklist.ps1
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

## Ortam Gereksinimleri

- Go 1.23+ (öneri: güncel stabil sürüm)
- PowerShell (Windows testleri için)
- Node.js + npm (`wscat` ile WebSocket testi için)
- Docker Desktop (Docker adımı için, opsiyonel)

## Kurulum

```bash
go mod tidy
```

Swagger CLI kurmak için:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Swagger dosyalarını üretmek için:

```bash
swag init
```

## Uygulamayı Çalıştırma

```bash
go run main.go
```

Varsayılan port: `8090`  
Base URL: `http://localhost:8090`

## Değerlendirme İçin Hızlı Test

Projeyi kontrol edecek kişi için en hızlı doğrulama adımları:

1. Uygulamayı çalıştır:

```bash
go mod tidy
swag init
go run main.go
```

2. Swagger üzerinden test et:

- `http://localhost:8090/swagger/index.html`

3. Otomatik checklist scriptini çalıştır (PowerShell):

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\run-ai-checklist.ps1
```

4. Docker ile test et (opsiyonel):

```bash
docker compose up --build -d
docker compose ps
```

## Swagger

Sunucu çalışırken:

- `http://localhost:8090/swagger/index.html`

## Kimlik Doğrulama ve Yetkilendirme

- Auth mekanizması: JWT Bearer Token
- Korumalı endpointler `/api` altında
- Öğretmen yetkisi gereken işlemlerde `TeacherOnly` middleware kullanılır

Authorization header formatı:

```text
Authorization: Bearer <TOKEN>
```

## Ana Endpointler

### Auth

- `POST /api/auth/register`
- `POST /api/auth/login`

### Courses

- `GET /api/courses` (pagination/filter/sort)
- `GET /api/courses/:id`
- `POST /api/courses` (teacher)
- `PUT /api/courses/:id` (owner teacher)
- `DELETE /api/courses/:id` (owner teacher)

### Lessons

- `GET /api/courses/:id/lessons`
- `POST /api/courses/:id/lessons` (teacher owner)
- `POST /api/lessons/:id/complete`

### Quiz

- `GET /api/lessons/:id/quiz`
- `POST /api/lessons/:id/quiz` (teacher)
- `POST /api/quiz/:id/submit` (otomatik puanlama)

### Progress

- `GET /api/my/progress`

### WebSocket

- `GET /ws/classroom/:courseId`

## Pagination / Filtering / Sorting

`GET /api/courses` endpointi şu query parametrelerini destekler:

- `page` (default: 1)
- `limit` (default: 10)
- `category`
- `sort` (ör: `title asc`, default: `created_at desc`)

Başarılı response alanları:

- `data`
- `page`
- `limit`
- `total`

## Rate Limiting

IP bazlı rate limiting uygulanır:

- 5 istek/saniye
- burst: 10

Limit aşıldığında:

- HTTP `429 Too Many Requests`
- `{"error":"Çok fazla istek gönderdiniz, lütfen bekleyin"}`

## WebSocket Testi (PowerShell)

Önce iki kullanıcı için token alıp iki ayrı PowerShell terminalinde bağlanın:

```powershell
wscat.cmd -c "ws://localhost:8090/ws/classroom/1" -H "Authorization: Bearer <TOKEN>"
```

Mesaj gönderimi:

```json
{"text":"Merhaba sinif!"}
```

## Docker

Build ve ayağa kaldırma:

```bash
docker compose up --build -d
```

Log izleme:

```bash
docker compose logs -f
```

Durdurma:

```bash
docker compose down
```

## Otomatik Test Scripti (22 Senaryo)

PowerShell test scripti:

- `scripts/run-ai-checklist.ps1`

Çalıştırma:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\run-ai-checklist.ps1
```

Script çıktısı:

- `PASS / FAIL / SKIP` bazlı tablo
- Toplam sonuç özeti (`PASS`, `FAIL`, `SKIP`)

## Notlar

- Parola alanı JSON response'larda gizlidir (`json:"-"`).
- Quiz doğru cevap alanı response'da gizlidir (`json:"-"`).
- Response yapısı genel olarak:
  - Hata: `{"error":"..."}`
  - Başarı: `{"message":"...","data":...}`

