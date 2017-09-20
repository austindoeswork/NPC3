class Kranch {
  constructor () {
    return;
  }
}

var Game = new Kranch();

function SetupGame () {
  const insertId = 'insert';
  const boardId = 'board';
  let boardObject = new Board(insertId, boardId);
  let boardElement = boardObject.board;


  boardElement.addEventListener('mouseover', ShowTooltip);
  boardElement.addEventListener('mouseout', HideTooltip);

  window.addEventListener('mousemove', (e) => {
    let tt = document.getElementById('tooltip');
    tt.style.bottom = window.innerHeight - e.clientY.toString() + 'px';
    tt.style.left = e.clientX.toString() + 'px';
  }, false);
}

function ClearClicked () {
  if (game.clickedx != -1 && game.clickedy != -1) {
    let index = game.clickedy * game.width + game.clickedx;
    game.tilelist[index].style.border = "1px solid grey";

    game.clickedx = -1;
    game.clickedy = -1;
  }
}
