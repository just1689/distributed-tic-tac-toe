package model

import "testing"

func TestGame_CheckForWinner(t *testing.T) {
	g := NewGame()
	g.Board = [][]string{
		{"", "X", ""},
		{"", "X", ""},
		{"", "X", ""},
	}
	hasWinner := g.CheckForWinner()
	if !hasWinner {
		t.Fail()
	}

}

func TestGame_CheckForWinner3(t *testing.T) {
	g := NewGame()
	g.Board = [][]string{
		{"", "", "X"},
		{"", "X", ""},
		{"X", "", ""},
	}
	hasWinner := g.CheckForWinner()
	if !hasWinner {
		t.Fail()
	}

}

func TestGame_CheckForWinner4(t *testing.T) {
	g := NewGame()
	g.Board = [][]string{
		{"0", "", "X"},
		{"", "0", ""},
		{"X", "0", ""},
	}
	hasWinner := g.CheckForWinner()
	if hasWinner {
		t.Fail()
	}

}
