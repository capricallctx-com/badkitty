// SPDX-License-Identifier: MIT
/*
   * Bad Kitty is a simple web server that can serve static files and reverse proxy requests to other servers.
   * It is designed to be a simple, easy to use, and easy to configure web server.

   * Contributors can add copyright here (not necessary - but a good idea).

   Copyright (c) 2024 - Caprica LLC
*/
package main

import (
	"log"
	"os/user"
)

func IsNotEmpty(s string) bool {
	return len(s) > 0
}
func AmIRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Printf("Error getting current user: %s", err.Error())
		return true
	}

	if currentUser.Uid == "0" {
		logger.Warn("WARNING: You are running this program as root. This can be a security risk.")
		return true
	}
	return false
}
