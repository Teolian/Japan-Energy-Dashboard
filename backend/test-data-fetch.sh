#!/bin/bash
# Test script for validating data fetching
# Usage: ./test-data-fetch.sh [YYYY-MM-DD]

set -e

DATE=${1:-$(date +%Y-%m-%d)}
echo "🧪 Testing data fetch for $DATE"
echo "================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Build binaries
echo "🔨 Building fetch binaries..."
go build -o fetch-demand cmd/fetch-demand-http/main.go
go build -o fetch-jepx cmd/fetch-jepx-http/main.go
go build -o fetch-reserve cmd/fetch-reserve-http/main.go
echo ""

# Create output directories
mkdir -p public/data/jp/tokyo
mkdir -p public/data/jp/kansai
mkdir -p public/data/jp/jepx
mkdir -p public/data/jp/system

# Test each source
ERRORS=0

# 1. Tokyo Demand
echo "📊 Testing Tokyo demand..."
if ./fetch-demand -area tokyo -date $DATE --use-http 2>&1 | tee /tmp/tokyo-fetch.log; then
  FILE="public/data/jp/tokyo/demand-$DATE.json"
  if [ -f "$FILE" ]; then
    SIZE=$(stat -f%z "$FILE" 2>/dev/null || stat -c%s "$FILE" 2>/dev/null)
    if [ $SIZE -gt 100 ]; then
      if jq empty "$FILE" 2>/dev/null; then
        POINTS=$(jq -r '.series | length' "$FILE" 2>/dev/null)
        echo -e "${GREEN}✓ Tokyo demand: $POINTS data points ($SIZE bytes)${NC}"

        # Check for expected values
        MIN_DEMAND=$(jq -r '[.series[].demand_mw] | min' "$FILE")
        MAX_DEMAND=$(jq -r '[.series[].demand_mw] | max' "$FILE")
        echo "  └─ Demand range: $MIN_DEMAND - $MAX_DEMAND MW"

        if [ "$POINTS" -ne 24 ]; then
          echo -e "${YELLOW}  ⚠️  Expected 24 points, got $POINTS${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      else
        echo -e "${RED}✗ Tokyo demand: Invalid JSON${NC}"
        ERRORS=$((ERRORS + 1))
      fi
    else
      echo -e "${RED}✗ Tokyo demand: File too small ($SIZE bytes)${NC}"
      ERRORS=$((ERRORS + 1))
    fi
  else
    echo -e "${RED}✗ Tokyo demand: File not created${NC}"
    ERRORS=$((ERRORS + 1))
  fi
else
  echo -e "${YELLOW}⚠️  Tokyo demand: Fetch failed (may not be available yet)${NC}"
fi
echo ""

# 2. Kansai Demand
echo "📊 Testing Kansai demand..."
if ./fetch-demand -area kansai -date $DATE --use-http 2>&1 | tee /tmp/kansai-fetch.log; then
  FILE="public/data/jp/kansai/demand-$DATE.json"
  if [ -f "$FILE" ]; then
    SIZE=$(stat -f%z "$FILE" 2>/dev/null || stat -c%s "$FILE" 2>/dev/null)
    if [ $SIZE -gt 100 ]; then
      if jq empty "$FILE" 2>/dev/null; then
        POINTS=$(jq -r '.series | length' "$FILE" 2>/dev/null)
        echo -e "${GREEN}✓ Kansai demand: $POINTS data points ($SIZE bytes)${NC}"

        MIN_DEMAND=$(jq -r '[.series[].demand_mw] | min' "$FILE")
        MAX_DEMAND=$(jq -r '[.series[].demand_mw] | max' "$FILE")
        echo "  └─ Demand range: $MIN_DEMAND - $MAX_DEMAND MW"

        if [ "$POINTS" -ne 24 ]; then
          echo -e "${YELLOW}  ⚠️  Expected 24 points, got $POINTS${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      else
        echo -e "${RED}✗ Kansai demand: Invalid JSON${NC}"
        ERRORS=$((ERRORS + 1))
      fi
    else
      echo -e "${RED}✗ Kansai demand: File too small ($SIZE bytes)${NC}"
      ERRORS=$((ERRORS + 1))
    fi
  else
    echo -e "${RED}✗ Kansai demand: File not created${NC}"
    ERRORS=$((ERRORS + 1))
  fi
else
  echo -e "${YELLOW}⚠️  Kansai demand: Fetch failed (may not be available yet)${NC}"
fi
echo ""

# 3. JEPX Tokyo
echo "💴 Testing JEPX Tokyo spot prices..."
if ./fetch-jepx -area tokyo -date $DATE --use-http 2>&1 | tee /tmp/jepx-tokyo-fetch.log; then
  FILE="public/data/jp/jepx/spot-tokyo-$DATE.json"
  if [ -f "$FILE" ]; then
    SIZE=$(stat -f%z "$FILE" 2>/dev/null || stat -c%s "$FILE" 2>/dev/null)
    if [ $SIZE -gt 100 ]; then
      if jq empty "$FILE" 2>/dev/null; then
        POINTS=$(jq -r '.price_yen_per_kwh | length' "$FILE" 2>/dev/null)
        echo -e "${GREEN}✓ JEPX Tokyo: $POINTS price points ($SIZE bytes)${NC}"

        MIN_PRICE=$(jq -r '[.price_yen_per_kwh[].price] | min' "$FILE")
        MAX_PRICE=$(jq -r '[.price_yen_per_kwh[].price] | max' "$FILE")
        AVG_PRICE=$(jq -r '[.price_yen_per_kwh[].price] | add / length' "$FILE")
        echo "  └─ Price range: $MIN_PRICE - $MAX_PRICE JPY/kWh (avg: $AVG_PRICE)"

        if [ "$POINTS" -ne 24 ]; then
          echo -e "${YELLOW}  ⚠️  Expected 24 points, got $POINTS${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      else
        echo -e "${RED}✗ JEPX Tokyo: Invalid JSON${NC}"
        ERRORS=$((ERRORS + 1))
      fi
    else
      echo -e "${RED}✗ JEPX Tokyo: File too small ($SIZE bytes)${NC}"
      ERRORS=$((ERRORS + 1))
    fi
  else
    echo -e "${RED}✗ JEPX Tokyo: File not created${NC}"
    ERRORS=$((ERRORS + 1))
  fi
else
  echo -e "${YELLOW}⚠️  JEPX Tokyo: Fetch failed (may not be available yet)${NC}"
fi
echo ""

# 4. JEPX Kansai
echo "💴 Testing JEPX Kansai spot prices..."
if ./fetch-jepx -area kansai -date $DATE --use-http 2>&1 | tee /tmp/jepx-kansai-fetch.log; then
  FILE="public/data/jp/jepx/spot-kansai-$DATE.json"
  if [ -f "$FILE" ]; then
    SIZE=$(stat -f%z "$FILE" 2>/dev/null || stat -c%s "$FILE" 2>/dev/null)
    if [ $SIZE -gt 100 ]; then
      if jq empty "$FILE" 2>/dev/null; then
        POINTS=$(jq -r '.price_yen_per_kwh | length' "$FILE" 2>/dev/null)
        echo -e "${GREEN}✓ JEPX Kansai: $POINTS price points ($SIZE bytes)${NC}"

        MIN_PRICE=$(jq -r '[.price_yen_per_kwh[].price] | min' "$FILE")
        MAX_PRICE=$(jq -r '[.price_yen_per_kwh[].price] | max' "$FILE")
        AVG_PRICE=$(jq -r '[.price_yen_per_kwh[].price] | add / length' "$FILE")
        echo "  └─ Price range: $MIN_PRICE - $MAX_PRICE JPY/kWh (avg: $AVG_PRICE)"

        if [ "$POINTS" -ne 24 ]; then
          echo -e "${YELLOW}  ⚠️  Expected 24 points, got $POINTS${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      else
        echo -e "${RED}✗ JEPX Kansai: Invalid JSON${NC}"
        ERRORS=$((ERRORS + 1))
      fi
    else
      echo -e "${RED}✗ JEPX Kansai: File too small ($SIZE bytes)${NC}"
      ERRORS=$((ERRORS + 1))
    fi
  else
    echo -e "${RED}✗ JEPX Kansai: File not created${NC}"
    ERRORS=$((ERRORS + 1))
  fi
else
  echo -e "${YELLOW}⚠️  JEPX Kansai: Fetch failed (may not be available yet)${NC}"
fi
echo ""

# 5. Reserve Data
echo "🔋 Testing reserve capacity..."
if ./fetch-reserve -date $DATE --use-http 2>&1 | tee /tmp/reserve-fetch.log; then
  FILE="public/data/jp/system/reserve-$DATE.json"
  if [ -f "$FILE" ]; then
    SIZE=$(stat -f%z "$FILE" 2>/dev/null || stat -c%s "$FILE" 2>/dev/null)
    if [ $SIZE -gt 100 ]; then
      if jq empty "$FILE" 2>/dev/null; then
        REGIONS=$(jq -r '.reserves | length' "$FILE" 2>/dev/null)
        echo -e "${GREEN}✓ Reserve data: $REGIONS regions ($SIZE bytes)${NC}"

        # List regions
        jq -r '.reserves[] | "  └─ \(.region): \(.reserve_percent)%"' "$FILE"

        if [ "$REGIONS" -ne 10 ]; then
          echo -e "${YELLOW}  ⚠️  Expected 10 regions, got $REGIONS${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      else
        echo -e "${RED}✗ Reserve data: Invalid JSON${NC}"
        ERRORS=$((ERRORS + 1))
      fi
    else
      echo -e "${RED}✗ Reserve data: File too small ($SIZE bytes)${NC}"
      ERRORS=$((ERRORS + 1))
    fi
  else
    echo -e "${RED}✗ Reserve data: File not created${NC}"
    ERRORS=$((ERRORS + 1))
  fi
else
  echo -e "${YELLOW}⚠️  Reserve data: Fetch failed (may not be available yet)${NC}"
fi
echo ""

# Summary
echo "================================"
if [ $ERRORS -eq 0 ]; then
  echo -e "${GREEN}✓ All tests passed!${NC}"
  exit 0
else
  echo -e "${RED}✗ Found $ERRORS error(s)${NC}"
  echo ""
  echo "Check logs in /tmp/*-fetch.log for details"
  exit 1
fi
