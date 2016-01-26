# RIP

> converts **r**egex **i**nto **p**attern

### Usage:
```
rip [REGEX] [PATTERN]
```

### Demo:

![](http://i.imgur.com/YaBAlRQ.gif)

### Pattern:

Output pattern variables are prefixed with `$`.

* `$debug` - dump all variables
* `$line` - the whole line
* `$0` - the whole match
* `$#` - capture group number
* `$name` - capture group name

If the pattern is omitted, it defaults to `$debug`.

### Defaults:

Invoking rip without any arguments is equivalent to doing

```
$ rip '.*' '$debug'
```

### FAQ

Q. Why would I use this instead of `sed`?
A. `$debug` makes it easier to incrementally build up your regex.

Q. How do I install it?
A. `go get github.com/icholy/rip`
