"use strict";

/**
 * Callback used by sendAjaxRequest.
 *
 * @callback ajaxCallBack
 * @param {?object} error
 * @param {?object} response
 */

/**
 *
 * @param {string} url
 * @param {string} data
 * @param {ajaxCallBack} callback
 */
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