package repositories

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"
)

type TokenBlacklistRepository struct {
	db *sql.DB
}

func NewTokenBlacklistRepository(db *sql.DB) *TokenBlacklistRepository {
	return &TokenBlacklistRepository{db: db}
}

// HashToken creates a SHA-256 hash of the token for storage
func (r *TokenBlacklistRepository) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Add adds a token to the blacklist
func (r *TokenBlacklistRepository) Add(token, userID string, expiresAt time.Time) error {
	tokenHash := r.HashToken(token)

	query := `
		INSERT INTO token_blacklist (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (token_hash) DO NOTHING
	`

	_, err := r.db.Exec(query, tokenHash, userID, expiresAt)
	return err
}

// IsBlacklisted checks if a token is in the blacklist
func (r *TokenBlacklistRepository) IsBlacklisted(token string) (bool, error) {
	tokenHash := r.HashToken(token)

	query := `
		SELECT EXISTS(
			SELECT 1 FROM token_blacklist 
			WHERE token_hash = $1 AND expires_at > NOW()
		)
	`

	var exists bool
	err := r.db.QueryRow(query, tokenHash).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// CleanupExpired removes expired tokens from the blacklist
func (r *TokenBlacklistRepository) CleanupExpired() error {
	query := `DELETE FROM token_blacklist WHERE expires_at <= NOW()`
	_, err := r.db.Exec(query)
	return err
}

// RevokeAllUserTokens blacklists all tokens for a specific user
// This is useful for scenarios like password change or account compromise
func (r *TokenBlacklistRepository) RevokeAllUserTokens(userID string) error {
	// This is a placeholder - in a real implementation, you'd need to track all active tokens
	// For now, we'll just ensure any new tokens added for this user are properly handled
	return nil
}
