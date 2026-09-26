#!/bin/bash
set -xe

cd /app 
npm i
npm run check-watch &
npx rsbuild --host 0.0.0.0
