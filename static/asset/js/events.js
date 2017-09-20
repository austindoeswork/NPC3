function ShowTooltip () {
  document.getElementById('tooltip').classList.toggle('hidden', false);
}

function HideTooltip () {
  document.getElementById('tooltip').classList.toggle('hidden', true);
}

function TileHover () {
  let tSpot = Game.map[i];

  // Make sure there's a troop on this spot before showing the tooltip
  if (tSpot && ('Troop' in tSpot)) {
    let tt = document.getElementById('tooltip');
    tt.style.borderColor = troopcolor(tSpot.Troop);

    // Populate da tooltip with da info
    tt.innerHTML = tSpot.Troop.Nickname + ' the ' + tSpot.Troop.Info.Name;
    tt.innerHTML += '<p class="description">' + tSpot.Troop.Info.Description + '</p>';
    tt.innerHTML += '<p class="quote">' + tSpot.Troop.Info.Quote + '</p>';

    ShowTooltip();
  }
}

function TileExit () {
  // Bye bye, birdie
  HideTooltip();
}

function AddListeners () {
  document.getElementById('create').onclick = CreateGame;
}
