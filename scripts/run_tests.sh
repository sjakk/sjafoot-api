#!/bin/bash


set -e

echo "Running Go tests..."

export SJAFOOT_DB_DSN=${SJAFOOT_DB_DSN:-"postgres://sjafoot_user:yourpassword@localhost/sjafoot_test?sslmode=disable"}

cd ../cmd/api

go test -v

echo "Tests passed successfully!"
