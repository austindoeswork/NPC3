let Gam = null;

// Helper for making a move
function MakeMove(troop, toX, toY) {
  var msg = {
    'Type': 'MOVE',
    'Troop': troop,
    'X': toX,
    'Y': toY
  };

  wssend(msg);
}

// Helper for ending a turn
function EndTurn () {
  var msg = {
    'Type': 'END',
  };

  wssend(msg);
}

class Kranch {
  constructor () {
    this.board = this.InitBoard('insert', 'board');

    this.clickedX = -1;
    this.clickedY = -1;

    this.ClearMap();

    this.player = '';
    this.troops = [];
    this.boulders = [];
  }

  InitBoard (insertId, boardId) {
    let boardObject = new Board(insertId, boardId);
    let boardElement = boardObject.boardHolder;

    boardElement.addEventListener('mouseover', ShowTooltip);
    boardElement.addEventListener('mouseout', HideTooltip);

    window.addEventListener('mousemove', (e) => {
      let tt = document.getElementById('tooltip');
      tt.style.bottom = window.innerHeight - e.clientY.toString() + 'px';
      tt.style.left = e.clientX.toString() + 'px';
    }, false);

    return boardObject;
  }

  get width () {
    return this.board.width;
  }

  get height () {
    return this.board.height;
  }

  get tileList () {
    return this.board.tiles;
  }

  // Clear the selected tile
  ClearClicked () {
    if (this.clickedX != -1 && this.clickedY != -1) {
      let index = this.clickedY * this.width + this.clickedX;
      this.tileList[index].classList.toggle('selected', false);

      this.clickedX = -1;
      this.clickedY = -1;
    }
  }

  // Reset tile elements
  ClearMap () {
    this.map = [];
    for (let i = 0; i < this.width * this.height; i++) {
      this.map.push({'Type':-1});
    }
  }

  ClearTiles () {
    for (let i = 0; i < this.tileList.length; i++) {
      this.tileList[i].innerHTML = ' ';
    }
  }

  // Read / Parse the state

  Update (state) {
    console.log('got update')

    this.ClearMap();
    this.ClearTiles();

    this.player = state.You;

    console.log(state.You);
    console.log(state.Turn % 2);
    if (state.You == state.Turn % 2) {
      document.getElementById('etb').disabled = false;
      document.getElementById('etb').innerHTML = 'End Turn';
    } else {
      document.getElementById('etb').disabled = true;
      document.getElementById('etb').innerHTML = 'Enemy Turn';
    }

    this.troops = state.Troops;
    this.boulders = state.Boulders;

    this.ParseTroops();
    this.ParseBoulders();

    if (state.Status == -1) {
      if (state.Winner == this.player) {
        UpdateStatus('You win! Game over.', 'green');
      } else {
        UpdateStatus('You win! Game over.', 'red');
      }
    }
  }

  ParseTroops () {
    for (let i = 0; i < this.troops.length; i++) {
      for (let j = 0; j < this.troops[i].length; j++) {
        let t = this.troops[i][j];
        let x = t.X;
        let y = t.Y;
        let index = y * this.width + x;

        if (t.Step < t.Info.Mv) {
          this.tileList[index].innerHTML = '<b class="glow">' + t.Info.ShortName + '</b>';
        } else {
          this.tileList[index].innerHTML = t.Info.ShortName;
        }

        this.tileList[index].innerHTML += ' <p class="tileinfo">'+ t.Info.Atk+ ' ' + t.HP + '</p>';
        this.map[index].Troop = t;
        this.map[index].Index = j;
      }
    }
  }

  ParseBoulders () {
    for (var i = 0; i < this.boulders.length; i++) {
      let x = this.boulders[i].X;
      let y = this.boulders[i].Y;
      let index = y * this.width + x;

      this.tileList[index].innerHTML = '* <p class="tileinfo">' + this.boulders[i].HP + '</p>';
      this.map[index].Type = 2;
      this.map[index].Index = i;
    }
  }
}
