package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func main() {
	dsn := "host=" + os.Getenv("DB_HOST") +
		// " user=" + os.Getenv("DB_USER") +
		" user=tenant_user" +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)

	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	http.Handle("GET /tenants/{tenant_id}/projects", setTenantMiddleware(http.HandlerFunc(getProjects)))
	slog.Info("Server is running on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Failed to start server: %v", err)
	}
}

func setTenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: 本来的には認証後にcontextにtenantIDを設定し、その値を実行時パラメータに設定する
		tenantID := r.PathValue("tenant_id")
		if tenantID == "" {
			http.Error(w, "Tenant ID is required", http.StatusBadRequest)
			return
		}

		// テナントIDを設定
		db.Exec(fmt.Sprintf("SET app.tenant_id = '%s'", tenantID))

		next.ServeHTTP(w, r)

		// リクエストが終了したらクリア
		db.Exec("RESET app.tenant_id")
	})
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	// tenantID := r.PathValue("id")
	var projects []Project

	// db.Debug().Where("projects.tenant_id = ?", tenantID).Find(&projects)

	// 同時にリクエストを受け付けても問題ないか確認するために30秒待つ
	time.Sleep(30 * time.Second)

	// WHERE句を付けないと全てのプロジェクトが取得される
	db.Debug().Find(&projects)

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
