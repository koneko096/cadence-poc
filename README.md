# POC Temporal

Simulating ride hailing booking workflow using Temporal.io 

## Requirements
- Go 1.22+
- JDK 17+
- [Temporal CLI](https://temporal.io/setup/install-temporal-cli)
- [Buf](https://buf.build/docs/cli/quickstart/)

## Running
1. Generate stubs `buf generate`
2. Run Temporal server `temporal dev`
3. Run the tests on the servers & workers `./run.sh`
4. Open temporal dev server to see workflow running on http://localhost:8233