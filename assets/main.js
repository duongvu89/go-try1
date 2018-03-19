'use strict';
setInterval(function(){
	$('#latest-pairs').load('/pairs');
}, 5000);

function displaySum() {
	$('#sum').load('/sum');
}

function displayMedian() {
	$('#median').load('/median');
}