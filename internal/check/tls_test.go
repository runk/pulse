package check

import (
	"crypto/tls"
	"testing"
	"time"
)

func TestTLSCheckValidate(t *testing.T) {
	scenarios := []struct {
		name     string
		check    TLSCheck
		expected string
		port     uint16
	}{
		{"host is required", TLSCheck{}, "host is required to perform tls check", 0},
		{"default port", TLSCheck{Host: "example.com"}, "", 443},
		{"custom port", TLSCheck{Host: "example.com", Port: 8443}, "", 8443},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := scenario.check.Validate()
			actual := ""
			if err != nil {
				actual = err.Error()
			}

			if actual != scenario.expected {
				t.Errorf("Want '%s' but got '%s'", scenario.expected, actual)
			}

			if scenario.port != 0 && scenario.check.Port != scenario.port {
				t.Errorf("Want port %d but got %d", scenario.port, scenario.check.Port)
			}
		})
	}
}

func TestTLSVersionName(t *testing.T) {
	scenarios := []struct {
		name     string
		version  uint16
		expected string
	}{
		{"tls 1.0", tls.VersionTLS10, "TLS 1.0"},
		{"tls 1.1", tls.VersionTLS11, "TLS 1.1"},
		{"tls 1.2", tls.VersionTLS12, "TLS 1.2"},
		{"tls 1.3", tls.VersionTLS13, "TLS 1.3"},
		{"unknown", 0x9999, "TLS 0x9999"},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actual := tlsVersionName(scenario.version)
			if actual != scenario.expected {
				t.Errorf("Want '%s' but got '%s'", scenario.expected, actual)
			}
		})
	}
}

func TestTLSCipherSuiteName(t *testing.T) {
	scenarios := []struct {
		name     string
		id       uint16
		expected string
	}{
		{
			"known",
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		},
		{"unknown", 0x9999, "0x9999"},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actual := tlsCipherSuiteName(scenario.id)
			if actual != scenario.expected {
				t.Errorf("Want '%s' but got '%s'", scenario.expected, actual)
			}
		})
	}
}

func TestDaysUntil(t *testing.T) {
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)

	scenarios := []struct {
		name     string
		future   time.Time
		expected int
	}{
		{"same time", now, 0},
		{"less than a day", now.Add(23 * time.Hour), 0},
		{"one day", now.Add(24 * time.Hour), 1},
		{"two days", now.Add(48 * time.Hour), 2},
		{"expired", now.Add(-24 * time.Hour), -1},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actual := daysUntil(scenario.future, now)
			if actual != scenario.expected {
				t.Errorf("Want %d but got %d", scenario.expected, actual)
			}
		})
	}
}
