// Copyright 2018 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"code.wolfmud.org/WolfMUD.git/attr"
	"code.wolfmud.org/WolfMUD.git/text"

	"strings"
)

// Syntax: $ACT <description>
func init() {
	addHandler(act{}, "$act")
}

type act cmd

func (act) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What was it that you wanted to act out?")
		return
	}

	who := attr.FindName(s.actor).TheName("Something")
	who = text.TitleFirst(who)

	msg := strings.Join(s.input, " ")

	s.msg.Actor.SendInfo(who, " ", msg)

	// Don't notify observers if it's crowded
	if !s.where.Crowded() {
		s.msg.Observer.SendInfo(who, " ", msg)
	}

	s.ok = true
}
