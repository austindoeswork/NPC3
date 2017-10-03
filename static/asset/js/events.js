function ShowTooltip () {
  document.getElementById('tooltip').classList.toggle('hidden', false);
}

function HideTooltip () {
  document.getElementById('tooltip').classList.toggle('hidden', true);
}

function TileHover (i) {
  let tSpot = Game.map[i];

  // Make sure there's a troop on this spot before showing the tooltip
  if (tSpot && ('Troop' in tSpot)) {
    let tt = document.getElementById('tooltip');
    tt.style.borderColor = TroopColor(tSpot.Troop);

    // Populate da tooltip with da info
    document.getElementById('ttName').innerHTML = tSpot.Troop.Nickname + ' the ' + tSpot.Troop.Info.Name;
    document.getElementById('ttDesc').innerHTML = tSpot.Troop.Info.Description;
    document.getElementById('ttQuote').innerHTML = tSpot.Troop.Info.Quote;

    ShowTooltip();
  } else {
    HideTooltip();
  }
}

function TileExit () {
  // Bye bye, birdie
  HideTooltip();
}

// Hover a tile or make a move
function ClickTile (x, y) {
  if (document.getElementById('etb').disabled) {
    return;
  }

  if (Game.clickedX == -1 || Game.clickedY == -1) {
    // If no tile is already selected, select the clicked tile
    let index = y * Game.width + x;
    Game.clickedX = x;
    Game.clickedY = y;
    Game.tileList[index].classList.toggle('selected', true);
  } else {
    // Else this click is a move input
    let index = Game.clickedY * Game.width + Game.clickedX;
    if (typeof Game.map[index].Troop != 'undefined') {
      if (Game.map[index].Troop.Owner == Game.player) {
        MakeMove(Game.map[index].Index, x, y);
      }
    }

    // Unhighlight the other tile
    Game.ClearClicked();
  }
}

function AddListeners () {
  document.getElementById('create').onclick = CreateGame;
  document.getElementById('etb').onclick = () => {
    if (this.disabled) {
      return false;
    } else {
      EndTurn();
    }
  }
}
