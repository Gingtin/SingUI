#!/usr/bin/env bash

# Quick run script for local development or testing
set -e

echo "Starting SingUI Server..."
go run main.go -p 2096 -d data/singbox_ui.db
