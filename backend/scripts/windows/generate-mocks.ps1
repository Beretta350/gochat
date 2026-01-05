# Script to generate mocks for interfaces using mockery

Write-Output "Generating mocks for interfaces..."

# Ensure mocks directory exists
if (-not (Test-Path -Path "mocks")) {
    New-Item -ItemType Directory -Path "mocks"
}

# Generate mock for wsadapter interfaces
Write-Output "Generating mocks for wsadapter..."
mockery --dir=internal/app/adapters/wsadapter --all --output=./mocks --outpkg=mocks

# Generate mock for websocket service interfaces
Write-Output "Generating mocks for websocket/service..."
mockery --dir=internal/app/websocket/service --all --output=./mocks --outpkg=mocks

# Generate mock for kafkaclient interfaces
Write-Output "Generating mocks for kafkaclient..."
mockery --dir=internal/app/kafkaclient --all --output=./mocks --outpkg=mocks

# Generate mock for kafkafactory interfaces
Write-Output "Generating mocks for kafkafactory..."
mockery --dir=./pkg/kafkafactory --all --output=./mocks --outpkg=mocks

Write-Output "Mock generation completed successfully!" 