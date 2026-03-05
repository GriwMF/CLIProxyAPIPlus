#!/bin/bash
# Extracts the Cline refresh token from the VSCode extension's secret storage file.

SECRETS_FILE="$HOME/.cline/data/secrets.json"

if [ ! -f "$SECRETS_FILE" ]; then
    echo "Error: Cline secrets file not found at $SECRETS_FILE"
    echo "You may need to log into Cline in VSCode first."
    exit 1
fi

TOKEN=$(python3 -c "import sys,json; d=json.load(sys.stdin); print(json.loads(d.get('cline:clineAccountId','{}')).get('refreshToken',''))" < "$SECRETS_FILE")

if [ -z "$TOKEN" ] || [ "$TOKEN" == "None" ]; then
    echo "Error: Refresh token not found in secrets file."
    exit 1
fi

echo "$TOKEN"
