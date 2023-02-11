// fixadated is a daemon for a collaborative date finding tool.
//
//	Copyright (C) 2023 Daniel Steinhauer <d.steinhauer@mailbox.org>
//
//	This program is free software: you can redistribute it and/or modify
//	it under the terms of the GNU Affero General Public License as
//	published by the Free Software Foundation, either version 3 of the
//	License, or (at your option) any later version.
//
//	This program is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even the implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU Affero General Public License for more details.
//
//	You should have received a copy of the GNU Affero General Public License
//	along with this program.  If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"os"
	"os/signal"
)

func GetSignalChannel() chan os.Signal {
	ret := make(chan os.Signal, 1)
	signal.Notify(ret, os.Interrupt, os.Kill)
	return ret
}
