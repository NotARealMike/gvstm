package correctness

import (
	"fmt"
	. "github.com/NotARealMike/gvstm/stm"
	"github.com/NotARealMike/gvstm/stmbench7/interfaces"
	"github.com/NotARealMike/gvstm/stmbench7/internal"
	"strconv"
	"strings"
)

type InvariantFailedError struct {
	message string
}

func (e *InvariantFailedError) Error() string {
	return e.message
}

func newInvFailedError(obj interface{}, id int, message string) error {
	withID := ""
	if id > 0 {
		withID = " with id = " + strconv.Itoa(id)
	}
	return &InvariantFailedError{
		message: fmt.Sprintf("Invariant violated! Object %T%s:\n  %s", obj, withID, message),
	}
}
func newUnexpectedIntError(obj interface{}, id int, attribute string, expected, found int) error {
	return newInvFailedError(obj, id, fmt.Sprintf("attribute %s --> expected: %d, found: %d", attribute, expected, found))
}

func newOutOfBoundError(obj interface{}, id int, attribute string, min, max, found int) error {
	return newInvFailedError(obj, id, fmt.Sprintf("attribute %s --> expected value in [%d,%d], found: %d", attribute, min, max, found))
}

func newUnexpectedStringError(obj interface{}, id int, attribute, expected, found string) error {
	return newInvFailedError(obj, id, fmt.Sprintf("attribute %s --> expected: %s, found: %s", attribute, expected, found))
}

func newIndexError(index, message string) error {
	return &InvariantFailedError{message: fmt.Sprintf("Index %s:\n  %s", index, message)}
}

func checkIndexValue(index string, traversedObject, indexedObject interface{}, elementID int) error {
	if indexedObject == nil {
		return newIndexError(index, fmt.Sprintf("element with id = %d in the data structure but not the index", elementID))
	}
	if indexedObject != traversedObject {
		return newIndexError(index, fmt.Sprintf("element with id = %d is different in data structure and index", elementID))
	}
	return nil
}

func checkAllTraversed(tx Transaction, index interfaces.Index, traversedSet map[interface{}]struct{}, indexName string) error {
	for _, key := range index.GetKeys(tx) {
		if obj, ok := traversedSet[index.Get(tx, key)]; !ok {
			return newIndexError(indexName, fmt.Sprintf("index contains too many elements (element %v)", obj))
		}
	}
	return nil
}

func checkValidType(typ string) bool {
	prefix := "type #"
	if !strings.HasPrefix(typ, prefix) {
		return false
	}
	typeNum, err := strconv.ParseInt(strings.TrimPrefix(typ, prefix), 10, 0)
	if err != nil {
		return false
	}
	if typeNum < 0 || int(typeNum) >= internal.NumTypes {
		return false
	}
	return true
}
