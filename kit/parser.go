package kit

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const kitFileName = "kit.yml"
const kitRefFileName = "kitref.yml"

func ParseKitFile(filePath string) (Kit, error) {
	kit := newKit()

	home, _ := os.UserHomeDir()
	if strings.HasPrefix(filePath, "~/") {
		filePath = filepath.Join(home, filePath[2:])
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return kit, &KitFileNotFoundError{}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return kit, &InvalidKitFileError{}
	}

	err = yaml.Unmarshal(data, &kit)
	if err != nil {
		return kit, &InvalidKitFileError{}
	}

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

	return kit, nil
}

func ParseKitFromKitRef(kitRef KitRef) (Kit, error) {
	kit, err := ParseKitFile(kitRef.Path)
	if err != nil {
		return Kit{}, err
	}

	kit.Ref = &kitRef
	return kit, nil
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

	return "", &NoContextKitFoundError{}
}

func FindContextKit() (Kit, error) {
	kitFilePath, err := FindKitFile()
	if err != nil {
		return Kit{}, err
	}

	contextKit, err := ParseKitFile(kitFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return contextKit, nil
}

func GetKitRefList() (KitRefList, error) {
	kitRefList := newKitRefList()

	kitDir := GetOrMakeKitDir()
	kitRefPath := path.Join(kitDir, kitRefFileName)

	if _, err := os.Stat(kitRefPath); os.IsNotExist(err) {
		content := []byte("")
		writeErr := os.WriteFile(kitRefPath, content, 0644)
		check(writeErr)
		return newKitRefList(), nil
	}

	data, err := os.ReadFile(kitRefPath)
	if err != nil {
		return kitRefList, &InvalidKitRefFileError{}
	}

	err = yaml.Unmarshal(data, &kitRefList)
	if err != nil {
		return kitRefList, &InvalidKitRefFileError{}
	}

	return kitRefList, nil
}

func GetOrMakeKitDir() string {
	home, _ := os.UserHomeDir()
	kitDirPath := path.Join(home, "/.kit")

	if _, err := os.Stat(kitDirPath); os.IsNotExist(err) {
		makeErr := os.Mkdir(kitDirPath, 0o755)
		check(makeErr)
	}

	return kitDirPath
}

func GetOrMakeKitExecDir() string {
	kitDir := GetOrMakeKitDir()
	kitExecDirPath := path.Join(kitDir, "/exec")

	if _, err := os.Stat(kitExecDirPath); os.IsNotExist(err) {
		makeErr := os.Mkdir(kitExecDirPath, 0o755)
		check(makeErr)
	}

	return kitExecDirPath
}

func GetUserKits() []Kit {
	var userKits []Kit

	kitRefList, err := GetKitRefList()
	if err != nil {
		return userKits
	}

	for _, kitRef := range kitRefList.References {
		kit, err := ParseKitFromKitRef(kitRef)
		if err != nil {
			continue
		}

		userKits = append(userKits, kit)
	}
	return userKits
}

func FindUserKit(name string) (Kit, error) {
	userKits := GetUserKits()

	for _, kit := range userKits {
		if kit.Ref.Alias == name || kit.Name == name {
			return kit, nil
		}
	}
	return Kit{}, &NoMatchingKitError{}
}

func GetGlobalUserKits() []Kit {
	var globalUserKits []Kit
	userKits := GetUserKits()

	for _, kit := range userKits {
		if kit.Ref.Global {
			globalUserKits = append(globalUserKits, kit)
		}
	}

	return globalUserKits
}
