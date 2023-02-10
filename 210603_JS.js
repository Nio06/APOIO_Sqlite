//****************************************
//--- Video & Resource Package specific JS
//****************************************

let surp = {};

let taskTableRows = 2;

//const url = "127.0.0.1:5558";
var username = "";
var password = "";
var token = "";
var catArrayReceived = false;
var catArray = new Array; // categories Array gets saved once and doesn't have to be sent anymore after first time
var category = {
  catName:"",
  taskName:"",
  status:"",
  dueDate:"",
  spredicted:"",
  upredicted:"",
  current:""
};
var paused = false;
var clickedID = new Array;
var running = false;

//*******************************************************
// Properties
//*******************************************************

//*******************************************************
// Non-Event-Triggered Methods
//*******************************************************

//*******************************************************
// Event-Triggered Methods
//*******************************************************

surp.openCategory = function (event) {
  let cat = event.target.parentNode;
  let catID = cat.id; // gets id of clicked category
  let catNum = event.target.innerHTML;
  if (catID[0] == "c") {
    let tid = "";
    for (let i = 1; i < catID.length; i++){
      tid = tid + catID[i];
    }
    catID = "";
    if (catArray[parseInt(tid)].length < 2) {
      catID = catArray[parseInt(tid)];
      clickedID = [];
      clickedID.push(parseInt(tid));
      let payload =
        "Msideapp#1" + "\x1F" + username + "\x1F" + catNum + "\x1F" + "catsk" + "\x1F" + token; // pass through category name instead of password
      const formdata = new FormData();
      formdata.append("accAction", payload);
      //	const formdata = new FormData(accountForm);
      fetch(url, {
        method: "POST",
        body: formdata,
      })
        .then((response) => response.text())
        .then((data) => {
          console.log(data);
          let json_obj = JSON.parse(data);
          var taskArray = new Array;
          taskArray = catArray[parseInt(tid)];
          for(var i in json_obj) {
            var object = Object.create(category);

            let taskName = json_obj[i].Taskname;
            let catName = json_obj[i].Catname;
            let dueDate = json_obj[i].Duedate;
            let status = json_obj[i].Status;
            if (status == "0") {
              status = "Incomplete";
            }
            let spredicted = json_obj[i].Spredicted;
            let upredicted = json_obj[i].Upredicted;
            let current = json_obj[i].Current;
            object.taskName = taskName;
            object.catName = catName;
            object.dueDate = dueDate;
            object.status = status;
            object.spredicted = spredicted;
            object.upredicted = upredicted;
            object.current = current;
            taskArray.push(object);
          }

          if (catArray[parseInt(tid)][0].catName == category.catName) {
            catArray[parseInt(tid)] = [];
            catArray[parseInt(tid)] = taskArray;
          }
          taskArray = [];
          surp.fillTaskTable(parseInt(tid));
        })
        .catch((error) => console.error(error));
        formdata.delete("accAction");
    } else {
      console.log("was already open once before");
      surp.fillTaskTable(parseInt(tid));
    }
    $("#categories").style.display = "none";
    $("#tasks").style.display = "block";
    sc.showOneScreen("#categories");
  }
}

surp.deleteTaskTable = function () {
  if ($("#taskCatTitle").innerHTML != "") {
    let rowNr = $("#taskTable").rows.length;
    for (var x = rowNr - 1; x > 0; x--) {
        $("#taskTable").deleteRow(x);
        let tempID = "#pt" + x.toString();
        if ($(tempID) != null){
          $(tempID).remove();
        }
    }
    $("#taskCatTitle").innerHTML = "";
  }
  //$("#taskTable").remove();
}

surp.fillTaskTable = function (tid) {
  surp.deleteTaskTable();
  sc.showOneScreen("#tasks");
  console.log("now tasktable should be visible");
  $("#categories").style.display = "none";
  var table = $("#taskTable");
  $("#taskCatTitle").innerHTML = catArray[tid][0].catName;
	for (let i = 1; i < catArray[tid].length; i++) {
			// create new table row with empty cells
			var row = table.insertRow(i);
			var c0 = row.insertCell(0);
			var c1 = row.insertCell(1);
      var c2 = row.insertCell(2);

			// fill first cell with task name
			c0.innerHTML = catArray[tid][i].taskName;
			c1.innerHTML = catArray[tid][i].dueDate;
      c2.innerHTML = catArray[tid][i].status;
      var node = document.createElement("P");             // Create a <p> node
      node.id = "pt" + i;
      var textnode = document.createTextNode(catArray[tid][i].taskName);         // Create a text node
      node.appendChild(textnode);                              // Append the text to <li>
      $("#dropTask").appendChild(node);     // Append <p> to <div> with id="dropTask"
      c0.style.backgroundColor = "white";
      c1.style.backgroundColor = "white";
      c2.style.backgroundColor = "white";
      c0.style.color = "black";
      c1.style.color = "black";
      c2.style.color = "black";
			row.id = "tt" + i; // gives current row in current Table an id for click eventListener later
      row.class = "A";
	}
  $("#dropTask").style.maxHeight = "300px";
  $("#dropTask").style.overflow = "auto";
  console.log($("#taskTable"));
}

surp.submitTime = function () {
  let currentTimeString = watch.innerHTML;
  let currentTime = surp.timerGetInt(currentTimeString);
  $("#start").style.display = "inline";
  $("#pause").style.display = "none";

  clearInterval(timer);
  $("#tasks").style.display = "block";
  $("#timerpage").style.display = "none";
  $("#optionsBar").style.display = "block";

  millisecound = 0;

  let dateTimer = new Date(millisecound);

  watch.innerHTML =
    ("0" + dateTimer.getUTCHours()).slice(-2) +
    ":" +
    ("0" + dateTimer.getUTCMinutes()).slice(-2) +
    ":" +
    ("0" + dateTimer.getUTCSeconds()).slice(-2) +
    ":" +
    ("0" + dateTimer.getUTCMilliseconds()).slice(-3, -1);

  watch.innerHTML = "00:00:00:00";
  let payload =
      "Msideapp#1" + "\x1F" + username + "\x1F" + $("#stopwatchCat").innerHTML + "\x1F" + "done" + "\x1F" + token + "\x1F" + $("#stopwatchTask").innerHTML + "\x1F" + currentTime.toString();
  const formdata = new FormData();
  formdata.append("accAction", payload);
  //	const formdata = new FormData(accountForm);
  fetch(url, {
      method: "POST",
      body: formdata,
  })
  .then((response) => response.text())
  .then((data) => {
        console.log(data);
  })
  .catch((error) => console.error(error));
  formdata.delete("accAction");
  catArray[clickedID[0]][clickedID[1]].status = "Complete";
  catArray[clickedID[0]][clickedID[1]].current = currentTime.toString();
  let row = $("#taskTable").rows[clickedID[1]];
  row.cells[2].innerHTML = "Complete";
  let tempVar = clickedID[0];
  clickedID = [];
  clickedID[0] = tempVar;
  running = false;
}

surp.submitTask = function () {
  /*var table = $("#taskTable");
  // create new table row with empty cells
  var row = table.insertRow(1);
  var c0 = row.insertCell(0);
  var c1 = row.insertCell(1);
  var c2 = row.insertCell(2);

  row.class = "A";
  row.id = "tt" + taskTableRows;
  taskTableRows++;

  c0.innerHTML = $("#taskNameInput").value;
  c1.innerHTML = $("#taskDateInput").value;
  c2.innerHTML = "Incomplete";*/
  let catName = "";
  let newCat = false;
  if ($("#taskCatInputSet").style.display != "none"){
    catName = $("#taskCatInputSet").innerHTML;
  } else {
    catName = $("#taskCatInput").value;
    newCat = true;
  }
  let obj = Object.create(category);

  let taskName = $("#taskNameInput").value;
  let dueDate = $("#taskDateInput").value;
  let spredicted = "0";
  let upredicted = ($("#predMinsInp").value * 60 + $("#predHoursInp").value * 3600);
  let current = "0";
  let status = "Incomplete";
  obj.taskName = taskName;
  obj.catName = catName;
  obj.dueDate = dueDate;
  obj.status = status;
  obj.spredicted = spredicted;
  obj.upredicted = upredicted;
  obj.current = current;

  if (!newCat) {
    catArray[clickedID[0]].push(obj);
  } else {
    let tempArray = new Array;
    tempArray.push(obj);
    catArray.push(tempArray);
    tempArray = [];
  }

  let payload =
      "Msideapp#1" + "\x1F" + username + "\x1F" + taskName + "\x1F"
      + "addTask" + "\x1F" + token + "\x1F" + catName + "\x1F" + dueDate
      + "\x1F" + status + "\x1F" + spredicted + "\x1F" + upredicted + "\x1F" + current;
  const formdata = new FormData();
  formdata.append("accAction", payload);
  //	const formdata = new FormData(accountForm);
  fetch(url, {
      method: "POST",
      body: formdata,
  })
  .then((response) => response.text())
  .then((data) => {
        console.log(data);
  })
  .catch((error) => console.error(error));
  formdata.delete("accAction");

  if (!newCat){
    surp.deleteTaskTable();
    surp.fillTaskTable(clickedID[0]);
    $("#tasks").style.display = "block";
  } else {
    surp.deleteCatTable();
    surp.fillCatTable();
    $("#categories").style.display = "block";
  }

  $("#taskentry").style.display = "none";
  $("#optionsBar").style.display = "block";
  $("#taskNameInput").value = "";
  $("#predMinsInp").value = "";
  $("#predHoursInp").value = "";
  $("#taskCatInput").value = "";
  const d = new Date();
  let day = d.getUTCDate();
  var dateDay = "";
  if (day < 10) {
    dateDay = "0" + day.toString();
  } else {
    dateDay = day.toString();
  }
  let month = d.getUTCMonth() + 1;
  let year = d.getUTCFullYear();
  if (month < 10) {
    $("#taskDateInput").value = year.toString() + "-0" + month.toString() + "-" + dateDay;
  } else {
    $("#taskDateInput").value = year.toString() + "-" + month.toString() + "-" + dateDay;
  }
  $("#predSpan").innerHTML = "";
}

surp.filter = function () {
  $("#filter").options[0].disabled = true;
  /*
  var e = $("#sort").value;
  var arr = new Array;
  var arr1 = new Array;

  if (e == "nfa") { // First Name Ascending
      let n = students.length;
      for (let i = 1; i < n; i++) {
          // Choosing the first element in our unsorted subarray
          arr = students[i];
          // The last element of our sorted subarray
          let j = i - 1;
          arr1 = students[j];
          while ((j > -1) && (arr[1] < arr1[1])) {
              students[j + 1] = students[j];
              j--;
              arr1 = students[j];
          }
          students[j + 1] = arr;
      }*/
    var e = $("#filter").value;
    let tempArray1 = new Array;
    let tempArray2 = new Array;
    let n = catArray.length;
    if (e == "1") {
      //name Ascending
      for (let i = 1; i < n; i++) {
        // Choosing the first element in our unsorted subarray
        tempArray1 = catArray[i];
        // The last element of our sorted subarray
        let j = i - 1;
        tempArray2 = catArray[j];
        while ((j > -1) && (tempArray1[0].catName < tempArray2[0].catName)) {
            catArray[j + 1] = catArray[j];
            j--;
            tempArray2 = catArray[j];
        }
        catArray[j + 1] = tempArray1;
      }
    } else if (e == "2") {
      //name Ascending
      for (let i = 1; i < n; i++) {
        // Choosing the first element in our unsorted subarray
        tempArray1 = catArray[i];
        // The last element of our sorted subarray
        let j = i - 1;
        tempArray2 = catArray[j];
        while ((j > -1) && (tempArray1[0].catName > tempArray2[0].catName)) {
            catArray[j + 1] = catArray[j];
            j--;
            tempArray2 = catArray[j];
        }
        catArray[j + 1] = tempArray1;
      }
    }
    surp.deleteCatTable();
    surp.fillCatTable();
}

surp.openStopwatch = function (event) {
  let tr = event.target.parentNode;
  let tid = tr.id;
  let temporaryTid = tid[0] + tid[1];
  if (temporaryTid == "tt") {
    var tempID = "";
    for (let i = 2; i < tid.length; i++){
      tempID = tempID + tid[i];
    }
    var i = 0;
    while (i < catArray.length){
      if (catArray[i][0].catName == $("#taskCatTitle").innerHTML){
        break;
      } else {
        i++;
      }
    }
    //clickedID.push(i);
    clickedID.push(parseInt(tempID));
    if (($("#stopwatchTask").innerHTML == catArray[i][clickedID[1]].taskName && $("#stopwatchCat").innerHTML == catArray[i][clickedID[1]].catName) ||
                          ($("#stopwatch").innerHTML == "00:00:00:00")){
      if (catArray[i][clickedID[1]].status == "Paused" || catArray[i][clickedID[1]].status == "Complete"){
        let curTime = parseInt(catArray[clickedID[0]][clickedID[1]].current);
        let h = 0;
        let min = 0;
        let sec = "";
        if (curTime > 3600) {
          h = curTime / 3600;
        } else {
          h = 0;
        }
        let hours = "";
        if (h < 10){
          hours = "0" + h.toString();
        } else {
          hours = h.toString();
        }
        if ((curTime - h * 3600) > 60) {
          min = (curTime - h * 3600) / 60;
        } else {
          min = 0;
        }
        let mins = "";
        if (min < 10){
          mins = "0" + min.toString();
        } else {
          mins = min.toString();
        }
        if ((curTime - h * 3600 - min * 60) < 10) {
          sec = (curTime - h * 3600 - min * 60).toString();
          sec = "0" + sec;
        } else {
          sec = (curTime - h * 3600 - min * 60).toString();
        }
        $("#stopwatch").innerHTML = hours + ":" + mins + ":" + sec + ":" + "00";
        if (catArray[i][clickedID[1]].status == "Complete"){
          $("#start").style.display = "none";
          $("#pause").style.display = "none";
          $("#timesubmit").style.display = "none";
        }
      } else {
        $("#stopwatch").innerHTML = "00" + ":" + "00" + ":" + "00" + ":" + "00";
      }
      $("#tasks").style.display = "none";
      $("#timerpage").style.display = "block";
      $("#stopwatchTask").innerHTML = catArray[i][clickedID[1]].taskName;
      $("#stopwatchCat").innerHTML = catArray[i][clickedID[1]].catName;
      $("#stopwatchPTime").innerHTML = catArray[i][clickedID[1]].upredicted;
      $("#optionsBar").style.display = "none";
    } else {
      alert("Another Task is still running!");
    }
  }
}

surp.openTaskEntry = function () {
  $("#tasks").style.display = "none";
  $("#categories").style.display = "none";
  $("#taskentry").style.display = "block";
  $("#optionsBar").style.display = "none";
  if (clickedID.length > 0){
    $("#taskCatInputSet").style.display = "inline";
    $("#taskCatInputSet").innerHTML = "";
    $("#taskCatInputSet").innerHTML = catArray[clickedID[0]][0].catName;
    $("#taskCatInput").style.display = "none";
  } else {
    $("#taskCatInput").style.display = "inline";
    $("#taskCatInputSet").style.display = "none";
  }
}

surp.timeStart = function () {
  $("#start").style.display = "none";
  $("#pause").style.display = "inline";
  watch.style.color = "#0f62fe";
  clearInterval(timer);
  if (catArray[clickedID[0]][clickedID[1]].status != "Paused") {
    catArray[clickedID[0]][clickedID[1]].current = "0";
    millisecound = 0;
  } else {
    millisecound = parseInt(catArray[clickedID[0]][clickedID[1]].current);
    millisecound = millisecound * 1000;
  }
  timer = setInterval(() => {
    millisecound += 10;

    let dateTimer = new Date(millisecound);

    watch.innerHTML =
      ("0" + dateTimer.getUTCHours()).slice(-2) +
      ":" +
      ("0" + dateTimer.getUTCMinutes()).slice(-2) +
      ":" +
      ("0" + dateTimer.getUTCSeconds()).slice(-2) +
      ":" +
      ("0" + dateTimer.getUTCMilliseconds()).slice(-3, -1);
  }, 10);
  paused = false;
  running = true;
}

surp.timePause = function () {
  let currentTimeString = watch.innerHTML;
  let currentTime = surp.timerGetInt(currentTimeString);
  $("#start").style.display = "inline";
  $("#pause").style.display = "none";
  watch.style.color = "red";
  paused = true;
  clearInterval(timer);
  let payload =
      "Msideapp#1" + "\x1F" + username + "\x1F" + $("#stopwatchCat").innerHTML + "\x1F" + "pause" + "\x1F" + token + "\x1F" + $("#stopwatchTask").innerHTML + "\x1F" + currentTime.toString();
  const formdata = new FormData();
  formdata.append("accAction", payload);
  //	const formdata = new FormData(accountForm);
  fetch(url, {
      method: "POST",
      body: formdata,
  })
  .then((response) => response.text())
  .then((data) => {
        console.log(data);
  })
  .catch((error) => console.error(error));
  formdata.delete("accAction");
  catArray[clickedID[0]][clickedID[1]].status = "Paused";
  catArray[clickedID[0]][clickedID[1]].current = currentTime.toString();
  let row = $("#taskTable").rows[clickedID[1]];
  row.cells[2].innerHTML = "Paused";
  running = false;
}

surp.timerGetInt = function (currentTimeString) {
  let h = parseInt(currentTimeString[0] + currentTimeString[1]);
  let min = parseInt(currentTimeString[3] + currentTimeString[4]);
  let sec = parseInt(currentTimeString[6] + currentTimeString[7]);
  if (h != 0) {
    h = h * 3600;
  }
  if (min != 0) {
    min = min * 60;
  }
  return h + min + sec;
}

surp.taskStartedTime = function () {
  $("#tasks").style.display = "block";
  $("#timerpage").style.display = "none";
  $("#optionsBar").style.display = "block";
  if (!running) {
    millisecound = 0;

    let dateTimer = new Date(millisecound);

    watch.innerHTML =
      ("0" + dateTimer.getUTCHours()).slice(-2) +
      ":" +
      ("0" + dateTimer.getUTCMinutes()).slice(-2) +
      ":" +
      ("0" + dateTimer.getUTCSeconds()).slice(-2) +
      ":" +
      ("0" + dateTimer.getUTCMilliseconds()).slice(-3, -1);

    watch.innerHTML = "00:00:00:00";
    $("#stopwatchTask").innerHTML = "";
    $("#stopwatchCat").innerHTML = "";
    $("#stopwatchPTime").innerHTML = "";
    $("#start").style.display = "inline";
    $("#pause").style.display = "none";
    $("#timesubmit").style.display = "inline";
  } else {
    let row = $("#taskTable").rows[clickedID[1]];
    row.cells[2].innerHTML = "in work";
  }
  let tempVar = clickedID[0];
  clickedID = [];
  clickedID[0] = tempVar;
}

surp.accountChangePassword = function () {
  var newPwd = prompt("Enter your new password:", "New Password");
  if (newPwd == null) {
    return;
  }
  var newPwd2 = prompt("Confirm new password:", "New Password");
  if (newPwd2 == null) {
    return;
  }
  if (newPwd != newPwd2) {
    alert("Passwords don't match");
    return;
  }
  let payload =
    "Msideapp#1" + "\x1F" + username + "\x1F" + newPwd + "\x1F" + "chpwd" + "\x1F" + token;
  const formdata = new FormData();
  formdata.append("accAction", payload);
  //	const formdata = new FormData(accountForm);
  fetch(url, {
    method: "POST",
    body: formdata,
  })
    .then((response) => response.text())
    .then((data) => {
      console.log(data);
        if (data == 1) {
        alert("Password has been changed!");
        return;
      }
      else {
        alert("Session has timed out, please login again!");
        $("#login").style.display = "block";
        $("#categories").style.display = "none";
        $("#tasks").style.display = "none";
        $("#timerpage").style.display = "none";
        $("#taskentry").style.display = "none";
        $("#welcome").style.display = "none";
        $("#optionsBar").style.display = "none";
        $("#usernameLogin").value = "";
        $("#passwordLogin").value = "";
        return;
      }
    })
    .catch((error) => console.error(error));
  formdata.delete("accAction");
}

surp.accountLogout = function () {
  if (confirm("Logging out")){
    let payload =
      "Msideapp#1" + "\x1F" + username + "\x1F" + password + "\x1F" + "lgout";
    const formdata = new FormData();
    formdata.append("accAction", payload);
    //	const formdata = new FormData(accountForm);
    fetch(url, {
      method: "POST",
      body: formdata,
    })
      .then((response) => response.text())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => console.error(error));
    formdata.delete("accAction");
    $("#login").style.display = "block";
    $("#categories").style.display = "none";
    $("#tasks").style.display = "none";
    $("#timerpage").style.display = "none";
    $("#taskentry").style.display = "none";
    $("#welcome").style.display = "none";
    $("#optionsBar").style.display = "none";
    $("#usernameLogin").value = "";
    $("#passwordLogin").value = "";
    surp.deleteCatTable();
    surp.deleteTaskTable();
    username = "";
    password = "";
    token = "";
    catArrayReceived = false;
    catArray = new Array; // categories Array gets saved once and doesn't have to be sent anymore after first time
    category = {
      catName:"",
      taskName:"",
      status:"",
      dueDate:"",
      spredicted:"",
      upredicted:"",
      current:""
    };
    paused = false;
    clickedID = [];
    running = false;
  }
}

surp.accountLogin = function () {
  username = $("#usernameLogin").value;
  password = $("#passwordLogin").value;
  const d = new Date();
  let day = d.getUTCDate();
  var dateDay = "";
  if (day < 10) {
    dateDay = "0" + day.toString();
  } else {
    dateDay = day.toString();
  }
  let month = d.getUTCMonth() + 1;
  let year = d.getUTCFullYear();
  if (month < 10) {
    $("#taskDateInput").value = year.toString() + "-0" + month.toString() + "-" + dateDay;
  } else {
    $("#taskDateInput").value = year.toString() + "-" + month.toString() + "-" + dateDay;
  }
  if (username == "" && password == "") {
    $("#usernameLogin").style.background = "red";
    $("#passwordLogin").style.background = "red";
  } else if (username == "") {
    $("#usernameLogin").style.background = "red";
    $("#passwordLogin").style.background = "white";
  } else if (password == "") {
    $("#passwordLogin").style.background = "red";
    $("#usernameLogin").style.background = "white";
  } else {
    let payload =
      "Msideapp#1" + "\x1F" + username + "\x1F" + password + "\x1F" + "login";
    const formdata = new FormData();
    formdata.append("accAction", payload);
    //	const formdata = new FormData(accountForm);
    fetch(url, {
      method: "POST",
      body: formdata,
    })
      .then((response) => response.text())
      .then((data) => {
        console.log("data:", data)
        if (data.charAt(0) == 1) {
          // do something
           token = data;
          formdata.delete("accAction");
          $("#welcome").style.display = "block";
          $("#optionsBar").style.display = "block";
          $("#usernameLogin").style.background = "white";
          $("#passwordLogin").style.background = "white";
          $("#login").style.display = "none";
          return token;
        } else {
          $("#incpwd").innerHTML = "Incorrect Username or Password";
        }
      })
      .catch((error) => console.error(error));
  }
}

surp.loginShowPassword = function () {
  var x = $("#passwordLogin");
  if (x.type === "password") {
    x.type = "text";
  } else {
    x.type = "password";
  }
}

surp.loginForgotPw = function () {
  var email = prompt(
    "Enter your Email Adress. You will shortly receive an email.",
    "Email"
  );
  let payload =
    "Msideapp#1" + "\x1F" + username + "\x1F" + password + "\x1F" + "fgtpwd";
  const formdata = new FormData();
  formdata.append("accAction", payload);
  //	const formdata = new FormData(accountForm);
  fetch(url, {
    method: "POST",
    body: formdata,
  })
    .then((response) => response.text())
    .then((data) => {
      console.log(data);
    })
    .catch((error) => console.error(error));
  formdata.delete("accAction");
}

surp.loginSignup = function () {
	var username = prompt("Enter your username:", "Username");
  if (username == null) {
    return;
  }
  var password = prompt("Enter new password:", "New Password");
  if (password == null) {
    return;
  }
  let payload =
    "Msideapp#1" + "\x1F" + username + "\x1F" + password + "\x1F" + "create";
  const formdata = new FormData();
  formdata.append("accAction", payload);
  fetch(url, {
    method: "POST",
    body: formdata,
  })
    .then((response) => response.text())
    .then((data) => {
      if (data == 1) {
      alert("Login Successful")
    } else if (data == 2) {
      alert("Username is taken.")
    }
      console.log(data);
    })
    .catch((error) => console.error(error));
  formdata.delete("accAction");
}

surp.showCatTable = function () {
  clickedID = [];
  $("#welcome").style.display="none";
  $("#categories").style.display="block";
  $("#tasks").style.display="none";
  $("#timerpage").style.display="none";
  $("#taskentry").style.display="none";
  if (catArrayReceived == false) {
    let payload = "Msideapp#1" + "\x1F" + username + "\x1F" + password + "\x1F" + "catab" + "\x1F" + token;
    const formdata = new FormData();
    formdata.append('accAction', payload);
  //	const formdata = new FormData(accountForm);
    fetch(url,{

        method:"POST",
        body:formdata,
    }).then(
        response => response.text()
    ).then(
        (data) => {
          let json_obj = JSON.parse(data);
          catArray = [];


          var a = new Array;
          for(var i in json_obj) {
            let catname = json_obj[i].Catname;
            let object = Object.create(category);
            object.catName = catname;
            a.push(object);
            catArray.push(a);
            catname = "";
            a = [];
          }
          surp.fillCatTable();
          clickedID = [];
        }
    ).catch(
        error => console.error(error)
    )
    formdata.delete('accAction');
  }
  }

surp.fillCatTable = function() {
  var table = $("#categoryTable");
	for (let i = 0; i < catArray.length; i++) {
			// create new table row with empty cells
			var row = table.insertRow(i + 1);
			var c0 = row.insertCell(0);


			// fill first cell with category name
			c0.innerHTML = catArray[i][0].catName;
      var node = document.createElement("P");             // Create a <p> node
      node.id = "pc" + i;
      var textnode = document.createTextNode(catArray[i][0].catName);         // Create a text node
      node.appendChild(textnode);                              // Append the text to <li>
      $("#dropCat").appendChild(node);     // Append <p> to <div> with id="dropTask"
      c0.style.backgroundColor = "white";
      c0.style.color = "black";
			row.id = "c" + i; // gives current row in current Table an id for click eventListener later
      row.class = "A";
	}
  category.catName = "";
  category.taskName = "";
  category.status = "";
  category.dueDate = "";
  category.spredicted = "";
  category.upredicted = "";
  category.current = "";
	catArrayReceived = true;
  $("#dropCat").style.maxHeight = "300px";
  $("#dropCat").style.overflow = "auto";
}

surp.cancelTaskCreate = function () {
  $("#taskentry").style.display = "none";
  $("#tasks").style.display = "block";
  $("#optionsBar").style.display = "block";
  $("#taskNameInput").value = "";
  $("#predMinsInp").value = "";
  $("#predHoursInp").value = "";
  $("#taskCatInput").value = "";
  const d = new Date();
  let day = d.getUTCDate();
  var dateDay = "";
  if (day < 10) {
    dateDay = "0" + day.toString();
  } else {
    dateDay = day.toString();
  }
  let month = d.getUTCMonth() + 1;
  let year = d.getUTCFullYear();
  if (month < 10) {
    $("#taskDateInput").value = year.toString() + "-0" + month.toString() + "-" + dateDay;
  } else {
    $("#taskDateInput").value = year.toString() + "-" + month.toString() + "-" + dateDay;
  }
  $("#predSpan").innerHTML = "";
}

surp.deleteCatTable = function () {
  let rowNr = $("#categoryTable").rows.length;
  for (var i = rowNr - 1; i > 0; i--) {
      $("#categoryTable").deleteRow(i);
  }
  for (let i = rowNr - 2; i > -1; i--) {
    let tempID = "#pc" + i.toString();
    $(tempID).remove();
  }
}

surp.showPredicted = function () {
  let payload =
    "Msideapp#1" + "\x1F" + username + "\x1F" + catArray[clickedID[0]][0].catName + "\x1F" + "guess" + "\x1F" + token + "\x1F" + ($("#predMinsInp").value * 60 + $("#predHoursInp").value * 3600);
  const formdata = new FormData();
  formdata.append("accAction", payload);
  fetch(url, {
    method: "POST",
    body: formdata,
  })
    .then((response) => response.text())
    .then((data) => {
      console.log(data);
      $("#predSpan").innerHTML = data;
    })
    .catch((error) => console.error(error));
  formdata.delete("accAction");
}

surp.deleteTask = function (event) {
  let taskID = event.target.id;
  let taskNr = "";
  for (let i = 2; i < taskID.length; i++){
    taskNr = taskNr + taskID[i];
  }
  let tID = parseInt(taskNr);
  if (confirm("Delete Task: " + catArray[clickedID[0]][tID].taskName)){
    let payload =
      "Msideapp#1" + "\x1F" + username + "\x1F" + catArray[clickedID[0]][tID].catName + "\x1F"
                  + "delTask" + "\x1F" + token  + "\x1F" + catArray[clickedID[0]][tID].taskName;
    surp.deleteTaskArray(tID);
    const formdata = new FormData();
    formdata.append("accAction", payload);
    fetch(url, {
      method: "POST",
      body: formdata,
    })
      .then((response) => response.text())
      .then((data) => {
        console.log(data);
        surp.deleteTaskTable();
        surp.fillTaskTable(clickedID[0]);
      })
      .catch((error) => console.error(error));
    formdata.delete("accAction");
  }
}

surp.deleteCat = function (event) {
  let catID = event.target.id;
  let catNr = "";
  for (let i = 2; i < catID.length; i++){
    catNr = catNr + catID[i];
  }
  let cID = parseInt(catNr);
  if (confirm("Delete Category: " + catArray[cID][0].catName)) {
    let payload =
      "Msideapp#1" + "\x1F" + username + "\x1F" + catArray[cID][0].catName + "\x1F" + "delCat" + "\x1F" + token;
    surp.deleteCatArray(cID);
    const formdata = new FormData();
    formdata.append("accAction", payload);
    fetch(url, {
      method: "POST",
      body: formdata,
    })
      .then((response) => response.text())
      .then((data) => {
        console.log(data);
        surp.deleteCatTable();
        surp.fillCatTable();
      })
      .catch((error) => console.error(error));
    formdata.delete("accAction");
  }
}

surp.deleteTaskArray = function (tid) {
  let tempArray = new Array;
  for (let i = 0; i < catArray[clickedID[0]].length; i++) {
    if (i != tid) {
      tempArray.push(catArray[clickedID[0]][i]);
    }
  }
  catArray[clickedID[0]] = [];
  catArray[clickedID[0]] = tempArray;
}

surp.deleteCatArray = function (cid) {
  let tempArray = new Array;
  for (let i = 0; i < catArray.length; i++) {
    if (i != cid) {
      tempArray.push(catArray[i]);
    }
  }
  catArray = [];
  catArray = tempArray;
}

surp.taskFilterAll = function () {
  //all
  surp.deleteTaskTable();
  surp.fillTaskTable(clickedID[0]);
}

surp.taskFilterIncomplete = function () {
  //incomplete
  surp.deleteTaskTable();
  surp.fillTaskTable(clickedID[0]);
  var table = $("#taskTable");
	for (let i = catArray[clickedID[0]].length - 1; i > 0; i--) {
    if (catArray[clickedID[0]][i].status != "Incomplete"){
			// create new table row with empty cells
			table.deleteRow(i);
      let tempID = "#pt" + i.toString();
      $(tempID).remove();
    }
	}
}

surp.taskFilterPaused = function () {
  //paused
  surp.deleteTaskTable();
  surp.fillTaskTable(clickedID[0]);
  var table = $("#taskTable");
	for (let i = catArray[clickedID[0]].length - 1; i > 0; i--) {
    if (catArray[clickedID[0]][i].status != "Paused"){
			// create new table row with empty cells
			table.deleteRow(i);
      let tempID = "#pt" + i.toString();
      $(tempID).remove();
    }
	}
}

surp.taskFilterComplete = function () {
  //complete
  surp.deleteTaskTable();
  surp.fillTaskTable(clickedID[0]);
  var table = $("#taskTable");
	for (let i = catArray[clickedID[0]].length - 1; i > 0; i--) {
    if (catArray[clickedID[0]][i].status != "Complete"){
			// create new table row with empty cells
			table.deleteRow(i);
      let tempID = "#pt" + i.toString();
      $(tempID).remove();
    }
	}
}

//*******************************************************
// Setting things up
//*******************************************************

$("#categoryTable").addEventListener("click", surp.openCategory);
$("#taskTable").addEventListener("click", surp.openStopwatch);
$("#taskAdd").addEventListener("click", surp.openTaskEntry);
$("#timesubmit").addEventListener("click", surp.submitTime);
$("#tasksubmit").addEventListener("click", surp.submitTask);
$("#filter").addEventListener("change", surp.filter);
$("#start").addEventListener("click", surp.timeStart);
$("#pause").addEventListener("click", surp.timePause);
$("#addButton").addEventListener("click", surp.openTaskEntry);
$("#taskStartedTime").addEventListener("click", surp.taskStartedTime);
$("#accChPwd").addEventListener("click", surp.accountChangePassword);
$("#accLogout").addEventListener("click", surp.accountLogout);
$("#submitLogin").addEventListener("click", surp.accountLogin);
$("#loginShowPassword").addEventListener("click", surp.loginShowPassword);
$("#forgotPwLogin").addEventListener("click", surp.loginForgotPw);
$("#signupLogin").addEventListener("click", surp.loginSignup);
$("#catButton").addEventListener("click", surp.showCatTable);
$("#taskcancel").addEventListener("click", surp.cancelTaskCreate);
$("#showPredicted").addEventListener("click", surp.showPredicted);
$("#dropTask").addEventListener("click", surp.deleteTask);
$("#dropCat").addEventListener("click", surp.deleteCat);
/*$("#taskFilterAll").addEventListener("click", surp.taskFilterAll);
$("#taskFilterIncomplete").addEventListener("click", surp.taskFilterIncomplete);
$("#taskFilterPaused").addEventListener("click", surp.taskFilterPaused);
$("#taskFilterComplete").addEventListener("click", surp.taskFilterComplete);*/
