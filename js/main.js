var app = angular.module('gofortune', [])

app.controller('cardsController', cardsController)

var cards = [
	//{file:'2C.png'},
	//{file:'2D.png'},
	//{file:'2H.png'},
	//{file:'2S.png'},
	//{file:'3C.png'},
	//{file:'3D.png'},
	//{file:'3H.png'},
	//{file:'3S.png'},
	//{file:'4C.png'},
	//{file:'4D.png'},
	//{file:'4H.png'},
	//{file:'4S.png'},
	//{file:'5C.png'},
	//{file:'5D.png'},
	//{file:'5H.png'},
	//{file:'5S.png'},
	//{file:'6C.png'},
	//{file:'6D.png'},
	//{file:'6H.png'},
	//{file:'6S.png'},
	//{file:'7D.png'},
	//{file:'7H.png'},
	//{file:'7S.png'},
	//{file:'7C.png'},
	//{file:'8C.png'},
	//{file:'8D.png'},
	//{file:'8H.png'},
	//{file:'8S.png'},
	//{file:'9C.png'},
	//{file:'9D.png'},
	//{file:'9H.png'},
	{file:'9S.png'},
	{file:'10C.png'},
	{file:'10D.png'},
	{file:'10H.png'},
	{file:'10S.png'},
	{file:'JC.png'},
	{file:'JD.png'},
	{file:'JH.png'},
	{file:'JS.png'},
	{file:'QC.png'},
	{file:'QD.png'},
	{file:'QH.png'},
	{file:'QS.png'},
	{file:'KC.png'},
	{file:'KD.png'},
	{file:'KH.png'},
	{file:'KS.png'},
	{file:'AC.png'},
	{file:'AD.png'},
	{file:'AH.png'},
	{file:'AS.png'}
]

function cardsController($scope) {
	$scope.cards = cards
}


