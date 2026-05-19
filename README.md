# pulse

Pulse is a small command-line health check runner. A policy file describes one or more checks, Pulse runs them, and each check can include assertions for the expected result.

Policy files are JSON.

## Basic Usage

Build the binary:

```sh
go build -o pulse
```

Run a policy:

```sh
./pulse run example/policy-basic.json
```

Run checks with a custom concurrency level:

```sh
./pulse run --concurrency 8 example/policy-multi.json
```

## Policy Shape

A policy has a name and a list of checks:

```json
{
  "name": "basic policy",
  "checks": [
    {
      "type": "http",
      "url": "https://example.com"
    }
  ]
}
```

Each check must include a `type`. Supported types are:

- `http`
- `dns`
- `tls`

## Examples

### HTTP Status And Body

```json
{
  "name": "http example",
  "checks": [
    {
      "type": "http",
      "url": "https://example.com",
      "assertions": [
        {
          "statusCode": {
            "between": [200, 299]
          },
          "body": {
            "contains": "Example Domain"
          },
          "headers": {
            "Content-Type": {
              "matches": "text/html.*"
            }
          }
        }
      ]
    }
  ]
}
```

### DNS TXT And MX

```json
{
  "name": "dns example",
  "checks": [
    {
      "type": "dns",
      "host": "github.com",
      "assertions": [
        {
          "txt": {
            "contains": "docusign=087098e3-3d46-47b7-9b4e-8a23028154cd"
          },
          "mx": {
            "contains": "github-com.mail.protection.outlook.com.",
            "length": {
              "equals": 1
            }
          }
        }
      ]
    }
  ]
}
```

### List Matching

List matchers are useful for DNS records. `contains` and `notContains` match exact list elements. `any` and `all` apply a string matcher to list elements.

```json
{
  "type": "dns",
  "host": "example.com",
  "assertions": [
    {
      "ns": {
        "length": {
          "gte": 1
        },
        "any": {
          "matches": ".*\\.iana-servers\\.net\\."
        }
      }
    }
  ]
}
```

### TLS Certificate And Protocols

```json
{
  "name": "tls example",
  "checks": [
    {
      "type": "tls",
      "host": "example.com",
      "assertions": [
        {
          "daysRemaining": {
            "gte": 14
          },
          "supportedVersions": {
            "contains": "TLS 1.3"
          },
          "supportedCiphers": {
            "notContains": "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
          }
        }
      ]
    }
  ]
}
```

## Check Reference

### HTTP Check

```json
{
  "type": "http",
  "url": "https://example.com",
  "method": "GET",
  "body": "base64-encoded request body",
  "assertions": []
}
```

Fields:

- `type`: must be `http`.
- `url`: required. Must use `http` or `https`.
- `method`: optional. Defaults to `GET`. Supported methods are `DELETE`, `GET`, `HEAD`, `OPTIONS`, `PATCH`, `PUT`, and `TRACE`.
- `body`: optional request body bytes. In JSON, Go unmarshals `[]byte` from a base64-encoded string.
- `assertions`: optional list of HTTP assertion objects.

HTTP assertion object:

```json
{
  "statusCode": {
    "equals": 200
  },
  "body": {
    "contains": "Example Domain"
  },
  "headers": {
    "Content-Type": {
      "matches": "text/html.*"
    }
  }
}
```

HTTP assertion fields:

- `statusCode`: number matcher applied to the response status code.
- `body`: string matcher applied to the response body.
- `headers`: object whose keys are header names and whose values are string matchers.

### DNS Check

```json
{
  "type": "dns",
  "host": "github.com",
  "assertions": []
}
```

Fields:

- `type`: must be `dns`.
- `host`: required.
- `assertions`: list of DNS assertion objects.

DNS assertion object:

```json
{
  "cname": {
    "equals": "example.com."
  },
  "mx": {
    "contains": "mail.example.com."
  },
  "txt": {
    "any": {
      "contains": "v=spf1"
    }
  },
  "ns": {
    "length": {
      "gte": 1
    }
  },
  "a": {
    "contains": "example.com."
  }
}
```

DNS assertion fields:

- `cname`: string matcher applied to `net.LookupCNAME(host)`.
- `mx`: string list matcher applied to MX record hosts.
- `txt`: string list matcher applied to TXT records.
- `ns`: string list matcher applied to NS record hosts.
- `a`: string list matcher applied to `net.LookupAddr(host)` results.

Note: the current `a` implementation uses reverse lookup via `net.LookupAddr`, so it expects an address-like host and returns names. For forward A/AAAA address lookup, the implementation should use `net.LookupHost` or `net.LookupIP`.

### TLS Check

```json
{
  "type": "tls",
  "host": "example.com",
  "port": 443,
  "assertions": []
}
```

Fields:

- `type`: must be `tls`.
- `host`: required. Used as both the TCP host and TLS server name.
- `port`: optional. Defaults to `443`.
- `assertions`: optional list of TLS assertion objects.

TLS assertion object:

```json
{
  "daysRemaining": {
    "gte": 14
  },
  "supportedVersions": {
    "contains": "TLS 1.3"
  },
  "supportedCiphers": {
    "notContains": "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
  }
}
```

TLS assertion fields:

- `daysRemaining`: number matcher applied to the leaf certificate's remaining lifetime in whole days.
- `supportedVersions`: string list matcher applied to supported TLS protocol versions.
- `supportedCiphers`: string list matcher applied to supported TLS cipher suite names.

TLS version names are:

- `TLS 1.0`
- `TLS 1.1`
- `TLS 1.2`
- `TLS 1.3`

Cipher names use Go's `crypto/tls` cipher suite names, such as `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`.

Note: Go can force individual cipher suites for TLS 1.0 through TLS 1.2. TLS 1.3 cipher suites are not configurable in the same way, so Pulse reports the negotiated TLS 1.3 cipher when TLS 1.3 is supported.

## Matcher Reference

### Number Matcher

Used for numeric values, such as HTTP status codes and list lengths.

```json
{
  "equals": 200,
  "notEquals": 500,
  "gt": 199,
  "gte": 200,
  "lt": 300,
  "lte": 299,
  "in": [200, 201, 204],
  "notIn": [500, 502, 503],
  "between": [200, 299]
}
```

Fields:

- `equals`: value must equal this number.
- `notEquals`: value must not equal this number.
- `gt`: value must be greater than this number.
- `gte`: value must be greater than or equal to this number.
- `lt`: value must be less than this number.
- `lte`: value must be less than or equal to this number.
- `in`: value must be one of these numbers.
- `notIn`: value must not be one of these numbers.
- `between`: value must be between two numbers, inclusive. The array must have exactly two values: `[min, max]`.

Number matcher fields are optional and can be combined.

### String Matcher

Used for single string values, such as HTTP response bodies, headers, and CNAME records.

```json
{
  "equals": "exact value",
  "notEquals": "forbidden value",
  "contains": "substring",
  "matches": "regular expression"
}
```

Fields:

- `equals`: string must exactly equal this value.
- `notEquals`: string must not exactly equal this value.
- `contains`: string must contain this substring.
- `matches`: string must match this Go regular expression.

String matching is case-sensitive. `matches` uses Go's regular expression syntax.

### String List Matcher

Used for lists of strings, such as DNS TXT, MX, and NS records.

```json
{
  "contains": "exact list element",
  "notContains": "forbidden list element",
  "length": {
    "gte": 1
  },
  "any": {
    "contains": "substring"
  },
  "all": {
    "matches": "^[a-z0-9.-]+\\.$"
  }
}
```

Fields:

- `contains`: list must contain this exact string as an element.
- `notContains`: list must not contain this exact string as an element.
- `length`: number matcher applied to the list length.
- `any`: string matcher that must match at least one element.
- `all`: string matcher that must match every element.

`contains` is exact element matching. To search within list elements, use `any` with a string matcher:

```json
{
  "txt": {
    "any": {
      "contains": "v=spf1"
    }
  }
}
```

## Full Schema

This is an informal JSON schema matching the current implementation:

```json
{
  "name": "string",
  "checks": [
    {
      "type": "http | dns | tls"
    }
  ]
}
```

HTTP check:

```json
{
  "type": "http",
  "url": "string, required",
  "method": "string, optional",
  "body": "base64 string, optional",
  "assertions": [
    {
      "statusCode": "NumberMatcher",
      "body": "StringMatcher",
      "headers": {
        "Header-Name": "StringMatcher"
      }
    }
  ]
}
```

DNS check:

```json
{
  "type": "dns",
  "host": "string, required",
  "assertions": [
    {
      "cname": "StringMatcher",
      "mx": "StringListMatcher",
      "txt": "StringListMatcher",
      "ns": "StringListMatcher",
      "a": "StringListMatcher"
    }
  ]
}
```

TLS check:

```json
{
  "type": "tls",
  "host": "string, required",
  "port": "number, optional",
  "assertions": [
    {
      "daysRemaining": "NumberMatcher",
      "supportedVersions": "StringListMatcher",
      "supportedCiphers": "StringListMatcher"
    }
  ]
}
```

Number matcher:

```json
{
  "equals": "number",
  "notEquals": "number",
  "gt": "number",
  "gte": "number",
  "lt": "number",
  "lte": "number",
  "in": ["number"],
  "notIn": ["number"],
  "between": ["number", "number"]
}
```

String matcher:

```json
{
  "equals": "string",
  "notEquals": "string",
  "contains": "string",
  "matches": "string"
}
```

String list matcher:

```json
{
  "contains": "string",
  "notContains": "string",
  "length": "NumberMatcher",
  "any": "StringMatcher",
  "all": "StringMatcher"
}
```
