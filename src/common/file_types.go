package common

type FileType struct {
	Extension  string
	FileFormat string
}

var TfFileType = FileType{Extension: ".tf", FileFormat: "tf"}
