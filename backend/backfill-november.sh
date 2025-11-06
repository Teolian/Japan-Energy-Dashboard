#!/bin/bash
# Backfill missing November data locally
# Run from backend directory

set -e

DATES=("2025-11-01" "2025-11-02" "2025-11-04")

echo "🔄 Backfilling missing November data..."
echo ""

for DATE in "${DATES[@]}"; do
  echo "📅 Fetching data for $DATE"

  # Tokyo demand
  ./fetch-demand -area tokyo -date $DATE --use-http \
    -output ../frontend/public/data/jp/tokyo/demand-$DATE.json || echo "⚠️ Tokyo failed"

  # JEPX Tokyo
  ./fetch-jepx -area tokyo -date $DATE --use-http \
    -output ../frontend/public/data/jp/jepx/spot-tokyo-$DATE.json || echo "⚠️ JEPX Tokyo failed"

  # JEPX Kansai
  ./fetch-jepx -area kansai -date $DATE --use-http \
    -output ../frontend/public/data/jp/jepx/spot-kansai-$DATE.json || echo "⚠️ JEPX Kansai failed"

  # Reserve
  ./fetch-reserve -date $DATE --use-http \
    -output ../frontend/public/data/jp/system/reserve-$DATE.json || echo "⚠️ Reserve failed"

  echo "✅ $DATE complete"
  echo ""
done

echo "📊 Backfill complete. Files created:"
ls -lh ../frontend/public/data/jp/tokyo/demand-2025-11-0*.json 2>/dev/null || true
ls -lh ../frontend/public/data/jp/jepx/spot-*-2025-11-0*.json 2>/dev/null || true
ls -lh ../frontend/public/data/jp/system/reserve-2025-11-0*.json 2>/dev/null || true
