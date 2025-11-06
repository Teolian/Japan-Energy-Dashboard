// Package storage provides data persistence layer for energy data.
package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/teo/aversome/backend/pkg/database"
)

// DataStorage handles CRUD operations for energy data
type DataStorage struct {
	db *database.DB
}

// NewDataStorage creates a new data storage instance
func NewDataStorage(db *database.DB) *DataStorage {
	return &DataStorage{db: db}
}

// DataRecord represents a stored data record
type DataRecord struct {
	ID        int
	DataType  string    // 'demand', 'jepx', 'reserve'
	Area      *string   // nullable for system-wide data
	Date      time.Time
	Data      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SaveData stores or updates energy data
func (s *DataStorage) SaveData(dataType string, area *string, date time.Time, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	query := `
		INSERT INTO energy_data (data_type, area, date, data, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (data_type, area, date)
		DO UPDATE SET
			data = EXCLUDED.data,
			updated_at = NOW()
	`

	_, err = s.db.Exec(query, dataType, area, date, jsonData)
	if err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

// GetData retrieves energy data by type, area, and date
func (s *DataStorage) GetData(dataType string, area *string, date time.Time) (json.RawMessage, error) {
	query := `
		SELECT data
		FROM energy_data
		WHERE data_type = $1
		  AND (area = $2 OR (area IS NULL AND $2 IS NULL))
		  AND date = $3
	`

	var data json.RawMessage
	err := s.db.QueryRow(query, dataType, area, date).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("data not found for %s/%v/%s", dataType, area, date.Format("2006-01-02"))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	return data, nil
}

// GetLatestData retrieves the most recent data for a given type and area
func (s *DataStorage) GetLatestData(dataType string, area *string) (json.RawMessage, time.Time, error) {
	query := `
		SELECT data, date
		FROM energy_data
		WHERE data_type = $1
		  AND (area = $2 OR (area IS NULL AND $2 IS NULL))
		ORDER BY date DESC
		LIMIT 1
	`

	var data json.RawMessage
	var date time.Time
	err := s.db.QueryRow(query, dataType, area).Scan(&data, &date)
	if err == sql.ErrNoRows {
		return nil, time.Time{}, fmt.Errorf("no data found for %s/%v", dataType, area)
	}
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to get latest data: %w", err)
	}

	return data, date, nil
}

// ListDates returns all available dates for a data type
func (s *DataStorage) ListDates(dataType string, area *string, limit int) ([]time.Time, error) {
	query := `
		SELECT DISTINCT date
		FROM energy_data
		WHERE data_type = $1
		  AND (area = $2 OR (area IS NULL AND $2 IS NULL))
		ORDER BY date DESC
		LIMIT $3
	`

	rows, err := s.db.Query(query, dataType, area, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list dates: %w", err)
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("failed to scan date: %w", err)
		}
		dates = append(dates, date)
	}

	return dates, nil
}

// DeleteOldData removes data older than the specified duration
func (s *DataStorage) DeleteOldData(olderThan time.Duration) (int64, error) {
	cutoffDate := time.Now().Add(-olderThan)

	query := `DELETE FROM energy_data WHERE date < $1`

	result, err := s.db.Exec(query, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetDataCount returns total number of records
func (s *DataStorage) GetDataCount() (int64, error) {
	var count int64
	err := s.db.QueryRow(`SELECT COUNT(*) FROM energy_data`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get data count: %w", err)
	}
	return count, nil
}
