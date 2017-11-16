//SPDX-License-Identifier: Apache-2.0

var loan = require('./controller.js');

module.exports = function(app){

  app.get('/get_loan/:id', function(req, res){
    loan.get_loan(req, res);
  });
  app.get('/add_loan/:loan', function(req, res){
    loan.add_loan(req, res);
  });
  app.get('/get_all_loan', function(req, res){
    console.log('check point C');
    loan.get_all_loan(req, res);
  });
  app.get('/change_lender/:lender', function(req, res){
    loan.change_lender(req, res);
  });
}
