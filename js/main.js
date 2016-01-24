var app = angular.module('gofortune', [])

app.controller('cardsController', cardsController)

function cardsController($scope, $http) {
	$scope.http = $http
	$scope.deal = deal
	$scope.startOver = startOver
	$scope.init = init
	$scope.init()
	$scope.fortune = document.getElementById('fortune')
}

function init() {
	var scope = this
	scope.count = 0
	scope.step = 1
	var response = function(data, status) {
		if (data.data.Error)
			alert(data.data.Error)
		scope.cards = data.data.Cards
		//console.log(data)
	}
	scope.http.post('/init', null).then(response, errorResponse)
}


function errorResponse(data, status) {
	alert('error: '+ data+', '+status)
	console.log(data)
} 

function deal(row) {
	var scope = this
	var response = function(data, status) {
		if (data.data.Error) {
			alert(data.data.Error)
			return
		}
		if (data.data.Card) {
			scope.step = 3
			scope.yourCard = data.data.Card

			scope.fortune.innerHTML = 'loading ...'
			var req = {Card: scope.yourCard}
			var f = function() {
				var showFortune = function(d) {
					scope.fortune.innerHTML = d.Tweet
				}
				scope.http.post('/fortune', JSON.stringify(req)).success(showFortune)
			}
			setTimeout(f, 1000)
			return
		}
		//console.log(data)
		scope.row1 = data.data.Row1
		scope.row2 = data.data.Row2
		scope.row3 = data.data.Row3
		scope.count++
		scope.step = 2
	}
	var cards = []
	if (scope.step == 2) {
		for (var i=0; i<scope.row1.length; i++)
			cards.push({Image: scope.row1[i].Image})
		for (var i=0; i<scope.row2.length; i++)
			cards.push({Image: scope.row2[i].Image})
		for (var i=0; i<scope.row3.length; i++)
			cards.push({Image: scope.row3[i].Image})
	} else {
		cards = scope.cards
	}
	var data = {
		Cards: cards,
		Row: row,
		Count: scope.count,
	}
	scope.http.post('/deal', JSON.stringify(data)).then(response, errorResponse)
}

function startOver() {
	this.init()
}
