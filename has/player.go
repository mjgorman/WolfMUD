// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package has

import (
	"io"
)

// Player is used to represent an actual player.
//
// Its default implementation is the attr.Player type.
type Player interface {
	Attribute

	// Player should implement a standard Write method to send data back to the
	// associated client.
	io.Writer
}
