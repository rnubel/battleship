package battleship

import (
	"testing"
)

func TestBoard(t *testing.T) {
	b := Board{width: 10, height: 10}

	if b.width != 10 || b.height != 10 {
		t.Log("Height and width were not set correctly")
		t.Fail()
	}
}

func TestHits(t *testing.T) {
	b := Board{width: 5, height: 5}

  b.recordHit(0, 2)

	if len(b.hits) != 1 {
		t.Log("Hit was not added correctly.")
		t.Fail()
	}
}

func TestShips(t *testing.T) {
  b := Board{width: 10, height: 10}

  ok, err := b.placeShip(0, 0, 10, true)
  if !ok || err != "" {
    t.Log("Was not able to place a 10-size ship at 0,0")
    t.Fail()
  }

  ok, err = b.placeShip(0, 0, 11, false)
  if ok || err != "out_of_bounds" {
    t.Log("Error was not returned properly")
    t.Fail()
  }
}

func TestHittingShips(t *testing.T) {
  b := Board{width: 10, height: 10}
  b.placeShip(0, 0, 10, true)

  hit, _ := b.recordShot(4, 4)
  if hit {
    t.Log("4, 4 returned 'hit' when it should not have")
    t.Fail()
  }

  hit, _ = b.recordShot(4, 0)
  if !hit {
    t.Log("4, 0 returned 'miss' when it should not have")
    t.Fail()
  }
}
