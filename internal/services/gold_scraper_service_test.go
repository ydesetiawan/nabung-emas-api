package services

import (
	"nabung-emas-api/internal/models"
	"testing"
)

func TestParsePrice(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      int64
		wantError bool
	}{
		{
			name:      "Standard Indonesian format",
			input:     "Rp1.234.567",
			want:      1234567,
			wantError: false,
		},
		{
			name:      "Without Rp prefix",
			input:     "1.234.567",
			want:      1234567,
			wantError: false,
		},
		{
			name:      "Plain number",
			input:     "1234567",
			want:      1234567,
			wantError: false,
		},
		{
			name:      "With IDR",
			input:     "IDR 1.234.567",
			want:      1234567,
			wantError: false,
		},
		{
			name:      "Small amount",
			input:     "Rp500",
			want:      500,
			wantError: false,
		},
		{
			name:      "Large amount",
			input:     "Rp123.456.789",
			want:      123456789,
			wantError: false,
		},
		{
			name:      "Empty string",
			input:     "",
			want:      0,
			wantError: true,
		},
		{
			name:      "Invalid format",
			input:     "abc",
			want:      0,
			wantError: true,
		},
		{
			name:      "Negative number",
			input:     "-1000",
			want:      0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePrice(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("parsePrice() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("parsePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectCategory(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		want        models.GoldCategory
	}{
		{
			name:        "Standard gold bar",
			productName: "Logam Mulia 1 gram",
			want:        models.GoldCategoryEmasBatangan,
		},
		{
			name:        "Gift series",
			productName: "Emas Batangan Gift Series 5 gram",
			want:        models.GoldCategoryEmasBatanganGiftSeries,
		},
		{
			name:        "Idul Fitri edition",
			productName: "Emas Batangan Selamat Idul Fitri 10 gram",
			want:        models.GoldCategoryEmasBatanganSelamatIdulFitri,
		},
		{
			name:        "Lebaran edition",
			productName: "Emas Lebaran 25 gram",
			want:        models.GoldCategoryEmasBatanganSelamatIdulFitri,
		},
		{
			name:        "Imlek edition",
			productName: "Emas Batangan Imlek 50 gram",
			want:        models.GoldCategoryEmasBatanganImlek,
		},
		{
			name:        "Chinese New Year",
			productName: "Gold Bar Chinese New Year 100 gram",
			want:        models.GoldCategoryEmasBatanganImlek,
		},
		{
			name:        "Batik Series III",
			productName: "Emas Batangan Batik Seri III 5 gram",
			want:        models.GoldCategoryEmasBatanganBatikSeriIII,
		},
		{
			name:        "Batik general",
			productName: "Emas Batik 10 gram",
			want:        models.GoldCategoryEmasBatanganBatikSeriIII,
		},
		{
			name:        "Pure silver",
			productName: "Perak Murni 100 gram",
			want:        models.GoldCategoryPerakMurni,
		},
		{
			name:        "Silver",
			productName: "Silver 50 gram",
			want:        models.GoldCategoryPerakMurni,
		},
		{
			name:        "Heritage silver",
			productName: "Perak Heritage 250 gram",
			want:        models.GoldCategoryPerakHeritage,
		},
		{
			name:        "Liontin Batik",
			productName: "Liontin Batik Seri III",
			want:        models.GoldCategoryLiontinBatikSeriIII,
		},
		{
			name:        "Pendant",
			productName: "Gold Pendant 1 gram",
			want:        models.GoldCategoryLiontinBatikSeriIII,
		},
		{
			name:        "Unknown product defaults to emas batangan",
			productName: "Unknown Product",
			want:        models.GoldCategoryEmasBatangan,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := detectCategory(tt.productName)
			if got != tt.want {
				t.Errorf("detectCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Multiple spaces",
			input: "Logam  Mulia   1  gram",
			want:  "Logam Mulia 1 gram",
		},
		{
			name:  "Newlines and tabs",
			input: "Logam\nMulia\t1 gram",
			want:  "Logam Mulia 1 gram",
		},
		{
			name:  "Leading and trailing spaces",
			input: "  Logam Mulia 1 gram  ",
			want:  "Logam Mulia 1 gram",
		},
		{
			name:  "Mixed whitespace",
			input: "  Logam\n\tMulia  1  gram  ",
			want:  "Logam Mulia 1 gram",
		},
		{
			name:  "Already clean",
			input: "Logam Mulia 1 gram",
			want:  "Logam Mulia 1 gram",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanText(tt.input)
			if got != tt.want {
				t.Errorf("cleanText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanPrice(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Standard format with Rp and dots",
			input: "Rp1.234.567",
			want:  "1234567",
		},
		{
			name:  "With IDR",
			input: "IDR 1.234.567",
			want:  "1234567",
		},
		{
			name:  "With commas",
			input: "Rp1,234,567",
			want:  "1234567",
		},
		{
			name:  "Mixed separators",
			input: "Rp 1.234,567",
			want:  "1234567",
		},
		{
			name:  "Already clean",
			input: "1234567",
			want:  "1234567",
		},
		{
			name:  "With extra spaces",
			input: "  Rp  1.234.567  ",
			want:  "1234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanPrice(tt.input)
			if got != tt.want {
				t.Errorf("cleanPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
