#!/usr/bin/env sh

set -euf -o pipefail

command -v curl > /dev/null 2>&1 || {
  echo 'Missing required dependency curl'
  exit 1
}

command -v jq > /dev/null 2>&1 || {
  echo 'Missing required dependency jq'
  exit 1
}

curl -u "$GRAFANA_AM_API_KEY" "$GRAFANA_AM_API_URL/aggregations/recommendations?verbose=false" \
  | jq 'sort_by(.metric)' \
  > recommendations.json
