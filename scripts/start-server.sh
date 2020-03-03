#!/bin/bash

echo "Start babel & Nodejs ========"
yarn start

echo "Start Go Run ========"
CompileDaemon \
    -color=true \
    -graceful-kill=true \
    -pattern="^(\.env.+|\.env)|(.+\.go|.+\.c)$" \
    -build="go build -mod=vendor -o ./go-ag-grid ./internal/..." \
    -command="./go-ag-grid"
echo ""
echo "End Go Run ********"
