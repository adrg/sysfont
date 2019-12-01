package sysfont

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/xdg"
)

type Finder struct {
	fonts []*Font
}

type FinderOpts struct {
	Extensions []string
}

func NewFinder(opts *FinderOpts) *Finder {
	if opts == nil {
		opts = &FinderOpts{Extensions: []string{".ttf", ".ttc", ".otf"}}
	}

	var fonts []*Font
	walker := func(filename string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Check file extension.
		if extensions := opts.Extensions; len(extensions) > 0 {
			extension := filepath.Ext(strings.ToLower(filename))
			if !strutil.SliceContains(extensions, extension) {
				return nil
			}
		}

		// Attempt to identify fonts by filename.
		matches := fontRegistry.matchFontsByFilename(filename)
		if len(matches) == 0 {
			matches = []*Font{&Font{Filename: filename}}
		}

		fonts = append(fonts, matches...)
		return nil
	}

	// Traverse OS font directories.
	for _, dir := range xdg.FontDirs {
		filepath.Walk(dir, walker)
	}

	return &Finder{
		fonts: fonts,
	}
}

func (f *Finder) List() []*Font {
	fonts := make([]*Font, len(f.fonts))
	for i, font := range f.fonts {
		fonts[i] = &(*font)
	}

	return fonts
}

func (f *Finder) Match(query string) *Font {
	font := fontRegistry.matchFont(query, f.fonts)
	if font == nil {
		font = f.findAlternative(query)
	}
	if font == nil {
		return nil
	}

	return &(*font)
}

func (f *Finder) findAlternative(query string) *Font {
	// Identify font family.
	family, _ := fontRegistry.matchFamily(query)

	// Identify alternate fonts based on the matched family.
	alternatives := fontRegistry.getAlternatives(family, f.fonts)

	// Identify best alternative.
	var maxScore float64
	var maxScoreFont *Font

	for _, font := range alternatives {
		if score := getFontStyleScore(query, font.Name); score > maxScore {
			maxScore = score
			maxScoreFont = font
		}
	}

	return maxScoreFont
}
