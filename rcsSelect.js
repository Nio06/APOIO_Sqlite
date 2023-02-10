//***********************************
//--- rcsSelect.js

// Short calls for selector functions
//***********************************



//*******************************************************
// Select items using *querySelector or *querySelectorAll
//*******************************************************

//-----------------------
//	document.querySelector
//-----------------------
$ = function(selectors)
{
return document.querySelector(selectors);
}

//--------------------------
//	document.querySelectorAll
//--------------------------
$_ = function(selectors)
{
return document.querySelectorAll(selectors);
}

//----------------------
// 	element.querySelector

//	eselect: selector for the element
//----------------------
$$ = function(elm, selectors)
{
return $(elm).querySelector(selectors);
}

//-------------------------
//	element.querySelectorAll

//	eselect: selector for the element
//-------------------------
$$_ = function(elm, selectors)
{
return $(elm).querySelectorAll(selectors);
}

