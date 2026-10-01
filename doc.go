// Package takeoffmetric embeds the US construction estimating constants
// published by TakeoffMetric and exposes them as typed, read-only records. It
// has no dependencies and makes no network calls.
//
// The constants are the values the calculators on https://takeoffmetric.com/
// use to turn dimensions into order quantities: bag yields, aggregate and soil
// densities, swell and compaction ratios, rebar weights, dressed lumber sizes,
// block counts, shingle coverage, asphalt unit weights and air properties. Each
// record gives the value, its unit, an SI conversion, the condition it applies
// to and the source it was read in.
//
// The package is a lookup library and computes nothing. [All] returns every
// record, [ByID] returns one, [Search] finds records by words and [Sources]
// lists the publications behind them. The calculators themselves run on the
// site: the bag yields and unit weight feed the concrete calculator at
// https://takeoffmetric.com/concrete/concrete-calculator/, the loose and bank
// densities feed the gravel calculator at
// https://takeoffmetric.com/earthwork/gravel-calculator/, and the dressed
// lumber sizes feed the board foot calculator at
// https://takeoffmetric.com/framing/board-foot-calculator/.
//
// These are estimating values. Densities change with moisture, gradation and
// compaction, and bag yields and coverage hold for the named product line only.
// Confirm with your plans, your supplier and your local code before ordering or
// building. How the calculators use the values is described at
// https://takeoffmetric.com/methodology/.
//
// Data licence: the embedded table is licensed under CC BY 4.0, attribution
// TakeoffMetric; each value remains attributed to its own source. The Go code
// is MIT licensed.
package takeoffmetric
