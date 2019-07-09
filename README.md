# Geo server

Uses new go modules so requires `go 1.12`. You must also place this project somewhere outside of your `GOPATH` for go to recognize it is a module. In GoLand ensure `Go > Go modules (vgo)` is enabled.

## Running
`go run server/server.go`

Hosts GraphQL Playground at `http://localhost:8080`.
