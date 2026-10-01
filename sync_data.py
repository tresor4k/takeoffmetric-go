"""Copy the published TakeoffMetric constants CSV into data/ (embedded by the package).

Usage: python sync_data.py [source.csv]
Default source: ../TakeoffMetrics/open-data/construction-constants/constants.csv
Then run: go vet ./... && go test ./...
"""
import csv
import pathlib
import shutil
import sys

HERE = pathlib.Path(__file__).resolve().parent
DEFAULT = HERE.parent / "TakeoffMetrics" / "open-data" / "construction-constants" / "constants.csv"
HEADER = ("id,category,material_or_item,quantity,value,unit,value_si,unit_si,condition,source_publisher,"
          "source_title,source_url,source_type,date_opened,used_by_tool_url,engine_file,engine_symbol,notes")

src = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else DEFAULT
with open(src, encoding="utf-8", newline="") as f:
    rows = list(csv.reader(f))
if ",".join(rows[0]) != HEADER:
    sys.exit(f"unexpected header in {src}: {rows[0]}")
dst = HERE / "data" / "constants.csv"
shutil.copyfile(src, dst)
print(f"synced {len(rows) - 1} rows from {src} to {dst}")
