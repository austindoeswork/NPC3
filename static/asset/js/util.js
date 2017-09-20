const BOARD_HEIGHT = 7;
const BOARD_WIDTH = 7;

/*
  Take in an x- and y-coordinate, output the rank and file
  e.g. (5, 3) => 'e4'
*/
function XyToRf (x, y) {
  const letters = 'abcdefghijklmnopqrstuvwxyz';

  // Id of the square
  const id = letters[x-1] + (BOARD_HEIGHT-y);
}
