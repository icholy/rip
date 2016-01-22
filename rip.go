package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"text/template"
)

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

func compileTemplate(tmplStr string, vars []string) (*template.Template, error) {
	tstr, err := replaceVars(tmplStr, func(name string) (string, error) {
		if name == "line" {
			return "{{.Line}}", nil
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

func main() {

	flag.Parse()

	args := flag.Args()

	pattern := ".*"
	if len(args) > 0 {
		pattern = args[0]
	}

	output := "$0"
	if len(args) > 1 {
		output = args[1]
	}

	// compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println(err)
		return
	}

	// compile the template
	templ, err := compileTemplate(output, re.SubexpNames())
	if err != nil {
		fmt.Printf("failed to parse template: %s\n", err)
		return
	}

	// apply transformation
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var (
			line    = scanner.Text()
			matches = re.FindStringSubmatch(line)
		)
		if len(matches) == 0 {
			continue
		}
		if err := templ.Execute(os.Stdout, struct {
			Matches []string
			Line    string
		}{
			Matches: matches,
			Line:    line,
		}); err != nil {
			fmt.Printf("failed to populate template: %s\n", err)
			return
		}
		os.Stdout.WriteString("\n")
	}
}
