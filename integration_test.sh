#!/bin/bash
set -e

# Integration test for haide
# Tests the full workflow in a temporary repository

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMP_DIR=$(mktemp -d)
TEMP_CONFIG_DIR=$(mktemp -d)
export HAIDE_HOME="$TEMP_CONFIG_DIR"

echo "Testing haide integration..."
echo "Temp repo: $TEMP_DIR"
echo "Temp config: $TEMP_CONFIG_DIR"

# Build haide
echo "Building haide..."
cd "$SCRIPT_DIR"
go build -o haide ./cmd/haide
HAIDE_BIN="$SCRIPT_DIR/haide"

# Create test repository
cd "$TEMP_DIR"
git init test-repo
cd test-repo

# Test 1: Initialize haide
echo -e "\n[Test 1] Initialize haide"
"$HAIDE_BIN" init
if [ ! -f .git/info/exclude ]; then
    echo "FAIL: .git/info/exclude not created"
    exit 1
fi
if ! grep -q "BEGIN HAIDE" .git/info/exclude; then
    echo "FAIL: Haide markers not found"
    exit 1
fi
echo "PASS"

# Test 2: Check config was created
echo -e "\n[Test 2] Config creation"
if [ ! -f "$HAIDE_HOME/config.ini" ]; then
    echo "FAIL: Config file not created"
    exit 1
fi
echo "PASS"

# Test 3: Add project-specific pattern
echo -e "\n[Test 3] Add pattern"
$HAIDE_BIN add "test-pattern.md"
if ! grep -q "test-pattern.md" .git/info/exclude; then
    echo "FAIL: Pattern not added to exclude file"
    exit 1
fi
if ! grep -q "test-pattern.md" "$HAIDE_HOME/config.ini"; then
    echo "FAIL: Pattern not added to config"
    exit 1
fi
echo "PASS"

# Test 4: Info command
echo -e "\n[Test 4] Info command"
$HAIDE_BIN info > /tmp/haide-info.txt
if ! grep -q "test-pattern.md" /tmp/haide-info.txt; then
    echo "FAIL: Info doesn't show pattern"
    exit 1
fi
echo "PASS"

# Test 5: Add global pattern
echo -e "\n[Test 5] Add global pattern"
$HAIDE_BIN add-global "*.global"
if ! grep -q "*.global" "$HAIDE_HOME/config.ini"; then
    echo "FAIL: Global pattern not added"
    exit 1
fi
echo "PASS"

# Test 6: Update to apply global pattern
echo -e "\n[Test 6] Update command"
yes | $HAIDE_BIN update
if ! grep -q "*.global" .git/info/exclude; then
    echo "FAIL: Global pattern not applied after update"
    exit 1
fi
echo "PASS"

# Test 7: Remove pattern
echo -e "\n[Test 7] Remove pattern"
$HAIDE_BIN remove "test-pattern.md"
if grep -q "test-pattern.md" .git/info/exclude; then
    echo "FAIL: Pattern still in exclude file after removal"
    exit 1
fi
echo "PASS"

# Test 8: Tracked file override
echo -e "\n[Test 8] Tracked file override"
# Clean haide first to start fresh
yes | $HAIDE_BIN clean
# Add and track CLAUDE.md before initializing haide
echo "test" > CLAUDE.md
git add -f CLAUDE.md  # Force add despite ignore
git config user.email "test@test.com"
git config user.name "Test User"
git commit -m "Add CLAUDE.md"
# Now initialize haide - should detect tracked file
$HAIDE_BIN init <<< "n"  # Answer 'n' to untrack prompt
if ! grep -q "+CLAUDE.md" "$HAIDE_HOME/config.ini"; then
    echo "FAIL: Override not added for tracked file"
    exit 1
fi
echo "PASS"

# Test 9: Clean command
echo -e "\n[Test 9] Clean command"
yes | $HAIDE_BIN clean
if grep -q "BEGIN HAIDE" .git/info/exclude; then
    echo "FAIL: Haide markers still present after clean"
    exit 1
fi
echo "PASS"

# Test 10: Dry-run doesn't modify
echo -e "\n[Test 10] Dry-run"
rm .git/info/exclude
$HAIDE_BIN init --dry-run
if [ -f .git/info/exclude ]; then
    echo "FAIL: Dry-run created exclude file"
    exit 1
fi
echo "PASS"

# Cleanup
cd /
rm -rf "$TEMP_DIR" "$TEMP_CONFIG_DIR"

echo -e "\n✓ All integration tests passed!"
