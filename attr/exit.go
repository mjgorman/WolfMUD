// Copyright 2015 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package attr

import (
	"code.wolfmud.org/WolfMUD.git/has"
)

// Constants for direction indexes. These can be used for the Link, AutoLink,
// Unlink and AutoUnlink methods. If these constants are modified probably need
// to update the Return function as well.
const (
	North byte = iota
	Northeast
	East
	Southeast
	South
	Southwest
	West
	Northwest
	Up
	Down
)

// directionNames is a lookup table for direction indexes to direction strings.
// When listing available exits they will be presented in the order they are in
// in this array.
var directionNames = [...]string{
	North:     "north",
	Northeast: "northeast",
	East:      "east",
	Southeast: "southeast",
	South:     "south",
	Southwest: "southwest",
	West:      "west",
	Northwest: "northwest",
	Up:        "up",
	Down:      "down",
}

// directionIndex is a lookup table for direction strings to direction indexes.
var directionIndex = map[string]byte{
	"N":         North,
	"NORTH":     North,
	"NE":        Northeast,
	"NORTHEAST": Northeast,
	"E":         East,
	"EAST":      East,
	"SE":        Southeast,
	"SOUTHEAST": Southeast,
	"S":         South,
	"SOUTH":     South,
	"SW":        Southwest,
	"SOUTHWEST": Southwest,
	"W":         West,
	"WEST":      West,
	"NW":        Northwest,
	"NORTHWEST": Northwest,
	"U":         Up,
	"UP":        Up,
	"D":         Down,
	"DOWN":      Down,
}

// Exits implements an attribute describing exits for the eight compass points
// north, northeast, east, southeast, south, southwest, west and northwest as
// well as the directions up and down and where they lead to. Exits are usually
// in pairs, for example one north and one back south although you can have one
// way exits or exits back that leads somewhere else.
type Exits struct {
	Attribute
	exits [len(directionNames)]has.Thing
}

// Some interfaces we want to make sure we implement
var (
	_ has.Exits = &Exits{}
)

// NewExits returns a new Exits attribute with no exits set. Exits should be
// added to the attribute using the Link and AutoLink methods. The reason exits
// cannot be set during initialisation like most other attributes is that all
// 'locations' have to be setup before they can all be linked together.
func NewExits() *Exits {
	return &Exits{Attribute{}, [len(directionNames)]has.Thing{}}
}

// FindExits searches the attributes of the specified Thing for attributes that
// implement has.Exits returning the first match it finds or nil otherwise.
func FindExits(t has.Thing) has.Exits {
	for _, a := range t.Attrs() {
		if a, ok := a.(has.Exits); ok {
			return a
		}
	}
	return nil
}

func (e *Exits) Dump() []string {
	buff := []byte{}
	for i, e := range e.exits {
		if e != nil {
			buff = append(buff, ", "...)
			buff = append(buff, directionNames[i]...)
			buff = append(buff, ": "...)
			if a := FindName(e); a != nil {
				buff = append(buff, a.Name()...)
			}
		}
	}
	if len(buff) > 0 {
		buff = buff[2:]
	}
	return []string{DumpFmt("%p %[1]T -> %s", e, buff)}
}

// Return calculates the opposite/return direction for the direction given.
// This is handy for calculating things like normal exits where if you go north
// you return by going back south. It is also useful for implementing ranged
// weapons, thrown weapons and spells. For example if you fire a bow west the
// person will see the arrow come from the east (from their perspective).
func Return(direction byte) byte {
	if direction < Up {
		return direction ^ 1<<2
	}
	return direction ^ 1
}

// Link links the given exit direction to the given Thing. If the given
// direction was already linked the exit will be overwritten - in effect the
// same as unlinking the exit first and then relinking it.
func (e *Exits) Link(direction byte, to has.Thing) {
	e.exits[direction] = to
}

// AutoLink links the given exit, calculates the opposite return exit and links
// that automatically as well - as long as the to Thing has an Exits attribute.
func (e *Exits) AutoLink(direction byte, to has.Thing) {
	e.Link(direction, to)
	if E := FindExits(to); E != nil {
		E.Link(Return(direction), e.Parent())
	}
}

// Unlink removes the exit for the given direction. It does not matter if the
// given direction was not linked in the first place.
func (e *Exits) Unlink(direction byte) {
	e.exits[direction] = nil
}

// AutoUnlink unlinks the given exit, calculates the opposite return exit and
// unlinks that automatically as well.
//
// BUG(diddymus): Does not check that exit A links to B and B links back to A.
// For example a maze may have an exit going North from A to B but going South
// from B takes you to C instead of back to A as would be expected!
func (e *Exits) AutoUnlink(direction byte) {
	to := e.exits[direction]
	e.Unlink(direction)

	if to == nil {
		return
	}

	if E := FindExits(to); E != nil {
		E.Unlink(Return(direction))
	}
}

// List will return a string listing the exits you can see. For example:
//
//	You can see exits east, southeast and south.
//
func (e *Exits) List() string {

	// Note we can tell the difference between l=0 initially and l=0 when the
	// last location was North by looking at the count c. If c is zero we have
	// not found any exits. If c is not zero then l=0 represents North.
	var (
		buff = make([]byte, 0, 1024) // buffer for direction list
		l    = 0                     // direction index of last exit found
		c    = 0                     // count of useable (linked) exits found
	)

	for i, e := range e.exits {
		switch {
		case e == nil:
			continue
		case c > 1:
			buff = append(buff, ", "...)
			fallthrough
		case c > 0:
			buff = append(buff, directionNames[l]...)
		}
		c++
		l = i
	}

	switch c {
	case 0:
		return "You can see no immediate exits from here."
	case 1:
		return "The only exit you can see from here is " + directionNames[l] + "."
	default:
		return "You can see exits " + string(buff) + " and " + directionNames[l] + "."
	}
}

// Move relocates a thing from it's current inventory to the inventory of the
// thing found following the given direction's exit. Note we use Thing a lot
// here as a location can be anything with an inventory - with or without
// exits!
//
// TODO: Is this the right place for this? It mostly deals with inventories so
// maybe it should go there? Really I guess this is 'glue' and should go into
// the cmd package as part of the move command itself?
func (e *Exits) Move(t has.Thing, direction string) (msg string, ok bool) {

	d, valid := directionIndex[direction]

	if !valid {
		msg = "You wanted to go which way!?"
		return
	}

	if e.exits[d] == nil {
		msg = "You can't go " + directionNames[d] + " from here!"
		return
	}

	from := FindInventory(e.Parent())
	if from == nil {
		msg = "You are not sure where you are, let alone where you are going."
		return
	}

	to := FindInventory(e.exits[d])
	if to == nil {
		msg = "For some odd reason you can't go " + directionNames[d] + "."
		return
	}

	if what := from.Remove(t); what == nil {
		msg = "Something stops you from leaving here!"
		return
	}

	to.Add(t)

	return "", true
}
