package executor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Adit0507/sql-query-optimizer/internal/catalog"
	"github.com/Adit0507/sql-query-optimizer/internal/plan"
)

type Row map[string]interface{} //row of data

type Iterator interface { //result set iterator
	Next() (Row, bool)
	Close()
}

type Executor struct {
	catalog *catalog.Catalog
}

func NewExecutor(cat *catalog.Catalog) *Executor {
	return &Executor{
		catalog: cat,
	}
}

func (e *Executor) Execute(plan plan.LogicalPlan) ([]Row, error) {
	iter, err := e.executeNode(plan)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var results []Row
	for {
		row, ok := iter.Next()
		if !ok {
			break
		}
		results = append(results, row)
	}

	return results, nil
}

func (e *Executor) executeNode(node plan.LogicalPlan) (Iterator, error) {
	switch n := node.(type) {
	case *plan.LogicalScan:
		return e.executeScan(n)

	default:
		return nil, fmt.Errorf("unsupported plan node: %T", node)
	}
}

type scanIterator struct {
	rows  []Row
	index int
}

func (s *scanIterator) Next() (Row, bool) {
	if s.index >= len(s.rows) {
		return nil, false
	}
	row := s.rows[s.index]
	s.index++

	return row, true
}
func (s *scanIterator) Close() {}

func (e *Executor) executeScan(scan *plan.LogicalScan) (Iterator, error) {
	data, err := os.ReadFile(scan.Table.DataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var rows []Row
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	return &scanIterator{rows: rows, index: 0}, nil
}

type filterIterator struct {
	input     Iterator
	predicate plan.Expr
}

func (f *filterIterator) Next() (Row, bool) {
	for {
		row, ok := f.input.Next()
		if !ok {
			return nil, false
		}

		result, err := evaluateExpr(f.predicate, row)
		if err != nil {
			continue
		}

		if boolResult, ok := result.(bool); ok && boolResult {
			return row, true
		}
	}
}
func (f *filterIterator) Close() {
	f.input.Close()
}

func (e *Executor) executeFilter(filter *plan.LogicalFilter) (Iterator, error) {
	input, err := e.executeNode(filter.Input)
	if err != nil {
		return nil, err
	}

	return &filterIterator{
		input:     input,
		predicate: filter.Predicate,
	}, nil
}

type projectIterator struct {
	input       Iterator
	projections []plan.Expr
	columnNames []string
}

func (p *projectIterator) Next() (Row, bool) {
	row, ok := p.input.Next()
	if !ok {
		return nil, false
	}

	ans := make(Row)
	for i, expr := range p.projections {
		value, err := evaluateExpr(expr, row)
		if err != nil {
			continue
		}

		ans[p.columnNames[i]] = value
	}

	return ans, true
}

func (p *projectIterator) Close() {
	p.input.Close()
}

func (e *Executor) executeProject(proj *plan.LogicalProject) (Iterator, error) {
	input, err := e.executeNode(proj.Input)
	if err != nil {
		return nil, err
	}

	return &projectIterator{
		input:       input,
		projections: proj.Projections,
		columnNames: proj.ColumnNames,
	}, nil
}

// join iterator
type joinIterator struct {
	left      Iterator
	right     Iterator
	condition plan.Expr
	joinType  plan.JoinType
	leftRow   Row
	rightRows []Row
	rightIdx  int
}

func (j *joinIterator) Next() (Row, bool) {
	for {
		if j.leftRow == nil {
			row, ok := j.left.Next()
			if !ok {
				return nil, false
			}

			j.leftRow = row
			j.rightIdx = 0
		}

		if j.rightIdx >= len(j.rightRows) {
			j.leftRow = nil
			continue
		}
		rightRow := j.rightRows[j.rightIdx]
		j.rightIdx++

		// combining rows
		combined := make(Row)
		for k, v := range j.leftRow {
			combined[k] = v
		}
		for k, v := range rightRow {
			combined[k] = v
		}

		// evaluate comdition
		ans, err := evaluateExpr(j.condition, combined)
		if err != nil || ans != true {
			continue
		}

		return combined, true
	}
}
func (j *joinIterator) Close() {
	j.left.Close()
	j.right.Close()
}
