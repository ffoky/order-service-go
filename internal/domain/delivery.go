package domain

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func validateDelivery(delivery *Delivery) error {
	if strings.TrimSpace(delivery.Name) == "" {
		return ErrInvalidDeliveryName
	}

	if strings.TrimSpace(delivery.Phone) != "" && !isValidPhone(delivery.Phone) {
		return ErrInvalidPhone
	}

	if strings.TrimSpace(delivery.Email) != "" && !isValidEmail(delivery.Email) {
		return ErrInvalidEmail
	}

	if strings.TrimSpace(delivery.Zip) == "" {
		return ErrInvalidZip
	}

	if strings.TrimSpace(delivery.City) == "" {
		return ErrInvalidCity
	}

	if strings.TrimSpace(delivery.Address) == "" {
		return ErrInvalidAddress
	}

	return nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

func isValidPhone(phone string) bool {
	cleanPhone := strings.ReplaceAll(strings.ReplaceAll(phone, " ", ""), "-", "")
	return phoneRegex.MatchString(cleanPhone)
}
