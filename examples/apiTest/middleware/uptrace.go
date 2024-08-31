package middleware

import "github.com/uptrace/uptrace-go/uptrace"

func InitUptrace(serviceName, version, dsn, env string) {
	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		uptrace.WithDSN(dsn),
		uptrace.WithDeploymentEnvironment(env),
		uptrace.WithServiceName(serviceName),
		uptrace.WithServiceVersion(version),
	)
}
