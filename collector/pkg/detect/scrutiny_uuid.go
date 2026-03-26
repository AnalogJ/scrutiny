package detect

import (
	"github.com/gofrs/uuid/v5"
)

// Randomly generated UUID v4 namespace for Scrutiny
var ScrutinyNamespaceUUID = uuid.Must(uuid.FromString("3ea22b35-682b-49fb-a655-abffed108e48"))

// WWN's are not actually unique so we use Model Name and Serial Number
// to hopefully create something that is actually unique despite
// manufacturer laziness
func GenerateScrutinyUUID(modelName string, serialNumber string, wwn string) uuid.UUID {
	name := modelName + serialNumber + wwn
	return uuid.NewV5(ScrutinyNamespaceUUID, name)
}
