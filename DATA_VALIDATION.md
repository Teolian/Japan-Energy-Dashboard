# Data Validation Guide

Руководство по проверке правильности скачивания, парсинга и сохранения данных для каждого источника.

## 🎯 Цель

Перед автоматизацией нужно убедиться, что:
1. **Скачивание**: Данные корректно загружаются с источника
2. **Парсинг**: CSV правильно преобразуется в JSON
3. **Сохранение**: Данные соответствуют реальности

---

## 📊 1. TEPCO (Tokyo Demand)

### Источник данных
- **URL**: https://www.tepco.co.jp/forecast/html/download-j.html
- **Формат**: ZIP архив с CSV файлами (Shift-JIS кодировка)
- **Структура**: `YYYYMM_power_usage.zip` → `YYYYMMDD_power_usage.csv`
- **Единицы**: 万kW (10,000 kW = 10 MW)

### Тестовая команда

```bash
cd backend

# Fetch данных за сегодня
./fetch-demand -area tokyo -date $(date +%Y-%m-%d) --use-http

# Проверить результат
cat public/data/jp/tokyo/demand-$(date +%Y-%m-%d).json | jq
```

### Что проверять

**1. Структура JSON:**
```json
{
  "area": "tokyo",
  "date": "2025-11-05",
  "timezone": "Asia/Tokyo",
  "timescale": "hourly",
  "series": [
    {
      "ts": "2025-11-05T00:00:00+09:00",
      "demand_mw": 28500,
      "forecast_mw": 29000
    }
    // ... 24 точки
  ],
  "source": {
    "name": "TEPCO",
    "url": "..."
  }
}
```

**2. Проверка данных:**
- [ ] **24 точки** (одна на каждый час)
- [ ] Timestamps идут с 00:00 до 23:00
- [ ] Timezone: `+09:00` (JST)
- [ ] Demand значения разумные:
  - Ночь (00:00-06:00): ~25,000-30,000 MW
  - День (10:00-18:00): ~35,000-45,000 MW
  - Пик: ~13:00-14:00
- [ ] Forecast присутствует (может быть null для прошлых дат)

**3. Кодировка:**
```bash
# Проверить что Shift-JIS корректно преобразован
grep -a "DATE\|TIME\|実績\|予測" backend/internal/adapters/testdata/tepco-sample.csv
```

**4. Единицы измерения:**
```bash
# В CSV: 万kW (10^4 kW)
# В JSON: MW (должны быть умножены на 10)
# Пример: CSV=2850.0 → JSON=28500 MW
```

### Известные проблемы

❌ **TODO**: Проверить конвертацию единиц измерения
❌ **TODO**: Убедиться что парсер обрабатывает разные форматы дат
❌ **TODO**: Проверить корректность распаковки ZIP

### Ручная проверка

1. Открыть: https://www.tepco.co.jp/forecast/html/download-j.html
2. Скачать последний ZIP вручную
3. Распаковать и посмотреть CSV
4. Сравнить с результатом fetch-demand

---

## 📊 2. Kansai Electric (Kansai Demand)

### Источник данных
- **URL**: https://www.kansai-td.co.jp/denkiyoho/download.html
- **Формат**: CSV (предположительно UTF-8 или Shift-JIS)
- **Единицы**: kW или MW (нужно проверить)

### Тестовая команда

```bash
cd backend

# Fetch данных за сегодня
./fetch-demand -area kansai -date $(date +%Y-%m-%d) --use-http

# Проверить результат
cat public/data/jp/kansai/demand-$(date +%Y-%m-%d).json | jq
```

### Что проверять

**1. Структура аналогична TEPCO**

**2. Проверка данных:**
- [ ] 24 точки (hourly)
- [ ] Demand значения разумные:
  - Ночь: ~12,000-15,000 MW (меньше чем Tokyo)
  - День: ~17,000-22,000 MW
  - Пик: ~13:00-14:00

**3. Кодировка:**
```bash
# Проверить кодировку CSV
file backend/internal/adapters/testdata/kansai-sample.csv
```

### Известные проблемы

❌ **TODO**: Проверить формат CSV (разделители, заголовки)
❌ **TODO**: Единицы измерения kW vs MW
❌ **TODO**: Наличие forecast данных

### Ручная проверка

1. Открыть: https://www.kansai-td.co.jp/denkiyoho/
2. Скачать CSV
3. Сравнить с результатом

---

## 💴 3. JEPX Spot Prices (Tokyo/Kansai)

### Источник данных
- **URL**: https://www.jepx.jp/
- **Формат**: CSV
- **Единицы**: JPY/kWh (日本円/キロワット時)

### Тестовая команда

```bash
cd backend

# Fetch Tokyo spot prices
./fetch-jepx -area tokyo -date $(date +%Y-%m-%d) --use-http
cat public/data/jp/jepx/spot-tokyo-$(date +%Y-%m-%d).json | jq

# Fetch Kansai spot prices
./fetch-jepx -area kansai -date $(date +%Y-%m-%d) --use-http
cat public/data/jp/jepx/spot-kansai-$(date +%Y-%m-%d).json | jq
```

### Что проверять

**1. Структура JSON:**
```json
{
  "date": "2025-11-05",
  "area": "tokyo",
  "timescale": "hourly",
  "price_yen_per_kwh": [
    {
      "ts": "2025-11-05T00:00:00+09:00",
      "price": 24.5
    }
    // ... 24 точки
  ],
  "source": {
    "name": "JEPX",
    "url": "https://www.jepx.jp/"
  },
  "meta": {
    "min_price": 18.2,
    "max_price": 42.8,
    "avg_price": 28.4
  }
}
```

**2. Проверка данных:**
- [ ] 24 точки (hourly)
- [ ] Цены разумные:
  - Ночь (00:00-06:00): ~15-25 JPY/kWh (низкий спрос)
  - День (10:00-20:00): ~25-45 JPY/kWh (высокий спрос)
  - Пик: обычно 18:00-20:00 (вечерний пик)
- [ ] Нет отрицательных цен (обычно)
- [ ] Meta корректные (min < avg < max)

**3. Сравнение Tokyo vs Kansai:**
```bash
# Tokyo обычно дороже на 5-10%
paste \
  <(jq -r '.price_yen_per_kwh[].price' public/data/jp/jepx/spot-tokyo-2025-11-05.json) \
  <(jq -r '.price_yen_per_kwh[].price' public/data/jp/jepx/spot-kansai-2025-11-05.json) \
  | awk '{print "Tokyo:", $1, "Kansai:", $2, "Diff:", $1-$2}'
```

### Известные проблемы

❌ **TODO**: Проверить формат даты в CSV
❌ **TODO**: Убедиться что цены в JPY/kWh (не JPY/MWh!)
❌ **TODO**: Проверить day-ahead vs spot цены

### Ручная проверка

1. Открыть: https://www.jepx.jp/market/index.html
2. Найти Spot Market → Day-Ahead
3. Сравнить цены по часам

---

## 🔋 4. OCCTO Reserve Capacity

### Источник данных
- **URL**: https://www.occto.or.jp/
- **Формат**: CSV
- **Охват**: 10 регионов Японии

### Тестовая команда

```bash
cd backend

# Fetch reserve data
./fetch-reserve -date $(date +%Y-%m-%d)
cat public/data/jp/system/reserve-$(date +%Y-%m-%d).json | jq
```

### Что проверять

**1. Структура JSON:**
```json
{
  "date": "2025-11-05",
  "reserves": [
    {
      "region": "hokkaido",
      "demand_mw": 3500,
      "capacity_mw": 5000,
      "reserve_percent": 42.86
    },
    {
      "region": "tokyo",
      "demand_mw": 35000,
      "capacity_mw": 42000,
      "reserve_percent": 20.0
    }
    // ... всего 10 регионов
  ],
  "source": {
    "name": "OCCTO",
    "url": "https://www.occto.or.jp/"
  }
}
```

**2. Проверка данных:**
- [ ] **10 регионов**: hokkaido, tohoku, tokyo, chubu, hokuriku, kansai, chugoku, shikoku, kyushu, okinawa
- [ ] Reserve percent = (capacity - demand) / capacity * 100
- [ ] Reserve обычно 5-30% (критично если <5%)
- [ ] Demand соответствует данным TEPCO/Kansai

**3. Расчёт reserve margin:**
```bash
# Проверить формулу
jq -r '.reserves[] | "\(.region): demand=\(.demand_mw) capacity=\(.capacity_mw) reserve=\(.reserve_percent)%"' \
  public/data/jp/system/reserve-2025-11-05.json
```

### Известные проблемы

❌ **TODO**: Проверить названия регионов (английские vs японские)
❌ **TODO**: Убедиться в правильности формулы reserve_percent
❌ **TODO**: Проверить что все 10 регионов присутствуют

### Ручная проверка

1. Открыть: https://www.occto.or.jp/kyoukei/toritukumi/system_reserve.html
2. Сравнить reserve margins по регионам

---

## 🧪 Комплексная проверка

### 1. Fetch всех данных за одну дату

```bash
#!/bin/bash
# test-all-data.sh

DATE=$(date +%Y-%m-%d)
echo "🧪 Testing data fetch for $DATE"

cd backend

# Build binaries
echo "🔨 Building fetch binaries..."
go build -o fetch-demand cmd/fetch-demand-http/main.go
go build -o fetch-jepx cmd/fetch-jepx-http/main.go
go build -o fetch-reserve cmd/fetch-reserve-http/main.go

# Fetch all data
echo "📊 Fetching Tokyo demand..."
./fetch-demand -area tokyo -date $DATE --use-http

echo "📊 Fetching Kansai demand..."
./fetch-demand -area kansai -date $DATE --use-http

echo "💴 Fetching JEPX Tokyo..."
./fetch-jepx -area tokyo -date $DATE --use-http

echo "💴 Fetching JEPX Kansai..."
./fetch-jepx -area kansai -date $DATE --use-http

echo "🔋 Fetching reserve data..."
./fetch-reserve -date $DATE

# Validate
echo ""
echo "✅ Validation Results:"
echo "===================="

for file in public/data/jp/*/*.json public/data/jp/*/*/*.json; do
  if [ -f "$file" ]; then
    size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
    if [ $size -gt 100 ]; then
      echo "✓ $file ($size bytes)"

      # Check if valid JSON
      if jq empty "$file" 2>/dev/null; then
        points=$(jq -r '.series | length' "$file" 2>/dev/null || jq -r '.reserves | length' "$file" 2>/dev/null)
        echo "  └─ Data points: $points"
      else
        echo "  └─ ⚠️  Invalid JSON!"
      fi
    else
      echo "⚠️  $file is too small ($size bytes)"
    fi
  fi
done
```

### 2. Сравнение с реальными данными

```bash
# Compare Tokyo demand from TEPCO website vs our data
echo "Comparison: TEPCO website vs our fetch"
echo "======================================"

# Manual: Open https://www.tepco.co.jp/forecast/html/index-j.html
# Get current demand value

# Our data:
jq -r '.series[-1] | "Our data: \(.ts) demand=\(.demand_mw) MW"' \
  public/data/jp/tokyo/demand-$(date +%Y-%m-%d).json
```

### 3. Проверка графиков

```bash
# Generate simple ASCII charts for visual inspection
pip install termgraph

# Tokyo demand chart
jq -r '.series[] | "\(.ts | split("T")[1] | split(":")[0]) \(.demand_mw)"' \
  public/data/jp/tokyo/demand-2025-11-05.json | termgraph

# JEPX price chart
jq -r '.price_yen_per_kwh[] | "\(.ts | split("T")[1] | split(":")[0]) \(.price)"' \
  public/data/jp/jepx/spot-tokyo-2025-11-05.json | termgraph
```

---

## 📋 Checklist перед автоматизацией

### TEPCO (Tokyo)
- [ ] ZIP архив скачивается
- [ ] Shift-JIS → UTF-8 работает
- [ ] CSV парсится без ошибок
- [ ] Единицы 万kW → MW конвертируются
- [ ] 24 точки в JSON
- [ ] Значения соответствуют реальности
- [ ] График выглядит корректно (пик днём)

### Kansai Electric
- [ ] CSV скачивается
- [ ] Кодировка определена
- [ ] Формат CSV корректный
- [ ] Единицы измерения верные
- [ ] 24 точки в JSON
- [ ] Значения меньше Tokyo (~50%)
- [ ] График выглядит корректно

### JEPX (Tokyo/Kansai)
- [ ] CSV скачивается для обеих зон
- [ ] Цены в JPY/kWh
- [ ] 24 точки (day-ahead)
- [ ] Ночные цены ниже дневных
- [ ] Пик цен вечером (18:00-20:00)
- [ ] Tokyo дороже Kansai на 5-10%
- [ ] Meta (min/max/avg) корректные

### OCCTO Reserve
- [ ] Данные за все 10 регионов
- [ ] Названия регионов корректные
- [ ] Reserve margin formula верная
- [ ] Tokyo/Kansai demand совпадают с TEPCO/Kansai
- [ ] Reserve обычно 5-30%

---

## 🐛 Reporting Issues

Когда нашли проблему:

1. **Создать issue** с описанием:
   ```markdown
   ## Problem
   [Описание проблемы]

   ## Source
   [TEPCO / Kansai / JEPX / OCCTO]

   ## Steps to Reproduce
   ./fetch-demand -area tokyo -date 2025-11-05 --use-http

   ## Expected
   [Что должно быть]

   ## Actual
   [Что получилось]

   ## Files
   - CSV: [приложить или показать пример]
   - JSON: [приложить вывод]
   ```

2. **Пометить в этом файле** проблемные источники

3. **Не запускать автоматизацию** пока не исправлено

---

## 📚 Полезные ссылки

- TEPCO Forecast: https://www.tepco.co.jp/forecast/html/index-j.html
- Kansai Electric: https://www.kansai-td.co.jp/denkiyoho/
- JEPX Market: https://www.jepx.jp/market/index.html
- OCCTO Reserve: https://www.occto.or.jp/kyoukei/toritukumi/system_reserve.html
