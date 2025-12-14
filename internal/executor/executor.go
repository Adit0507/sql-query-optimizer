package executor

import (
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

func ( e*Executor) Execute(plan plan.LogicalPlan) ([]Row, error) {
	iter, err := e.executeNode(plan)
	if err != nil {
		return  nil, err
	}
	defer iter.Close()

	var results []Row
	for {
		
	}


	return  results, nil
}

func (e *Executor) executeNode(plan plan.LogicalPlan) (Iterator, error) {}