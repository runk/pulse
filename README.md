# pulse

Pulse is a monitoring tool that allows you to define checks in a JSON/YAML configuration file. Each check can be of a specific type (e.g., HTTP, DNS) and includes assertions to validate the expected outcomes.

It's first and foremost a command line monitoring tool, but it can also be used for testing and validation purposes. The configuration file is flexible and allows you to define various checks with different assertions based on your needs.

## Assertions

All checks are defined in the `assert` section of the configuration file. Each check consists of a domain-specific field, a matcher, and an expected value. The structure is as follows:

```yaml
assert:
  <domain-specific-field>:
    <matcher>: <expected>
```

With vocabulary for matchers as follows:

```yaml
equals: value
notEquals: value
contains: value
matches: regex
exists: true
in: [a, b, c]
gt: number
gte: number
lt: number
lte: number
between: [min, max]
all: [...]
any: [...]
not: ...
```

Domain specific fields vary depending on the check type, consult with the reference.

## Check Types

### HTTP Check

### DNS Check
