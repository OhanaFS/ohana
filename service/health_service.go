package service

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type healthStatus string

const (
	HealthStatusGood     healthStatus = "good"
	HealthStatusDegraded healthStatus = "degraded"
	HealthStatusBad      healthStatus = "bad"
)

// ClusterHealth contains various information on the cluster's health.
type ClusterHealth struct {
	Status   healthStatus  `json:"status"`
	Uptime   time.Duration `json:"uptime"`
	DBStatus sql.DBStats   `json:"db_status"`
}

// Health is the service responsible for monitoring cluster health.
type Health interface {
	Status() *ClusterHealth
}

// health is the implementation for the Health service.
type health struct {
	createdAt time.Time
	db        *gorm.DB
}

// NewHealth creates a new instance of the Health service.
func NewHealth(db *gorm.DB) Health {
	return &health{
		createdAt: time.Now(),
		db:        db,
	}
}

// Status returns the status of the cluster. For now, it will always return
// good.
func (h *health) Status() *ClusterHealth {

	db, err := h.db.DB()

	var dbStats sql.DBStats

	if err == nil {
		dbStats = db.Stats()
	}

	return &ClusterHealth{
		Status:   HealthStatusGood,
		Uptime:   time.Now().Sub(h.createdAt),
		DBStatus: dbStats,
	}
}
