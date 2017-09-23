const BOARD_HEIGHT = 7;
const BOARD_WIDTH = 7;

/*
  Take in an x- and y-coordinate, output the rank and file
  e.g. (5, 3) => 'e4'
*/
function XyToRf (x, y) {
  const letters = 'abcdefghijklmnopqrstuvwxyz';

  // Id of the square
  const id = letters[x] + (BOARD_HEIGHT-y);
  return id;
}

/*
  Decide what color to use for a troop (based on team)
*/
function TroopColor (t) {
  if (t.Owner == 0) {
    if (!t.CanAct) {
      return 'var(--my-team-inactive)';
    } else {
      return 'var(--my-team-color)';
    }
  }

  if (t.Owner == 1) {
    if (!t.CanAct) {
      return 'var(--their-team-inactive)';
    } else {
      return 'var(--their-team-color)';
    }
  }
}


/*
  Change the text and color of the status message
*/
function UpdateStatus (text, className) {
  let statusEl = document.getElementById('status');
  statusEl.innerHTML = text;
  statusEl.className = className;
}
