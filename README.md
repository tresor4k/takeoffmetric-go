# takeoffmetric-go

Dependency-free Go package that embeds the US construction estimating constants published by
[TakeoffMetric](https://takeoffmetric.com/) as typed, read-only records: bag yields, aggregate and soil densities,
swell and compaction ratios, rebar weights, dressed lumber sizes, block counts, shingle coverage, asphalt unit weights
and air properties. These are published constants with their sources, not original research: each record gives the
value, its unit, an SI conversion, the condition it applies to and the source it was read in.

```go
import "github.com/tresor4k/takeoffmetric-go"

c, ok := takeoffmetric.ByID("CON-001")   // QUIKRETE 40 lb bag yield
rows := takeoffmetric.Search("bank density")
all := takeoffmetric.All()
src := takeoffmetric.Sources()
```

The package is a lookup library and computes nothing. To turn dimensions into quantities, use the calculators on the
site, for example the [concrete calculator](https://takeoffmetric.com/concrete/concrete-calculator/). How they use
these values: [methodology](https://takeoffmetric.com/methodology/).

Estimating values only: confirm with your plans, your supplier and your local code before ordering or building.

## Updating the data

`python sync_data.py` copies `../TakeoffMetrics/open-data/construction-constants/constants.csv` (or the file given as
first argument) into `data/`. Then run `go vet ./... && go test ./...`.

## Licence

- Code: MIT, see LICENSE.
- Constants (`data/`): CC BY 4.0, see data/LICENSE. Attribution: TakeoffMetric. Each value remains attributed to its
  own source through the `source_publisher` and `source_url` fields.
- The calculation engine of takeoffmetric.com is not included in this module.
