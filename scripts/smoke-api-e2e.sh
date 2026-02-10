#!/usr/bin/env bash

set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://127.0.0.1:8085}"
SMOKE_DB_HOST="${SMOKE_DB_HOST:-mysql}"
SMOKE_DB_PORT="${SMOKE_DB_PORT:-3306}"
SMOKE_DB_NAME="${SMOKE_DB_NAME:-goreport}"
SMOKE_DB_USER="${SMOKE_DB_USER:-root}"
SMOKE_DB_PASSWORD="${SMOKE_DB_PASSWORD:-root}"

TOKEN=""
DATASOURCE_ID=""
DATASET_ID=""
REPORT_ID=""

log() {
  printf '[smoke] %s\n' "$1"
}

call_api() {
  local method="$1"
  local path="$2"
  local data="${3:-}"
  local auth_header=()

  if [[ -n "$TOKEN" ]]; then
    auth_header=(-H "Authorization: Bearer ${TOKEN}")
  fi

  if [[ -n "$data" ]]; then
    curl -sS -X "$method" "${API_BASE_URL}${path}" \
      -H 'Content-Type: application/json' \
      "${auth_header[@]}" \
      -d "$data"
  else
    curl -sS -X "$method" "${API_BASE_URL}${path}" \
      "${auth_header[@]}"
  fi
}

extract_json_field() {
  local json="$1"
  local expr="$2"
  python3 -c "import json,sys; data=json.loads(sys.stdin.read()); value=${expr}; print(value if value is not None else '')" <<<"$json"
}

assert_success() {
  local json="$1"
  local stage="$2"
  local success
  success="$(extract_json_field "$json" "data.get('success')")"
  if [[ "$success" != "True" ]]; then
    log "${stage} failed: ${json}"
    exit 1
  fi
}

cleanup() {
  set +e

  if [[ -n "$REPORT_ID" ]]; then
    call_api DELETE "/api/v1/jmreport/delete?id=${REPORT_ID}" >/dev/null
  fi

  if [[ -n "$DATASET_ID" ]]; then
    call_api DELETE "/api/v1/datasets/${DATASET_ID}" >/dev/null
  fi

  if [[ -n "$DATASOURCE_ID" ]]; then
    call_api DELETE "/api/v1/datasource/${DATASOURCE_ID}" >/dev/null
  fi
}

trap cleanup EXIT

log "checking backend health"
health_json="$(call_api GET '/health')"
if [[ -z "$health_json" ]]; then
  log 'health check failed: empty response'
  exit 1
fi

log "login"
login_json="$(call_api POST '/api/v1/auth/login' '{"username":"admin","password":"admin123"}')"
assert_success "$login_json" 'login'
TOKEN="$(extract_json_field "$login_json" "data.get('result',{}).get('token','')")"
if [[ -z "$TOKEN" ]]; then
  log 'login failed: token is empty'
  exit 1
fi

suffix="$(date +%s)"
ds_name="smoke-ds-${suffix}"
dataset_name="smoke-dataset-${suffix}"
report_name="smoke-report-${suffix}"

log "create datasource"
create_ds_payload="$(cat <<EOF
{"name":"${ds_name}","type":"mysql","host":"${SMOKE_DB_HOST}","port":${SMOKE_DB_PORT},"database":"${SMOKE_DB_NAME}","username":"${SMOKE_DB_USER}","password":"${SMOKE_DB_PASSWORD}"}
EOF
)"
create_ds_json="$(call_api POST '/api/v1/datasource/create' "$create_ds_payload")"
assert_success "$create_ds_json" 'create datasource'
DATASOURCE_ID="$(extract_json_field "$create_ds_json" "data.get('result',{}).get('id','')")"

log "test datasource connection"
test_ds_json="$(call_api POST '/api/v1/datasource/test' "$create_ds_payload")"
assert_success "$test_ds_json" 'test datasource'

log "create dataset"
create_dataset_payload="$(cat <<EOF
{"name":"${dataset_name}","type":"sql","datasourceId":"${DATASOURCE_ID}","config":{"query":"SELECT id, username FROM users"}}
EOF
)"
create_dataset_json="$(call_api POST '/api/v1/datasets' "$create_dataset_payload")"
assert_success "$create_dataset_json" 'create dataset'
DATASET_ID="$(extract_json_field "$create_dataset_json" "data.get('result',{}).get('id','')")"

log "preview dataset"
preview_dataset_json="$(call_api GET "/api/v1/datasets/${DATASET_ID}/preview")"
assert_success "$preview_dataset_json" 'preview dataset'

log "create report"
create_report_payload="$(cat <<EOF
{"name":"${report_name}","type":"report","config":{"cells":[{"row":0,"col":0,"text":"Smoke OK"}]}}
EOF
)"
create_report_json="$(call_api POST '/api/v1/jmreport/create' "$create_report_payload")"
assert_success "$create_report_json" 'create report'
REPORT_ID="$(extract_json_field "$create_report_json" "data.get('result',{}).get('id','')")"

log "preview report"
preview_report_json="$(call_api POST '/api/v1/jmreport/preview' "{\"id\":\"${REPORT_ID}\",\"params\":{}}")"
assert_success "$preview_report_json" 'preview report'

html_len="$(extract_json_field "$preview_report_json" "len(data.get('result',{}).get('html',''))")"
if [[ "$html_len" == "0" ]]; then
  log 'preview report failed: html is empty'
  exit 1
fi

log 'smoke e2e passed'
