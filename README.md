# RIP

> sed and grep had a baby

### Usage:
```
rip [REGEX] [TEMPLATE]
```

### Template:

Template variables are prefixed with `$`.

* `$line` - the whole line
* `$#` - capture group number
* `$name` - capture group name

### Examples:

```
$ ls | rip '^(?P<first_char>.)' '$first_char'
$ ss -p | rip 'users:\(\("(\w+)' '$1'
```
