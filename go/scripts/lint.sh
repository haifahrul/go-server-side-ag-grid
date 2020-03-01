#!/bin/bash

echo "Start linting========"
CompileDaemon -build="echo true" \
    -color="true" \
    -command="golint ./internal/... " \
    -command="go run ./internal/... "

echo ""
echo "End Linting********"
