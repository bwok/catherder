"use strict";

var editObj = new function(){
	var errorArea, dateContainer, adminhash, descrElem, emailElem, notifyElem;

	this.init = function(){
		errorArea = document.getElementById('errorArea');
		dateContainer = document.getElementById('dateContainer');
		descrElem = document.getElementById("description");
		emailElem = document.getElementById("adminEmail");
		notifyElem = document.getElementById("notifications");

		var params = new URLSearchParams(window.location.search.substring(1));
		adminhash = params.get("id");

		if(adminhash === null){
			var startDate = new Date();
			startDate.setHours(0, 0, 0, 0);
			dateTool.init(dateContainer, startDate.valueOf(), []);
		} else{
			getMeetUp();
		}

		document.getElementById('saveButt').addEventListener('click', function(){
			saveMeetUp();
		});

		document.getElementById('cancelButt').addEventListener('click', function(){
			window.location.href = window.location.origin
		});

	};

	/**
	 * Shows an error message.
	 * @param {string} errMsg
	 */
	function showError(errMsg){
		errorArea.textContent = errMsg;
		errorArea.classList.remove("hidden");
	}

	/**
	 * Clears the error message
	 */
	function clearError(){
		errorArea.textContent = "";
		errorArea.classList.add("hidden");
	}

	/**
	 * Gets the admin version of the meetup data.
	 */
	function getMeetUp(){
		sendAjaxRequest("/api/getadminmeetup", JSON.stringify({adminhash: adminhash}), function(error, response){
			if(error !== null){
				showError(error.toString());
			} else if(response.error !== ""){
				showError(response.error);
			} else{
				emailElem.value = response.result.adminemail;
				notifyElem.checked = response.result.sendalerts;
				descrElem.value = response.result.description;

				if(response.result.dates.length === 0){
					var startDate = new Date();
					startDate.setHours(0, 0, 0, 0);
					dateTool.init(dateContainer, startDate.valueOf(), response.result.dates);
				} else{
					dateTool.init(dateContainer, response.result.dates[0], response.result.dates);
				}
			}
		});
	}

	/**
	 * Saves the meetup values
	 */
	function saveMeetUp(){
		clearError();

		var args = {
			adminhash: adminhash,
			description: descrElem.value,
			dates: dateTool.getDates(),
			users: [],
			adminemail: emailElem.value,
			sendalerts: notifyElem.checked
		};

		sendAjaxRequest("/api/updatemeetup", JSON.stringify(args), function(error, response){
			if(error !== null){
				showError(error.toString());
			} else if(response.error !== ""){
				showError(response.error);
			} else{
				var link = document.getElementById("userLink");
				link.href = window.location.origin + "/view?id=" + encodeURIComponent(response.result.userhash);
				link.textContent = window.location.origin + "/view?id=" + encodeURIComponent(response.result.userhash);

				link = document.getElementById("adminLink");
				link.href = window.location.origin + "/edit?id=" + encodeURIComponent(response.result.adminhash);
				link.textContent = window.location.origin + "/edit?id=" + encodeURIComponent(response.result.adminhash);

				document.getElementById("linkArea").classList.remove("hidden");

				// meetup created, remove the edit area just leaving the links.
				document.body.removeChild(document.querySelector(".editArea"));
			}
		});
	}

};


document.onreadystatechange = function(){
	if(document.readyState === "complete"){
		editObj.init();
	}
};
