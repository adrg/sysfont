package sysfont_test

import (
	"fmt"

	"github.com/adrg/sysfont"
)

func ExampleNewFinder() {
	// Create a new font finder using the default options.
	finder := sysfont.NewFinder(nil)

	// Create a new finder which only searches for TTF files.
	finder = sysfont.NewFinder(&sysfont.FinderOpts{
		Extensions: []string{".ttf"},
	})

	// Create a new finder that searches for fonts only in the current directory.
	finder = sysfont.NewFinder(&sysfont.FinderOpts{
		SearchPaths: []string{"."},
	})

	// List detected fonts.
	for _, font := range finder.List() {
		fmt.Println(font.Family, font.Name, font.Filename)
	}
}

func ExampleFinder_List() {
	finder := sysfont.NewFinder(nil)

	for _, font := range finder.List() {
		fmt.Println(font.Family, font.Name, font.Filename)
	}
}

func ExampleFinder_Match() {
	finder := sysfont.NewFinder(nil)

	terms := []string{
		"AmericanTypewriter",
		"AmericanTypewriter-Bold",
		"Andale",
		"Arial",
		"Arial Bold",
		"Arial-BoldItalicMT",
		"ArialMT",
		"Baskerville",
		"Candara",
		"Corbel",
		"Gill Sans",
		"Hoefler Text Bold",
		"Impact",
		"Palatino",
		"Symbol",
		"Tahoma",
		"Times",
		"Times Bold",
		"Times BoldItalic",
		"Times Italic Bold",
		"Times Roman",
		"Verdana",
		"Verdana-Italic",
		"Webdings",
		"ZapfDingbats",
	}

	for _, term := range terms {
		font := finder.Match(term)
		if font == nil {
			// Match should always return a font. However, it is safer to check.
			continue
		}

		fmt.Printf("%-30s -> %-30s (%s)\n", term, font.Name, font.Filename)
	}
}
