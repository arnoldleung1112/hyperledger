// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();
	
	$scope.queryAllLoan = function(){
		console.log('check point A');
		appFactory.queryAllLoan(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_loan = array;
		});
	}

	$scope.queryLoan = function(){

		var id = $scope.loan_id;

		appFactory.queryLoan(id, function(data){
			$scope.query_loan = data;

			if ($scope.query_loan == "Could not locate loan"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordloan = function(){

		appFactory.recordloan($scope.loan, function(data){
			$scope.create_loan = data;
			$("#success_create").show();
		});
	}

	$scope.changeLender = function(){

		appFactory.changeLender($scope.lender, function(data){
			$scope.change_lender = data;
			if ($scope.change_lender == "Error: no loan found"){
				$("#error_holder").show();
				$("#success_holder").hide();
			} else{
				$("#success_holder").show();
				$("#error_holder").hide();
			}
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllLoan = function(callback){
		console.log('check point B');
    	$http.get('/get_all_loan/').success(function(output){
			callback(output)
		});
	}

	factory.queryLoan = function(id, callback){
    	$http.get('/get_loan/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordloan = function(data, callback){

		data.location = data.longitude + ", "+ data.latitude;

		var loan = data.id + "-" + data.location + "-" + data.timestamp + "-" + data.lender + "-" + data.vessel;

    	$http.get('/add_loan/'+loan).success(function(output){
			callback(output)
		});
	}

	factory.changeLender = function(data, callback){

		var lender = data.id + "-" + data.name;

    	$http.get('/change_lender/'+lender).success(function(output){
			callback(output)
		});
	}

	return factory;
});


