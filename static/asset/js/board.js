class Board {
  constructor (insert="boardWrap", board="board") {
    /* TODO: there is a bug where if you have to make either board or boardWrap,
             css transforms aren't applied properly. Therefore, just make sure they
             already exist before you make a board. This bug will probably never
             get fixed, at least not by me.
     */
    this.insert = document.getElementById(insert);

    // Create insertion point if it doesn't exist yet
    if (this.insert == null) {
      this.insert = document.createElement('div');
      this.insert.id = insert;
      document.body.appendChild(this.insert);
    }

    this.boardHolder = document.getElementById(board);
    // Create board if it doesn't exist yet
    if (this.boardHolder == null) {
      this.boardHolder = document.createElement('div');
      this.boardHolder.id = board;
      document.body.appendChild(this.board);
    }

    this.height = BOARD_HEIGHT;
    this.width = BOARD_WIDTH;

    this.tiles = [];

    // TODO: i and j might be mixed up here
    // Create the dom elements for all squares
    for (let j = 0; j < this.height; j++) {
      let col = document.createElement('div');
      col.className = 'column';
      for (let i = 0; i < this.width; i++) {
        let sq = document.createElement('div');
        sq.className = 'sq';
        sq.id = XyToRf(i, j)

        col.appendChild(sq);

        sq.onmouseenter = TileHover.bind(sq, j*this.width + i);
        sq.onmouseleave = TileExit.bind(sq);
        sq.onclick = ClickTile.bind(sq, i, j);

        this.tiles.push(sq)
      }

      this.boardHolder.appendChild(col);
    }
  }
}

function PlaceSprite (troop, player, id) {
  // Check which side
  const row = id[0];
  const col = Number(id[1]);

  const el = document.getElementById(id);
  const r = el.getBoundingClientRect();
  console.log(el);
  console.log(r);

  let img = document.createElement('img');
  // img.src = './graphics/cannibal/portrait.png';
  img.src = './asset/graphics/';

  img.className = 'troopa';

  if (Number(player) == 1) {
    img.src += 'p2/';
  } else {
    // Player 1, mirror the pieces
    img.style.transform = 'scaleX(-1)';
    img.src += 'p1/';
  }

  img.src += troop + '.png';

  let bottom = window.scrollY + r.bottom;

  let left = (r.left + r.right)/2;
  left -= (r.right - r.left)/2;

  img.style.bottom = bottom + 'px';
  img.style.left = left + 'px';

  img.style.width = r.width + 'px';
  img.style.height = 'auto';

  img.style.zIndex = (1000 - col) + '' ;

  document.body.insertBefore(img, document.getElementById('boardWrap'));

  return img;
}

function ClearSprites () {
  // Delete all sprites from the dom
  let els = document.getElementsByClassName('troopa');
  while (els.length > 0) {
    els[0].parentNode.removeChild(els[0]);
  }
}
