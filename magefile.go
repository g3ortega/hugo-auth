// +build mage

package main

import (
	"github.com/g3ortega/hugo-auth"
)

// Import docs from different sources
func UpdateContent() {
	hugo_auth.UpdateContent()
}
