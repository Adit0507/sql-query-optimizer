package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Adit0507/sql-query-optimizer/internal/catalog"
	"github.com/Adit0507/sql-query-optimizer/internal/executor"
	"github.com/Adit0507/sql-query-optimizer/internal/parser"
	"github.com/Adit0507/sql-query-optimizer/internal/plan"
)

func main() {
	fmt.Println("SQL Query optimizer")

	cat := catalog.NewCatalog()
	if err := cat.LoadFromFile("catalog.json"); err != nil {
		fmt.Printf("Error loading catalog: %v\n", err)
		return
	}
	fmt.Println("Catalog loaded")
	
	planner := plan.NewPlanner(cat)
	exec := executor.NewExecutor(cat)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\nEnter SQL queries (type 'exit' to quit, 'help' for commands):")
	
	for {
		fmt.Println("\ngoquery> ")
		if !scanner.Scan(){
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == ""{
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}
		if input == "help" {
			printHelp()
			continue
		}

		if strings.HasPrefix(input, "EXPLAIN"){
			query := strings.TrimPrefix(input, "EXPLAIN ")
			executeExplain(query, planner)
			
			continue
		}

		executeQuery(input, planner, exec)
	}
}

func executeQuery(query string, planner *plan.Planner, exec *executor.Executor) {
	p := parser.NewParser(query)
	stmt := p.Parse()

	if len(p.Errors()) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range p.Errors() {
			fmt.Printf("  - %s\n", err)
		}
		return
	}

	// creatin logical plan
	logicalPlan, err := planner.CreateLogicalPlan(stmt)
	if err != nil {
		fmt.Printf("Planning error: %v\n", err)
		return
	}

	results, err := exec.Execute(logicalPlan)
	if err != nil {
		fmt.Printf("Execution error: %v")
		return
	}

	displayResults(results)
}

func executeExplain(query string, planner *plan.Planner) {
	p := parser.NewParser(query)
	stmt := p.Parse()

	if len(p.Errors()) > 0 {
		fmt.Println("Parse errors:")

		for _, err := range p.Errors() {
			fmt.Printf("  - %s\n", err)
		}
		return
	}

	// logical plan
	logicalPlan, err := planner.CreateLogicalPlan(stmt)
	if err != nil {
		fmt.Printf("Planning error: %v\n", err)
		return
	}

	fmt.Println("\nLogical Plan:")
	fmt.Println("-------------")
	plan.PrintPlan(logicalPlan, 0)
}

func displayResults(results []executor.Row) {
	if len(results) == 0 {
		fmt.Println("(0 rows)")
		return
	}

	var columns []string
	for col := range results[0] {
		columns = append(columns, col)
	}

	fmt.Println()
	for i, col := range columns {
		fmt.Printf("%-20s", col)
		if i < len(columns)-1 {
			fmt.Print("| ")
		}
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", len(columns)*22))

	for _, row := range results {
		for i, col := range columns {
			val := row[col]
			fmt.Printf("%-20v", val)
			if i < len(columns)-1 {
				fmt.Print("| ")
			}
		}
	}

	fmt.Printf("\n(%d rows)\n", len(results))
}

func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  SELECT ...           - Execute a SELECT query")
	fmt.Println("  EXPLAIN SELECT ...   - Show query execution plan")
	fmt.Println("  help                 - Show this help message")
	fmt.Println("  exit/quit            - Exit the program")
	fmt.Println("\nExample queries:")
	fmt.Println("  SELECT * FROM users")
	fmt.Println("  SELECT name, email FROM users WHERE age > 25")
	fmt.Println("  SELECT users.name, orders.amount FROM users JOIN orders ON users.id = orders.user_id")
}
