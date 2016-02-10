package main

import (
	"bufio"
	"bytes"
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
	Count   int
}

func (d *TemplateData) Debug() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n%s\n", d.Line, strings.Repeat("-", len(d.Line)))
	for i, v := range d.Vars {
		if len(v) == 0 {
			fmt.Fprintf(&buf, "%%%d = %s\n", i, d.Matches[i])
		} else {
			fmt.Fprintf(&buf, "%%%s = %s\n", v, d.Matches[i])
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
			return 0, fmt.Errorf("%%%s exceedes the number of subexpressions", name)
		}
		return i, nil
	}
	for i := 0; i < len(vars); i++ {
		if vars[i] == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("%%%s does not correspond to any subexpression", name)
}

func replaceVars(s string, f func(string) (string, error)) (string, error) {
	var (
		regex   = regexp.MustCompile(`%(:?([\w\d]+)|{([\w\d]+)})`)
		matches = regex.FindAllStringSubmatchIndex(s, -1)
		index   = 0
		buffer  bytes.Buffer
	)
	for _, m := range matches {
		var name string
		if m[4] != -1 {
			name = s[m[4]:m[5]]
		} else {
			name = s[m[6]:m[7]]
		}
		if !isEscaped(s, m[0]) {
			buffer.WriteString(s[index:m[0]])
			if replace, err := f(name); err != nil {
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

func compilePattern(pattern string, vars []string) (*template.Template, error) {
	tstr, err := replaceVars(pattern, func(name string) (string, error) {
		switch name {
		case "line":
			return "{{.Line}}", nil
		case "debug":
			return "{{.Debug}}", nil
		case "count":
			return "{{.Count}}", nil
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
	defaultRegex   = ".*"
	defaultPattern = "%0"
)

func init() {
	if envPattern, ok := os.LookupEnv("RIP_PATTERN"); ok {
		defaultPattern = envPattern
	}
	if envRegex, ok := os.LookupEnv("RIP_REGEX"); ok {
		defaultRegex = envRegex
	}
}

func main() {

	args := os.Args[1:]

	regex := defaultRegex
	if len(args) > 0 {
		regex = args[0]
	}

	pattern := defaultPattern
	if len(args) > 1 {
		pattern = args[1]
	}

	// compile the regex
	re, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println(err)
		return
	}

	// compile the template
	vars := re.SubexpNames()
	templ, err := compilePattern(pattern, vars)
	if err != nil {
		fmt.Printf("failed to parse template: %s\n", err)
		return
	}

	// apply transformation
	scanner := bufio.NewScanner(os.Stdin)
	var count int
	for scanner.Scan() {
		var (
			line    = scanner.Text()
			matches = re.FindAllStringSubmatch(line, -1)
		)
		for _, match := range matches {
			if len(match) == 0 {
				continue
			}
			count++
			if err := templ.Execute(os.Stdout, &TemplateData{
				Matches: match,
				Line:    line,
				Count:   count,
				Vars:    vars,
			}); err != nil {
				fmt.Printf("failed to populate template: %s\n", err)
				return
			}
			os.Stdout.WriteString("\n")
		}
	}
}
