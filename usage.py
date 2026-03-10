#!/usr/bin/env python3
# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "requests",
#     "python-dotenv",
# ]
# ///
"""Daily Copilot usage snapshot extractor"""
import requests
import json
import os
from datetime import datetime
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

GITHUB_TOKEN = os.getenv("GITHUB_TOKEN")
OUTPUT_DIR = "./copilot_usage_data"

if not GITHUB_TOKEN:
    raise ValueError("❌ GITHUB_TOKEN not found in .env file!")

os.makedirs(OUTPUT_DIR, exist_ok=True)

headers = {
    "Authorization": f"Bearer {GITHUB_TOKEN}",
    "Accept": "application/json",
}

# Call the internal Copilot user API
response = requests.get(
    "https://api.github.com/copilot_internal/user",
    headers=headers
)

if response.status_code == 200:
    data = response.json()
    snapshot = {
        "timestamp": datetime.utcnow().isoformat(),
        "date": datetime.utcnow().strftime("%Y-%m-%d"),
        "raw_response": data,
    }

    filename = f"{OUTPUT_DIR}/usage_{datetime.utcnow().strftime('%Y%m%d_%H%M%S')}.json"
    with open(filename, "w") as f:
        json.dump(snapshot, f, indent=2)

    print(f"✅ Saved: {filename}")
    print(f"📊 Response preview:\n{json.dumps(data, indent=2)[:500]}")
else:
    print(f"❌ Error: {response.status_code}")
    print(f"Response: {response.text}")
