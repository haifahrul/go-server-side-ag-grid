#!/bin/bash

echo "Start linting========"
CompileDaemon -build="echo true" \
    -color="true" \
    -command="golint ./internal/..."

echo ""
echo "End Linting********"
