var ws;
var xhr = new XMLHttpRequest();

function getgames() {
	xhr.open( "GET", "//austindoes.work/games", true ); 
    xhr.send();
	xhr.onreadystatechange = processgames;

	window.setInterval(function(){
  	/// call your function here
	xhr.open( "GET", "//austindoes.work/games", true ); 
    xhr.send();
	xhr.onreadystatechange = processgames;
	}, 5000);
}

function processgames(e) {
	if (xhr.readyState == 4 && xhr.status == 200) {
		var response = JSON.parse(xhr.responseText);

		let glist  = document.getElementById("gamelist")
		glist.innerHTML = "";

		for (var i = 0; i < response.length; i++) {
			let gamerow = document.createElement("div");
			let touchrow = rowclicked;
			gamerow.onclick = touchrow.bind(gamerow, response[i].Name);
			gamerow.className = " gamerow ";

			glist.appendChild(gamerow);
			gamerow.innerHTML = response[i].Name + " -- " + response[i].Players.toString();
		}
	}
}

function addlisteners() {   
	document.getElementById("join").onclick = wsjoin;
	document.getElementById("close").onclick = function(evt) {
        if (!ws) {
			return false;
        }
	    setstatus("DISCONNECTED", "label label-danger");
	    ws.close();
        return false;
	};
}

function rowclicked(gamename) {
    var wspath = "ws://" + location.hostname + ":" + location.port;
    var wsroute = "/ws";
    wsopen(wspath, wsroute, gamename, handlemessage);
}

function wsjoin() {
    var wspath = "ws://" + location.hostname + ":" + location.port;
    var wsroute = "/ws";
    var gname = document.getElementById("gamename").value;
    wsopen(wspath, wsroute, gname, handlemessage);
}

function handlemessage(msg) {
	let cmd = JSON.parse(msg);
	if (cmd.Type == "PROMPT") {
		setstatus(cmd.Message, "");
		if (cmd.Message == "YOUR TURN") {
			let audio = new Audio('asset/sound/block.mp3');
			audio.play();
		}
	} else if (cmd.Type == "ACK") {
		if (cmd.Success) {
			setstatus(cmd.Message, "label label-info");
		    if (cmd.Message.includes("healed")) {
					let audio = new Audio('asset/sound/heal.mp3');
					audio.play();
					console.log("HEALED");
			}
		    else {
					let audio = new Audio('asset/sound/chess.mp3');
					audio.play();
					console.log("OTERH");
			}
		} else {
			setstatus(cmd.Message, "label label-danger");
		}
	} else if (cmd.Type == "STATE") {
		update(cmd);
	}
}

function wsopen(wspath, wsroute, gname, handlemessage) {
    if (ws) {
		return false;
    }
    if (gname != "") {
		ws = new WebSocket(wspath + wsroute + "?game=" + gname);
    } else {
		ws = new WebSocket(wspath + wsroute);
    }

    ws.onopen = function(evt) {
	    setstatus("CONNECTED", "label label-info");
    }
    ws.onclose = function(evt) {
		setstatus("DISCONNECTED", "label label-danger");
		ws = null;
    }
    ws.onmessage = function(evt) {
		handlemessage(evt.data);
    }
    ws.onerror = function(evt) {
		setstatus("ERROR", "label label-danger");
    }
    return true;
}

function send(input) {
    if (!ws) {
			return false;
    }
    ws.send(JSON.stringify(input));
    return true;
}

function setstatus(statusstring, className) {
    s = document.getElementById("status");
    s.className = className;
    s.innerHTML = statusstring;
}
