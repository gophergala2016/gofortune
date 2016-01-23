var app = angular.module('gofortune', [])

app.controller('cardsController', cardsController)

function cardsController($scope, $http) {
	$scope.step = 1
	$scope.count = 1
	$scope.http = $http
	$scope.deal = deal
	$scope.step2 = step2
	var response = function(data, status) {
		if (data.data.Error)
			alert(data.data.Error)
		$scope.cards = data.data.Cards
		//console.log(data)
	}
	$http.post('/init', null).then(response, errorResponse)
}


function errorResponse(data, status) {
	alert('error: '+ data+', '+status)
	console.log(data)
} 

function step2() {
	var scope = this
	var cards = scope.cards
	var row0 = []
	var row1 = []
	var row2 = []
	for (var i=0; i<cards.length; i++) {
		if (i<7)
			row0.push(cards[i])
		if (i>6 && i<14)
			row1.push(cards[i])
		if (i>13)
			row2.push(cards[i])
	}
	scope.row0 = row0
	scope.row1 = row1
	scope.row2 = row2
	scope.step = 2
}

function deal(row) {
	var scope = this
	var response = function(data, status) {
		if (data.data.Error) {
			alert(data.data.Error)
			return
		}
		//console.log(data)
		scope.row0 = data.data.Row0
		scope.row1 = data.data.Row1
		scope.row2 = data.data.Row2
		scope.count++
	}
	var cards = []
	for (var i=0; i<scope.cards.length; i++) {
		cards.push({Image: scope.cards[i].Image})
	}
	var data = {
		Cards: cards,
		Row: row,
		Count: scope.count,
	}
	scope.http.post('/deal', JSON.stringify(data)).then(response, errorResponse)
}
