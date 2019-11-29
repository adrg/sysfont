package sysfont

import (
	"strings"
	"unicode"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func cleanQuery(query string) string {
	return strings.Join(strings.FieldsFunc(strings.ToLower(query), func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}), " ")
}

func extractFamily(query string) string {
	family := cleanQuery(query)
	for _, fontStyle := range fontStyles {
		family = strings.Replace(family, fontStyle, "", -1)
	}

	return cleanQuery(family)
}

func extractStyles(query string) string {
	query = cleanQuery(query)

	var matched []string
	for _, fontStyle := range fontStyles {
		if strings.Contains(query, fontStyle) {
			matched = append(matched, fontStyle)
		}
	}

	return strings.Join(matched, " ")
}

func getFamilyScore(query, family string) float64 {
	return strutil.Similarity(query, cleanQuery(family), metrics.NewJaroWinkler())
}

func getFontScore(query, queryFamily string, font *Font) float64 {
	query, queryFamily = strings.ToLower(query), strings.ToLower(queryFamily)
	name, family := strings.ToLower(font.Name), strings.ToLower(font.Family)

	jw := metrics.NewJaroWinkler()
	return (strutil.Similarity(query, name, jw) +
		strutil.Similarity(queryFamily, family, jw)) / 2
}

func getFontStyleScore(query, font string) float64 {
	sd := metrics.NewSorensenDice()
	sd.NgramSize = 4

	fontStyles := extractStyles(font)
	similarity := strutil.Similarity(extractStyles(query), fontStyles, sd)
	if similarity == 0 {
		if len(fontStyles) == 0 || strings.Contains(fontStyles, "regular") {
			return 0.1
		}
	}

	return similarity
}

func in(q string, terms ...string) bool {
	for _, term := range terms {
		if q == term {
			return true
		}
	}

	return false
}
