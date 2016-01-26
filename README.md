# RIP

> converts **r**egex **i**nto **p**attern

### Usage:
```
rip [REGEX] [PATTERN]
```

### Pattern:

Output pattern variables are prefixed with `$`.

* `$debug` - dump all variables
* `$line` - the whole line
* `$0` - the whole match
* `$#` - capture group number
* `$name` - capture group name

If the pattern is omitted, it defaults to `$debug`.

### Examples:

```
$ ls | rip '^(?P<first_char>.)' '$first_char'
$ ss -p | rip 'users:\(\("(\w+)' '$1'
```
