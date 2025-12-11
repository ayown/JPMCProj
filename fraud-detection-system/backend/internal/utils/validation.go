package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`\+?[0-9]{10,15}`)
	urlRegex   = regexp.MustCompile(`https?://[^\s]+`)
)

// ValidateEmail checks if an email is valid
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if a password meets requirements
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// ExtractURLs extracts all URLs from text
func ExtractURLs(text string) []string {
	urls := urlRegex.FindAllString(text, -1)
	if urls == nil {
		return []string{}
	}
	return urls
}

// ExtractPhoneNumbers extracts phone numbers from text
func ExtractPhoneNumbers(text string) []string {
	phones := phoneRegex.FindAllString(text, -1)
	if phones == nil {
		return []string{}
	}
	return phones
}

// HasUrgentWords checks if text contains urgent/panic words
func HasUrgentWords(text string) bool {
	urgentWords := []string{
		"urgent", "immediately", "expire", "suspended", "blocked",
		"verify now", "act now", "limited time", "last chance",
		"mandatory", "required", "within 24 hours", "account closed",
	}

	lowerText := strings.ToLower(text)
	for _, word := range urgentWords {
		if strings.Contains(lowerText, word) {
			return true
		}
	}
	return false
}

// CountUrgentWords counts urgent words in text
func CountUrgentWords(text string) int {
	urgentWords := []string{
		"urgent", "immediately", "expire", "suspended", "blocked",
		"verify now", "act now", "limited time", "last chance",
		"mandatory", "required", "within 24 hours", "account closed",
	}

	count := 0
	lowerText := strings.ToLower(text)
	for _, word := range urgentWords {
		count += strings.Count(lowerText, word)
	}
	return count
}

// HasKYCKeywords checks if text contains KYC-related keywords
func HasKYCKeywords(text string) bool {
	kycWords := []string{
		"kyc", "know your customer", "verification", "verify",
		"update details", "update information", "pan card",
		"aadhaar", "identity", "documents",
	}

	lowerText := strings.ToLower(text)
	for _, word := range kycWords {
		if strings.Contains(lowerText, word) {
			return true
		}
	}
	return false
}

// HasBankNames checks if text contains bank names
func HasBankNames(text string) bool {
	bankNames := []string{
		"hdfc", "icici", "sbi", "axis", "kotak", "yes bank",
		"pnb", "bank of baroda", "canara", "union bank",
		"idbi", "indian bank", "rbi", "reserve bank",
	}

	lowerText := strings.ToLower(text)
	for _, bank := range bankNames {
		if strings.Contains(lowerText, bank) {
			return true
		}
	}
	return false
}

// CalculateSpecialCharRatio calculates the ratio of special characters
func CalculateSpecialCharRatio(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	specialCount := 0
	for _, char := range text {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			specialCount++
		}
	}

	return float64(specialCount) / float64(len(text))
}

// CalculateCapitalRatio calculates the ratio of capital letters
func CalculateCapitalRatio(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	capitalCount := 0
	letterCount := 0
	for _, char := range text {
		if unicode.IsLetter(char) {
			letterCount++
			if unicode.IsUpper(char) {
				capitalCount++
			}
		}
	}

	if letterCount == 0 {
		return 0
	}

	return float64(capitalCount) / float64(letterCount)
}

// CalculateNumberRatio calculates the ratio of numbers
func CalculateNumberRatio(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	numberCount := 0
	for _, char := range text {
		if unicode.IsNumber(char) {
			numberCount++
		}
	}

	return float64(numberCount) / float64(len(text))
}

// SanitizeInput sanitizes user input
func SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	// Trim whitespace
	input = strings.TrimSpace(input)
	return input
}

// MaskPII masks personally identifiable information
func MaskPII(text string) string {
	// Mask phone numbers
	text = phoneRegex.ReplaceAllStringFunc(text, func(match string) string {
		if len(match) > 4 {
			return match[:2] + strings.Repeat("*", len(match)-4) + match[len(match)-2:]
		}
		return match
	})

	// Mask email addresses
	text = emailRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := strings.Split(match, "@")
		if len(parts) == 2 && len(parts[0]) > 2 {
			return parts[0][:2] + strings.Repeat("*", len(parts[0])-2) + "@" + parts[1]
		}
		return match
	})

	return text
}

