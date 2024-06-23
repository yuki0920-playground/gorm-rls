package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func main() {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tenants/{id}/projects", getProjects)

	slog.Info("Server is running on port 8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Failed to start server: %v", err)
	}
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("id")
	var projects []Project

	db.Debug().Where("projects.tenant_id = ?", tenantID).Find(&projects)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

type Tenant struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Projects []Project `json:"projects" gorm:"foreignKey:TenantID"`
}

type Project struct {
	TenantID string `json:"tenant_id"`
	Tenant   Tenant `json:"tenant" gorm:"foreignKey:TenantID"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}
