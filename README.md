# RIP

> sed and grep had a baby

### Usage:
```
rip [REGEX] [TEMPLATE]
```

### Template:

Template variables are prefixed with `$`.

* `$debug` - dump all variables
* `$line` - the whole line
* `$0` - the whole match
* `$#` - capture group number
* `$name` - capture group name

If the template is omitted, it defaults to `$debug`.

### Examples:

```
$ ls | rip '^(?P<first_char>.)' '$first_char'
$ ss -p | rip 'users:\(\("(\w+)' '$1'
```
