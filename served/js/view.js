"use strict";


var viewObj = new function(){
	var errorArea, userhash, columnCont;
	var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
	var days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];


	/**
	 * Initialise any bits that need initialising.
	 */
	this.init = function(){
		var params = new URLSearchParams(window.location.search.substring(1));
		userhash = params.get("id");
		errorArea = document.getElementById('errorArea');
		columnCont = document.querySelector(".columnsContainer");

		if(userhash === null){
			showError("No id argument was found in the URL.");
		} else {
			document.querySelector(".shareLink").textContent = window.location.origin + "/view?id=" + encodeURIComponent(userhash);
		}



		document.getElementById("saveButt").addEventListener("click", function(){
			addUser();
		});

		refreshDateGrid();
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
	 * Gets the username and checked checkboxes, and requests the backend add the user.
	 */
	function addUser(){
		clearError();

		var userName = document.querySelector(".username").value;
		var checkedDates = document.querySelectorAll(".newuser:checked");
		var dates = [];
		for(var i = 0; i < checkedDates.length; i++){
			dates.push(parseInt(checkedDates[i].name, 10))
		}

		var args = {
			username: userName,
			userhash: userhash,
			dates: dates
		};

		sendAjaxRequest("/api/updateuser", JSON.stringify(args), function(error, response){
			if(error !== null){
				showError(error.toString());
			} else if(response.error !== ""){
				showError(response.error);
			} else{
				refreshDateGrid();
			}
		});
	}

	/**
	 * Clears and redraws the date grid.
	 */
	function refreshDateGrid(){
		columnCont.innerHTML = '<div class="nameColumn"><div class="dummyBox"></div></div>';		// Reset container on each refresh

		sendAjaxRequest("/api/getusermeetup", JSON.stringify({userhash: userhash}), function(error, response){
			if(error !== null){
				showError(error.toString());
			} else if(response.error !== ""){
				showError(response.error);
			} else{
				clearError();
				var i;
				var usersArray = response.result.users;

				/*
				Create users column
				 */
				var nameColumn = columnCont.querySelector(".nameColumn");
				for(i = 0; i < usersArray.length; i++){
					var userDiv = document.createElement("div");
					userDiv.classList.add("row");
					userDiv.textContent = usersArray[i].name;
					nameColumn.appendChild(userDiv);
				}
				nameColumn.insertAdjacentHTML("beforeend", '<div class="row"><input class="username" type="text" name="username" placeholder="New user..."></div>');

				/*
				Create date columns
				 */
				var datesArray = response.result.dates;

				for(i = 0; i < datesArray.length; i++){
					var dateColumn = document.createElement("div");
					dateColumn.classList.add("dateColumn");

					// Generate the Month/date/dayofweek header
					dateColumn.innerHTML = '<div class="dateBox"><span></span><span class="date"></span><span></span></div>';
					var date = new Date(datesArray[i]);
					var spans = dateColumn.querySelectorAll("span");
					spans[0].textContent = months[date.getMonth()];
					spans[1].textContent = date.getDate().toString(10);
					spans[2].textContent = days[date.getDay()];


					// Generate existing users checkbox rows
					for(var usrIndex = 0; usrIndex < usersArray.length; usrIndex++){
						var row = document.createElement("div");
						row.classList.add("row");
						row.classList.add("rowUnavailable");

						var checkbox = document.createElement("input");
						checkbox.type = "checkbox";
						checkbox.disabled = true;
						checkbox.checked = false;

						for(var dateIndex = 0; dateIndex < usersArray[usrIndex].dates.length; dateIndex++){
							if(datesArray[i] === usersArray[usrIndex].dates[dateIndex]){
								checkbox.checked = true;
								row.classList.remove("rowUnavailable");
								row.classList.add("rowAvailable");
								break;
							}
						}

						row.appendChild(checkbox);
						dateColumn.appendChild(row);
					}

					dateColumn.insertAdjacentHTML("beforeend", '<div class="row"><input type="checkbox" class="newuser" name="' + datesArray[i] + '"></div></div>');
					columnCont.appendChild(dateColumn);
				}
			}
		});
	}
};


document.onreadystatechange = function(){
	if(document.readyState === "complete"){
		viewObj.init();
	}
};
