# Go test naming convention linter

```sh
go install github.com/devnev/testnamestyle/cmd/go-testnamestyle
go-testnamestyle -rule 'Test(?:[^_]+)?(?:_[a-z][^_]*)?' all
```

* no rules are configured by default
* `-rule` may be specified multiple times to add further rules
* if any rule matches, the test name is accepted
* rules are regexp patterns and must match the entire test name
* for rules with one capture group, the captured string must match a struct type
  or function declaration in the package
* for rules with two capture groups, the first captured string must match a
  struct type declaration in the package, and the second captured string must
  match a method of that type
