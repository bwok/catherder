function sendAjaxRequest(url, data, callback) {
	// TODO missing a lot of missing bits here. Hello IE11 :/
	var oReq = new XMLHttpRequest();
	oReq.addEventListener("load", function () {
		callback(JSON.parse(this.responseText), null);
	});
	oReq.open("POST", url);
	oReq.setRequestHeader("Content-Type", "application/json");
	oReq.send(data);
}