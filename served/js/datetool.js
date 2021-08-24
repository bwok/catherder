"use strict";
/**
 * A scrollbox with selectable dates. Scrollable within the allowable javascript date ranges.
 */
var dateTool = new function(){
	var currDate, parentContainer, dateScrollCont;
	var selectedDates = [];		// UTC timestamps of selected dates, stored as numbers not strings.
	var numDateElements = 10;	// The number of visible date elements in the tool. Scrolls left and right by this many.
	var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
	var days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
	var minDate = 0;
	var maxDate = 8640000000000000; // A.D. 275760

	/**
	 * Initialises the date tool. Throws an exception if the input parameters are invalid.
	 * @param {HTMLElement} parentElement
	 * @param {Number} startDate
	 * @param {Array} dateArray    array of selected dates. unix millisecond timestamps.
	 */
	this.init = function(parentElement, startDate, dateArray){
		if(parentElement instanceof Node === false || document.contains(parentElement) === false){
			throw "dateTool init(): parentContainer is not a document node."
		} else{
			parentContainer = parentElement;
		}
		if(Number.isInteger(startDate) === false || startDate < minDate || startDate > maxDate){
			throw "dateTool init(): startDate is not a valid unix timestamp."
		}
		if(Array.isArray(dateArray) === false){
			throw "dateTool init(): dateArray is not an array."
		}

		currDate = startDate;
		selectedDates = dateArray;
		parentContainer = parentElement;
		createSkeleton();
	};

	/**
	 * Returns the selected dates as an array of timestamps (number type).
	 * @returns {Array.<number>}
	 */
	this.getDates = function(){
		return selectedDates;
	};

	/**
	 * Creates the tool html skeleton in the parentElement.
	 */
	function createSkeleton(){
		parentContainer.innerHTML = '<span class="scrollBox"><svg xmlns="http://www.w3.org/2000/svg" width="5" height="10" viewBox="0 0 5 10"><path d="M 0,5 5,10 5,0 Z"></path></svg></span>' +
			'<div class="dateScrollCont"></div>' +
			'<span class="scrollBox"><svg xmlns="http://www.w3.org/2000/svg" width="5" height="10" viewBox="0 0 5 10"><path d="M 0,0 0,10 5,5 Z"></path></svg></span>';

		dateScrollCont = parentContainer.querySelector(".dateScrollCont");
		var aScrollElems = parentContainer.querySelectorAll(".scrollBox");

		// scroll left
		aScrollElems[0].addEventListener("click", function(){
			scrollElements(-1);
		});
		// scroll right
		aScrollElems[1].addEventListener("click", function(){
			scrollElements(1);
		});

		makeElements(1);
	}

	/**
	 * Scrolls elements left or right. Each element
	 * @param {number} direction -1 scrolls left, 1 scrolls right
	 */
	function scrollElements(direction){
		var startDate = currDate;

		if(direction === -1 && minDate < startDate){
			makeElements(-1);
		} else if(direction === 1 && maxDate > startDate + (numDateElements * 86400000)){	// numDays*num milliseconds in a day
			makeElements(1)
		}
	}

	/**
	 * Makes the date elements for the tool
	 * @param {number }direction    -1 to prepend nodes, 1 to append them
	 */
	function makeElements(direction){

		if(dateScrollCont.children.length !== 0){
			currDate = currDate + (numDateElements * 86400000) * direction;
		}

		var startDate = currDate;
		dateScrollCont.innerHTML = '';

		for(var i = 0; i < numDateElements; i++){
			var parentSpan = document.createElement("span");
			parentSpan.classList.add("dateBox");
			parentSpan.setAttribute("data-date", startDate);

			// Highlight previously selected dates on scroll
			if(selectedDates.indexOf(startDate) >= 0){
				parentSpan.classList.add("selectedDate");
			}

			// On click add or remove date from the selectedDates array
			parentSpan.addEventListener("click", function(){
				var thisDate = parseInt(this.getAttribute("data-date"), 10);
				var arrIndex = selectedDates.indexOf(thisDate);

				if(arrIndex < 0){
					selectedDates.push(thisDate);
					this.classList.add("selectedDate");
				} else{
					selectedDates.splice(arrIndex, 1);
					this.classList.remove("selectedDate");
				}
				// Keep the date array sorted
				selectedDates.sort(function(a, b){
					return a - b;
				});
			});

			var dateObj = new Date(startDate);
			var monthSpan = document.createElement("span");
			monthSpan.textContent = months[dateObj.getMonth()];

			var dateSpan = document.createElement("span");
			dateSpan.classList.add("date");
			dateSpan.textContent = dateObj.getDate().toString(10);

			var daySpan = document.createElement("span");
			daySpan.textContent = days[dateObj.getDay()];

			parentSpan.appendChild(monthSpan);
			parentSpan.appendChild(dateSpan);
			parentSpan.appendChild(daySpan);
			dateScrollCont.appendChild(parentSpan);

			startDate = startDate + 86400000;
		}
	}
};