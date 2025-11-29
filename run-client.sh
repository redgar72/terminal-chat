#!/bin/bash
cd "$(dirname "$0")/client"
go run main.go "$1"
