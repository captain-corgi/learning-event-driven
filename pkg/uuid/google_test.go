package uuid

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewGoogle(t *testing.T) {
	tests := []struct {
		name string
		want int // expected length of UUID string
	}{
		{
			name: "generates valid UUID",
			want: 36, // UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGoogle()

			// Test length
			if len(got) != tt.want {
				t.Errorf("NewGoogle() length = %v, want %v", len(got), tt.want)
			}

			// Test format - should be parseable by google/uuid
			if _, err := uuid.Parse(got); err != nil {
				t.Errorf("NewGoogle() generated invalid UUID: %v", err)
			}

			// Test uniqueness by generating multiple UUIDs
			second := NewGoogle()
			if got == second {
				t.Errorf("NewGoogle() generated duplicate UUIDs: %v", got)
			}

			// Test format pattern (8-4-4-4-12 with hyphens)
			parts := strings.Split(got, "-")
			if len(parts) != 5 {
				t.Errorf("NewGoogle() format invalid, expected 5 parts separated by hyphens, got %d", len(parts))
			}

			expectedLengths := []int{8, 4, 4, 4, 12}
			for i, part := range parts {
				if len(part) != expectedLengths[i] {
					t.Errorf("NewGoogle() part %d length = %d, want %d", i, len(part), expectedLengths[i])
				}
			}
		})
	}
}

func TestParseGoogle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid UUID v4",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			want:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "valid UUID v1",
			input:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			want:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			wantErr: false,
		},
		{
			name:    "valid UUID uppercase",
			input:   "550E8400-E29B-41D4-A716-446655440000",
			want:    "550e8400-e29b-41d4-a716-446655440000", // should normalize to lowercase
			wantErr: false,
		},
		{
			name:    "valid UUID without hyphens",
			input:   "550e8400e29b41d4a716446655440000",
			want:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "nil UUID",
			input:   "00000000-0000-0000-0000-000000000000",
			want:    "00000000-0000-0000-0000-000000000000",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - too short",
			input:   "550e8400-e29b-41d4-a716",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - too long",
			input:   "550e8400-e29b-41d4-a716-446655440000-extra",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "550g8400-e29b-41d4-a716-446655440000",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong separators",
			input:   "550e8400_e29b_41d4_a716_446655440000",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - missing hyphens",
			input:   "550e8400e29b41d4a716446655440000extra",
			want:    "",
			wantErr: true,
		},
		{
			name:    "whitespace",
			input:   "   ",
			want:    "",
			wantErr: true,
		},
		{
			name:    "special characters",
			input:   "!@#$%^&*()_+-=[]{}|;:,.<>?",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGoogle(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGoogle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ParseGoogle() = %v, want %v", got, tt.want)
			}

			// If no error expected, verify the result is a valid UUID
			if !tt.wantErr && got != "" {
				if _, parseErr := uuid.Parse(got); parseErr != nil {
					t.Errorf("ParseGoogle() returned invalid UUID: %v", parseErr)
				}
			}
		})
	}
}

func TestMustParseGoogle(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantPanic bool
	}{
		{
			name:      "valid UUID v4",
			input:     "550e8400-e29b-41d4-a716-446655440000",
			want:      "550e8400-e29b-41d4-a716-446655440000",
			wantPanic: false,
		},
		{
			name:      "valid UUID v1",
			input:     "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			want:      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			wantPanic: false,
		},
		{
			name:      "valid UUID uppercase",
			input:     "550E8400-E29B-41D4-A716-446655440000",
			want:      "550e8400-e29b-41d4-a716-446655440000",
			wantPanic: false,
		},
		{
			name:      "valid UUID without hyphens",
			input:     "550e8400e29b41d4a716446655440000",
			want:      "550e8400-e29b-41d4-a716-446655440000",
			wantPanic: false,
		},
		{
			name:      "nil UUID",
			input:     "00000000-0000-0000-0000-000000000000",
			want:      "00000000-0000-0000-0000-000000000000",
			wantPanic: false,
		},
		{
			name:      "empty string should panic",
			input:     "",
			want:      "",
			wantPanic: true,
		},
		{
			name:      "invalid format should panic",
			input:     "invalid-uuid",
			want:      "",
			wantPanic: true,
		},
		{
			name:      "invalid characters should panic",
			input:     "550g8400-e29b-41d4-a716-446655440000",
			want:      "",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustParseGoogle() should have panicked for input: %v", tt.input)
					}
				}()
				MustParseGoogle(tt.input)
			} else {
				got := MustParseGoogle(tt.input)
				if got != tt.want {
					t.Errorf("MustParseGoogle() = %v, want %v", got, tt.want)
				}

				// Verify the result is a valid UUID
				if _, err := uuid.Parse(got); err != nil {
					t.Errorf("MustParseGoogle() returned invalid UUID: %v", err)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewGoogle(b *testing.B) {
	for b.Loop() {
		NewGoogle()
	}
}

func BenchmarkParseGoogle(b *testing.B) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	b.ResetTimer()

	for b.Loop() {
		ParseGoogle(validUUID)
	}
}

func BenchmarkMustParseGoogle(b *testing.B) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	b.ResetTimer()

	for b.Loop() {
		MustParseGoogle(validUUID)
	}
}
