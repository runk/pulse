package check

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/runk/pulse/internal/assertion"
)

const defaultTLSPort = 443

type TLSAssertion struct {
	DaysRemaining     *assertion.NumberMatcher     `json:"daysRemaining,omitempty"`
	SupportedVersions *assertion.StringListMatcher `json:"supportedVersions,omitempty"`
	SupportedCiphers  *assertion.StringListMatcher `json:"supportedCiphers,omitempty"`
}

type TLSCheck struct {
	Host       string         `json:"host"`
	Port       uint16         `json:"port,omitempty"`
	Assertions []TLSAssertion `json:"assertions,omitempty"`
}

func (TLSCheck) Type() string { return "tls" }

func (c *TLSCheck) Validate() error {
	if c.Host == "" {
		return errors.New("host is required to perform tls check")
	}

	if c.Port == 0 {
		c.Port = defaultTLSPort
	}

	return nil
}

func (c TLSCheck) Run(ctx context.Context) error {
	addr := net.JoinHostPort(c.Host, strconv.Itoa(int(c.Port)))
	needsDaysRemaining := len(c.Assertions) == 0
	needsSupportedVersions := false
	needsSupportedCiphers := false

	for _, assertion := range c.Assertions {
		needsDaysRemaining = needsDaysRemaining || assertion.DaysRemaining != nil
		needsSupportedVersions = needsSupportedVersions || assertion.SupportedVersions != nil
		needsSupportedCiphers = needsSupportedCiphers || assertion.SupportedCiphers != nil
	}

	var daysRemaining int
	if needsDaysRemaining {
		conn, err := dialTLS(ctx, addr, c.Host, 0, 0, nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		state := conn.ConnectionState()
		if len(state.PeerCertificates) == 0 {
			return errors.New("tls connection did not return peer certificates")
		}

		daysRemaining = daysUntil(state.PeerCertificates[0].NotAfter, time.Now())
	}

	var err error
	var supportedVersions []string
	if needsSupportedVersions {
		supportedVersions, err = supportedTLSVersions(ctx, addr, c.Host)
		if err != nil {
			return err
		}
	}

	var supportedCiphers []string
	if needsSupportedCiphers {
		supportedCiphers, err = supportedTLSCiphers(ctx, addr, c.Host)
		if err != nil {
			return err
		}
	}

	errs := []error{}
	for _, assertion := range c.Assertions {
		if assertion.DaysRemaining != nil {
			if err := assertion.DaysRemaining.Match(daysRemaining); err != nil {
				errs = append(errs, err)
			}
		}

		if assertion.SupportedVersions != nil {
			if err := assertion.SupportedVersions.Match(supportedVersions); err != nil {
				errs = append(errs, err)
			}
		}

		if assertion.SupportedCiphers != nil {
			if err := assertion.SupportedCiphers.Match(supportedCiphers); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func dialTLS(ctx context.Context, addr, serverName string, minVersion, maxVersion uint16, cipherSuites []uint16) (*tls.Conn, error) {
	dialer := &tls.Dialer{
		NetDialer: &net.Dialer{},
		Config: &tls.Config{
			ServerName:         serverName,
			InsecureSkipVerify: true,
			MinVersion:         minVersion,
			MaxVersion:         maxVersion,
			CipherSuites:       cipherSuites,
		},
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		_ = conn.Close()
		return nil, fmt.Errorf("wanted *tls.Conn, but got %T", conn)
	}

	return tlsConn, nil
}

func supportedTLSVersions(ctx context.Context, addr, serverName string) ([]string, error) {
	versions := []uint16{
		tls.VersionTLS10,
		tls.VersionTLS11,
		tls.VersionTLS12,
		tls.VersionTLS13,
	}

	supported := []string{}
	for _, version := range versions {
		conn, err := dialTLS(ctx, addr, serverName, version, version, nil)
		if err != nil {
			// This error could be one of two
			// 1. version didn't work - we accept those as normal testing behaviour
			// 2. context errors - should exit immediately
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			continue
		}

		supported = append(supported, tlsVersionName(version))
		_ = conn.Close()
	}

	return supported, nil
}

func supportedTLSCiphers(ctx context.Context, addr, serverName string) ([]string, error) {
	seen := map[uint16]bool{}
	supported := []string{}

	for _, suite := range append(tls.CipherSuites(), tls.InsecureCipherSuites()...) {
		for _, version := range []uint16{tls.VersionTLS10, tls.VersionTLS11, tls.VersionTLS12} {
			conn, err := dialTLS(ctx, addr, serverName, version, version, []uint16{suite.ID})
			if err != nil {
				// This error could be one of two
				// 1. version/suite didn't work - we accept those as normal testing behaviour
				// 2. context errors - should exit immediately
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				continue
			}

			state := conn.ConnectionState()
			if !seen[state.CipherSuite] {
				seen[state.CipherSuite] = true
				supported = append(supported, tlsCipherSuiteName(state.CipherSuite))
			}

			_ = conn.Close()
		}
	}

	conn, err := dialTLS(ctx, addr, serverName, tls.VersionTLS13, tls.VersionTLS13, nil)
	if err == nil {
		state := conn.ConnectionState()
		if !seen[state.CipherSuite] {
			supported = append(supported, tlsCipherSuiteName(state.CipherSuite))
		}
		_ = conn.Close()
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return supported, nil
}

func daysUntil(future, now time.Time) int {
	return int(future.Sub(now).Hours() / 24)
}

func tlsVersionName(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("TLS 0x%04x", version)
	}
}

func tlsCipherSuiteName(id uint16) string {
	for _, suite := range append(tls.CipherSuites(), tls.InsecureCipherSuites()...) {
		if suite.ID == id {
			return suite.Name
		}
	}

	return fmt.Sprintf("0x%04x", id)
}
