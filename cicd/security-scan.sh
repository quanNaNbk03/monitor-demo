#!/bin/bash

REPORT_ROOT="security-reports"
EXCLUDE_DIRS="security-reports" "cicd" "frontend" "test" "docs"
JENKINS_CONTAINER_ID="jenkins"
WS_PATH=$WORKSPACE

rm -rf $REPORT_ROOT
mkdir -p $REPORT_ROOT/{gosec,gitleaks,semgrep,trivy}

echo "========================================================="
echo "STARTING FULL-PROJECT SECURITY SCANNING (ROOT LEVEL)"
echo "========================================================="

echo "[1/4] Running Gitleaks (Full Scan)..."
docker run --rm --network host --volumes-from $JENKINS_CONTAINER_ID \
    -w "$WS_PATH" zricethezav/gitleaks:latest \
    detect --source "." \
    --report-format sarif \
    --report-path "$REPORT_ROOT/gitleaks/all-secrets.sarif" \
    --no-git --exit-code 0

echo "[2/4] Running Gosec (Recursive Scan)..."
docker run --rm --network host --volumes-from $JENKINS_CONTAINER_ID \
    -w "$WS_PATH" securego/gosec:latest \
    -fmt sarif -out "$REPORT_ROOT/gosec/all-go-issues.sarif" \
    -exclude-dir="frontend,test,docs" \
    "./..." > /dev/null 2>&1

echo "[3/4] Running Semgrep (Full Scan)..."
docker run --rm --network host --volumes-from $JENKINS_CONTAINER_ID \
    -w "$WS_PATH" semgrep/semgrep:latest \
    semgrep scan --config=p/default --sarif \
    --output="$REPORT_ROOT/semgrep/all-sast.sarif" \
    --exclude="frontend" --exclude="test" --exclude="docs" --exclude="security-reports" .

echo "[4/4] Running Trivy (Filesystem Scan)..."
docker run --rm --network host --volumes-from $JENKINS_CONTAINER_ID \
    -w "$WS_PATH" aquasec/trivy:latest \
    fs --format sarif --output "$REPORT_ROOT/trivy/all-deps.sarif" \
    --severity HIGH,CRITICAL .

sudo chown -R $(id -u):$(id -g) $REPORT_ROOT

echo "========================================================="
echo "SCAN COMPLETED. Reports are in $REPORT_ROOT"
echo "========================================================="