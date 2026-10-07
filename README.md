# takeoffmetric-go

Dependency-free Go package that embeds version 1.1.0 of the US construction estimating constants published by
[TakeoffMetric](https://takeoffmetric.com/) as typed, read-only records: bag yields, aggregate and soil densities,
swell and compaction ratios, rebar weights, dressed lumber sizes, block counts, shingle coverage, asphalt unit weights
and air properties. These are published constants with their sources, not original research: each record gives the
value, its unit, an SI conversion, the condition it applies to and the source it was read in.

55 of the 148 records carry no value: their publishers do not allow republication. Those records keep the quantity, its
unit and the link to the source; `Withheld` is true, `Value` and `ValueSI` are nil, and `WithheldHosts` lists the domains
concerned. Read the value at `Source.URL`.

```go
import "github.com/tresor4k/takeoffmetric-go"

c, ok := takeoffmetric.ByID("CON-007")   // c.Value is a *float64, nil when c.Withheld
rows := takeoffmetric.Search("bank density")
all := takeoffmetric.All()
src := takeoffmetric.Sources()
```

The package is a lookup library and computes nothing. To turn dimensions into quantities, use the calculators on the
site, for example the [concrete calculator](https://takeoffmetric.com/concrete/concrete-calculator/). How they use
these values: [methodology](https://takeoffmetric.com/methodology/).

Estimating values only: confirm with your plans, your supplier and your local code before ordering or building.

## Updating the data

`python sync_data.py` copies `../TakeoffMetrics-open-data/us-construction-estimating-constants/constants.csv` (or the
file given as first argument) into `data/`. It accepts only the 1.1.0 header (with `values_withheld` and
`withheld_hosts`). Then run `go vet ./... && go test ./...`.

## Licence

- Code: MIT, see LICENSE.
- Constants (`data/`), see data/LICENSE:
  - the compilation (selection of rows, SI conversions, conditions, notes, column structure): CC BY 4.0, attribution
    TakeoffMetric Editorial;
  - values from federal publications: works of the US government, in the public domain;
  - any other value stays under the terms of its publisher and is not relicensed by CC BY 4.0; each row names its
    publisher and links its source (`source_publisher`, `source_title`, `source_url`);
  - 55 rows carry no value (`values_withheld` = `true`): their publishers do not allow republication.
- Cite as: TakeoffMetric Editorial, *US Construction Estimating Constants*, version 1.1.0, 2026, takeoffmetric.com.
- The calculation engine of takeoffmetric.com is not included in this module.
