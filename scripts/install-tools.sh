#!/bin/bash
set -e
go install github.com/projectdiscovery/katana/cmd/katana@latest
go install github.com/projectdiscovery/httpx/cmd/httpx@latest
echo "katana and httpx installed successfully"
