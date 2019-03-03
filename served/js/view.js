"use strict";


function init() {
	var params = new URLSearchParams(window.location.search.substring(1));
	var userhash = params.get("id");

	if(userhash === null){
		console.log("error: no userhash found in url string.");
		return;
	}

	sendAjaxRequest("/api/getusermeetup", JSON.stringify({userhash: userhash}), function (error, response) {
		if(error !== null){
			errorArea.textContent = error;
			errorArea.classList.remove("hidden");
		}
		else if(response.error !== ""){
			errorArea.textContent = response.error;
			errorArea.classList.remove("hidden");
		} else {
			console.log(response)
		}
	});

	

	/*
	var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
	var days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
	var dateColumns = document.querySelectorAll(".dateColumn");

	for (var i = 0; i < dateColumns.length; i++) {
		var dateBox = dateColumns[i].querySelector(".dateBox");
		var dateString = dateColumns[i].querySelector("input:enabled[type=checkbox]").getAttribute("name");
		if( dateString != null ){
			var unixDate = parseInt(dateString, 10);
			if(isNaN(unixDate) === false){
				var date = new Date(unixDate);
				var spans = dateBox.querySelectorAll("span");
				spans[0].textContent = months[date.getMonth()];
				spans[1].textContent = date.getDate().toString(10);
				spans[2].textContent = days[date.getDay()]
			}
		}
	}
	*/


}


document.onreadystatechange = function () {
	if (document.readyState === "complete") {
		init();
	}
};
