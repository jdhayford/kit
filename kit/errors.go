package kit

type KitFileNotFoundError struct {
	error
}

func (e *KitFileNotFoundError) Error() string {
	return "Kit file not found"
}

type InvalidKitFileError struct {
	error
}

func (e *InvalidKitFileError) Error() string {
	return "Kit file is invalid"
}

type InvalidKitRefFileError struct {
	error
}

func (e *InvalidKitRefFileError) Error() string {
	return "kitref file is invalid"
}

type NoMatchingKitError struct {
	error
}

func (e *NoMatchingKitError) Error() string {
	return "No matching user kit found"
}

type NoContextKitFoundError struct {
	error
}

func (e *NoContextKitFoundError) Error() string {
	return "No kit.yml file found"
}
