package takeoffmetric

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed data/constants.csv
var constantsCSV string

// Constant is one row of the constants table. Value and ValueSI are the
// numbers of the row; ValueText keeps the value as the source prints it.
type Constant struct {
	ID             string    // stable row id, for example "CON-001"
	Category       string    // concrete, rebar, aggregates, earthwork, ...
	MaterialOrItem string    // the product or material the value is for
	Quantity       string    // what the value measures
	Value          float64   // the value in US customary units
	ValueText      string    // Value exactly as written in the table
	Unit           string    // unit of Value
	ValueSI        *float64  // SI conversion of Value; nil when the table leaves it empty
	UnitSI         string    // unit of ValueSI
	Condition      string    // the condition the value applies to (loose, bank, compacted, ...)
	Source         Source    // the publication the value was read in
	DateOpened     time.Time // date the source was opened and the value read
	UsedByToolURL  []string  // TakeoffMetric calculator pages that use the value
	EngineFile     []string  // site source files holding the value (not part of this module)
	EngineSymbol   []string  // symbols in those files, one per EngineFile entry
	Notes          string
}

// Source is the publication a constant was read in.
type Source struct {
	Publisher string
	Title     string
	URL       string
	Type      string // federal_public_domain, manufacturer_datasheet, trade_association, supplier or other
}

var (
	once      sync.Once
	constants []Constant
	byID      map[string]int
	loadErr   error
)

func load() {
	rows, err := csv.NewReader(strings.NewReader(constantsCSV)).ReadAll()
	if err != nil {
		loadErr = err
		return
	}
	byID = make(map[string]int, len(rows))
	for i, r := range rows[1:] {
		c, err := parse(r)
		if err != nil {
			loadErr = fmt.Errorf("data/constants.csv row %d: %w", i+2, err)
			return
		}
		byID[c.ID] = len(constants)
		constants = append(constants, c)
	}
}

func parse(r []string) (Constant, error) {
	if len(r) != 18 {
		return Constant{}, fmt.Errorf("%d fields, want 18", len(r))
	}
	v, err := strconv.ParseFloat(r[4], 64)
	if err != nil {
		return Constant{}, fmt.Errorf("value: %w", err)
	}
	var si *float64
	if r[6] != "" {
		x, err := strconv.ParseFloat(r[6], 64)
		if err != nil {
			return Constant{}, fmt.Errorf("value_si: %w", err)
		}
		si = &x
	}
	d, err := time.Parse("2006-01-02", r[13])
	if err != nil {
		return Constant{}, fmt.Errorf("date_opened: %w", err)
	}
	return Constant{
		ID: r[0], Category: r[1], MaterialOrItem: r[2], Quantity: r[3],
		Value: v, ValueText: r[4], Unit: r[5], ValueSI: si, UnitSI: r[7], Condition: r[8],
		Source:        Source{Publisher: r[9], Title: r[10], URL: r[11], Type: r[12]},
		DateOpened:    d,
		UsedByToolURL: split(r[14]), EngineFile: split(r[15]), EngineSymbol: split(r[16]),
		Notes: r[17],
	}, nil
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "; ")
}

func table() []Constant {
	once.Do(load)
	if loadErr != nil {
		panic("takeoffmetric: embedded data is invalid: " + loadErr.Error())
	}
	return constants
}

// All returns every constant, in table order. The slice is a copy.
func All() []Constant {
	return append([]Constant(nil), table()...)
}

// ByID returns the constant with the given id, for example "CON-001".
func ByID(id string) (Constant, bool) {
	t := table()
	i, ok := byID[id]
	if !ok {
		return Constant{}, false
	}
	return t[i], true
}

// Search returns the constants whose id, category, material, quantity,
// condition, publisher or notes contain every word of query, ignoring case.
// An empty query returns nothing.
func Search(query string) []Constant {
	words := strings.Fields(strings.ToLower(query))
	if len(words) == 0 {
		return nil
	}
	var out []Constant
	for _, c := range table() {
		text := strings.ToLower(strings.Join([]string{c.ID, c.Category, c.MaterialOrItem,
			c.Quantity, c.Condition, c.Source.Publisher, c.Notes}, " "))
		match := true
		for _, w := range words {
			if !strings.Contains(text, w) {
				match = false
				break
			}
		}
		if match {
			out = append(out, c)
		}
	}
	return out
}

// Sources returns each distinct source of the table once, sorted by
// publisher, then title, then URL.
func Sources() []Source {
	seen := map[Source]bool{}
	var out []Source
	for _, c := range table() {
		if !seen[c.Source] {
			seen[c.Source] = true
			out = append(out, c.Source)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Publisher != b.Publisher {
			return a.Publisher < b.Publisher
		}
		if a.Title != b.Title {
			return a.Title < b.Title
		}
		return a.URL < b.URL
	})
	return out
}

// CSV returns the embedded table as published (UTF-8, comma-separated),
// licensed CC BY 4.0.
func CSV() string {
	return constantsCSV
}
