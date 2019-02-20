"use strict";

function init () {
	var errorArea = document.getElementById('errorArea');
	var dateContainer = document.getElementById('dateContainer');

	document.getElementById('saveButt').addEventListener('click', function () {
		errorArea.textContent = "";
		errorArea.classList.add("hidden");

		var args = {
			description: document.getElementById("description").value,
			dates: [],
			admin: {
				email: document.getElementById("adminEmail").value,
				alerts: document.getElementById("notifications").checked
			}
		};

		var selectedDates = dateTool.getDates();
		for (var i = 0; i < selectedDates.length; i++) {
			args.dates.push({date: selectedDates[i], users: []});
		}

		sendAjaxRequest("create", JSON.stringify(args), function (error, response) {
			if(error !== null){
				errorArea.textContent = error;
				errorArea.classList.remove("hidden");
			}
			else if(response.error !== ""){
				errorArea.textContent = response.error;
				errorArea.classList.remove("hidden");
			} else {
				var link = document.getElementById("userLink");
				link.href = window.location.origin + "/view?id=" + encodeURIComponent(response.result.userlink);
				link.textContent = window.location.origin + "/view?id=" + encodeURIComponent(response.result.userlink);

				link = document.getElementById("adminLink");
				link.href = window.location.origin + "/admin?id=" + encodeURIComponent(response.result.adminlink);
				link.textContent = window.location.origin + "/admin?id=" + encodeURIComponent(response.result.adminlink);

				document.getElementById("linkArea").classList.remove("hidden");

				// meetup created, remove the edit area just leaving the links.
				document.body.removeChild(document.querySelector(".editArea"));
			}
		});
	});
	document.getElementById('cancelButt').addEventListener('click', function () {
		window.location.href = window.location.origin
	});

	dateTool.init(dateContainer, new Date());
}




document.onreadystatechange = function () {
	if (document.readyState === "complete") {
		init();
	}
};
