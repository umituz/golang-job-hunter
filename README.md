# Golang Job Hunter 🚀

A powerful job scraper API built with **Golang** and **Gin framework** that automatically finds and filters tech jobs containing specific keywords like **Flutter**, **Laravel**, and **Golang**. Perfect for developers looking for their next opportunity!

## ✨ Features

- **Smart Job Filtering**: Automatically detects jobs containing Flutter, Laravel, Golang keywords
- **RemoteOK Integration**: Scrapes jobs from RemoteOK API with real-time data
- **RESTful API**: Clean and well-documented endpoints
- **Duplicate Prevention**: Smart duplicate detection prevents data redundancy  
- **Keyword Matching**: Advanced text analysis for accurate job categorization
- **SQLite Database**: Lightweight, zero-configuration database
- **CORS Enabled**: Ready for frontend integration
- **Pagination Support**: Efficient data handling for large result sets
- **Search Functionality**: Advanced search with location and keyword filters

## 🛠️ Tech Stack

- **Backend**: [Golang](https://golang.org/) with [Gin Framework](https://gin-gonic.com/)
- **Database**: SQLite with [GORM](https://gorm.io/) ORM
- **Job Source**: [RemoteOK API](https://remoteok.io/)
- **HTTP Client**: Native Go http package with JSON parsing

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/umituz/golang-job-hunter.git
   cd golang-job-hunter/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Test the API**
   ```bash
   curl http://localhost:8081/health
   ```

## 📚 API Documentation

### Base URL
```
http://localhost:8081
```

### Endpoints

#### Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "ok",
  "message": "Job Hunter API is running",
  "time": "2025-08-20T14:06:15.467655+03:00"
}
```

#### Get Filtered Jobs
```http
GET /api/v1/jobs?page=1&limit=20
```
**Response:**
```json
{
  "jobs": [
    {
      "id": 73,
      "title": "Fullstack Golang Engineer",
      "company": "Discern",
      "description": "...",
      "location": "Stamford, United States",
      "salary": "50000",
      "url": "https://remoteok.io/remote-jobs/...",
      "source": "RemoteOK",
      "remote": false,
      "keywords": "golang,go",
      "hasKeywords": true,
      "createdAt": "2025-08-20T14:07:16.486995+03:00",
      "updatedAt": "2025-08-20T14:07:16.486995+03:00"
    }
  ],
  "total": 75,
  "page": 1,
  "limit": 20,
  "pages": 4
}
```

#### Search Jobs
```http
POST /api/v1/search
Content-Type: application/json

{
  "keywords": ["golang", "flutter"],
  "location": "remote",
  "page": 1,
  "limit": 10
}
```

#### Get Single Job
```http
GET /api/v1/jobs/{id}
```

#### Scrape New Jobs
```http
POST /api/v1/scrape
```
**Response:**
```json
{
  "message": "Scraping completed",
  "scraped": 75,
  "saved": 75,
  "duplicates": 0,
  "time": "2025-08-20T14:07:16.519624+03:00"
}
```

#### Get Statistics
```http
GET /api/v1/stats
```
**Response:**
```json
{
  "totalJobs": 75,
  "jobsWithKeywords": 75,
  "keywordsPercentage": 100,
  "lastUpdated": "2025-08-20T14:07:38.382134+03:00"
}
```

## 🏗️ Architecture

```
├── main.go                 # Application entry point
├── models/
│   └── job.go             # Job data model with keyword detection
├── database/
│   └── database.go        # Database operations and queries  
├── scrapers/
│   └── remoteok.go        # RemoteOK API integration
├── go.mod                 # Go module dependencies
└── jobs.db               # SQLite database (auto-created)
```

## 🔄 How It Works

1. **Scraping**: The application fetches job data from RemoteOK API
2. **Filtering**: Each job is analyzed for Flutter, Laravel, Golang keywords
3. **Storage**: Filtered jobs are stored in SQLite with automatic keyword tagging
4. **API**: Clean REST endpoints provide access to filtered job data
5. **Search**: Advanced search capabilities with pagination and filters

## 🎯 Keyword Detection

The system automatically detects these keywords in job titles, descriptions, and company names:
- **Flutter** - Mobile app development
- **Laravel** - PHP web framework  
- **Golang/Go** - Backend development

## 📊 Example Results

Recent scraping session found:
- **75 total jobs** scraped from RemoteOK
- **100% keyword match rate** (only relevant jobs stored)
- Popular positions: Senior Golang Engineer, Flutter Developer, Full-stack Engineer
- Salary range: $50k - $140k+

## 🔧 Development

### Project Structure
```go
// Job Model with automatic keyword detection
type Job struct {
    ID          uint      `json:"id" gorm:"primarykey"`
    Title       string    `json:"title" gorm:"not null"`
    Company     string    `json:"company" gorm:"not null"`
    Description string    `json:"description" gorm:"type:text"`
    Location    string    `json:"location"`
    Salary      string    `json:"salary"`
    URL         string    `json:"url" gorm:"unique;not null"`
    Source      string    `json:"source"`
    Remote      bool      `json:"remote"`
    Keywords    string    `json:"keywords"`
    HasKeywords bool      `json:"hasKeywords"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### Adding New Job Sources

1. Create a new scraper in `/scrapers/`
2. Implement the scraping logic
3. Add to the scraper endpoint in `main.go`

### Database Schema

The application uses SQLite with GORM for automatic migrations:
- **jobs** table with indexes on keywords and timestamps
- **Unique constraints** on job URLs to prevent duplicates
- **Automatic keyword detection** via GORM hooks

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [RemoteOK API](https://remoteok.io/api) - Job data source
- [Gin Framework](https://gin-gonic.com/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [SQLite](https://sqlite.org/) - Database

## 📧 Contact

Umit Uz - [@umituz](https://github.com/umituz)

Project Link: [https://github.com/umituz/golang-job-hunter](https://github.com/umituz/golang-job-hunter)

---

⭐ **Star this repository if it helped you find your dream job!** ⭐