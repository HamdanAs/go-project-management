package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(value interface{}) error {
	// Implementasi Scan untuk mengonversi nilai database ke UUIDArray
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("failed to parse value: unsupported data type")
	}

	str = str[1 : len(str)-1] // Menghapus kurung kurawal
	parts := strings.Split(str, ",")

	*a = make(UUIDArray, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(strings.Trim(part, `"`)) // menghapus spasi dan tanda kutip
		if part == "" {
			continue
		}

		u, err := uuid.Parse(part)
		if err != nil {
			return fmt.Errorf("invalid uuid in array: %v", err)
		}

		*a = append(*a, u)
	}

	return nil
}

func (a UUIDArray) Value() (driver.Value, error) {
	// Implementasi Value untuk mengonversi UUIDArray ke format yang dapat disimpan di database
	if len(a) == 0 {
		return "{}", nil
	}

	parts := make([]string, 0, len(a))
	for _, part := range a {
		parts = append(parts, fmt.Sprintf(`"%s"`, part.String()))
	}

	return fmt.Sprintf("{%s}", strings.Join(parts, ",")), nil
}

func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
