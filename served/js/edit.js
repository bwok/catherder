"use strict";

function init () {
	var errorArea = document.getElementById('errorArea');
	var dateContainer = document.getElementById('dateContainer');

	document.getElementById('saveButt').addEventListener('click', function () {
		errorArea.textContent = "";
		errorArea.classList.add("hidden");

		var args = {
			description: document.getElementById("description").value,
			dates: dateTool.getDates(),
			users: [],
			adminemail: document.getElementById("adminEmail").value,
			sendalerts: document.getElementById("notifications").checked

		};

		sendAjaxRequest("/api/updatemeetup", JSON.stringify(args), function (error, response) {
			if(error !== null){
				errorArea.textContent = error;
				errorArea.classList.remove("hidden");
			}
			else if(response.error !== ""){
				errorArea.textContent = response.error;
				errorArea.classList.remove("hidden");
			} else {
				var link = document.getElementById("userLink");
				link.href = window.location.origin + "/view?id=" + encodeURIComponent(response.result.userhash);
				link.textContent = window.location.origin + "/view?id=" + encodeURIComponent(response.result.userhash);

				link = document.getElementById("adminLink");
				link.href = window.location.origin + "/admin?id=" + encodeURIComponent(response.result.adminhash);
				link.textContent = window.location.origin + "/admin?id=" + encodeURIComponent(response.result.adminhash);

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
