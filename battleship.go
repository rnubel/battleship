package battleship

type Coord struct {
	x int
	y int
}

type Ship struct {
	startPos   Coord
	horizontal bool
	size       int
	status     []bool
}

type Board struct {
	width  int
	height int
	hits   []Coord
	misses []Coord
	ships  []Ship
}

func (s Ship) isHit(x, y int) bool {
	return (s.horizontal && x < s.startPos.x+s.size && y == s.startPos.y ||
		!s.horizontal && y < s.startPos.y+s.size && x == s.startPos.x)
}

func (b *Board) recordHit(x, y int) {
	c := Coord{x: x, y: y}
	b.hits = append(b.hits, c)
}

func (b *Board) recordMiss(x, y int) {
	c := Coord{x: x, y: y}
	b.misses = append(b.misses, c)
}

func (b *Board) recordShot(x, y int) (hit bool, err string) {
	if x < 0 || y < 0 || x >= b.width || y >= b.height {
		hit, err = false, "out_of_bounds"
		return
	}

	for i := range b.ships {
		hit = hit || b.ships[i].isHit(x, y)
	}

	if hit {
		b.recordHit(x, y)
	} else {
		b.recordMiss(x, y)
	}

	return
}

func (b *Board) placeShip(x, y, size int, horizontal bool) (ok bool, err string) {
	c := Coord{x: x, y: y}
	s := Ship{startPos: c, horizontal: horizontal, size: size}

	if x < 0 || y < 0 || (horizontal && x+size > b.width) || (!horizontal && y+size > b.height) {
		ok, err = false, "out_of_bounds"
		return
	}

	b.ships = append(b.ships, s)

	ok, err = true, ""
	return
}
