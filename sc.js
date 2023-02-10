//*************************************************
// SC.JS: CODE FOR MANAGING SHARED CATEGORIES
//*************************************************

let sc = {};


//----------------------------------------------------
// sc.showOneScreen (display none for all, then display one
//									showMe is the one to show.
//									# should be the first character
//----------------------------------------------------

let screens = [	"#welcome",
								"#categories",
								"#tasks",
								"#timerpage",
								"#taskentry",
								"#recvScr",
								"#shareScr"];

sc.showOneScreen = function(showMe)	{

for (let i = 0; i < screens.length; i++)	{
	$(screens[i]).style.display = "none";
	}

$(showMe).style.display = "block";
}



//----------------------------------------------------
// Show or hide the main options bar
//----------------------------------------------------

sc.hideMainOpts = function() {$("#optionsBar").style.display="none";}
sc.showMainOpts = function() {$("#optionsBar").style.display="block";}


//*************************************************
// FETCH CALLBACKS
//*************************************************
sc.shareCatFCB = function(data)	{

let myExpCats = [];

if (data !== "")	//: No data means no categories for this user but here,
									//	there is data to parse
	{
	let myCats = JSON.parse(data),
			textList = "";

	for(var i in json_obj) {
            let catname = json_obj[i].Catname;
            let object = Object.create(category);
            object.catName = catname;
            a.push(object);
            catArray.push(a);
            catname = "";
            a = [];
          }

$("#sharetestView").innerText = "";

	for (let i = 0; i < myCats.length; i++)
		{
		console.log(myCats[i]);
		$("#sharetestView").innerText += myCats[i] + "<br>";
		}
	}
}



//*************************************************
// EVENT LISTENER FUNCTIONS
//*************************************************

sc.shareCat = function()	{

sc.showOneScreen("#shareScr");
let payload = util.makePayload("Msideapp#1", username, password, "catab", token);
console.log("Payload: " + payload);
util.Fetch(url, payload, sc.shareCatFCB);
}

//*************************************************

sc.toExportScr = function()	{
$("#expScr").style.display = "none";
sc.hideMainOpts();
$("#addExpScr").style.display = "block";
}

//*************************************************

sc.createNewExpCat = function()	{
let newCatName = $("#newExpName").value,
		exportList = $("#exportList").value,
		payload;


//--- BUILD THE PAYLOAD & FETCH
payload = util.makePayload("Msideapp#1", username, newCatName, "addShareCat", token);
payload += util.appendLinesSep(exportList);
util.Fetch(url, payload, sc.createNewExpCatFCB);


alert("Payload: " + payload);
}

//*************************************************
// EVENT LISTENERS
//*************************************************

//: Get to the Export Categories screen from the topline menu
$("#shareButton").addEventListener("click", sc.shareCat);

