#!/usr/bin/env sh

set -euf -o pipefail

command -v jq > /dev/null 2>&1 || {
  echo 'Missing required dependency jq'
  exit 1
}

command -v terraform > /dev/null 2>&1 || {
  echo 'Missing required dependency terraform'
  exit 1
}

# In case this is the first time, initialize terraform (it's safe to run this multiple times).
terraform init

# Auto-apply the latest recommendations.
terraform apply -auto-approve

# Output those recommendations as something human-readable.
terraform output -json rules \
  | jq 'map({metric, match_type, aggregations, drop, drop_labels, keep_labels, aggregation_delay, aggregation_interval}) | sort_by(.metric)' \
  > rules.json
