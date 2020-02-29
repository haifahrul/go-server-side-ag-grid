import {Grid} from 'ag-grid-community';
import 'ag-grid-enterprise';

import "ag-grid-community/dist/styles/ag-grid.css";
import "ag-grid-community/dist/styles/ag-theme-balham.css";

var cacheBlockSize = 2;

const gridOptions = {

    rowModelType: 'serverSide',

    columnDefs: [
        {field: 'athlete'},
        {field: 'country', hide: true},
        {field: 'sport', hide: true},
        // {field: 'country', rowGroup: true, hide: true},
        // {field: 'sport', rowGroup: true, hide: true},
        {field: 'year', filter: 'number', filterParams: {newRowsAction: 'keep'}},
        {field: 'gold', aggFunc: 'sum'},
        {field: 'silver', aggFunc: 'sum'},
        {field: 'bronze', aggFunc: 'sum'},
    ],

    defaultColDef: {
        sortable: true
    },

    pagination: true,
    paginationPageSize: cacheBlockSize,

    // debug: true,
    cacheBlockSize: cacheBlockSize,
    // maxBlocksInCache: cacheBlockSize,
    // purgeClosedRowNodes: true,
    // maxConcurrentDatasourceRequests: 2,
    // blockLoadDebounceMillis: 1000
};

function onPageSizeChanged() {
    document.getElementById('page-size').addEventListener("change", function(e) {
        var value = document.getElementById('page-size').value;
        console.log(value);
        cacheBlockSize = value;
        gridOptions.api.paginationSetPageSize(Number(value));
        gridOptions.api.setServerSideDatasource(datasource);
    })
}

const gridDiv = document.querySelector('#myGrid');
new Grid(gridDiv, gridOptions);

onPageSizeChanged()

const datasource = {
    getRows(params) {

         fetch('./olympicWinners/', {
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

gridOptions.api.setServerSideDatasource(datasource);