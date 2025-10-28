package adapters

import (
	"os"
	"strings"
	"testing"
)

func TestJEPXAdapter_ParseCSV(t *testing.T) {
	tests := []struct {
		name        string
		csvPath     string
		date        string
		area        string
		wantErr     bool
		wantHours   int
		wantPrice0  float64 // Price at hour 0
		wantPrice23 float64 // Price at hour 23
	}{
		{
			name:        "valid CSV with Tokyo prices",
			csvPath:     "testdata/jepx-sample.csv",
			date:        "2025-10-23",
			area:        "tokyo",
			wantErr:     false,
			wantHours:   24,
			wantPrice0:  24.32,
			wantPrice23: 25.60,
		},
		{
			name:        "valid CSV with Kansai prices",
			csvPath:     "testdata/jepx-sample.csv",
			date:        "2025-10-23",
			area:        "kansai",
			wantErr:     false,
			wantHours:   24,
			wantPrice0:  23.15,
			wantPrice23: 24.80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.csvPath)
			if err != nil {
				t.Fatalf("Failed to open test CSV: %v", err)
			}
			defer f.Close()

			adapter := NewJEPXAdapter()
			resp, err := adapter.ParseCSV(f, tt.date, tt.area)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Validate response structure
			if resp.Date != tt.date {
				t.Errorf("Date = %v, want %v", resp.Date, tt.date)
			}

			if resp.Area != tt.area {
				t.Errorf("Area = %v, want %v", resp.Area, tt.area)
			}

			if resp.Timescale != "hourly" {
				t.Errorf("Timescale = %v, want hourly", resp.Timescale)
			}

			// Validate number of price points
			if len(resp.PriceYenPerKwh) != tt.wantHours {
				t.Errorf("Got %d price points, want %d", len(resp.PriceYenPerKwh), tt.wantHours)
			}

			// Validate first price point
			if len(resp.PriceYenPerKwh) > 0 {
				first := resp.PriceYenPerKwh[0]
				if first.Price != tt.wantPrice0 {
					t.Errorf("First price = %v, want %v", first.Price, tt.wantPrice0)
				}
				// Validate timestamp format
				if !strings.Contains(first.Timestamp, tt.date) {
					t.Errorf("Timestamp %v doesn't contain date %v", first.Timestamp, tt.date)
				}
				if !strings.Contains(first.Timestamp, "+09:00") {
					t.Errorf("Timestamp %v doesn't have Asia/Tokyo offset (+09:00)", first.Timestamp)
				}
			}

			// Validate last price point
			if len(resp.PriceYenPerKwh) > 0 {
				last := resp.PriceYenPerKwh[len(resp.PriceYenPerKwh)-1]
				if last.Price != tt.wantPrice23 {
					t.Errorf("Last price = %v, want %v", last.Price, tt.wantPrice23)
				}
			}

			// Validate source attribution
			if resp.Source.Name != "JEPX" {
				t.Errorf("Source name = %v, want JEPX", resp.Source.Name)
			}
			if resp.Source.URL == "" {
				t.Error("Source URL is empty")
			}
		})
	}
}

func TestJEPXAdapter_normalizeDate(t *testing.T) {
	adapter := NewJEPXAdapter()

	tests := []struct {
		input string
		want  string
	}{
		{"2025-10-23", "2025-10-23"},
		{"2025/10/23", "2025-10-23"},
		{"20251023", "2025-10-23"},
		{" 2025-10-23 ", "2025-10-23"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := adapter.normalizeDate(tt.input)
			if got != tt.want {
				t.Errorf("normalizeDate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestJEPXAdapter_buildTimestamp(t *testing.T) {
	adapter := NewJEPXAdapter()

	tests := []struct {
		date        string
		hour        int
		wantHour    string
		wantOffset  string
	}{
		{"2025-10-23", 0, "T00:00:00", "+09:00"},
		{"2025-10-23", 12, "T12:00:00", "+09:00"},
		{"2025-10-23", 23, "T23:00:00", "+09:00"},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			got := adapter.buildTimestamp(tt.date, tt.hour)
			if !strings.Contains(got, tt.date) {
				t.Errorf("Timestamp %v doesn't contain date %v", got, tt.date)
			}
			if !strings.Contains(got, tt.wantHour) {
				t.Errorf("Timestamp %v doesn't contain hour %v", got, tt.wantHour)
			}
			if !strings.Contains(got, tt.wantOffset) {
				t.Errorf("Timestamp %v doesn't contain offset %v", got, tt.wantOffset)
			}
		})
	}
}
