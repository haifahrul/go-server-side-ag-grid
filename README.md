[<img src="https://img.shields.io/github/issues/haifahrul/go-server-side-ag-grid"/>](https://github.com/haifahrul/go-server-side-ag-grid/issues)
[<img src="https://img.shields.io/github/forks/haifahrul/go-server-side-ag-grid"/>](https://github.com/haifahrul/go-server-side-ag-grid/network/members)
![stars](https://img.shields.io/github/stars/haifahrul/go-server-side-ag-grid)
![Hacktoberfest](https://img.shields.io/badge/Hacktoberfest-magenta)

<img src="https://github.com/haifahrul/go-server-side-ag-grid/blob/master/data/sample.png?raw=true" style="height: 200px; width: auto;"/>

# ag-Grid Server-Side 

- [Ag Grid](https://www.ag-grid.com/)
- Node.js Example
- Go Example
- MySQL
- PostgreSQL

A reference implementation showing how to perform server-side operations using ag-Grid with api server node.js and go.

![](https://github.com/ag-grid/ag-grid/blob/latest/packages/ag-grid-docs/src/nodejs-server-side-operations/app-arch.png "")

Reff: for full details see: http://ag-grid.com/nodejs-server-side-operations/

## Pre requested

- copy file `.env-example` and rename to `.env`
- setup your credential database such as username, password and etc
- Install `go get github.com/githubnemo/CompileDaemon`
- Install `go get golang.org/x/lint/golint`

## How to Use

- Clone the project
- run `yarn install`
- start with `yarn start`
- open browser at `localhost:4000`

If you want to start the angular and api node.js
- start with `yarn dev` to run angular & node.js server with nodemon

### Run Go with Makefile
- `make goget` To run `go get & go mod vendor`
- `make gorun` To run `go run main.go`
