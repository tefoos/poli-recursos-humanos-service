package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (c *Client) ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

func (c *Client) ValidateDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func (c *Client) ReadEmailInput(prompt string) string {
	for {
		email := c.ReadInput(prompt)
		if c.ValidateEmail(email) {
			return email
		}
		fmt.Printf("Email inválido. Use formato: usuario@dominio.com\n")
	}
}

func (c *Client) ReadDateInput(prompt string) string {
	for {
		date := c.ReadInput(prompt)
		if c.ValidateDate(date) {
			return date
		}
		fmt.Printf("Fecha inválida. Use formato: YYYY-MM-DD (ej: 1990-01-15)\n")
	}
}

func (c *Client) ReadPositiveFloatInput(prompt string) float64 {
	for {
		input := c.ReadInput(prompt)
		val, err := strconv.ParseFloat(input, 64)
		if err != nil || val <= 0 {
			fmt.Printf("Debe ser un número mayor a 0\n")
			continue
		}
		return val
	}
}

func (c *Client) ReadRangeFloatInput(prompt string, min, max float64) float64 {
	for {
		input := c.ReadInput(prompt)
		val, err := strconv.ParseFloat(input, 64)
		if err != nil || val < min || val > max {
			fmt.Printf("Debe ser un número entre %.2f y %.2f\n", min, max)
			continue
		}
		return val
	}
}

func (c *Client) ReadValidatedName(prompt string, required bool) *string {
	for {
		input := c.ReadInput(prompt)
		trimmed := strings.TrimSpace(input)

		if !required && trimmed == "" {
			return nil
		}

		if required && trimmed == "" {
			fmt.Printf("Este campo es requerido\n")
			continue
		}

		if len(trimmed) > 50 {
			fmt.Printf("Máximo 50 caracteres permitidos\n")
			continue
		}

		return &trimmed
	}
}
