package services

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// StringUtils provides utility functions for string manipulation
type StringUtils struct{}

// NewStringUtils creates a new string utilities instance
func NewStringUtils() *StringUtils {
	return &StringUtils{}
}

// CleanFormatting removes formatting characters from a string
func (s *StringUtils) CleanFormatting(text string) string {
	// Remove color codes and other formatting
	text = strings.ReplaceAll(text, "{", "")
	text = strings.ReplaceAll(text, "}", "")
	text = strings.ReplaceAll(text, "<", "")
	text = strings.ReplaceAll(text, ">", "")
	return text
}

// MD5 returns the MD5 hash of a string
func (s *StringUtils) MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// IsValidEmail checks if a string is a valid email address
func (s *StringUtils) IsValidEmail(email string) bool {
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// SanitizeInput removes potentially dangerous characters from input
func (s *StringUtils) SanitizeInput(input string) string {
	// Remove common injection characters
	input = strings.ReplaceAll(input, "'", "")
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.ReplaceAll(input, ";", "")
	input = strings.ReplaceAll(input, "--", "")
	return input
}

// TruncateString truncates a string to a maximum length
func (s *StringUtils) TruncateString(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}
