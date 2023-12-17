package sysfont

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/xdg"
)

// Finder is used to identify installed fonts. It can match fonts based on user
// queries and suggest alternative fonts if the requested fonts are not found.
type Finder struct {
	fonts []*Font
}

// FinderOpts contains options for configuring a font finder.
type FinderOpts struct {
	// Extensions controls which types of font files the finder reports.
	Extensions []string

	// SearchPaths is a list of paths to search for fonts.
	SearchPaths []string
}

// NewFinder returns a new font finder. If the opts parameter is nil, default
// options are used.
//
// Default options:
//   Extensions: []string{".ttf", ".ttc", ".otf"}
//   SearchPaths: xdg.FontDirs
//
// NOTE: See https://github.com/adrg/xdg#other-directories for more information
// about the default search paths.
func NewFinder(opts *FinderOpts) *Finder {
	if opts == nil {
		opts = &FinderOpts{Extensions: []string{".ttf", ".ttc", ".otf"}}
	}

	if len(opts.SearchPaths) == 0 {
		opts.SearchPaths = xdg.FontDirs
	}

	var fonts []*Font
	walker := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
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
			basename := filepath.Base(filename)
			family := cleanQuery(strings.TrimSuffix(basename, filepath.Ext(basename)))
			matches = append(matches, &Font{Family: family, Filename: filename})
		}

		fonts = append(fonts, matches...)
		return nil
	}

	// Traverse OS font directories.
	for _, dir := range opts.SearchPaths {
		if err := filepath.Walk(dir, walker); err != nil {
			continue
		}
	}

	return &Finder{
		fonts: fonts,
	}
}

// List returns the list of installed fonts. The finder attempts to identify
// the name and family of the returned fonts. If identification is not possible,
// only the filename field will be filled.
func (f *Finder) List() []*Font {
	fonts := make([]*Font, 0, len(f.fonts))
	for _, font := range f.fonts {
		fonts = append(fonts, font.clone())
	}

	return fonts
}

// Match attempts to identify the best matching installed font based on the
// specified query. If no close match is found, alternative fonts are searched.
// If no alternative font is found, a suitable default font is returned.
func (f *Finder) Match(query string) *Font {
	font := fontRegistry.matchFont(query, f.fonts)
	if font == nil {
		font = f.findAlternative(query)
	}

	return font.clone()
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
