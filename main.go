/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/yourusername/k8s-controller-tutorial/cmd"

// Version will be injected at build time
var version = "dev"

func main() {
	cmd.Execute()
}
