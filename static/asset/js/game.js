// Assumes there is a ws
// TODO fix that
var game = {};
game.tilelist = [];
game.troops = [];
game.boulders = [];
game.map = [];
game.width = 12;
game.height = 7;
game.player = -1;
game.clickedx = -1;
game.clickedy = -1;



var tileHover = (function(i) {
	let troopa = game.map[i];
	if (troopa && ("Troop" in troopa)) {
		let tt = document.getElementById("tooltip");
		tt.style.opacity = "1";
		tt.style.borderColor = troopcolor(troopa.Troop);
		tt.innerHTML = troopa.Troop.Nickname + " the " + troopa.Troop.Info.Name;
		tt.innerHTML += "<p class='description'>" + troopa.Troop.Info.Description + "</p>";
		tt.innerHTML += "<p class='quote'>" + troopa.Troop.Info.Quote + "</p>";
	}
});

var tileUnhover = (function() {
	// jello birdie
	let tt = document.getElementById("tooltip");
	tt.style.opacity = 0;
});

var setupgame = (function() {

	$("#board").mouseleave(clearclicked);
	var dad = document.getElementById("board");

	window.addEventListener('mousemove', (function(e) {
		let tt = document.getElementById("tooltip");
		tt.style.bottom = window.innerHeight - e.clientY.toString() + "px";
		tt.style.left = e.clientX.toString() + "px";
	}), false);

	for (var i = 0; i < game.width * game.height; i++) {
		let tile = document.createElement("div");
		tile.className += " box "
		let touchtile = boxclicked;

		tile.onmouseenter = tileHover.bind(tile, i);
		tile.onmouseleave = tileUnhover.bind(tile);
		tile.onclick = touchtile.bind(tile, i % game.width, Math.floor(i/game.width));
		tile.innerHTML = " "

		dad.appendChild(tile);
		game.tilelist.push(tile);
	}

	document.getElementById("endturn").onclick = endturn;
});

var boxclicked = (function(x, y) {
	if (game.clickedx == -1 || game.clickedy == -1) {
		let index = y * game.width + x;
		game.clickedx = x;
		game.clickedy = y;
		game.tilelist[index].style.border = "1px solid #00ff99";
	} else {
		let index = game.clickedy * game.width + game.clickedx;
		if (typeof game.map[index].Troop != "undefined") {
			if (game.map[index].Troop.Owner == game.player) {
				makemove(game.map[index].Index, x, y);
			}
		}
		clearclicked();
	}
});

var clearclicked = (function() {
	if (game.clickedx != -1 && game.clickedy != -1) {
		let index = game.clickedy * game.width + game.clickedx;
		game.tilelist[index].style.border = "1px solid grey";

		game.clickedx = -1;
		game.clickedy = -1;
	}
});

var clearmap = (function() {
	game.map = [];
	for (var i = 0; i < game.width * game.height; i++) {
		game.map.push({"Type":-1});
	}
});

var cleartiles = (function() {
	for (var i = 0; i < game.tilelist.length; i++) {
		game.tilelist[i].innerHTML = " ";
		game.tilelist[i].style.backgroundColor = "white";
	}
});

var makemove = (function(troop, tox, toy) {
	var msg = {
		"Type": "MOVE",
		"Troop": troop,
		"X": tox,
		"Y": toy
	};
	send(msg);
});

var endturn = (function() {
	var msg = {
		"Type": "END"
	};
	send(msg);
});

var update = (function(state) {
    console.log("got update")

	clearmap();
	cleartiles();

	game.player = state.You;

	if (state.You == state.Turn % 2) {
		document.getElementById("endturn").className = "btn btn-sm btn-warning";
		document.getElementById("endturn").innerHTML = "End Turn";
	}
	else {
		document.getElementById("endturn").className = "btn btn-sm btn-default";
		document.getElementById("endturn").innerHTML = "Enemy Turn";
	}
	game.troops = state.Troops;
	game.boulders = state.Boulders;
	parsetroops();
	parseboulders();

	if (state.Status == -1) {
		if (state.Winner == game.player) {
			setstatus("YOU WIN. GAME OVER");
		} else {
			setstatus("YOU LOSE. GAME OVER");
		}
	}
});

var parseboulders = (function(){
		for (var i = 0; i < game.boulders.length; i++) {
				let x = game.boulders[i].X;
				let y = game.boulders[i].Y;
				let index = y * game.width + x;
				game.tilelist[index].innerHTML = "*" +
					" <p class='tileinfo'>" + game.boulders[i].HP + "</p>"

				game.map[index].Type = 2;
				game.map[index].Index = i;
		}
});

var troopcolor = (function(t){
	if (t.Owner == 0) {
		if (!t.CanAct) {
			return "#9898b3"
		} else {
			return "lightblue"
		}
	}
   	if (t.Owner == 1) {
		if (!t.CanAct) {
			return "#936c6c" //greyish red
		} else {
			return "pink"
		}
	}

});

var parsetroops = (function(){
		for (var i = 0; i < game.troops.length; i++) {
				for (var j = 0; j < game.troops[i].length; j++) {
						let t = game.troops[i][j];
						let x = t.X;
						let y = t.Y;
						let index = y * game.width + x;
						if (t.Step < t.Info.Mv) {
							game.tilelist[index].innerHTML = "<b class='glow'>" + t.Info.ShortName + "</b>";
						}
						else {
							game.tilelist[index].innerHTML = t.Info.ShortName;
						}
						game.tilelist[index].innerHTML += " <p class='tileinfo'>"+ t.Info.Atk+ " " + t.HP + "</p>"

						game.tilelist[index].style.backgroundColor = troopcolor(t);
						game.map[index].Troop = t;
						game.map[index].Index = j;
				}
		}
});
