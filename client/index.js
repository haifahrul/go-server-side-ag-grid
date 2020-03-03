import {Grid} from 'ag-grid-community';
import 'ag-grid-enterprise';

import "ag-grid-community/dist/styles/ag-grid.css";
import "ag-grid-community/dist/styles/ag-theme-balham.css";

var defaultPage = 10;

const gridOptions = {

    rowModelType: 'serverSide',

    columnDefs: [
        {field: 'athlete', filter: 'text', filterParams: {newRowsAction: 'keep'}},
        {field: 'country', filter: 'text', filterParams: {newRowsAction: 'keep'}},
        // {field: 'sport', hide: true},
        // {field: 'country', rowGroup: true, hide: true},
        // {field: 'sport', rowGroup: true, hide: true},
        {field: 'year', filter: 'number', filterParams: {newRowsAction: 'keep'}},
        // {field: 'gold', aggFunc: 'sum'},
        // {field: 'silver', aggFunc: 'sum'},
        // {field: 'bronze', aggFunc: 'sum'},

        // {field: 'athlete'},
        // {field: 'country', rowGroup: true, hide: true},
        // {field: 'sport', rowGroup: true, hide: true},
        // {field: 'year', filter: 'number', filterParams: {newRowsAction: 'keep'}},
        // {field: 'gold', aggFunc: 'sum'},
        // {field: 'silver', aggFunc: 'sum'},
        // {field: 'bronze', aggFunc: 'sum'},
    ],

    defaultColDef: {
        filter: 'agSetColumnFilter',
        sortable: true,
        enableRowGroup: true
    },
    sideBar: true,
    rowDragManaged: true,
    rowGroupPanelShow: 'always',
    floatingFilter: true,
    // pagination: true,
    paginationPageSize: defaultPage,

    // debug: true,
    cacheBlockSize: 100,
    // maxBlocksInCache: cacheBlockSize,
    purgeClosedRowNodes: true,
    maxConcurrentDatasourceRequests: 2,
    blockLoadDebounceMillis: 1000
};

function onPageSizeChanged() {
    document.getElementById('page-size').addEventListener("change", function(e) {
        var value = document.getElementById('page-size').value;
        console.log(value);
        // cacheBlockSize = value;
        gridOptions.api.paginationSetPageSize(Number(value));
        gridOptions.api.setServerSideDatasource(datasource);
    })
}
onPageSizeChanged()

// API NODE.js
const gridDivNode = document.querySelector('#myGrid');
new Grid(gridDivNode, gridOptions);
const datasourceNode = {
    getRows(params) {
        fetch('./nodeOlympicWinners/', {
            method: 'post',
            body: JSON.stringify(params.request),
            headers: {"Content-Type": "application/json; charset=utf-8"}
        })
        .then(httpResponse => httpResponse.json())
        .then(response => {
            params.successCallback(response.rows, response.lastRow);
        })
        .catch(error => {
            console.error(error);
            params.failCallback();
        })
    }
};
gridOptions.api.setServerSideDatasource(datasourceNode);

// API Golang - DB MySQL
const gridDivGoMySQL = document.querySelector('#myGridGoMySQL');
new Grid(gridDivGoMySQL, gridOptions);
const datasourceGoMySQL = {
    getRows(params) {
        fetch('./goOlympicWinnersMySQL/', {
            method: 'post',
            body: JSON.stringify(params.request),
            headers: {"Content-Type": "application/json; charset=utf-8"}
        })
        .then(httpResponse => httpResponse.json())
        .then(response => {
            if (response.rows && response.rows.length > 0) {
                params.successCallback(response.rows, response.lastRow);
            } else {
                params.successCallback([], 0);
            }
        })
        .catch(error => {
            console.error(error);
            params.failCallback();
        })
    }
};
gridOptions.api.setServerSideDatasource(datasourceGoMySQL);

// API Golang - DB PostgreSQL
const gridDivGoPostgreSQL = document.querySelector('#myGridGoPostgreSQL');
new Grid(gridDivGoPostgreSQL, gridOptions);
const datasourceGoPostgreSQL = {
    getRows(params) {
        fetch('./goOlympicWinnersPostgreSQL/', {
            method: 'post',
            body: JSON.stringify(params.request),
            headers: {"Content-Type": "application/json; charset=utf-8"}
        })
        .then(httpResponse => httpResponse.json())
        .then(response => {
            if (response.rows && response.rows.length > 0) {
                params.successCallback(response.rows, response.lastRow);
            } else {
                params.successCallback([], 0);
            }
        })
        .catch(error => {
            console.error(error);
            params.failCallback();
        })
    }
};
gridOptions.api.setServerSideDatasource(datasourceGoPostgreSQL);