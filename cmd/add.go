/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	editor "github.com/davidtaing/scriptcli/internal/editor"
)

var (
	scriptName = "helloworld"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if scriptName == "" {
			scriptName, err = promptScriptName()
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		var openEditor = promptOpenEditor()

		if openEditor {
			if Editor == "" {
				Editor = promptSelectEditor()
			}

			path, _ := createScript(scriptName)
			editor.OpenScriptInEditor(path, Editor)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&scriptName, "script", "s", "", "script file name")
}

func promptScriptName() (string, error) {
	var result string
	var err error

	prompt := promptui.Prompt{
		Label: "What would you like to name your new script?",
	}

	result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	if result == "" {
		return "", errors.New("New script name was not provided. Exiting")
	}

	return result, nil
}

func promptOpenEditor() bool {
	p := promptui.Prompt{
		Label:     "Open in editor?",
		IsConfirm: true,
	}

	result, _ := p.Run()

	return result == "y" || result == "Y"
}

func promptSelectEditor() string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label: "Which text editor would you like to use?",
			Items: editor.ValidEditors,
		}

		index, result, err = prompt.Run()
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	fmt.Printf("You choose %s\n", result)
	return result
}

func createScript(name string) (string, error) {
	var path = fmt.Sprintf("bin/%s.sh", name)

	err := os.WriteFile(path, []byte(`#!/bin/bash
echo "This Script was generated by Script CLI"`), 0755)

	if err != nil {
		fmt.Println("Error creating new script:", err)
		return "", err
	}

	fmt.Println("Created new script at", path)

	return path, nil
}
