#!/bin/bash

echo "Start Go Run ========"
CompileDaemon -build="echo true" \
    -color="true" \
    -command="golint ./internal/... " \
    -command="go run ./internal/... "

echo ""
echo "End Go Run ********"
