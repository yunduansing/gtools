package jaeger

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"testing"
	"time"
)

func TestNewJaeger(t *testing.T) {
	NewJaeger(context.Background(),"http://localhost:14268/api/traces","jaeger-test","dev","1")
	//time.Sleep(time.Second*5)
	tr := otel.Tracer("controller")

	ctx, span := tr.Start(context.Background(), "call service")

	tracer(ctx)
	span.End()
	time.Sleep(time.Second*2)
}

func tracer(ctx context.Context)  {
	tr := otel.Tracer("service")
	_, span := tr.Start(ctx, "call dao")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()
}
