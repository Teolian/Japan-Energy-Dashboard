#!/bin/bash
# Generate estimated generation mix for all available dates

set -e

echo "Generating estimated generation mix data..."

# Find all demand files and generate corresponding generation files
for area in tokyo kansai; do
  echo ""
  echo "Processing area: $area"

  for demand_file in public/data/jp/$area/demand-*.json; do
    if [ -f "$demand_file" ]; then
      # Extract date from filename (demand-2025-11-03.json -> 2025-11-03)
      date=$(basename "$demand_file" | sed 's/demand-//' | sed 's/.json//')

      # Check if JEPX data exists
      jepx_file="public/data/jp/jepx/spot-$area-$date.json"
      if [ ! -f "$jepx_file" ]; then
        echo "⚠️  Skipping $date: JEPX data not found"
        continue
      fi

      # Generate generation mix
      echo "  Estimating generation for $date..."
      ./estimate-generation -area "$area" -date "$date" 2>&1 | grep -E "(Renewable|Carbon|Peak|Successfully)"
    fi
  done
done

echo ""
echo "✓ All generation mix files generated"
