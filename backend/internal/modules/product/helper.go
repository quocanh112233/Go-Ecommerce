package product

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// removeAccents removes Vietnamese diacritics
// "Áo thun" -> "Ao thun"
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: Mark, nonspacing
	}), norm.NFC)
	result, _, _ := transform.String(t, s)

	// Handle Đ/đ specifically
	result = strings.ReplaceAll(result, "đ", "d")
	result = strings.ReplaceAll(result, "Đ", "D")
	return result
}

// generatePrefix generates prefix from text
// "Áo thun trắng" -> "Ao thun trang" -> "ATT"
func generatePrefix(text string) string {
	text = removeAccents(text)
	words := strings.Fields(text)
	var prefix string
	for _, word := range words {
		runes := []rune(strings.ToUpper(word))
		if len(runes) > 0 {
			prefix += string(runes[0])
		}
	}
	return prefix
}

// generateSKU creates SKU from category, product names and variant ID
// Example: "Áo thun" + "Áo thun trắng" + 1 -> "ATATT1"
func generateSKU(categoryName, productName string, variantID uint) string {
	catPrefix := generatePrefix(categoryName)
	prodPrefix := generatePrefix(productName)
	return fmt.Sprintf("%s%s%d", catPrefix, prodPrefix, variantID)
}

// generateSlug creates URL-friendly slug from name
func generateSlug(name string) string {
	slug := removeAccents(name)
	slug = strings.ToLower(strings.TrimSpace(slug))
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric except hyphens
	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// calculateTotalStock sums stock from all variants
func calculateTotalStock(variants []ProductVariant) int {
	total := 0
	for _, v := range variants {
		total += v.Stock
	}
	return total
}
