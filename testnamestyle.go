package testnamestyle

import (
	"errors"
	"go/types"
	"regexp"

	"golang.org/x/tools/go/analysis"
)

var (
	rules         = []*regexp.Regexp{}
	TestNameStyle = &analysis.Analyzer{
		Name: "testnamestyle",
		Doc:  "check test naming style",
		Run:  run,
		Requires: []*analysis.Analyzer{
			TestFuncs,
		},
	}
)

func init() {
	TestNameStyle.Flags.Func("rule", "add a rule to check", func(s string) error {
		r, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		if r.NumSubexp() > 3 {
			return errors.New("at most two match groups may be used")
		}
		for _, name := range r.SubexpNames() {
			if name != "" {
				return errors.New("named match group are not permitted")
			}
		}
		rules = append(rules, r)
		return nil
	})
}

func run(pass *analysis.Pass) (interface{}, error) {
	testFuncs := pass.ResultOf[TestFuncs].(TestFuncList)
	if len(rules) == 0 {
		return nil, nil
	}

NextFunc:
	for _, fn := range testFuncs {
		name := fn.Name.Name
		for _, rule := range rules {
			if matchRule(rule, name, pass.Pkg) {
				continue NextFunc
			}
		}
		pass.ReportRangef(fn.Name, "test function %s does not match any rules", name)
	}
	return nil, nil
}

func matchRule(rule *regexp.Regexp, name string, pkg *types.Package) bool {
	match := rule.FindStringSubmatchIndex(name)
	if match == nil || match[0] != 0 || match[1] != len(name) {
		return false
	}

	part1 := string(rule.ExpandString(nil, "$1", name, match))
	part2 := string(rule.ExpandString(nil, "$2", name, match))
	if part1 == "" && part2 == "" {
		return true
	}

	switch decl := pkg.Scope().Lookup(part1).(type) {
	case *types.TypeName:
		namedType, _ := decl.Type().(*types.Named)
		structType, _ := namedType.Underlying().(*types.Struct)
		if structType == nil {
			return false
		}
		if part2 == "" {
			return true
		}
		funcObj, _, _ := types.LookupFieldOrMethod(namedType, true, pkg, part2)
		return funcObj != nil
	case *types.Func:
		return part2 == ""
	case nil:
		return false
	default:
		return false
	}
}
