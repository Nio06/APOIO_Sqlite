//*************************************************
// UTIL.JS: GENERAL AND LOCAL UTILITIES
//*************************************************
let util = {};

//-------------------------------------------------

util.makePayload = function()
//--- Variable number of items appended with \x1F as a separator
{ let txt = "";

if (arguments.length === 0)
	return txt;

for (let i = 0; i < arguments.length; i++)
	{
	txt += arguments[i];
	if ((i + 1) < arguments.length)
		txt += "\x1F";
	}
return txt;
}

//-------------------------------------------------

util.appendLinesSep = function(str)
//--- str is a multiline.  This returns str after replacing '\n' with \x1F
//		every place where '\n' appears.

{ let ary = str.split('\n'),
			txt = "";

for (let i = 0; i < ary.length; i++)
	{
	txt += "\x1F" + ary[i];
	}
return txt;
}

//-------------------------------------------------

util.Fetch = async function (url, payload, theJob) {
	const formdata = new FormData();
	formdata.append("accAction", payload);

//--- DO THE ACTUAL FETCH
fetch(url, {method: "POST", body: formdata})
.then(response => {
  if (response.status !== 200) {
    throw new Error(`Response.status was not 200, it was ${response.status}`);
  }
  return response.text();
})
.then(rtnData => {
	theJob(rtnData);
  })
.catch(e => {
  console.log('There has been a problem with your fetch operation: ' + e.message);
});
}

