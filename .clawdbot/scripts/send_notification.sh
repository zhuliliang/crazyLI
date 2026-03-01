#!/usr/bin/env bash
# Telegram notification helper.

set -euo pipefail

MESSAGE="$1"

TELEGRAM_BOT_TOKEN="8663455961:AAEZifuyg497cThOzbaHr_TQ6L0_msbKGP4"
TELEGRAM_CHAT_ID="1030611758"

curl -sS -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
  -d chat_id="${TELEGRAM_CHAT_ID}" \
  -d text="${MESSAGE}" >/dev/null
