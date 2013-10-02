package battleship

import (
	"testing"
)

func TestBoard(t *testing.T) {
	b := Board{Width: 10, Height: 10}

	if b.Width != 10 || b.Height != 10 {
		t.Log("Height and Width were not set correctly")
		t.Fail()
	}
}

func TestShips(t *testing.T) {
  b := Board{Width: 10, Height: 10}

  ok, err := b.PlaceShip(0, 0, 10, true)
  if !ok || err != "" {
    t.Error("Was not able to place a 10-size ship at 0,0")
  }

  ok, err = b.PlaceShip(0, 0, 11, false)
  if ok || err != "out_of_bounds" {
    t.Error("Error was not returned properly")
  }
}

func TestHittingShips(t *testing.T) {
  b := Board{Width: 10, Height: 10}
  b.PlaceShip(0, 0, 10, true)

  hit, _ := b.RecordShot(4, 4)
  if hit {
    t.Error("4, 4 returned 'hit' when it should not have")
  }

  hit, _ = b.RecordShot(4, 0)
  if !hit {
    t.Error("4, 0 returned 'miss' when it should not have")
  }
}

func TestDestroyingShips(t *testing.T) {
  b := Board{Width: 10, Height: 10}
  b.PlaceShip(0, 0, 2, true)

  b.RecordShot(0,0)
  b.RecordShot(1,0)

  if !b.ships[0].isDead {
    t.Error("A ship of 2 parts was not destroyed properly")
  }


  b.PlaceShip(3, 3, 1, true)

  b.RecordShot(3,3)

  if !b.ships[1].isDead {
    t.Error("A ship of 1 part was not destroyed properly")
  }
}

func TestLosingBoard(t *testing.T) {
  b := Board{Width: 10, Height: 10}
  b.PlaceShip(0, 0, 2, true)
  b.PlaceShip(2, 2, 2, false)

  b.RecordShot(0, 0)
  b.RecordShot(1, 0)
  b.RecordShot(2, 2)
  b.RecordShot(2, 3)

  if !b.IsLost() {
    t.Error("Board is not lost even after all ships are destroyed")
  }
}

func TestCreateGame(t *testing.T) {
  p1, p2 := Player{Identifier: "p1"}, Player{Identifier: "p2"}
  g := CreateGame(10, 10, p1, p2)

  if g.Board1 == nil || g.Board2 == nil || g.Player1 != p1 || g.Player2 != p2 {
    t.Error("Game was not set up with two boards and correct players")
  }

  if g.Phase != NOTSTARTED {
    t.Error("Game is not in correct phase (NOTSTARTED)")
  }
}

func TestGameTurns(t *testing.T) {
  p1, p2 := Player{Identifier: "p1"}, Player{Identifier: "p2"}
  g := CreateGame(10, 10, p1, p2)
  g.Start()

  if g.CurrentPlayer != p1 {
    t.Error("Player 1 did not start the game")
  }

  g.endTurn()

  if g.CurrentPlayer != p2 {
    t.Error("Turn did not switch to player 2")
  }
}

func TestPlacement(t *testing.T) {
  p1, p2 := Player{Identifier: "p1"}, Player{Identifier: "p2"}
  g := CreateGame(10, 10, p1, p2)
  g.Start()

  if g.Phase != PLACEMENT {
    t.Error("Game is not in the PLACEMENT phase")
  }

  loc := Coord{x: 0, y: 0}
  placement := Placement{loc: loc, size: 4, horizontal: true}
  turn := Turn{Player: p1, TurnType: PLACEMENT_TURN, Placement: placement}

  result := g.executeTurn(turn)

  if !result.Ok || result.Error != "" {
    t.Error("Result was not as expected")
  }

  if len(g.Board1.ships) != 1 {
    t.Error("Turn did not actually place a ship")
  }
}

func TestPlacementLimits(t *testing.T) {
  p1, p2 := Player{Identifier: "p1"}, Player{Identifier: "p2"}
  g := CreateGame(10, 10, p1, p2)
  g.ShipAllocation = []int{1}
  g.Start()

  loc := Coord{x: 0, y: 0}
  placement := Placement{loc: loc, size: 6, horizontal: true} // no size 6 ships should be allowed
  turn := Turn{Player: p1, TurnType: PLACEMENT_TURN, Placement: placement}

  result := g.executeTurn(turn)

  if result.Ok || result.Error != "cannot_place_ship_of_that_size" {
    t.Error("Invalid placement was not rejected properly")
  }

  loc = Coord{x: 0, y: 0}
  placement = Placement{loc: loc, size: 1, horizontal: true}
  turn = Turn{Player: p1, TurnType: PLACEMENT_TURN, Placement: placement}

  result = g.executeTurn(turn)

  if !result.Ok {
    t.Error("Valid placement was rejected")
  }

  result = g.executeTurn(turn)

  if result.Ok || result.Error != "cannot_place_ship_of_that_size" {
    t.Error("Duplicate placement was not rejected")
  }
}

func TestPhaseTransitions(t *testing.T) {
  p1, p2 := Player{Identifier: "p1"}, Player{Identifier: "p2"}
  g := CreateGame(10, 10, p1, p2)
  g.ShipAllocation = []int{1}
  g.Start()

  loc := Coord{x: 0, y: 0}
  placement := Placement{loc: loc, size: 1, horizontal: true}
  turn1 := Turn{Player: p1, TurnType: PLACEMENT_TURN, Placement: placement}
  g.SubmitTurn(turn1)

  turn2 := Turn{Player: p2, TurnType: PLACEMENT_TURN, Placement: placement}
  g.SubmitTurn(turn2)

  if g.Phase != BATTLE {
    t.Error("Phase did not change to BATTLE once both players had placed all available ships")
  }

  // Now, have player1 win it.
  salvo := Salvo{locs: []Coord{loc}}
  turn3 := Turn{Player: p1, TurnType: SALVO_TURN, Salvo: salvo}
  g.SubmitTurn(turn3)

  if g.Phase != FINISHED {
    t.Error("Phase did not change to FINISHED when player 1 won")
  }

}
