package battleship

type Coord struct {
	x int
	y int
}

type ShipPart struct {
  loc   Coord
  isHit bool
}

type Ship struct {
  parts   []ShipPart
  isDead  bool
}

type Board struct {
	Width  int
	Height int
	hits   []Coord
	misses []Coord
	ships  []Ship
}

func (s *Ship) covers(x, y int) bool {
  for i := range(s.parts) {
    if s.parts[i].loc.x == x && s.parts[i].loc.y == y {
      return true
    }
  }

  return false
}

func (s *Ship) RecordShot(x, y int) bool {
  for i := range(s.parts) {
    if s.parts[i].loc.x == x && s.parts[i].loc.y == y {
      s.parts[i].isHit = true
      s.isDead = s.checkIfDestroyed()
      return true
		}
  }

  return false
}

func (s *Ship) checkIfDestroyed() bool {
  for i := range(s.parts) {
    if !s.parts[i].isHit {
      return false
    }
  }

  return true
}

func (b *Board) recordHit(x, y int) {
	c := Coord{x: x, y: y}
	b.hits = append(b.hits, c)
}

func (b *Board) recordMiss(x, y int) {
	c := Coord{x: x, y: y}
	b.misses = append(b.misses, c)
}

// Public API follows.

func (b *Board) RecordShot(x, y int) (hit bool, err string) {
	if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
		hit, err = false, "out_of_bounds"
		return
	}

	for i := range b.ships {
		hit = hit || b.ships[i].RecordShot(x, y)
	}

	if hit {
		b.recordHit(x, y)
	} else {
		b.recordMiss(x, y)
	}

	return
}

func (b *Board) HitAt(x, y int) bool {
  for i := range(b.hits) {
    if b.hits[i].x == x && b.hits[i].y == y {
      return true
    }
  }
  return false
}

func (b *Board) MissAt(x, y int) bool {
  for i := range(b.misses) {
    if b.misses[i].x == x && b.misses[i].y == y {
      return true
    }
  }
  return false
}

func (b *Board) ShipAt(x, y int) bool {
  for i := range(b.ships) {
    if b.ships[i].covers(x, y) {
      return true
    }
  }
  return false
}

func (b *Board) NumShips() int {
  return len(b.ships)
}

func (b *Board) PlaceShip(x, y, size int, horizontal bool) (ok bool, err string) {
	s := Ship{}

	if x < 0 || y < 0 || (horizontal && x+size > b.Width) || (!horizontal && y+size > b.Height) {
		ok, err = false, "out_of_bounds"
		return
	}

  for i := 0; i < size; i++ {
    loc := Coord{x: x, y: y}
    if horizontal {
      loc.x += i
    } else {
      loc.y += i
    }

    part := ShipPart{loc: loc, isHit: false}
    s.parts = append(s.parts, part)
  }

	b.ships = append(b.ships, s)

	ok, err = true, ""
	return
}

func (b *Board) IsLost() (lost bool) {
  lost = true
  for i := range(b.ships) {
    lost = lost && b.ships[i].isDead
  }

  return
}
