package domain

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

// Delivery представляет информацию о доставке
// @Description Данные о доставке заказа
type Delivery struct {
	Name    string `json:"name" example:"Test Testov"`
	Phone   string `json:"phone" example:"+9720000000"`
	Zip     string `json:"zip" example:"2639809"`
	City    string `json:"city" example:"Kiryat Mozkin"`
	Address string `json:"address" example:"Ploshad Mira 15"`
	Region  string `json:"region" example:"Kraiot"`
	Email   string `json:"email" example:"test@gmail.com"`
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
