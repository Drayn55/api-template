package slugutil

import (
	"regexp"
	"strings"
)

func Generate(input string) string {
	if input == "" {
		return ""
	}

	numberMap := map[string]string{
		"0": "nol",
		"1": "satu",
		"2": "dua",
		"3": "tiga",
		"4": "empat",
		"5": "lima",
		"6": "enam",
		"7": "tujuh",
		"8": "delapan",
		"9": "sembilan",
	}

	for k, v := range numberMap {
		input = strings.ReplaceAll(input, k, " "+v+" ")
	}

	slug := strings.ToLower(strings.TrimSpace(input))
	slug = strings.ReplaceAll(slug, ">", "more-than")
	slug = strings.ReplaceAll(slug, "<", "less-than")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, "\\", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, ".", "-")

	re := regexp.MustCompile(`[^a-z0-9-]`)
	slug = re.ReplaceAllString(slug, "")

	reDash := regexp.MustCompile(`-+`)
	slug = reDash.ReplaceAllString(slug, "-")

	return strings.Trim(slug, "-")
}
