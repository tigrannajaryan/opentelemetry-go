package semconv

import "go.opentelemetry.io/otel/attribute"

type service struct {
	attrs []attribute.KeyValue
}

func NewService(serviceName string) service {
	return service{
		attrs: []attribute.KeyValue{ServiceNameKey.String(serviceName)},
	}
}

func (s service) WithServiceNamespace(str string) service {
	return service{
		attrs: append(s.attrs, ServiceNamespaceKey.String(str)),
	}
}
