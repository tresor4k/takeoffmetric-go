package takeoffmetric

import (
	"encoding/csv"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// source is the published table this module copies with sync_data.py.
var source = filepath.Join("..", "TakeoffMetrics-open-data", "us-construction-estimating-constants", "constants.csv")

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
		if r[18] == "true" {
			if !c.Withheld || c.Value != nil || c.ValueSI != nil || strings.Join(c.WithheldHosts, ";") != r[19] {
				t.Errorf("%s: withheld row exposes %+v", c.ID, c)
			}
			continue
		}
		if v, _ := strconv.ParseFloat(r[4], 64); c.Withheld || c.Value == nil || *c.Value != v {
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

func TestFederalValueKept(t *testing.T) {
	c, ok := ByID("CON-007")
	if !ok || c.Withheld || c.Source.Type != "federal_public_domain" || c.Value == nil || *c.Value != 140 || c.ValueText != "140" {
		t.Fatalf("ByID(CON-007) = %+v, %v", c, ok)
	}
	if c.ValueSI == nil || *c.ValueSI != 2243 || c.WithheldHosts != nil {
		t.Fatalf("CON-007 SI = %v, hosts %q", c.ValueSI, c.WithheldHosts)
	}
}

func TestWithheldHasNoValue(t *testing.T) {
	c, ok := ByID("CON-001")
	if !ok || !c.Withheld || c.Value != nil || c.ValueSI != nil || c.ValueText != "" || c.Unit != "ft3" || c.Category != "concrete" {
		t.Fatalf("ByID(CON-001) = %+v, %v", c, ok)
	}
	if len(c.WithheldHosts) != 1 || c.WithheldHosts[0] != "quikrete.com" {
		t.Fatalf("CON-001 WithheldHosts = %q", c.WithheldHosts)
	}
	n := 0
	for _, c := range All() {
		if c.Withheld {
			n++
			if c.Value != nil || c.ValueSI != nil || c.ValueText != "" {
				t.Errorf("%s is withheld but has a value", c.ID)
			}
		}
	}
	if n != 55 {
		t.Fatalf("%d withheld rows, want 55", n)
	}
}

func TestOldHeaderRefused(t *testing.T) {
	old := "id,category,material_or_item,quantity,value,unit,value_si,unit_si,condition,source_publisher," +
		"source_title,source_url,source_type,date_opened,used_by_tool_url,engine_file,engine_symbol,notes"
	if _, err := parse(strings.Split(old, ",")); err == nil {
		t.Fatal("parse accepted an 18-field 1.0.0 row")
	}
	py, err := exec.LookPath("python")
	if err != nil {
		t.Skip("python not found")
	}
	f := filepath.Join(t.TempDir(), "constants.csv")
	if err := os.WriteFile(f, []byte(old+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(py, "sync_data.py", f).CombinedOutput(); err == nil {
		t.Fatalf("sync_data.py accepted the 1.0.0 header: %s", out)
	}
}

func TestByID(t *testing.T) {
	c, ok := ByID("CON-001")
	if !ok {
		t.Fatal("CON-001 not found")
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
	if c, ok := ByID("HVC-002"); !ok || c.Withheld || c.Value == nil || c.ValueSI != nil {
		t.Fatalf("HVC-002 has a value and an empty value_si in the table, got %+v", c)
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
		if c.Withheld && c.Value != nil {
			t.Fatalf("Search returned a value for withheld row %s", c.ID)
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
