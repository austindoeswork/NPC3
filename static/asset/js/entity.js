function UUID4 () {
  const choices = '1234567890abcdef';
  let out = '';
  for (let i = 0; i < 32; i++) {
    out += choices[Math.floor(Math.random()*16)]
    if (i == 8 || i == 12 || i == 16 || i == 20) {
      out += '-';
    }
  }
  return out;
}

class Entity {
  constructor (type) {
    this.id = UUID4();

    if (typeof this.type == 'undefined') {
      // Will load this from server
      return;
    }

    this.Type = type;

    if (this.Type == 'air') {
      this.Representation = '.';
      this.Name = 'Noone';
      return;
    }

    // Load stats from template
    const template = DATA[Name];
    _.forEach(template, ((v, k) => {
      this[k] = v;
    }).bind(this));

    // Non-templated stats
    this.CurrentHP = this.Health;
    this.AttackBonus = 0;
    this.MovesLeft = this.MoveSpeed;
    this.CanAttack = true;

    // Create the element that will represent this troop
    this.Sprite = null;
  }

  MoveTo (x, y) {
    const letters = 'abcdefghijklmnopqrstuvwxyz';
    const height = BOARD_HEIGHT;

    // Id of the square
    const id = XyToRf(x, y);
  }

  PlaceSprite (id, offX, offY) {
    // Check which side
    const row = id[0];
    const col = Number(id[1]);
    const rside = ['e', 'f', 'g', 'h'];

    const el = document.getElementById(id);
    const r = el.getBoundingClientRect();

    let img = document.createElement('img');
    // img.src = './graphics/cannibal/portrait.png';
    img.src = './graphics/knight/idle.gif';

    img.style.position = 'fixed';

    if (rside.includes(row)) {
      img.style.transform = 'scaleX(-1)';
      offX *= -1;
    }

    let bottom = window.innerHeight - r.bottom;
    bottom += offY * r.height;

    let left = r.left;
    left += offX * r.width;

    img.style.bottom = bottom + 'px';
    img.style.left = left + 'px';

    img.style.width = r.width + 'px';
    img.style.height = 'auto';

    img.style.zIndex = (1000 - col) + '' ;

    document.body.insertBefore(img, document.getElementById('boardWrap'));
  }
}
