@echo off
REM Script to generate mocks for interfaces using mockery

echo Generating mocks for interfaces...

REM Ensure mocks directory exists
if not exist mocks mkdir mocks

REM Generate mock for wsadapter interfaces
echo Generating mocks for wsadapter...
mockery --dir=internal/app/adapters/wsadapter --all --output=./mocks --outpkg=mocks

REM Generate mock for websocket service interfaces
echo Generating mocks for websocket/service...
mockery --dir=internal/app/websocket/service --all --output=./mocks --outpkg=mocks

REM Generate mock for kafkaclient interfaces
echo Generating mocks for kafkaclient...
mockery --dir=internal/app/kafkaclient --all --output=./mocks --outpkg=mocks

REM Generate mock for kafkafactory interfaces
echo Generating mocks for kafkafactory...
mockery --dir=./pkg/kafkafactory --all --output=./mocks --outpkg=mocks

echo Mock generation completed successfully! 