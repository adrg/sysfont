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

func getFontStyleScore(query, font string) float64 {
	sd := metrics.NewSorensenDice()
	sd.NgramSize = 4

	fontStyles := extractStyles(font)
	queryStyles := extractStyles(query)

	similarity := strutil.Similarity(queryStyles, fontStyles, sd)
	for _, sf := range []string{"mono", "sans", "serif"} {
		if strings.Contains(queryStyles, sf) && strings.Contains(fontStyles, sf) {
			similarity += 0.1
		}
	}
	if similarity == 0 {
		if len(fontStyles) == 0 || strings.Contains(fontStyles, "regular") {
			similarity += 0.1
		}
	}

	return similarity
}
