// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/g3ortega/hugo-auth"
)

// Import docs from different sources
func UpdateContent() {
	hugo_auth.UpdateContent()
}
