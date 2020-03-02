import webpack from 'webpack';
import webpackMiddleware from 'webpack-dev-middleware';
import webpackConfig from '../webpack.config.js';
import express from 'express';
import bodyParser from 'body-parser';
import fetch from 'node-fetch'
import OlympicWinnersService from './olympicWinnersService';

const app = express();
app.use(webpackMiddleware(webpack(webpackConfig)));
app.use(bodyParser.urlencoded({extended: false}));
app.use(bodyParser.json());

app.post('/nodeOlympicWinners', function (req, res) {
    OlympicWinnersService.getData(req.body, (rows, lastRow) => {
        res.json({rows: rows, lastRow: lastRow});
    });
});

app.post('/goOlympicWinnersSQL', function (req, res) {
    console.log('---BODY---', req.body)
    console.log('')
    fetch('http://localhost:8080/sql-olympic-winners', {
        method: 'post',
        body: JSON.stringify(req.body),
        headers: {"Content-Type": "application/json; charset=utf-8;"}
     })
     .then(httpResponse => httpResponse.json())
     .then(response => {
        res.json({rows: response.rows, lastRow: response.lastRow});
     })
     .catch(error => {
         console.error(error);
     })
});

app.listen(4000, () => {
    console.log('Started on localhost:4000');
});
