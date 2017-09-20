var xhr = new XMLHttpRequest();
var ws;

function GetGames () {
  xhr.open( 'GET', '//austindoes.work/games', true);
  xhr.send();
  xhr.onreadystatechange = ProcessGames;

  window.setInterval(() => {
    xhr.open('GET', '//austindoes.work/games', true);
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

  _.forEach (response, (game) => {
    let gameEl = document.createElement('div');
    gameEl.onclick = _onRowClick.bind(gamerow, game.Name, game.Players);
    gameEl.className += ' gameRow';

    gameEl.innerHTML = game.Name + ' -- ' + game.Players.toString();
    games.appendChild(gameEl);
  });
}


function _onRowClick (gameName, numPlayers) {
  if (ws) {
    ws.close();
  }

  const wspath = 'ws://' + location.hostname + ':' + location.port;
  let wsroute = '/ws';
  if (numplayers >= 2) {
    wsroute += 'watch';
  }

  wsopen(wspath, wsroute, gameName, _onWsMessage);
}

function AddListeners () {
  document.getElementById('join').onclick = JoinGame;
  document.getElementById('close').onclick = function(evt) {
    if (!ws) {
      return false;
    }

    UpdateStatus('DISCONNECTED', 'label label-danger');
    ws.close();
    return false;
  };
}


function JoinGame () {
  if (ws) {
    ws.close();
  }

  const wsPath = 'ws://' + location.hostname + ':' + location.port;
  const wsRoute = '/ws';
  const gameName = document.getElementById('gamename').value;
  wsopen(wsPath, wsRoute, _onWsMessage, gameName);
}

function _onWsMessage (cmd) {
  cmd = JSON.parse(cmd);

  const yourTurn = document.getElementById('status').innerHTML == 'YOUR TURN';
  const message = cmd.Message;

  if (cmd.Type == 'STATE') {
    if (!yourTurn) {
      UpdateStatus(message, '');
    }
    Game.Update(cmd);
  } else {
    // Sound effects
    if (message == 'YOUR TURN') {
      let audio = new Audio('asset/sound/block.mp3');
      audio.play();
    } else if (message.includes('healed')) {
      let audio = new Audio('asset/sound/heal.mp3');
      audio.play();
    } else {
      let audio = new Audio('asset/sound/chess.mp3');
      audio.play();
    }

    // Update the status
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
    ws = new WebSocket(path + route + '?game=' + gname);
  }

  ws.onopen = (evt) => {
    UpdateStatus('CONNECTED', 'label label-info');
  }

  ws.onclose = (evt) => {
    UpdateStatus('DISCONNECTED', 'label label-danger');
    console.log(evt);
    ws = null;
  }

  ws.onmessage = (evt) => {
    _onWsMessage(evt.data);
  }

  ws.onerror = function(evt) {
    UpdateStatus('ERROR', 'label label-danger');
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

function UpdateStatus (text, className) {
  let statusEl = document.getElementById('status');
  statusEl.innerHTML = text;
  statusEl.className = className;
}
