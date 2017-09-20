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
		let glist = document.getElementById("gamelist");
		glist.innerHTML = "";

		for (var i = 0; i < response.length; i++) {
			let gamerow = document.createElement("div");
			let rowjoin = rowclicked;
			gamerow.onclick = rowjoin.bind(gamerow, response[i].Name, response[i].Players);
			gamerow.className = " gamerow ";

			// let joinrow = document.createElement("button");
			// let rowjoin = joinclicked;
			// joinrow.onclick = rowjoin.bind(joinrow, response[i].Name);

			glist.appendChild(gamerow);
			gamerow.innerHTML = response[i].Name + " -- " + response[i].Players.toString();
			// gamerow.appendChild(joinrow);
			// gamerow.appendChild(watchrow);
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

function rowclicked(gamename, numplayers) {
    if (ws) {
		ws.close();
	}
		console.log("ROWCLICKED", gamename, numplayers);

    var wspath = "ws://" + location.hostname + ":" + location.port;
    var wsroute = "/ws";
    if (numplayers >= 2) {
		wsroute = "/wswatch";
    }

    wsopen(wspath, wsroute, gamename, handlemessage);
}

function wsjoin() {
    if (ws) {
		ws.close();
	}
    var wspath = "ws://" + location.hostname + ":" + location.port;
    var wsroute = "/ws";
    var gname = document.getElementById("gamename").value;
    wsopen(wspath, wsroute, gname, handlemessage);
}

function handlemessage(msg) {
	let cmd = JSON.parse(msg);
		console.log(cmd);
	
	if (cmd.Type == "STATE") {
		let currentmsg = document.getElementById("status").innerHTML;
		if (currentmsg  != "YOUR TURN") {
			setstatus(cmd.Message, "");
		}
	} else {
		messagenoise(cmd.Message);
		setstatus(cmd.Message, "");
	}
	if (cmd.Type == "STATE") {
		update(cmd);
	}	
}

function messagenoise(msg) {
	if (msg == "YOUR TURN") {
		let audio = new Audio('asset/sound/block.mp3');
		audio.play();
	}
	else if (msg.includes("healed")) {
		let audio = new Audio('asset/sound/heal.mp3');
		audio.play();
	}
	else {
		let audio = new Audio('asset/sound/chess.mp3');
		audio.play();
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
		console.log(evt);
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
