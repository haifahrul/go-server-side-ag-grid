#!/bin/bash

echo "Start linting========"
CompileDaemon -build="echo true" \
    -color="true" \
    -command="golint ./cmd/... " \
    -command="go run ./cmd/..."

echo ""
echo "End Linting********"
