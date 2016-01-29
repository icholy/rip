# RIP

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

If the pattern is omitted, it defaults to `$debug`.  
Variable names can be isolated with braces. `${debug}`

### Defaults:

Invoking rip without any arguments is equivalent to doing

```
$ rip '.*' '$debug'
```

You can change the default pattern by setting the `RIP_PATTERN` env variable.
Similarly, the default regex can be changed with the `RIP_REGEX` variable.


```
$ export RIP_PATTERN='$0'
```

### FAQ

**Q.** Why would I use this instead of `sed`?  
**A.** `$debug` makes it easier to incrementally build up your regex.

**Q.** How do I install it?  
**A.** Download [here](https://github.com/icholy/rip/releases)
