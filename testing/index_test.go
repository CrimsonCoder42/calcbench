// Package index_test keeps the three places that describe a function in agreement: the doc comment, the
// package doc.go and INDEX.md, plus doc/taxonomy.md against the tree (.claude/CLAUDE.md "Finding things").
// It reads files only; it imports no calcbench package.
package index_test

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

const root = ".."

// groupPrefix maps a top-level folder to its ID prefix.
var groupPrefix = map[string]string{
	"algo": "ALG", "num": "NUM", "unit": "UNT", "sci": "SCI", "eng": "ENG", "fin": "FIN", "med": "MED", "life": "LIF",
}

// calculatorGroups are the groups whose functions are calculators and so also carry Units: (go-functions.md).
var calculatorGroups = map[string]bool{"unit": true, "sci": true, "eng": true, "fin": true, "med": true, "life": true}

var (
	idRe       = regexp.MustCompile(`^([A-Z]{3})-([A-Z0-9]{1,6})-([0-9]{3})$`)
	labelRe    = regexp.MustCompile(`^([A-Z][A-Za-z]*):\s+(.*)$`)
	docEntryRe = regexp.MustCompile(`^\t([A-Z]{3}-[A-Z0-9]{1,6}-[0-9]{3})\s+(\S+)\s*(.*)$`)
	shortRe    = regexp.MustCompile(`ID short form: ([A-Z0-9]+)`)
	labels     = map[string]bool{"ID": true, "Time": true, "Space": true, "Source": true, "Tags": true,
		"Units": true, "Range": true, "Example": true, "Verified": true, "Implementation": true}
)

// entry is one ID'd declaration: a function, a method, a type or a constant.
type entry struct {
	pkg, file, name string
	line            int
	fields          map[string][]string
	order           []string
}

func (e entry) where() string { return fmt.Sprintf("%s:%d %s", e.file, e.line, e.name) }

func (e entry) id() string {
	if v := e.fields["ID"]; len(v) > 0 {
		return v[0]
	}
	return ""
}

// parseFields reads the labelled lines of a doc comment. A line indented by two or more spaces after a
// labelled line continues it (joined with one space).
func parseFields(doc string) (map[string][]string, []string) {
	f := map[string][]string{}
	var order []string
	cur := ""
	for _, ln := range strings.Split(doc, "\n") {
		if m := labelRe.FindStringSubmatch(ln); m != nil && labels[m[1]] {
			cur = m[1]
			f[cur] = append(f[cur], strings.TrimSpace(m[2]))
			order = append(order, cur)
			continue
		}
		if cur != "" && strings.HasPrefix(ln, "  ") && strings.TrimSpace(ln) != "" {
			l := f[cur]
			l[len(l)-1] += " " + strings.TrimSpace(ln)
			continue
		}
		cur = ""
	}
	return f, order
}

// tree is everything the checks need, read once.
type tree struct {
	entries    []entry             // every ID'd declaration
	uncovered  []string            // exported functions and methods with no ID and no covering ID (D-071)
	pkgs       map[string]bool     // packages on disk (folders holding non-test .go files)
	short      map[string]string   // package → SHORT from doc.go
	docIDs     map[string][]string // package → IDs listed in doc.go
	docLine    map[string]string   // ID → doc.go "Does" text
	docName    map[string]string   // ID → doc.go name
	readErrors []string
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

// resultTypes names the types a function returns (pointer and generic wrappers removed).
func resultTypes(fd *ast.FuncDecl) []string {
	var out []string
	if fd.Type.Results == nil {
		return out
	}
	for _, r := range fd.Type.Results.List {
		t := r.Type
		for {
			switch x := t.(type) {
			case *ast.StarExpr:
				t = x.X
				continue
			case *ast.IndexExpr:
				t = x.X
				continue
			case *ast.IndexListExpr:
				t = x.X
				continue
			case *ast.ArrayType:
				t = x.Elt
				continue
			}
			break
		}
		if id, ok := t.(*ast.Ident); ok {
			out = append(out, id.Name)
		}
	}
	return out
}

func readTree(t *testing.T) *tree {
	t.Helper()
	tr := &tree{pkgs: map[string]bool{}, short: map[string]string{}, docIDs: map[string][]string{},
		docLine: map[string]string{}, docName: map[string]string{}}
	fset := token.NewFileSet()
	type pending struct {
		pkg, recv, name, where string
		results                []string
	}
	var noID []pending
	idTypes := map[string]bool{}  // pkg.Type carrying an ID
	retTypes := map[string]bool{} // pkg.Type returned by an ID'd function or method
	for group := range groupPrefix {
		dir := filepath.Join(root, group)
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
				return nil
			}
			rel, _ := filepath.Rel(root, p)
			rel = filepath.ToSlash(rel)
			pkg := filepath.ToSlash(filepath.Dir(rel))
			tr.pkgs[pkg] = true
			af, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
			if err != nil {
				tr.readErrors = append(tr.readErrors, err.Error())
				return nil
			}
			if filepath.Base(p) == "doc.go" && af.Doc != nil {
				for _, ln := range strings.Split(af.Doc.Text(), "\n") {
					if m := docEntryRe.FindStringSubmatch(ln); m != nil {
						tr.docIDs[pkg] = append(tr.docIDs[pkg], m[1])
						tr.docName[m[1]] = m[2]
						tr.docLine[m[1]] = strings.TrimSpace(m[3])
					}
					if m := shortRe.FindStringSubmatch(ln); m != nil {
						tr.short[pkg] = m[1]
					}
				}
			}
			for _, decl := range af.Decls {
				switch dcl := decl.(type) {
				case *ast.GenDecl:
					for _, sp := range dcl.Specs {
						var doc *ast.CommentGroup
						var name *ast.Ident
						switch s := sp.(type) {
						case *ast.TypeSpec:
							doc, name = s.Doc, s.Name
						case *ast.ValueSpec:
							doc, name = s.Doc, s.Names[0]
						default:
							continue
						}
						if doc == nil && len(dcl.Specs) == 1 {
							doc = dcl.Doc
						}
						if doc == nil {
							continue
						}
						f, order := parseFields(doc.Text())
						if len(f["ID"]) == 0 {
							continue
						}
						if _, ok := sp.(*ast.TypeSpec); ok {
							idTypes[pkg+"."+name.Name] = true
						}
						tr.entries = append(tr.entries, entry{pkg: pkg, file: rel, name: name.Name,
							line: fset.Position(name.Pos()).Line, fields: f, order: order})
					}
				case *ast.FuncDecl:
					if !dcl.Name.IsExported() {
						continue
					}
					name, recv := dcl.Name.Name, ""
					if dcl.Recv != nil {
						recv = receiverType(dcl)
						if !ast.IsExported(recv) {
							continue
						}
						name = recv + "." + name
					}
					doc := ""
					if dcl.Doc != nil {
						doc = dcl.Doc.Text()
					}
					f, order := parseFields(doc)
					where := fmt.Sprintf("%s:%d %s", rel, fset.Position(dcl.Pos()).Line, name)
					if len(f["ID"]) == 0 {
						noID = append(noID, pending{pkg, recv, dcl.Name.Name, where, resultTypes(dcl)})
						continue
					}
					for _, rt := range resultTypes(dcl) {
						retTypes[pkg+"."+rt] = true
					}
					tr.entries = append(tr.entries, entry{pkg: pkg, file: rel, name: name,
						line: fset.Position(dcl.Pos()).Line, fields: f, order: order})
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}
	// D-071: a method is covered by its receiver type's ID or by a covered or ID'd function returning that
	// type; a constructor New<T> is covered by T's ID. Repeat until nothing new is covered.
	covered := make([]bool, len(noID))
	for changed := true; changed; {
		changed = false
		for i, p := range noID {
			if covered[i] {
				continue
			}
			if (p.recv != "" && (idTypes[p.pkg+"."+p.recv] || retTypes[p.pkg+"."+p.recv])) ||
				(p.recv == "" && strings.HasPrefix(p.name, "New") && idTypes[p.pkg+"."+strings.TrimPrefix(p.name, "New")]) {
				covered[i], changed = true, true
				for _, rt := range p.results {
					retTypes[p.pkg+"."+rt] = true
				}
			}
		}
	}
	for i, p := range noID {
		if !covered[i] {
			tr.uncovered = append(tr.uncovered, p.where)
		}
	}
	slices.SortFunc(tr.entries, func(a, b entry) int { return cmp.Compare(a.id(), b.id()) })
	return tr
}

// indexRow is one row of INDEX.md, cells unescaped.
type indexRow struct {
	line                                               int
	id, function, pkg, does, time, space, source, tags string
}

func (r indexRow) removed() bool { return strings.Contains(r.does, "REMOVED") }

// splitRow splits a markdown table row on "|" not preceded by a backslash and unescapes "\|".
func splitRow(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "|")
	s = strings.TrimSuffix(s, "|")
	var cells []string
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == '\\' && i+1 < len(s) && s[i+1] == '|':
			b.WriteByte('|')
			i++
		case s[i] == '|':
			cells = append(cells, strings.TrimSpace(b.String()))
			b.Reset()
		default:
			b.WriteByte(s[i])
		}
	}
	return append(cells, strings.TrimSpace(b.String()))
}

func readIndex(t *testing.T) []indexRow {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "INDEX.md"))
	if err != nil {
		t.Fatal(err)
	}
	var rows []indexRow
	for i, ln := range strings.Split(string(data), "\n") {
		c := splitRow(ln)
		if len(c) < 1 || !idRe.MatchString(c[0]) {
			continue
		}
		if len(c) != 8 {
			t.Errorf("INDEX.md:%d: %d cells, want 8", i+1, len(c))
			continue
		}
		rows = append(rows, indexRow{i + 1, c[0], c[1], c[2], c[3], c[4], c[5], c[6], c[7]})
	}
	return rows
}

// report fails the test with every offender on its own line.
func report(t *testing.T, what string, bad []string) {
	t.Helper()
	if len(bad) > 0 {
		t.Errorf("%s (%d):\n  %s", what, len(bad), strings.Join(bad, "\n  "))
	}
}

// (a) Every exported function has an ID line in the documented shape, or is covered by one (D-071).
func TestIndex_docCommentShape(t *testing.T) {
	tr := readTree(t)
	report(t, "parse errors", tr.readErrors)
	report(t, "exported functions with no ID", tr.uncovered)
	var bad []string
	for _, e := range tr.entries {
		id := e.id()
		m := idRe.FindStringSubmatch(id)
		group := strings.SplitN(e.pkg, "/", 2)[0]
		switch {
		case m == nil:
			bad = append(bad, e.where()+": ID "+id+" not <PREFIX>-<SHORT>-<NNN>")
			continue
		case m[1] != groupPrefix[group]:
			bad = append(bad, e.where()+": prefix "+m[1]+", want "+groupPrefix[group])
		case m[2] != tr.short[e.pkg]:
			bad = append(bad, e.where()+": SHORT "+m[2]+", doc.go says "+tr.short[e.pkg])
		}
		want := []string{"ID", "Time", "Space", "Source", "Tags"}
		if calculatorGroups[group] {
			want = append(want, "Units")
		}
		if group == "med" {
			want = append(want, "Range")
		}
		for _, l := range want {
			if len(e.fields[l]) == 0 || e.fields[l][0] == "" {
				bad = append(bad, e.where()+": missing "+l+":")
			}
		}
		// ID, Time, Space, Source, Tags appear in that order.
		pos := map[string]int{}
		for i, l := range e.order {
			if _, seen := pos[l]; !seen {
				pos[l] = i
			}
		}
		for i := 1; i < 5; i++ {
			a, b := want[i-1], want[i]
			if pa, ok := pos[a]; ok {
				if pb, ok := pos[b]; ok && pb < pa {
					bad = append(bad, e.where()+": "+b+": before "+a+":")
				}
			}
		}
		if len(e.fields["ID"]) > 1 {
			bad = append(bad, e.where()+": more than one ID line")
		}
	}
	report(t, "doc comments out of shape", bad)
}

// (e) No ID appears twice in code.
func TestIndex_noDuplicateIDs(t *testing.T) {
	tr := readTree(t)
	seen := map[string]string{}
	var bad []string
	for _, e := range tr.entries {
		if w, ok := seen[e.id()]; ok {
			bad = append(bad, e.id()+": "+w+" and "+e.where())
		}
		seen[e.id()] = e.where()
	}
	report(t, "IDs used twice in code", bad)
}

// (b) and (d) Every ID is in INDEX.md exactly once, and the row matches the doc comment word for word.
func TestIndex_indexAgreesWithCode(t *testing.T) {
	tr := readTree(t)
	rows := readIndex(t)
	byID := map[string][]indexRow{}
	for _, r := range rows {
		byID[r.id] = append(byID[r.id], r)
	}
	code := map[string]entry{}
	for _, e := range tr.entries {
		code[e.id()] = e
	}
	var missing, dup, stale, differ, unsorted []string
	for _, e := range tr.entries {
		switch n := len(byID[e.id()]); {
		case n == 0:
			missing = append(missing, e.id()+" "+e.where())
		case n > 1:
			dup = append(dup, e.id())
		}
	}
	for i, r := range rows {
		if i > 0 && rows[i-1].id >= r.id {
			unsorted = append(unsorted, fmt.Sprintf("INDEX.md:%d %s after %s", r.line, r.id, rows[i-1].id))
		}
		e, ok := code[r.id]
		if !ok {
			if !r.removed() {
				stale = append(stale, fmt.Sprintf("INDEX.md:%d %s (no code, not REMOVED)", r.line, r.id))
			}
			continue
		}
		if r.removed() {
			stale = append(stale, fmt.Sprintf("INDEX.md:%d %s is marked REMOVED but exists in code", r.line, r.id))
		}
		join := func(l string) string { return strings.Join(e.fields[l], "; ") }
		for _, c := range []struct{ col, got, want string }{
			{"Function", r.function, e.name}, {"Package", r.pkg, e.pkg}, {"Does", r.does, tr.docLine[r.id]},
			{"Time", r.time, join("Time")}, {"Space", r.space, join("Space")},
			{"Source", r.source, join("Source")}, {"Tags", r.tags, join("Tags")},
		} {
			if c.got != c.want {
				differ = append(differ, fmt.Sprintf("INDEX.md:%d %s %s: index %q, code %q", r.line, r.id, c.col, c.got, c.want))
			}
		}
	}
	report(t, "IDs missing from INDEX.md", missing)
	report(t, "IDs in INDEX.md more than once", dup)
	report(t, "INDEX.md rows with no code", stale)
	report(t, "INDEX.md cells that differ from the doc comment or doc.go", differ)
	report(t, "INDEX.md not sorted by ID", unsorted)
}

// (c) Every ID is listed exactly once in its own package's doc.go, under the same name.
func TestIndex_docGoAgreesWithCode(t *testing.T) {
	tr := readTree(t)
	rows := readIndex(t)
	removed := map[string]bool{}
	for _, r := range rows {
		if r.removed() {
			removed[r.id] = true
		}
	}
	code := map[string]entry{}
	for _, e := range tr.entries {
		code[e.id()] = e
	}
	var bad []string
	for _, e := range tr.entries {
		n := 0
		for _, id := range tr.docIDs[e.pkg] {
			if id == e.id() {
				n++
			}
		}
		if n != 1 {
			bad = append(bad, fmt.Sprintf("%s listed %d times in %s/doc.go", e.id(), n, e.pkg))
		}
		if tr.docName[e.id()] != e.name {
			bad = append(bad, fmt.Sprintf("%s: doc.go names %q, code %q", e.id(), tr.docName[e.id()], e.name))
		}
	}
	for pkg, ids := range tr.docIDs {
		for _, id := range ids {
			if e, ok := code[id]; ok && e.pkg != pkg {
				bad = append(bad, fmt.Sprintf("%s listed in %s/doc.go but lives in %s", id, pkg, e.pkg))
			}
			if _, ok := code[id]; !ok && !removed[id] {
				bad = append(bad, fmt.Sprintf("%s listed in %s/doc.go: no code and no REMOVED row", id, pkg))
			}
		}
	}
	for pkg := range tr.pkgs {
		if tr.short[pkg] == "" {
			bad = append(bad, pkg+"/doc.go: no \"ID short form:\" line")
		}
	}
	slices.Sort(bad)
	report(t, "doc.go lists out of step with the code", bad)
}

// (f) Every package on disk is a built row of doc/taxonomy.md with the same SHORT, and every built row exists.
func TestIndex_taxonomyAgreesWithTree(t *testing.T) {
	tr := readTree(t)
	data, err := os.ReadFile(filepath.Join(root, "doc", "taxonomy.md"))
	if err != nil {
		t.Fatal(err)
	}
	type row struct{ short, status string }
	tax := map[string]row{}
	pkgRe := regexp.MustCompile("^`([a-z0-9/]+)`$")
	for _, ln := range strings.Split(string(data), "\n") {
		c := splitRow(ln)
		if len(c) != 5 || groupPrefix[strings.SplitN(strings.Trim(c[1], "`"), "/", 2)[0]] == "" {
			continue
		}
		m := pkgRe.FindStringSubmatch(c[1])
		if m == nil {
			continue
		}
		tax[m[1]] = row{c[2], c[3]}
	}
	var bad []string
	for pkg := range tr.pkgs {
		r, ok := tax[pkg]
		switch {
		case !ok:
			bad = append(bad, pkg+": on disk, not in doc/taxonomy.md")
		case r.status != "built":
			bad = append(bad, pkg+": on disk, taxonomy status "+r.status)
		case r.short != tr.short[pkg]:
			bad = append(bad, pkg+": taxonomy SHORT "+r.short+", doc.go "+tr.short[pkg])
		}
	}
	for pkg, r := range tax {
		if r.status == "built" && !tr.pkgs[pkg] {
			bad = append(bad, pkg+": taxonomy says built, no package on disk")
		}
		if r.status != "built" && r.status != "planned" {
			bad = append(bad, pkg+": taxonomy status "+r.status+", want built or planned")
		}
	}
	slices.Sort(bad)
	report(t, "doc/taxonomy.md out of step with the tree", bad)
}

// TestIndex_counts logs the three counts the reports quote: functions | INDEX rows | doc.go entries.
func TestIndex_counts(t *testing.T) {
	tr := readTree(t)
	rows := readIndex(t)
	live := 0
	for _, r := range rows {
		if !r.removed() {
			live++
		}
	}
	docN := 0
	code := map[string]bool{}
	for _, e := range tr.entries {
		code[e.id()] = true
	}
	for _, ids := range tr.docIDs {
		for _, id := range ids {
			if code[id] {
				docN++
			}
		}
	}
	t.Logf("IDs in code %d | INDEX rows %d (+%d REMOVED) | doc.go entries %d | packages %d",
		len(tr.entries), live, len(rows)-live, docN, len(tr.pkgs))
	if len(tr.entries) != live || live != docN {
		t.Errorf("counts differ: code %d, INDEX %d, doc.go %d", len(tr.entries), live, docN)
	}
}
