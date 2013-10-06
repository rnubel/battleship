package battleship

import "fmt"

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

type GamePhase int
const (
  NOTSTARTED GamePhase = iota
  PLACEMENT
  BATTLE    // um, some go magic makes this work
  FINISHED
)

type Player struct {
  Identifier string
  Name       string
}

type Placement struct {
  loc         Coord
  size        int
  horizontal  bool
}

type Salvo struct {
  Locs        []Coord
}

type TurnType int
const (
  PLACEMENT_TURN TurnType = iota
  SALVO_TURN
)

type Turn struct {
  Player    Player
  TurnType  TurnType
  Placement Placement
  Salvo     Salvo
}

type Result struct {
  Ok    bool
  Error string
  Hits  int
}

type Game struct {
  Board1    *Board
  Board2    *Board
  Player1   Player
  Player2   Player

  ShipAllocation  []int
  SalvoAllocation map[int]int // ship size -> shots to fire
  Phase           GamePhase
  CurrentPlayer   Player
  Turns           []Turn
}

func (s *Ship) covers(x, y int) bool {
  for i := range(s.parts) {
    if s.parts[i].loc.x == x && s.parts[i].loc.y == y {
      return true
    }
  }

  return false
}

func (s *Ship) Size() int {
  return len(s.parts)
}

func (s *Ship) RecordShot(x, y int) bool {
  for i := range(s.parts) {
    if s.parts[i].loc.x == x && s.parts[i].loc.y == y {
      if s.parts[i].isHit {
        return false // already hit
      }
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

func (g *Game) canPlace(b *Board, s int) bool {
  allowedShips := make(map[int]int)
  for i := range(g.ShipAllocation) {
    allowedShips[g.ShipAllocation[i]]++
  }

  for i := range(b.ships) {
    allowedShips[b.ships[i].Size()]--
  }

  return allowedShips[s] > 0
}

func (g *Game) allShipsPlacedOnBoard(b *Board) bool {
  allowedShips := make(map[int]int)
  for i := range(g.ShipAllocation) {
    allowedShips[g.ShipAllocation[i]]++
  }

  for i := range(b.ships) {
    allowedShips[b.ships[i].Size()]--
  }

  for i := range(g.ShipAllocation) {
    if allowedShips[g.ShipAllocation[i]] > 0 {
      return false
    }
  }

  return true
}

func (g *Game) allShipsPlaced() bool {
  return g.allShipsPlacedOnBoard(g.Board1) && g.allShipsPlacedOnBoard(g.Board2)
}

func (g *Game) endTurn() {
  if g.CurrentPlayer == g.Player1 {
    g.CurrentPlayer = g.Player2
  } else {
    g.CurrentPlayer = g.Player1
  }

  if g.Phase == PLACEMENT && g.allShipsPlaced() {
    g.Phase = BATTLE
  } else if g.Phase == BATTLE && (g.Board1.IsLost() || g.Board2.IsLost()) {
    g.Phase = FINISHED
  }
}

func (g *Game) boardsForPlayer(p Player) (myBoard, theirBoard *Board) {
  if p.Identifier == g.Player1.Identifier {
    myBoard     = g.Board1
    theirBoard  = g.Board2
  } else {
    myBoard     = g.Board2
    theirBoard  = g.Board1
  }

  return
}

func (g *Game) executeTurn(t Turn) (r Result) {
  myBoard, theirBoard := g.boardsForPlayer(t.Player)

  switch t.TurnType {
  case PLACEMENT_TURN:
    if g.Phase != PLACEMENT {
      return Result{Ok: false, Error: "not_placement_phase"}
    }

    p := t.Placement
    if !g.canPlace(myBoard, p.size) {
      return Result{Ok: false, Error: "cannot_place_ship_of_that_size"}
    }

    ok, err := myBoard.PlaceShip(
      p.loc.x,
      p.loc.y,
      p.size,
      p.horizontal)

    return Result{Ok: ok, Error: err}
  case SALVO_TURN:
    if g.Phase != BATTLE {
      return Result{Ok: false, Error: "not_battle_phase"}
    }

    if t.Player.Identifier != g.CurrentPlayer.Identifier {
      return Result{Ok: false, Error: "not_your_turn"}
    }

    if len(t.Salvo.Locs) > g.ShotsPlayerCanFire(t.Player) {
      return Result{Ok: false, Error: fmt.Sprintf("too_many_shots__max_is_%d", g.ShotsPlayerCanFire(t.Player))}
    }

    s, ok := t.Salvo, true
    for i := range(s.Locs) {
      ok = ok && theirBoard.IsValidShot(s.Locs[i].x, s.Locs[i].y)
    }

    if !ok {
      return Result{Ok: false, Error: "salvo_out_of_bounds"}
    }

    hits := 0
    for i := range(s.Locs) {
      // without a transaction, too late to handle errors.
      hit, _ := theirBoard.RecordShot(s.Locs[i].x, s.Locs[i].y)
      if hit {
        hits++
      }
    }

    return Result{Ok: true, Hits: hits}
  default:
    return Result{Ok: false, Error: "invalid_turn_type"}
  }
}

// Public API follows.

func CreateGame(width, height int, p1, p2 Player) Game {
  b1, b2 := Board{Width: width, Height: height}, Board{Width: width, Height: height}
  ships := []int{2,3,3,4,5} // allocated ship sizes
  shots := map[int]int{1:1, 2:1, 3:1, 4:1, 5:1} // no variance here; just one shot per ship.
  g := Game{Board1: &b1,
            Board2: &b2,
            Player1: p1,
            Player2: p2,
            Phase: NOTSTARTED,
            CurrentPlayer: p1,
            ShipAllocation: ships,
            SalvoAllocation: shots}

  return g
}

func (b *Board) IsValidShot(x, y int) bool {
	return !(x < 0 || y < 0 || x >= b.Width || y >= b.Height)
}

func (b *Board) RecordShot(x, y int) (hit bool, err string) {
	if !b.IsValidShot(x, y) {
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

    if b.ShipAt(loc.x, loc.y) {
      return false, "collides_with_other_ship"
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

func (g *Game) Start() {
  g.CurrentPlayer = g.Player1
  g.Phase = PLACEMENT
}

func (g *Game) SubmitTurn(t Turn) (ok bool, err string, hits int) {
  result := g.executeTurn(t)

  if result.Ok {
    g.endTurn()
    return true, "", result.Hits
  } else {
    return false, result.Error, 0
  }
}

func (g *Game) SubmitPlacementTurn(playerID string, x, y, size int, horizontal bool) (bool, string) {
  player := Player{Identifier: playerID}
  loc := Coord{x: x, y: y}
  placement := Placement{loc: loc, size: size, horizontal: horizontal}

  ok, err, _ := g.SubmitTurn(Turn{Player: player, TurnType: PLACEMENT_TURN, Placement: placement})
  return ok, err
}

func (g *Game) SubmitSalvoTurn(playerID string, locs []Coord) (bool, string, int) {
  player := Player{Identifier: playerID}
  salvo := Salvo{Locs: locs}

  ok, err, hits := g.SubmitTurn(Turn{Player: player, TurnType: SALVO_TURN, Salvo: salvo})
  return ok, err, hits
}

func (g *Game) ShotsPlayerCanFire(p Player) (shots int) {
  myBoard, _ := g.boardsForPlayer(p)

  for i := range(myBoard.ships) {
    if !myBoard.ships[i].isDead {
      shots += g.SalvoAllocation[myBoard.ships[i].Size()]
    }
  }

  return
}

func NewCoord(x, y int) Coord {
  return Coord{x: x, y: y}
}
