/*
   Restfool-go

   Copyright (C) 2018 Carsten Seeger

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.

   @author Carsten Seeger
   @copyright Copyright (C) 2018 Carsten Seeger
   @license http://www.gnu.org/licenses/gpl-3.0 GNU General Public License 3
   @link https://github.com/cseeger-epages/rest-api-go-skeleton
*/

package restfool

import (
	"net/http"
)

// RestAPI contains api data
type RestAPI struct {
	Conf   config
	Routes []route
}

// Msg is the standard message type
type Msg struct {
	Message string `json:"message"`
}

// ErrMsg is the standard error message type
type ErrMsg struct {
	Error string `json:"error"`
}

// QueryStrings contains all possible query options
type QueryStrings struct {
	prettify bool
}

type pathList struct {
	Method      string
	Pattern     string
	Description interface{}
}

type pathPrefix struct {
	Prefix  string
	Handler http.Handler
}
