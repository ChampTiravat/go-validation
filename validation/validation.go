package validation

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultDataType string = "string"
	StringDataType  string = "string"
	IntDataType     string = "int"
	FloatDataType   string = "float"
	FileDataType    string = "file"

	TypeFormData       string = "form-data"
	TypeMultipartFile  string = "multipart-file"
	TypeFormURLEncoded string = "form-urlencoded"
	TypeURLParam       string = "url-param"
	TypeQueryString    string = "query-string"
)

type (
	Field struct {
		Name          string
		FieldType     string
		Value         interface{}
		File          *multipart.FileHeader
		FileName      string
		FileExtension string
		DataType      string
		DefaultValue  string
		ErrMsg        string
		Required      bool
	}

	File struct {
		Content   *multipart.File
		Name      string
		Extension string
	}

	Validator struct {
		Context      *gin.Context // TODO: Add support for other HTTP frameworks.
		Fields       map[string]Field
		CurrentField string
	}
)

func New(ctxPtr *gin.Context) Validator {
	return Validator{
		Context: ctxPtr,
		Fields:  map[string]Field{},
	}
}

func FromRequest(ctxPtr *gin.Context) Validator {
	return Validator{
		Context: ctxPtr,
		Fields:  map[string]Field{},
	}
}

// TODO: Validate struct
// func FromStruct(ctxPtr *gin.Context) Validator {
// 	return Validator{
// 		Context: ctxPtr,
// 		Fields:  map[string]Field{},
// 	}
// }

// TODO: Validate map
// func FromMap(ctxPtr *gin.Context) Validator {
// 	return Validator{
// 		Context: ctxPtr,
// 		Fields:  map[string]Field{},
// 	}
// }

// TODO: Add more string validation and string formating functions such as isEmail, isPhone, RegEx etc.
// Example: email := v.Form("email").Required().String().Format("email")

func GetFileExtension(fileName string) string {
	fileExtension := ""

	names := strings.Split(fileName, ".")
	if len(names) >= 2 {
		fileExtension = names[len(names)-1]
	}

	return strings.ToLower(fileExtension)
}

func (vPtr *Validator) AddField(fieldName string, fieldType string) Validator {
	var (
		err           error
		fieldValue    string
		fileName      string
		fileExtension string
		fileHeaderPtr *multipart.FileHeader
	)

	switch fieldType {
	case TypeQueryString:
		fieldValue = vPtr.Context.DefaultQuery(fieldName, "")
	case TypeURLParam:
		fieldValue = vPtr.Context.Param(fieldName)
	case TypeFormData:
		fieldValue = vPtr.Context.DefaultPostForm(fieldName, "")
	case TypeFormURLEncoded:
		fieldValue = vPtr.Context.DefaultPostForm(fieldName, "")
	case TypeMultipartFile:
		fileHeaderPtr, err = vPtr.Context.FormFile(fieldName)
		if err != nil {
			err = errors.New("file not found")
			fileHeaderPtr = nil
			fileExtension = ""
			fileName = ""
		} else {
			fileName = (*fileHeaderPtr).Filename
			fileExtension = GetFileExtension(fileName)
		}
	default:
		fieldValue = vPtr.Context.DefaultPostForm(fieldName, "")
	}

	var newField Field
	if err != nil {
		newField = Field{
			Name:          fieldName,
			FieldType:     fieldType,
			Value:         "",
			File:          nil,
			FileName:      fileName,
			FileExtension: fileExtension,
			DataType:      DefaultDataType,
			DefaultValue:  "",
			ErrMsg:        err.Error(),
			Required:      false,
		}
	} else {
		newField = Field{
			Name:          fieldName,
			FieldType:     fieldType,
			Value:         fieldValue,
			File:          fileHeaderPtr,
			FileName:      fileName,
			FileExtension: fileExtension,
			DataType:      DefaultDataType,
			DefaultValue:  "",
			ErrMsg:        "",
			Required:      false,
		}
	}

	(*vPtr).Fields[fieldName] = newField
	(*vPtr).CurrentField = fieldName

	return *vPtr
}

func (vPtr *Validator) Multipart(fieldName string) *Validator {
	(*vPtr).AddField(fieldName, TypeMultipartFile)
	return vPtr
}

func (vPtr *Validator) Form(fieldName string) *Validator {
	(*vPtr).AddField(fieldName, TypeFormURLEncoded)
	return vPtr
}

func (vPtr *Validator) FormData(fieldName string) *Validator {
	(*vPtr).AddField(fieldName, TypeFormData)
	return vPtr
}

func (vPtr *Validator) Query(fieldName string) *Validator {
	(*vPtr).AddField(fieldName, TypeQueryString)
	return vPtr
}

func (vPtr *Validator) Param(fieldName string) *Validator {
	(*vPtr).AddField(fieldName, TypeURLParam)
	return vPtr
}

func (vPtr *Validator) Required() *Validator {
	currentFieldName := (*vPtr).CurrentField

	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.Required = true

	textNotOK := tmpField.Value.(string) == ""
	fileNotOK := tmpField.File == nil

	if textNotOK && fileNotOK {
		msg := fmt.Sprintf(
			"`%s` must be specified in %s",
			tmpField.Name,
			tmpField.FieldType,
		)
		tmpField.ErrMsg = msg
	}

	(*vPtr).Fields[currentFieldName] = tmpField

	return vPtr
}

func (vPtr *Validator) Optional() *Validator {
	currentFieldName := (*vPtr).CurrentField

	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.Required = false

	(*vPtr).Fields[currentFieldName] = tmpField

	return vPtr
}

func (vPtr *Validator) Default(defaultValue interface{}) *Validator {
	currentFieldName := (*vPtr).CurrentField

	// TODO: If real value is empty, set value to default value
	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.DefaultValue = defaultValue.(string)

	(*vPtr).Fields[currentFieldName] = tmpField

	return vPtr
}

func (vPtr *Validator) Error(errMsg string) *Validator {
	currentFieldName := (*vPtr).CurrentField

	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.ErrMsg = errMsg

	(*vPtr).Fields[currentFieldName] = tmpField

	return vPtr
}

func (v Validator) CheckIfUnwantedFieldsExist() bool {
	return false
}

func (v Validator) CheckIfEmpty() bool {
	count := 0

	for _, field := range v.Fields {
		value := field.Value
		if value != nil {
			count++
		}
	}

	return count == 0
}

func (v Validator) Done() error {
	if v.CheckIfEmpty() {
		return errors.New("all inputs cannot be empty")
	}

	// TODO: Add a function to check if there're fields which're not predefined in the request.
	// if v.CheckIfUnwantedFieldsExist() {}

	for _, field := range v.Fields {
		if field.Required && field.ErrMsg != "" {
			return errors.New(field.ErrMsg)
		}
	}

	return nil
}

// TODO: Return both pointer to a file and filename
// file, filename := v.Multipart("document_file").Required().File()

func (vPtr *Validator) File() *File {
	currentFieldName := (*vPtr).CurrentField

	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.DataType = FileDataType

	filePtr := tmpField.File
	if filePtr == nil {
		msg := ""

		if tmpField.ErrMsg == "" {
			msg = fmt.Sprintf("%s must be %s", tmpField.Name, tmpField.DataType)
		} else {
			msg = fmt.Sprintf(" and the value must be %s", tmpField.DataType)
		}

		tmpField.ErrMsg = tmpField.ErrMsg + msg

		return nil
	}

	(*vPtr).Fields[currentFieldName] = tmpField

	fileContent, _ := (*filePtr).Open()

	file := &File{
		Content:   &fileContent,
		Name:      tmpField.FileName,
		Extension: tmpField.FileExtension,
	}

	return file
}

func (vPtr *Validator) String() string {
	currentFieldName := (*vPtr).CurrentField

	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.DataType = StringDataType

	strValue, assertionOK := tmpField.Value.(string)
	if !assertionOK {
		msg := ""

		if tmpField.ErrMsg == "" {
			msg = fmt.Sprintf("%s must be %s", tmpField.Name, tmpField.DataType)
		} else {
			msg = fmt.Sprintf(" and the value must be %s", tmpField.DataType)
		}

		tmpField.ErrMsg = tmpField.ErrMsg + msg
	}

	tmpField.Value = strValue
	(*vPtr).Fields[currentFieldName] = tmpField

	return strValue
}

func (vPtr *Validator) Int() int {
	currentFieldName := (*vPtr).CurrentField

	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.DataType = IntDataType

	value, err := strconv.ParseInt(tmpField.Value.(string), 10, 32)
	if err != nil {
		msg := ""

		if tmpField.ErrMsg == "" {
			msg = fmt.Sprintf("%s must be %s", tmpField.Name, tmpField.DataType)
		} else {
			msg = fmt.Sprintf(" and the value must be %s", tmpField.DataType)
		}

		tmpField.ErrMsg = tmpField.ErrMsg + msg
	}

	tmpField.Value = int(value)
	(*vPtr).Fields[currentFieldName] = tmpField

	return int(value)
}

func (vPtr *Validator) Float32() float32 {
	currentFieldName := (*vPtr).CurrentField

	tmpField := (*vPtr).Fields[currentFieldName]
	tmpField.DataType = FloatDataType

	value, err := strconv.ParseFloat(tmpField.Value.(string), 32)
	if err != nil {
		msg := ""

		if tmpField.ErrMsg == "" {
			msg = fmt.Sprintf("%s must be %s", tmpField.Name, tmpField.DataType)
		} else {
			msg = fmt.Sprintf(" and the value must be %s", tmpField.DataType)
		}

		tmpField.ErrMsg = tmpField.ErrMsg + msg
	}

	tmpField.Value = float32(value)
	(*vPtr).Fields[currentFieldName] = tmpField

	return float32(value)
}
