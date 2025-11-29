#!/bin/bash
set -euo pipefail
cd /app
npm install
npm run dev -- --host 0.0.0.0 --port "${NUXT_PORT:-3000}"
