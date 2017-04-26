var tilelist = [];
var troops = [
        [
            {
                "CanAct": true,
                "HP": 2,
                "Info": {
                    "Atk": 2,
                    "MaxHP": 2,
                    "Mv": 2,
                    "Name": "general",
                    "Rng": 1,
                    "Secondary": 0,
                    "ShortName": "g"
                },
                "Owner": 0,
                "Step": 0,
                "X": 0,
                "Y": 2
            },
            {
                "CanAct": true,
                "HP": 12,
                "Info": {
                    "Atk": 4,
                    "MaxHP": 12,
                    "Mv": 2,
                    "Name": "knight",
                    "Rng": 1,
                    "Secondary": 0,
                    "ShortName": "k"
                },
                "Owner": 0,
                "Step": 0,
                "X": 1,
                "Y": 3
            }
        ],
        [
            {
                "CanAct": true,
                "HP": 2,
                "Info": {
                    "Atk": 2,
                    "MaxHP": 2,
                    "Mv": 2,
                    "Name": "general",
                    "Rng": 1,
                    "Secondary": 0,
                    "ShortName": "g"
                },
                "Owner": 1,
                "Step": 0,
                "X": 9,
                "Y": 2
            },
            {
                "CanAct": true,
                "HP": 12,
                "Info": {
                    "Atk": 4,
                    "MaxHP": 12,
                    "Mv": 2,
                    "Name": "knight",
                    "Rng": 1,
                    "Secondary": 0,
                    "ShortName": "k"
                },
                "Owner": 1,
                "Step": 0,
                "X": 8,
                "Y": 2
            }
        ]
    ];
var boulders = [{X:0,Y:0},{X:4,Y:6}];

var width = 12;
var height = 7;

window.onload = function() {
	console.log("hello world");

	var dad = document.getElementById("board");
	for (var i = 0; i < width * height; i++) {
		let kid = document.createElement("div");
		kid.className += " box "
		let touchkid = (function(x, y) {
			console.log(x, y);
		});

		kid.onclick = touchkid.bind(kid, i % width, Math.floor(i/width));
		// kid.x = i % 12
		// kid.y = Math.floor(i/2);
		kid.innerHTML = " "

		dad.appendChild(kid);
		tilelist.push(kid);
	}

    parseboulders();
    parsetroops();
};

var parseboulders = (function(){
		for (var i = 0; i < boulders.length; i++) {
				let x = boulders[i].X;
				let y = boulders[i].Y;
				let index = y * width + x;
				tilelist[index].innerHTML = "*";
		}
});

var parsetroops = (function(){
		for (var i = 0; i < troops.length; i++) {
				for (var j = 0; j < troops[0].length; j++) {
						let t = troops[i][j];
						let x = t.X;
						let y = t.Y;
						let index = y * width + x;
						tilelist[index].innerHTML = t.Info.ShortName;
				}
		}
});
