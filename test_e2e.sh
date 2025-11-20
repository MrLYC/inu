#!/bin/bash
# End-to-end test for inu CLI
# This script tests the full anonymize -> restore workflow

set -e

echo "=== Inu CLI End-to-End Test ==="
echo

# Check required environment variables
if [ -z "$OPENAI_API_KEY" ] || [ -z "$OPENAI_MODEL_NAME" ]; then
    echo "Error: Required environment variables not set"
    echo "Please set:"
    echo "  export OPENAI_API_KEY='your-api-key'"
    echo "  export OPENAI_MODEL_NAME='gpt-4'"
    echo "  export OPENAI_BASE_URL='https://api.openai.com/v1'  # optional"
    exit 1
fi

# Clean up any previous test files
rm -f test_output.txt test_entities.yaml test_restored.txt

echo "Step 1: Anonymize text from file"
./bin/inu anonymize \
    --file test_input.txt \
    --output test_output.txt \
    --output-entities test_entities.yaml \
    --print-entities

echo
echo "Step 2: Verify anonymized output and entities files were created"
if [ ! -f test_output.txt ] || [ ! -f test_entities.yaml ]; then
    echo "Error: Output files not created"
    exit 1
fi
echo "✓ Files created successfully"

echo
echo "Step 3: Restore text using saved entities"
./bin/inu restore \
    --file test_output.txt \
    --entities test_entities.yaml \
    --output test_restored.txt \
    --print

echo
echo "Step 4: Compare restored text with original"
if diff -q test_input.txt test_restored.txt > /dev/null; then
    echo "✓ SUCCESS: Restored text matches original!"
else
    echo "✗ FAIL: Restored text differs from original"
    echo
    echo "Differences:"
    diff test_input.txt test_restored.txt || true
    exit 1
fi

echo
echo "=== All tests passed! ==="

# Clean up test files
rm -f test_output.txt test_entities.yaml test_restored.txt
