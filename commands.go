package main

import (
    "fmt"
    "strings"
)

type CommandProcessing func([]string) error
type Command struct {
    Info string
    Executor CommandProcessing
}

var (
    CommandTable = make(map[string]Command)
)

/* Return error if command not exist in CommandTable, nil otherwise */
func IsCommandExist(command string) error {
    var err error
    _, exist := CommandTable[command]

    if !exist {
        err = fmt.Errorf("Unknown command: '%s'", command)
    }

    return err
}

/* Search for command executor and call it, print error in terminal for any errors occured */
func ExecuteCommand(input string) {
    splitted := strings.Split(input, " ")

    if len(splitted) == 0 {
        return
    }

    err := IsCommandExist(splitted[0])

    if err != nil {
        Terminal.Print(err.Error())
        return
    }

    commandStruct, _ := CommandTable[splitted[0]]

    err = commandStruct.Executor(splitted[1:])

    if err != nil {
        Terminal.Print(err.Error())
    }
}

/* Write given command to list, overwrite command if exist */
func RegisterCommand(name string, obj Command) {
    CommandTable[name] = obj
}
