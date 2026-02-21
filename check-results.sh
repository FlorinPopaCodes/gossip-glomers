#!/usr/bin/env bash
set -euo pipefail

# Usage: ./check-results.sh <max-msgs-per-op> <max-median-ms> <max-max-ms> [results-file]
MAX_MSGS="${1:?Usage: $0 <max-msgs-per-op> <max-median-ms> <max-max-ms> [results-file]}"
MAX_MEDIAN="${2:?}"
MAX_MAX="${3:?}"
RESULTS="${4:-store/broadcast/latest/results.edn}"

if [ ! -f "$RESULTS" ]; then
  echo "FAIL: results file not found: $RESULTS"
  exit 1
fi

# Extract metrics from EDN
# :servers {:send-count ..., :msgs-per-op 3646.2449}
msgs_per_op=$(awk '/:servers/{found=1} found && /:msgs-per-op/{gsub(/[^0-9.]/,"",$0); print; exit}' "$RESULTS")
# :stable-latencies {0 0, 0.5 57, 0.95 118, 0.99 190, 1 203} (may span multiple lines)
latencies=$(awk '/:stable-latencies/{found=1} found{line=line $0} found && /\}/{gsub(/[{},]/, " ", line); print line; exit}' "$RESULTS")
median_latency=$(echo "$latencies" | awk '{for(i=1;i<=NF;i++) if($i=="0.5") print $(i+1)}')
max_latency=$(echo "$latencies" | awk '{print $NF}')

pass=true

echo ""
printf "%-25s %10s %10s %s\n" "Metric" "Actual" "Limit" "Result"
printf "%-25s %10s %10s %s\n" "-------------------------" "----------" "----------" "------"

check() {
  local name="$1" actual="$2" limit="$3"
  if awk "BEGIN{exit !($actual > $limit)}"; then
    printf "%-25s %10s %10s FAIL\n" "$name" "$actual" "<$limit"
    pass=false
  else
    printf "%-25s %10s %10s PASS\n" "$name" "$actual" "<$limit"
  fi
}

check "msgs-per-op (servers)" "$msgs_per_op" "$MAX_MSGS"
check "median latency (ms)" "$median_latency" "$MAX_MEDIAN"
check "max latency (ms)" "$max_latency" "$MAX_MAX"

echo ""
if $pass; then
  echo "All checks passed!"
else
  echo "Some checks FAILED"
  exit 1
fi
