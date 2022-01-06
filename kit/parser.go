package kit

import (
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const kitFileName = "kit.yml"

func ParseKitFile(filePath string) Kit {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		panic(errors.New("parseKitFile(): File doesnt exist"))
	}

	data, err := os.ReadFile(filePath)
	check(err)

	kit := newKit()
	yaml.Unmarshal(data, kit)
	check(err)

	// TODO: VALIDATION LOGIC

	// Assign key to command struct as Alias
	for k, v := range kit.Commands {
		// Manually set KitArgument.Name from map
		for argk, argv := range v.Arguments {
			argv.Name = argk
			if len(argv.Prompt) == 0 {
				argv.Prompt = "Provide value for argument: " + argv.Name
			}
			v.Arguments[argk] = argv
		}

		v.Alias = k
		kit.Commands[k] = v
	}

	return kit
}

func FindKitFile() (string, error) {
	filePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for filePath != "/" {
		// Look for .kit file
		kitFilePath := path.Join(filePath, kitFileName)
		if _, err := os.Stat(kitFilePath); !os.IsNotExist(err) {
			return kitFilePath, nil
		}

		// Navigate to parent dir
		dirPath := path.Dir(filePath)
		filePath = dirPath
	}

	return "", errors.New("no kit.yml file found")
}
