package bcrypt

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Hash simple password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Hash empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "Hash long password",
			password: strings.Repeat("a", 1000),
			wantErr:  false,
		},
		{
			name:     "Hash special characters",
			password: "!@#$%^&*()_+-=[]{}|;':\",./<>?`~",
			wantErr:  false,
		},
		{
			name:     "Hash unicode characters",
			password: "密码测试🔐",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if !strings.HasPrefix(hash, "$argon2id$") {
					t.Errorf("HashPassword() hash doesn't have correct prefix, got %s", hash)
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testPassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		want           bool
	}{
		{
			name:           "Correct password",
			password:       password,
			hashedPassword: hash,
			want:           true,
		},
		{
			name:           "Wrong password",
			password:       "wrongPassword",
			hashedPassword: hash,
			want:           false,
		},
		{
			name:           "Empty password with valid hash",
			password:       "",
			hashedPassword: hash,
			want:           false,
		},
		{
			name:           "Valid password with invalid hash",
			password:       password,
			hashedPassword: "invalid-hash",
			want:           false,
		},
		{
			name:           "Valid password with empty hash",
			password:       password,
			hashedPassword: "",
			want:           false,
		},
		{
			name:           "Valid password with malformed hash",
			password:       password,
			hashedPassword: "$argon2id$invalid",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyPassword(tt.password, tt.hashedPassword); got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
	}{
		{
			name: "Invalid hash - wrong number of parts",
			hash: "$argon2id$invalid",
		},
		{
			name: "Invalid hash - wrong algorithm",
			hash: "$bcrypt$v=19$m=65536,t=3,p=4$salt$hash",
		},
		{
			name: "Invalid hash - wrong version format",
			hash: "$argon2id$version=19$m=65536,t=3,p=4$salt$hash",
		},
		{
			name: "Invalid hash - invalid base64 salt",
			hash: "$argon2id$v=19$m=65536,t=3,p=4$!!!invalid!!!$hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := decodeHash(tt.hash)
			if err == nil {
				t.Errorf("decodeHash() expected error for hash: %s", tt.hash)
			}
		})
	}
}

func TestHashUniqueness(t *testing.T) {
	password := "samePassword"
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword() should generate unique hashes for the same password due to random salt")
	}

	if !VerifyPassword(password, hash1) {
		t.Error("First hash should verify correctly")
	}
	if !VerifyPassword(password, hash2) {
		t.Error("Second hash should verify correctly")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := HashPassword(password)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyPassword(password, hash)
	}
}
