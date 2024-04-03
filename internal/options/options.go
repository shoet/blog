package options

import (
	"fmt"
	"reflect"
)

type ListBlogOptions struct {
	IsPublic bool
	Limit    int64
}

const DefaultLimit int64 = 10
const DefaultIsPublic bool = false

var ErrNotPointer = fmt.Errorf("v is not pointer")
var ErrFieldNotFound = fmt.Errorf("field is not found")
var ErrDefaultValueUnmatchType = fmt.Errorf("defaultValue is type unmatched")

// NewListBlogOptionsはデフォルト値が設定されたListBlogOptionsを生成する
func NewListBlogOptions(isPublic *bool, limit *int64) (*ListBlogOptions, error) {
	option := new(ListBlogOptions)
	if err := SetDefault(option, "IsPublic", isPublic, DefaultIsPublic); err != nil {
		return nil, fmt.Errorf("failed to set default value IsPublic: %v", err)
	}
	if err := SetDefault(option, "Limit", limit, DefaultLimit); err != nil {
		return nil, fmt.Errorf("failed to set default value Limit: %v", err)
	}
	return option, nil
}

// SetDefaultは構造体のフィールドにsetValueがnilの場合defaultValueを設定する
func SetDefault[T any](v any, fieldName string, setValue T, defaultValue any) error {
	rV := reflect.ValueOf(v)
	rSetValue := reflect.ValueOf(setValue)
	rDefaultValue := reflect.ValueOf(defaultValue)

	// rVがポインタ型でない場合はエラーを返す
	if rV.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	// setValueがポインタ型でない場合はエラーを返す
	if rSetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("setValue is not pointer")
	}

	// ポインタの中身を取得する
	rV = rV.Elem()

	// フィールドが存在しない場合はエラーを返す
	rField := rV.FieldByName(fieldName)
	if !rField.IsValid() || !rField.CanSet() {
		return ErrFieldNotFound
	}

	// フィールドの型がdefaultValueの型と異なる場合はエラーを返す
	if rField.Type() != rDefaultValue.Type() {
		return ErrDefaultValueUnmatchType
	}

	if rSetValue.IsNil() {
		rField.Set(rDefaultValue)
	} else {
		rField.Set(rSetValue.Elem())
	}
	return nil
}

type ListTagsOptions struct {
	Limit int
}
