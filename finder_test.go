package sysfont

import "testing"

func TestFinder(t *testing.T) {
	f := NewFinder(&FinderOpts{SearchPaths: []string{"testdata/fonts"}})
	fonts := f.List()

	// We expect Arial twice as we get double matches due to the lowercase and capitalised versions.
	exp := []*Font{
		{Name: "Arial", Family: "Arial", Filename: "testdata/fonts/arial.ttf"},
		{Name: "Arial", Family: "Arial", Filename: "testdata/fonts/arial.ttf"},
		{Name: "Dingbats", Family: "Dingbats", Filename: "testdata/fonts/dingbats.ttc"},
		{Name: "Helvetica", Family: "Helvetica", Filename: "testdata/fonts/helvetica.otf"},
	}

	if len(fonts) != len(exp) {
		t.Errorf("Expected %d fonts, got %d", len(exp), len(fonts))
	}

	for i, font := range fonts {
		if font.Name != exp[i].Name {
			t.Errorf("Expected font name %s, got %s", exp[i].Name, font.Name)
		}
		if font.Family != exp[i].Family {
			t.Errorf("Expected font family %s, got %s", exp[i].Family, font.Family)
		}
		if font.Filename != exp[i].Filename {
			t.Errorf("Expected font filename %s, got %s", exp[i].Filename, font.Filename)
		}
	}
}
