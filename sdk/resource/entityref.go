package resource

import (
	"reflect"

	"go.opentelemetry.io/otel/attribute"
)

type resourceEntityRef struct {
	schemaUrl string

	// Defines the entity type, e.g "service", "k8s.pod", etc.
	typ string

	// Set of Resource attribute keys that identify the entity.
	id attribute.Set

	// Set of Resource attribute keys that describe the entity.
	attrs attribute.Set
}

var _ = map[resourceEntityRef]int{}

func (r *resourceEntityRef) SchemaUrl() string {
	return r.schemaUrl
}

func (r *resourceEntityRef) Type() string {
	return r.typ
}

func (r *resourceEntityRef) Id() (keys []string) {
	iter := r.id.Iter()
	for iter.Next() {
		keys = append(keys, string(iter.Attribute().Key))
	}
	return keys
}

func (r *resourceEntityRef) Attrs() (keys []string) {
	iter := r.attrs.Iter()
	for iter.Next() {
		keys = append(keys, string(iter.Attribute().Key))
	}
	return keys
}

func keysToSet(keys []string) attribute.Set {

}

type resourceEntityRefs struct {
	iface any
}

func (d resourceEntityRefs) reflectValue() reflect.Value {
	return reflect.ValueOf(d.iface)
}

// Valid returns true if this value refers to a valid Set.
func (d resourceEntityRefs) Valid() bool {
	return d.iface != nil
}

// Len returns the number of attributes in this set.
func (l *resourceEntityRefs) Len() int {
	if l == nil || !l.Valid() {
		return 0
	}
	return l.reflectValue().Len()
}

// Get returns the KeyValue at ordered position idx in this set.
func (l *resourceEntityRefs) Get(idx int) (resourceEntityRef, bool) {
	if l == nil || !l.Valid() {
		return resourceEntityRef{}, false
	}
	value := l.reflectValue()

	if idx >= 0 && idx < value.Len() {
		// Note: The Go compiler successfully avoids an allocation for
		// the interface{} conversion here:
		return value.Index(idx).Interface().(resourceEntityRef), true
	}

	return resourceEntityRef{}, false
}

func (d resourceEntityRefs) AsSlice() []resourceEntityRef {

}

var resourceEntityRefType = reflect.TypeOf(resourceEntityRef{})

func NewResourceEntityRefs(refs []resourceEntityRef) resourceEntityRefs {
	at := reflect.New(reflect.ArrayOf(len(refs), resourceEntityRefType)).Elem()
	for i, keyValue := range refs {
		*(at.Index(i).Addr().Interface().(*resourceEntityRef)) = keyValue
	}
	return resourceEntityRefs{iface: at.Interface()}
}
