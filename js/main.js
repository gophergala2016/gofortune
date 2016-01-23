var app = angular.module('gofortune', [])

app.controller('cardsController', cardsController)

function cardsController($scope, $http) {
	var response = function(data, status) {
		if (data.data.Error)
			alert(data.data.Error)
		$scope.cards = data.data.Cards
		//console.log(data)
	}
	$scope.cards = []
	$http.post('/init', null).then(response, errorResponse)
}


function errorResponse(data, status) {
	alert('error: '+ data+', '+status)
} 
