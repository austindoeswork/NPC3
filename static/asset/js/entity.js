const DATA = {
  King: {
    Representation: 'K',
    Name: 'Shirley',
    GraphicsDir: '',
    AttackPower: 2,
    Health: 5,
    MoveSpeed: 2,
    AttackRange: 1,
    OffsetY: 0,
    OffsetX: 0,
    Quote: 'If you die in the Kranch, you die in real life',
  },

  Ranger: {
    Representation: 'R',
    Name: 'Ranger',
    GraphicsDir: '',
    AttackPower: 2,
    Health: 8,
    MoveSpeed: 2,
    AttackRange: 2,
    OffsetY: 0,
    OffsetX: 0,
    Quote: 'To be one with the blistering desert... that is the way of my people',
  },

  Assassin: {
    Representation: 'A',
    Name: 'Assassin',
    GraphicsDir: '',
    AttackPower: 2,
    Health: 8,
    MoveSpeed: 3,
    AttackRange: 1,
    OffsetY: 0,
    OffsetX: 0,
    Quote: 'Would a real assassin give you a quote?'
  },

  Knight: {
    Representation: 'N',
    Name: 'Knight',
    GraphicsDir: '',
    AttackPower: 4,
    Health: 10,
    MoveSpeed: 2,
    AttackRange: 1,
    OffsetY: 0,
    OffsetX: 0,
    Quote: 'You look like the north end of a south-facing elephant',
  },

  Cannibal: {
    Representation: 'C',
    Name: 'Dune Worm',
    GraphicsDir: '',
    AttackPower: 3,
    Health: 6,
    MoveSpeed: 2,
    AttackRange: 1,
    OffsetX: 0.1,
    OffsetY: 0.3,
    Quote: 'GGLGLKGLGKLGGLKGL',
  },

  Healer: {
    Representation: 'H',
    Name: 'Medic',
    GraphicsDir: '',
    AttackPower: -2,
    Health: 8,
    MoveSpeed: 2,
    AttackRange: 1,
    OffsetY: 0,
    OffsetX: 0,
    Quote: 'I just like feeling important'
  },
};

const STARTS = {
  Player1: {
    King: {
      X: 2,
      Y: 1,
    }
  },

  Player2: {
    //
  },
};

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
