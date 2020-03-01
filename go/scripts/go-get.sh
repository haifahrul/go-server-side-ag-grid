#!/bin/bash

echo "Start Go Get ========"
cd internal
go get
go mod vendor
echo ""
echo "End Go Get ********"
