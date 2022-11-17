package schema

import (
	"gorm.io/datatypes"
)

type Schema struct {
	ID       string          `json:"id" gorm:"not null;column:id;primaryKey"`
	SchemaID string          `json:"name" gorm:"not null;column:schema_id;uniqueIndex"`
	Schema   *datatypes.JSON `json:"schema" gorm:"column:schema"`
}
