package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type TemplateData struct {
	Matches []string
	Line    string
	Vars    []string
}

func (d *TemplateData) Debug() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n%s\n", d.Line, strings.Repeat("-", len(d.Line)))
	for i, v := range d.Vars {
		if len(v) == 0 {
			fmt.Fprintf(&buf, "$%d = %s\n", i, d.Matches[i])
		} else {
			fmt.Fprintf(&buf, "$%s = %s\n", v, d.Matches[i])
		}
	}
	return buf.String()
}

func isEscaped(s string, pos int) bool {
	slashes := 0
	for i := pos - 1; i >= 0; i-- {
		if s[i] == '\\' {
			slashes++
		} else {
			break
		}
	}
	return slashes%2 != 0
}

func varToIndex(vars []string, name string) (int, error) {
	if i, err := strconv.Atoi(name); err == nil {
		if i >= len(vars) {
			return 0, fmt.Errorf("$%s exceedes the number of subexpressions", name)
		}
		return i, nil
	}
	for i := 0; i < len(vars); i++ {
		if vars[i] == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("$%s does not correspond to any subexpression", name)
}

func replaceVars(s string, f func(string) (string, error)) (string, error) {
	var (
		regex   = regexp.MustCompile(`\$[\w\d]+`)
		matches = regex.FindAllStringIndex(s, -1)
		index   = 0
		buffer  bytes.Buffer
	)
	for _, m := range matches {
		if !isEscaped(s, m[0]) {
			buffer.WriteString(s[index:m[0]])
			if replace, err := f(s[m[0]+1 : m[1]]); err != nil {
				return "", err
			} else {
				buffer.WriteString(replace)
			}
			index = m[1]
		}
	}
	buffer.WriteString(s[index:])
	return buffer.String(), nil
}

func compileTemplate(pattern string, vars []string) (*template.Template, error) {
	tstr, err := replaceVars(pattern, func(name string) (string, error) {
		switch name {
		case "line":
			return "{{.Line}}", nil
		case "debug":
			return "{{.Debug}}", nil
		}
		index, err := varToIndex(vars, name)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("{{index .Matches %d}}", index), nil
	})
	if err != nil {
		return nil, err
	}
	return template.New("").Parse(tstr)
}

var (
	help = flag.Bool("h", false, "show help")
)

func main() {

	flag.Parse()

	if *help {
		fmt.Println("rip [REGEX] [PATTERN]")
		return
	}

	args := flag.Args()

	regex := ".*"
	if len(args) > 0 {
		regex = args[0]
	}

	pattern := "$debug"
	if len(args) > 1 {
		pattern = args[1]
	}

	// compile the regex
	re, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// compile the pattern
	vars := re.SubexpNames()
	patternTemplate, err := compileTemplate(pattern, vars)
	if err != nil {
		fmt.Printf("failed to parse template: %s\n", err)
		return
	}

	empty := make([]string, len(vars))

	// apply transformation
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var (
			line    = scanner.Text()
			matches = re.FindStringSubmatch(line)
		)
		if len(matches) == 0 {
			matches = empty
		}
		if err := patternTemplate.Execute(os.Stdout, &TemplateData{
			Matches: matches,
			Line:    line,
			Vars:    vars,
		}); err != nil {
			fmt.Printf("failed to populate template: %s\n", err)
			return
		}
		os.Stdout.WriteString("\n")
	}
}
