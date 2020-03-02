# ag-Grid Server-Side 

- [Ag Grid](https://www.ag-grid.com/)
- Node.js Example
- Go Example

A reference implementation showing how to perform server-side operations using ag-Grid with api server node.js, go and MySQL.

![](https://github.com/ag-grid/ag-grid/blob/latest/packages/ag-grid-docs/src/nodejs-server-side-operations/app-arch.png "")

Reff: for full details see: http://ag-grid.com/nodejs-server-side-operations/

## Usage

- Clone the project
- run `yarn install`
- start with `yarn start` to run only angular
- open browser at `localhost:4000`

If you want to start the angular and api node.js
- start with `yarn dev` to run angular & node.js server

## GO Pre requested
- Install `go get github.com/githubnemo/CompileDaemon`
- Install `go get golang.org/x/lint/golint`

### Run Go with Makefile
- `cd go`
- `make goget` To run `go get & go mod vendor`
- `make gorun` To run `go run main.go`
