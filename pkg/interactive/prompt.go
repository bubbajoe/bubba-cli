package interactive

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/bubbajoe/bubba-cli/pkg/util"
	"github.com/c-bata/go-prompt"
)

var history = []string{"exit"}

var baseDirectory string = "./"

func StartInteractivePrompt() {
	setupBaseCommands()
	setupSearchCommand()
	setupEnvionmentVariableCmd()
	setupVsmCommand()
	for {
		in := prompt.Input(baseDirectory+"| bb> ", completer,
			prompt.OptionTitle("bubba-prompt"),
			prompt.OptionHistory(history),
			prompt.OptionPrefixTextColor(prompt.Red),
			prompt.OptionPreviewSuggestionTextColor(prompt.DarkGray),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))
		if in == "exit" {
			break
		}
		if in == "" {
			continue
		}
		tokens := strings.Split(in, " ")
		processInput(tokens[0], tokens)
	}
}

var commandFunctionMap = map[string]func([]string) error{
	"help": func(inputs []string) error {
		for k, d := range commandDefinitions {
			fmt.Printf("%s - %s\n", k, d)
		}
		return nil
	},
}

var commandDefinitions = map[string]string{}

func RegisterCommand(cmd string, def string, fn func([]string) error) {
	commandFunctionMap[cmd] = fn
	commandDefinitions[cmd] = def
}

func processInput(base string, inputs []string) {
	if processFunc, ok := commandFunctionMap[inputs[0]]; ok {
		err := processFunc(inputs[1:])
		if err != nil {
			fmt.Printf("'%s' error: %s", inputs[0], err)
		} else {
			history = append(history, base)
		}
	} else {
		fmt.Printf("Subcommand '%s' not found\n", inputs[0])
	}
}

func setupEnvionmentVariableCmd() {
	commandFunctionMap["env"] = func(inputs []string) error {
		if len(inputs) == 0 {
			for _, e := range os.Environ() {
				fmt.Println(e)
			}
			return nil
		}
		env := os.Getenv(inputs[0])
		fmt.Printf("env: %s='%s'\n", inputs[0], env)
		return nil
	}
	RegisterSuggestion("env", util.SliceMap(os.Environ(), func(s string) prompt.Suggest {
		return prompt.Suggest{Text: s}
	}))
}

func setupSearchCommand() {
	commandFunctionMap["search"] = func(inputs []string) error {

		return nil
	}
	RegisterSuggestion("env", util.SliceMap(os.Environ(), func(s string) prompt.Suggest {
		return prompt.Suggest{Text: s}
	}))
}

func setupVsmCommand() {
	commandFunctionMap["vsm"] = func(inputs []string) error {
		if len(inputs) == 0 {
			for _, e := range os.Environ() {
				fmt.Println(e)
			}
			return nil
		}
		env := os.Getenv(inputs[0])
		fmt.Printf("env: %s='%s'\n", inputs[0], env)
		return nil
	}
	RegisterSuggestion("env", util.SliceMap(os.Environ(), func(s string) prompt.Suggest {
		return prompt.Suggest{Text: s}
	}))
}

func setupBaseCommands() {
	dir, err := os.Getwd()
	if err != nil {
		// unable to set current working directory
	} else {
		baseDirectory = dir
	}

	RegisterCommand("pwd", "", func(inputs []string) error {
		fmt.Println("", baseDirectory)
		return nil
	})

	RegisterCommand("cd", "", func(inputs []string) error {
		if len(inputs) == 0 {
			return nil
		}
		loc := inputs[0]
		newDir := path.Join(baseDirectory, inputs[0])
		if strings.HasPrefix(loc, "/") {
			newDir = loc
		} else if strings.HasPrefix(loc, "~") {
			newDir = strings.Replace(loc, "~", os.Getenv("HOME"), 1)
		}
		_, err := os.ReadDir(newDir)
		if err != nil {
			return err
		}
		baseDirectory = newDir
		return nil
	})

	RegisterCommand("ls", "list files and directories", func(inputs []string) error {
		dirs, err := os.ReadDir(baseDirectory)
		if err != nil {
			return nil
		}
		for _, d := range dirs {
			if d.IsDir() {
				fmt.Printf("%s/\n", d.Name())
			} else {
				fmt.Printf("%s\n", d.Name())
			}
		}
		return nil
	})
	RegisterSuggestionFunc("cd", func(input string) []prompt.Suggest {
		dirs, err := os.ReadDir(baseDirectory)
		if err != nil {
			return []prompt.Suggest{}
		}
		return util.SliceFilter(dirs, func(de fs.DirEntry) *prompt.Suggest {
			if de.IsDir() {
				return &prompt.Suggest{Text: de.Name()}
			}
			return nil
		})
	})
}
