function init() {
	var errorArea = document.getElementById('errorArea');

	// Add click handlers to all date boxes.
	var dateBoxes = document.querySelectorAll(".dateBox");
	for(var i=0; i < dateBoxes.length; i++){
		dateBoxes[i].addEventListener("click", function(){
			if(this.classList.contains("selectedDate")){
				this.classList.remove("selectedDate");
			} else {
				this.classList.add("selectedDate");
			}
		});
	}

	document.getElementById("saveButt").addEventListener("click", function(){
		var dateBoxes = document.getElementById("dateContainer").querySelectorAll(".selectedDate");
		var dates = [];

		for(var i=0; i < dateBoxes.length; i++){
			dates.push({date: parseInt(dateBoxes[i].getAttribute("data-date"), 10), users: []});
		}

		var args = {
			description: document.getElementById("description").value,
			dates: dates,
			admin: {
				email: document.getElementById("adminEmail").value,
				alerts: document.getElementById("notifications").checked
			}
		};

		var currentUrl = new URL(window.location.href);
		var url = "adminsave?id=" + currentUrl.searchParams.get("id");

		sendAjaxRequest(url, JSON.stringify(args), function (response, error) {
			if(error !== null){
				errorArea.textContent = error;
				errorArea.classList.remove("hidden");
			}
			else if(response.error !== ""){
				errorArea.textContent = response.error;
				errorArea.classList.remove("hidden");
			} else {
				window.location.href = window.location.origin + "/view?id=" + encodeURIComponent(response.result);
			}
		});
	});

	document.getElementById("cancelButt").addEventListener("click", function(){
		window.location.href = window.location.origin + "/view?id=" + encodeURIComponent(document.getElementById("userHash").value);
	});

	document.getElementById("deleteButt").addEventListener("click", function(){
		if(confirm("Are you sure?") === false){
			return
		}

		var currentUrl = new URL(window.location.href);
		var url = "admindelete?id=" + currentUrl.searchParams.get("id");

		sendAjaxRequest(url, JSON.stringify({}), function (response, error) {
			if(error !== null){
				errorArea.textContent = error;
				errorArea.classList.remove("hidden");
			}
			else if(response.error !== ""){
				errorArea.textContent = response.error;
				errorArea.classList.remove("hidden");
			} else {
				window.location.href = window.location.origin;
			}
		});
	});
}


document.onreadystatechange = function () {
	if (document.readyState === "complete") {
		init();
	}
};
