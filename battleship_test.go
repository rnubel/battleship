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
