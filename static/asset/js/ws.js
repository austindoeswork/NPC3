var xhr = new XMLHttpRequest();
var ws;

const URL = '//' + location.hostname


// XHR, deal with game list

function GetGames () {
  xhr.open( 'GET', URL + '/games', true);
  xhr.send();
  xhr.onreadystatechange = ProcessGames;

  window.setInterval(() => {
    xhr.open('GET', URL + '/games', true);
    xhr.send();
    xhr.onreadystatechange = ProcessGames;
  }, 5000);
}

function ProcessGames (e) {
  if (xhr.readyState != 4 || xhr.status != 200) {
    return;
  }

  const response = JSON.parse(xhr.responseText);

  let games = document.getElementById('gameList');
  games.innerHTML = '';

  _.forEach(response, (game) => {
    let gameEl = document.createElement('div');
    gameEl.id = game.Name;
    gameEl.onclick = _onRowClick.bind(gameEl, game.Name, game.Players);
    gameEl.className += ' gameRow';

    gameEl.innerHTML = game.Name + ' -- ' + game.Players.toString();
    games.appendChild(gameEl);
  });
}


// Websocket handlers

function _onWsMessage (cmd) {
  cmd = JSON.parse(cmd);

  const yourTurn = document.getElementById('status').innerHTML == 'Your turn';
  let message = cmd.Message;

  console.log(cmd);

  if (cmd.Type == 'STATE') {
    if (!yourTurn) {
      UpdateStatus(message, '');
    }
    Game.Update(cmd);
  } else {
    // Sound effects
    if (message == 'YOUR TURN') {
      let audio = new Audio('asset/audio/block.mp3');
      audio.play();
    } else if (message.includes('healed')) {
      let audio = new Audio('asset/audio/heal.mp3');
      audio.play();
    } else {
      let audio = new Audio('asset/audio/chess.mp3');
      audio.play();
    }

    // Update the status
    message = message.toLowerCase();
    message = message.charAt(0).toUpperCase() + message.slice(1);
    UpdateStatus(message, '');
  }
}

function wsopen (path, route, msgHandler, gameName) {
  if (ws) {
    return false;
  }

  if (gameName == null) {
    ws = new WebSocket(path + route);
  } else {
    ws = new WebSocket(path + route + '?game=' + gameName);
  }

  ws.onopen = (evt) => {
    UpdateStatus('Connected', 'green');
  }

  ws.onclose = (evt) => {
    UpdateStatus('Disconnected', 'gray');
    console.log(evt);
    ws = null;
  }

  ws.onmessage = (evt) => {
    _onWsMessage(evt.data);
  }

  ws.onerror = function(evt) {
    UpdateStatus('Error :(', 'red');
  }

  return true;
}

function wssend (input) {
  if (!ws) {
    return false;
  }

  ws.send(JSON.stringify(input));
  return true;
}

// Helper for creating a game
function CreateGame () {
  if (ws) {
    ws.close();
  }

  const wsPath = 'ws://' + location.hostname + ':' + location.port;
  const wsRoute = '/ws';
  const gameName = document.getElementById('gameName').value;

  wsopen(wsPath, wsRoute, _onWsMessage, gameName);
}

function _onRowClick (gameName, numPlayers) {
  if (ws) {
    ws.close();
  }

  const wspath = 'ws://' + location.hostname + ':' + location.port;
  let wsroute = '/ws';
  if (numPlayers >= 2) {
    wsroute += 'watch';
  }

  wsopen(wspath, wsroute, _onWsMessage, gameName);
  // document.getElementById(gameName).classList.toggle('selected', true);
}
