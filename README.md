![](rip.png)

> Convert **r**egex **i**nto **p**attern

### Description:

Extract data from input using regular expressions.

### Usage:
```
rip [REGULAR EXPRESSION] [OUTPUT PATTERN]
```

### Demo:

![](http://i.imgur.com/1mpK75L.gif)

### Pattern:

Output pattern variables are prefixed with `$`.

* `$debug` - dump all variables
* `$line` - the whole line
* `$count` - current item number
* `$0` - the whole match
* `$#` - capture group number
* `$name` - capture group name

Variable names can be isolated with braces. `${debug}`.  
To insert a literal `$` into your pattern, escape it with a backslash.

### Defaults:

Invoking rip without any arguments is equivalent to doing

```
$ rip '.*' '$debug'
```

### Environment Variables

* `RIP_PATTERN` changes the default pattern (default=`$debug`).
* `RIP_REGEX` changes the default regex (default=`.*`).
* `RIP_PREFIX` changed the default variable prefix (default=`$`).

```
$ export RIP_PREFIX='%'
$ export RIP_PATTERN='%debug'
```

### FAQ

**Q.** Why would I use this instead of `sed`?  
**A.** `$debug` makes it easier to incrementally build up your regex.

**Q.** How do I install it?  
**A.** `go install github.com/icholy/rip@latest`
