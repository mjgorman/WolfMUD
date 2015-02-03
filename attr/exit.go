// Copyright 2015 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package attr

import (
	"code.wolfmud.org/WolfMUD-mini.git/has"
)

const (
	_              = iota
	N, NORTH uint8 = iota, iota
	NE, NORTHEAST
	E, EAST
	SE, SOUTHEAST
	S, SOUTH
	SW, SOUTHWEST
	W, WEST
	NW, NORTHWEST
	U, UP
	D, DOWN
)

var directionLongNames = [...]string{
	N:  "north",
	NE: "northeast",
	E:  "east",
	SE: "southeast",
	S:  "south",
	SW: "southwest",
	W:  "west",
	NW: "northwest",
	U:  "up",
	D:  "down",
}

var directionIndex = map[string]uint8{
	"N":         N,
	"NORTH":     N,
	"NE":        NE,
	"NORTHEAST": NE,
	"E":         E,
	"EAST":      E,
	"SE":        SE,
	"SOUTHEAST": SE,
	"S":         S,
	"SOUTH":     S,
	"SW":        SW,
	"SOUTHWEST": SW,
	"W":         W,
	"WEST":      W,
	"NW":        NW,
	"NORTHWEST": NW,
	"U":         U,
	"UP":        U,
	"D":         D,
	"DOWN":      D,
}

type exits struct {
	attribute
	exits [len(directionLongNames)]has.Thing
}

// Some interfaces we want to make sure we implement
var (
	_ has.Attribute = Exits()
	_ has.Exits     = Exits()
)

func Exits() *exits {
	return nil
}

func (*exits) New() *exits {
	return &exits{attribute{}, [len(directionLongNames)]has.Thing{}}
}

func (*exits) Find(t has.Thing) has.Exits {
	for _, a := range t.Attrs() {
		if a, ok := a.(has.Exits); ok {
			return a
		}
	}
	return nil
}

func (e *exits) Dump() []string {
	buff := []byte{}
	for i, e := range e.exits {
		if e != nil {
			buff = append(buff, ", "...)
			buff = append(buff, directionLongNames[i]...)
			buff = append(buff, ": "...)
			if a := Name().Find(e); a != nil {
				buff = append(buff, a.Name()...)
			}
		}
	}
	if len(buff) > 0 {
		buff = buff[2:]
	}
	return []string{DumpFmt("%p %[1]T -> %s", e, buff)}
}

func (e *exits) Link(direction uint8, to has.Thing) {
	e.exits[direction] = to
}

func (e *exits) Unlink(direction uint8) {
	e.exits[direction] = nil
}

func (e *exits) List() string {

	var (
		buff = make([]byte, 0, 1024) // buffer for direction list
		l    = 0                     // Last direction found
		c    = 0                     // Count of directions processed
	)

	for i, e := range e.exits {
		if e != nil {
			if l > 0 {
				if c > 1 {
					buff = append(buff, ", "...)
				}
				buff = append(buff, directionLongNames[l]...)
			}
			c++
			l = i
		}
	}

	switch c {
	case 0:
		return "You can see no immediate exits from here."
	case 1:
		return "The only exit you can see from here is " + directionLongNames[l] + "."
	default:
		return "You can see exits " + string(buff) + " and " + directionLongNames[l] + "."
	}
}

func (e *exits) Place(t has.Thing) {
	if a := Inventory().Find(e.Parent()); a != nil {
		a.Add(t)
	}
}

func (e *exits) Move(t has.Thing, cmd string) string {

	// Check direction is valid e.g. "N" or "NORTH"
	d := directionIndex[cmd]
	if d == 0 {
		return "You wanted to go which way!?"
	}

	if e.exits[d] == nil {
		return "You can't go " + directionLongNames[d] + " from here!"
	}

	from := Inventory().Find(e.Parent())
	if from == nil {
		return "You are not sure where you are, let alone where you are going."
	}

	to := Inventory().Find(e.exits[d])
	if to == nil {
		return "For some odd reason you can't go " + directionLongNames[d] + "."
	}

	if what := from.Remove(t); what == nil {
		return "Something stops you from leaving here!"
	}

	to.Add(t)

	return ""
}
