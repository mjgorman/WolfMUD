// Copyright 2012 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

// Package item implements a very basic, general item that can be picked up and
// carried around.
package item

import (
	"wolfmud.org/entities/location"
	"wolfmud.org/entities/thing"
	"wolfmud.org/utils/command"
	"wolfmud.org/utils/inventory"
	"wolfmud.org/utils/units"
)

// Item type is a default implemention of an item.
type Item struct {
	*thing.Thing
	weight units.Weight
}

// New allocates a new Item and returning a pointer reference to it.
func New(name string, aliases []string, description string, weight units.Weight) *Item {
	return &Item{
		Thing:  thing.New(name, aliases, description),
		weight: weight,
	}
}

// Process implements the command.Interface to handle Item specific
// commands.
func (i *Item) Process(cmd *command.Command) (handled bool) {

	// This specific item?
	if !i.IsAlias(cmd.Target) {
		return
	}

	switch cmd.Verb {
	case "DROP":
		handled = i.drop(cmd)
	case "WEIGH":
		handled = i.weigh(cmd)
	case "EXAMINE", "EXAM":
		handled = i.examine(cmd)
	case "GET":
		handled = i.get(cmd)
	case "JUNK":
		handled = i.junk(cmd)
	}

	return
}

// drop removes an Item from the command issuer's inventory and puts it into
// the inventory of the issuer's current location. For this to happen a few
// conditions must be true:
//
//	1. Issuer must be at some sort of location
//	2. Issuer must implement an inventory
//	3. Inventory must contain the requested item
//
func (i *Item) drop(cmd *command.Command) (handled bool) {
	if m, ok := cmd.Issuer.(location.Locateable); ok {
		if inv, ok := cmd.Issuer.(inventory.Interface); ok {
			if inv.Contains(i) {
				inv.Remove(i)
				cmd.Respond("You drop %s.", i.Name())
				cmd.Broadcast([]thing.Interface{cmd.Issuer}, "You see %s drop %s.", cmd.Issuer.Name(), i.Name())

				m.Locate().Add(i)

				handled = true
			}
		}
	}
	return
}

// weigh estimates the weight of the specified item.
func (i *Item) weigh(cmd *command.Command) (handled bool) {
	cmd.Respond("You estimate the weight of %s to be about %s.", i.Name(), i.weight)
	return true
}

// examine describes the specific item.
func (i *Item) examine(cmd *command.Command) (handled bool) {
	cmd.Respond("You examine %s. %s", i.Name(), i.Description())
	return true
}

// get removes an Item from the command issuer's current location and puts it
// into it's own inventory. For this to happen a few conditions must be true:
//
//	1. Issuer must be at some sort of location
//	2. Issuer must implement an inventory
//	3. Issuer's location must contain the requested item
//
func (i *Item) get(cmd *command.Command) (handled bool) {
	if m, ok := cmd.Issuer.(location.Locateable); ok {
		if inv, ok := cmd.Issuer.(inventory.Interface); ok {
			if l := m.Locate(); l.Contains(i) {
				l.Remove(i)
				cmd.Broadcast([]thing.Interface{cmd.Issuer}, "You see %s pick up %s.", cmd.Issuer.Name(), i.Name())

				inv.Add(i)
				cmd.Respond("You pickup %s.", i.Name())

				handled = true
			}
		}
	}
	return
}

// TODO: Implement junk command
func (i *Item) junk(cmd *command.Command) (handled bool) {
	cmd.Respond("Junk not implemented yet.")
	return true
}
