package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func main() {
	fmt.Println("main is called")

	dsn := "host=" + os.Getenv("DB_HOST") +
		// " user=" + os.Getenv("DB_USER") +
		" user=tenant_user" +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	connector := &connector{dsn: dsn, d: &stdlib.Driver{}}
	sqlDB := sql.OpenDB(connector)

	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	http.Handle("GET /tenants/{tenant_id}/projects", setTenantMiddleware(http.HandlerFunc(getProjects)))
	http.Handle("GET /tenants/{tenant_id}/project_tags/{id}", setTenantMiddleware(http.HandlerFunc(getProjectTag)))

	slog.Info("Server is running on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Failed to start server: %v", err)
	}
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	var projects []Project

	// コネクション作成、取得時にcontextからtenant_idを取得し実行時パラメータにセットするためにcontextを渡す
	// WHERE句を付けななくても、contextからtenant_idを取得して実行時パラメータにセットするため、organization_idによる絞り込みが行われる
	db.Debug().WithContext(r.Context()).Find(&projects)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func getProjectTag(w http.ResponseWriter, r *http.Request) {
	var tag ProjectTag

	// PreloadでRLSポリシーを適用しているテーブルを取得する場合は、WithContextでcontextを渡す必要がある
	db.Debug().WithContext(r.Context()).Preload("Project").Find(&tag, "id = ?", r.PathValue("id"))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

func setTenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.PathValue("tenant_id")
		if tenantID == "" {
			http.Error(w, "Tenant ID is required", http.StatusBadRequest)
			return
		}

		// NOTE: 認証後にcontextにtenant_idをセットする
		ctx := context.WithValue(r.Context(), "tenant_id", tenantID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	TenantID string `json:"tenant_id"`
	Tenant   Tenant `json:"tenant" gorm:"foreignKey:TenantID;references:ID"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

type ProjectTag struct {
	ProjectID string  `json:"project_id"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID;references:ID"`
	ID        string  `json:"id"`
}
