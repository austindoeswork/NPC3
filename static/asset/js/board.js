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
