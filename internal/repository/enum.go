package repository

type PrivateDataType string

const (
	PrivateDataTypeUnknown PrivateDataType = "UNKNOWN"
	PrivateDataTypeText    PrivateDataType = "TEXT"
	PrivateDataTypeFile    PrivateDataType = "FILE"
	PrivateDataTypeAuth    PrivateDataType = "AUTH"
	PrivateDataTypeBank    PrivateDataType = "BANK"
)

func (t PrivateDataType) ToString() string {
	return string(t)
}
