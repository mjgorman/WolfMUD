// Copyright 2015 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"code.wolfmud.org/WolfMUD.git/attr"
)

// Syntax: ( N | NORTH | NE | NORTHEAST | E | EAST | SE | SOUTHEAST | S | SOUTH
//				 | SW | SOUTHWEST | W | WEST | NW | NORTHWEST | U | UP | D | DOWN)
//
func init() {
	AddHandler(Move,
		"N", "NE", "E", "SE", "S", "SW", "W", "NW", "U", "D",
		"NORTH", "NORTHEAST", "EAST", "SOUTHEAST",
		"SOUTH", "SOUTHWEST", "WEST", "NORTHWEST",
		"UP", "DOWN",
	)
}

// TODO: Move does not support vetoes yet.
func Move(s *state) {

	from := s.where

	// A thing can only move itself if it knows where it is
	if from == nil {
		s.msg.Actor.Send("You are not sure where you are, let alone where you are going!")
		return
	}

	// Is where we are exitable?
	exits := attr.FindExits(from.Parent())
	if !exits.Found() {
		s.msg.Actor.Send("You can't see anywhere to go from here.")
		return
	}

	// Is direction a valid direction? Move could have been called directly by
	// another command just passing in the direction.
	direction, err := exits.NormalizeDirection(s.cmd)
	if err != nil {
		s.msg.Actor.Send("You wanted to go which way!?")
		return
	}

	wayToGo := exits.ToName(direction)

	// Find out where our exit leads to
	to := exits.LeadsTo(direction)
	if to == nil {
		s.msg.Actor.Send("You can't go ", wayToGo, " from here!")
		return
	}

	// Are we locking our destination yet? If not add it to the locks and simply
	// return. The parser will detect the locks have changed and reprocess the
	// command with the new locks held.
	if !s.CanLock(to) {
		s.AddLock(to)
		return
	}

	if from.Remove(s.actor) == nil {
		s.msg.Actor.Send("Something stops you from leaving here!")
		return
	}

	to.Add(s.actor)

	// Re-point where we are and re-alias observer
	s.where = to
	s.msg.Observer = s.msg.Observers[s.where]

	// Get actors name
	name := attr.FindName(s.actor).Name("someone")

	// Broadcast leaving and arrival notifications
	s.msg.Observers[from].Send("You see ", name, " go ", wayToGo, ".")
	s.msg.Observers[to].Send("You see ", name, " enter.")

	// Describe our destination
	s.scriptActor("LOOK")
}
