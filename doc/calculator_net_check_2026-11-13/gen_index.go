//go:build ignore

// gen_index regenerates INDEX.md from the code (D-070), the way testing/index_test.go reads it: every ID'd
// declaration's doc comment gives Time, Space, Source and Tags (repeated labels joined by "; "), the package
// doc.go line gives Does, a literal "|" is written "\|", REMOVED rows are kept as they stand, rows sorted by ID.
// The preamble (everything up to the table's |---| line) is kept unchanged.
//
// Run from the repo root: go run doc/calculator_net_check_2026-11-13/gen_index.go <out path>
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	idRe       = regexp.MustCompile(`^([A-Z]{3})-([A-Z0-9]{1,6})-([0-9]{3})$`)
	labelRe    = regexp.MustCompile(`^([A-Z][A-Za-z]*):\s+(.*)$`)
	docEntryRe = regexp.MustCompile(`^\t([A-Z]{3}-[A-Z0-9]{1,6}-[0-9]{3})\s+(\S+)\s*(.*)$`)
	labels     = map[string]bool{"ID": true, "Time": true, "Space": true, "Source": true, "Tags": true,
		"Units": true, "Range": true, "Example": true, "Verified": true, "Implementation": true}
	groups = []string{"algo", "num", "unit", "sci", "eng", "fin", "med", "life"}
)

func parseFields(doc string) map[string][]string {
	f := map[string][]string{}
	cur := ""
	for _, ln := range strings.Split(doc, "\n") {
		if m := labelRe.FindStringSubmatch(ln); m != nil && labels[m[1]] {
			cur = m[1]
			f[cur] = append(f[cur], strings.TrimSpace(m[2]))
			continue
		}
		if cur != "" && strings.HasPrefix(ln, "  ") && strings.TrimSpace(ln) != "" {
			l := f[cur]
			l[len(l)-1] += " " + strings.TrimSpace(ln)
			continue
		}
		cur = ""
	}
	return f
}

func receiverType(fd *ast.FuncDecl) string {
	t := fd.Recv.List[0].Type
	if s, ok := t.(*ast.StarExpr); ok {
		t = s.X
	}
	switch x := t.(type) {
	case *ast.IndexExpr:
		t = x.X
	case *ast.IndexListExpr:
		t = x.X
	}
	if id, ok := t.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

type entry struct {
	id, name, pkg string
	f             map[string][]string
}

func esc(s string) string { return strings.ReplaceAll(s, "|", `\|`) }

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: gen_index.go <out path>")
		os.Exit(2)
	}
	fset := token.NewFileSet()
	var entries []entry
	docLine := map[string]string{}
	for _, g := range groups {
		if _, err := os.Stat(g); err != nil {
			continue
		}
		err := filepath.WalkDir(g, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
				return nil
			}
			pkg := filepath.ToSlash(filepath.Dir(p))
			af, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			if filepath.Base(p) == "doc.go" && af.Doc != nil {
				for _, ln := range strings.Split(af.Doc.Text(), "\n") {
					if m := docEntryRe.FindStringSubmatch(ln); m != nil {
						docLine[m[1]] = strings.TrimSpace(m[3])
					}
				}
			}
			add := func(doc *ast.CommentGroup, name string) {
				if doc == nil {
					return
				}
				f := parseFields(doc.Text())
				if len(f["ID"]) == 0 {
					return
				}
				entries = append(entries, entry{f["ID"][0], name, pkg, f})
			}
			for _, decl := range af.Decls {
				switch dcl := decl.(type) {
				case *ast.GenDecl:
					for _, sp := range dcl.Specs {
						var doc *ast.CommentGroup
						var name string
						switch s := sp.(type) {
						case *ast.TypeSpec:
							doc, name = s.Doc, s.Name.Name
						case *ast.ValueSpec:
							doc, name = s.Doc, s.Names[0].Name
						default:
							continue
						}
						if doc == nil && len(dcl.Specs) == 1 {
							doc = dcl.Doc
						}
						add(doc, name)
					}
				case *ast.FuncDecl:
					if !dcl.Name.IsExported() {
						continue
					}
					name := dcl.Name.Name
					if dcl.Recv != nil {
						r := receiverType(dcl)
						if !ast.IsExported(r) {
							continue
						}
						name = r + "." + name
					}
					add(dcl.Doc, name)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	old, err := os.ReadFile("INDEX.md")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lines := strings.Split(string(old), "\n")
	var preamble []string
	rows := map[string]string{}
	inTable := false
	for _, ln := range lines {
		if !inTable {
			preamble = append(preamble, ln)
			if strings.HasPrefix(ln, "|---|") {
				inTable = true
			}
			continue
		}
		cells := strings.SplitN(strings.TrimPrefix(ln, "| "), " |", 2)
		if len(cells) == 2 && idRe.MatchString(cells[0]) && strings.Contains(ln, "REMOVED") {
			rows[cells[0]] = ln
		}
	}
	code := map[string]bool{}
	for _, e := range entries {
		if code[e.id] {
			fmt.Fprintln(os.Stderr, "duplicate ID in code:", e.id)
			os.Exit(1)
		}
		code[e.id] = true
		if _, ok := rows[e.id]; ok {
			fmt.Fprintln(os.Stderr, "REMOVED row exists in code:", e.id)
		}
		j := func(l string) string { return esc(strings.Join(e.f[l], "; ")) }
		rows[e.id] = fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |", e.id, esc(e.name), e.pkg,
			esc(docLine[e.id]), j("Time"), j("Space"), j("Source"), j("Tags"))
	}
	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var b strings.Builder
	b.WriteString(strings.Join(preamble, "\n") + "\n")
	for _, id := range ids {
		b.WriteString(rows[id] + "\n")
	}
	if err := os.WriteFile(os.Args[1], []byte(b.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("code IDs %d, rows %d (%d REMOVED)\n", len(entries), len(ids), len(ids)-len(entries))
}
