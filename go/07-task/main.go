package main

import (
	"github.com/ramin0/live/go/task/cmd"
)

/*
task is a CLI for managing your TODOs.

Usage:
  task [command]

Available Commands:
  add         Add a new task to your TODO list
  do          Mark a task on your TODO list as complete
  list        List all of your incomplete tasks

Use "task [command] --help" for more information about a command.
*/

func main() {
	cmd.Execute()
}
