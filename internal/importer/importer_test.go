package importer

import (
	"strings"
	"testing"
)

func TestImportReport(t *testing.T) {
	r, e := ReadJSON(strings.NewReader(`{"ID":"a","Kind":"x"}`))
	if e != nil || Validate(r) != nil {
		t.Fatal(e)
	}
}
