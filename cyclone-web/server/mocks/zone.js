/*jshint node:true*/

module.exports = function(app) {
  var express = require('express');
  var excerptsRouter = express.Router();

  var data = [
    {id: 1, name: 'zone-1'},
    {id: 2, name: 'zone-2'},
    {id: 3, name: 'zone-3'},
  ];

  excerptsRouter.get('/', function(req, res) {
    res.send(data);
  });

  excerptsRouter.post('/web', function(req, res) {
    // var resp = JSON.stringify(body);
    // res.status(201).end(resp);
  });

  excerptsRouter.get('/:id', function(req, res) {
    // var excerpt = null;
    // excerpts.forEach(function(item){
    //   if (item.id === parseInt(req.params.id)) {
    //     excerpt = item;
    //   }
    // });

    // if (excerpt) {
    //   res.send(excerpt);
    // } else {
    //   res.status(204).end({
    //     message: "record not found"
    //   });
    // }
  });

  excerptsRouter.put('/:id', function(req, res) {
    // res.send({
    //   'excerpts': {
    //     id: req.params.id
    //   }
    // });
  });

  excerptsRouter.delete('/:id', function(req, res) {
    // res.status(204).end();
  });

  // The POST and PUT call will not contain a request body
  // because the body-parser is not included by default.
  // To use req.body, run:

  //    npm install --save-dev body-parser

  // After installing, you need to `use` the body-parser for
  // this mock uncommenting the following line:
  //
  app.use('/api/zone', require('body-parser').json());
  app.use('/api/zone', excerptsRouter);
};
