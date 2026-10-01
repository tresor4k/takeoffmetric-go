package takeoffmetric

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// source is the published table this module copies with sync_data.py.
var source = filepath.Join("..", "TakeoffMetrics", "open-data", "construction-constants", "constants.csv")

func TestMatchesThePublishedTable(t *testing.T) {
	f, err := os.Open(source)
	if err != nil {
		t.Skipf("%s not found (module outside its source tree)", source)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	rows = rows[1:]
	all := All()
	if len(all) != len(rows) {
		t.Fatalf("All() has %d constants, %s has %d rows (run python sync_data.py)", len(all), source, len(rows))
	}
	for i, r := range rows {
		c := all[i]
		if c.ID != r[0] || c.ValueText != r[4] || c.Unit != r[5] || c.UnitSI != r[7] || c.Source.URL != r[11] {
			t.Errorf("row %d differs from the published table: %+v", i+2, c)
		}
		if v, _ := strconv.ParseFloat(r[4], 64); c.Value != v {
			t.Errorf("%s: Value %v, published %s", c.ID, c.Value, r[4])
		}
	}
}

func TestEmbeddedRowCount(t *testing.T) {
	rows, err := csv.NewReader(strings.NewReader(CSV())).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if got := len(All()); got != len(rows)-1 || got == 0 {
		t.Fatalf("All() = %d constants, CSV has %d data rows", got, len(rows)-1)
	}
}

func TestByID(t *testing.T) {
	c, ok := ByID("CON-001")
	if !ok || c.Value != 0.30 || c.ValueText != "0.30" || c.Unit != "ft3" || c.Category != "concrete" {
		t.Fatalf("ByID(CON-001) = %+v, %v", c, ok)
	}
	if c.ValueSI == nil || *c.ValueSI != 0.008495 || c.UnitSI != "m3" {
		t.Fatalf("CON-001 SI = %v %s", c.ValueSI, c.UnitSI)
	}
	if c.Source.Publisher != "QUIKRETE" || c.DateOpened.Format("2006-01-02") != "2026-09-19" {
		t.Fatalf("CON-001 source = %+v, opened %v", c.Source, c.DateOpened)
	}
	if len(c.UsedByToolURL) < 2 || c.UsedByToolURL[0] != "https://takeoffmetric.com/concrete/concrete-calculator/" {
		t.Fatalf("CON-001 UsedByToolURL = %q", c.UsedByToolURL)
	}
	if _, ok := ByID("NOPE-000"); ok {
		t.Fatal("ByID found an unknown id")
	}
}

func TestEmptySIStaysNil(t *testing.T) {
	if c, ok := ByID("HVC-002"); !ok || c.ValueSI != nil {
		t.Fatalf("HVC-002 has an empty value_si in the table, got %v", *c.ValueSI)
	}
}

func TestSearch(t *testing.T) {
	got := Search("QUIKRETE 40 lb")
	if len(got) == 0 {
		t.Fatal("Search found no QUIKRETE 40 lb row")
	}
	for _, c := range got {
		if !strings.Contains(strings.ToLower(c.MaterialOrItem+c.Source.Publisher), "quikrete") {
			t.Fatalf("unexpected match %s", c.ID)
		}
	}
	if Search("") != nil || Search("zzzz-no-such-word") != nil {
		t.Fatal("empty or unknown query returned results")
	}
}

func TestSources(t *testing.T) {
	src := Sources()
	if len(src) == 0 {
		t.Fatal("no sources")
	}
	seen := map[Source]bool{}
	for _, s := range src {
		if seen[s] || !strings.HasPrefix(s.URL, "http") {
			t.Fatalf("duplicate or bad source %+v", s)
		}
		seen[s] = true
	}
	for _, c := range All() {
		if !seen[c.Source] {
			t.Fatalf("%s source missing from Sources()", c.ID)
		}
	}
}

func TestAllIsACopy(t *testing.T) {
	a := All()
	a[0].ID = "changed"
	if All()[0].ID == "changed" {
		t.Fatal("All() exposes the internal table")
	}
}
