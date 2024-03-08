package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
)

func getTodoFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := usr.HomeDir + "/todos.txt"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
	
	return path
}

var fileName = getTodoFilePath()

func WriteToFile(todo string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.WriteString("[❎] " + todo); err != nil {
		log.Fatal(err)
	}
}

func markDone(id int) {
	todos := readFromFile()
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i, todo := range todos {
		if i == id {
			todo = strings.Replace(todo, "[❎]", "[✅]", 1)
		}
		if _, err := file.WriteString(todo + "\n"); err != nil {
			log.Fatal(err)
		}
	}

}

func updateTodo(id int, new string) {
	todos := readFromFile()
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	new = strings.TrimSpace(new)

	for i, todo := range todos {
		if i == id {
			status := todo[:5]
			todo = status + " " + new
		}
		if _, err := file.WriteString(todo + "\n"); err != nil {
			log.Fatal(err)
		}
	}
}

func readFromFile() []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var todos []string
	for scanner.Scan() {
		todos = append(todos, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return todos
}

func deleteTodo(id int) {
	todos := readFromFile()
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i, todo := range todos {
		if i != id {
			if _, err := file.WriteString(todo + "\n"); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	fmt.Println(color.BlueString("\n" +
		"██████╗  ██████╗     ██████╗  ██████╗ \n" +
		"██╔════╝ ██╔═══██╗    ██╔══██╗██╔═══██╗\n" +
		"██║  ███╗██║   ██║    ██║  ██║██║   ██║\n" +
		"██║   ██║██║   ██║    ██║  ██║██║   ██║\n" +
		"╚██████╔╝╚██████╔╝    ██████╔╝╚██████╔╝\n" +
		" ╚═════╝  ╚═════╝     ╚═════╝  ╚═════╝\n"))

	br := bufio.NewReader(os.Stdin)

	for {
		todosTable := table.NewWriter()
		todosTable.SetOutputMirror(os.Stdout)
		todosTable.AppendHeader(table.Row{"ID", "Status", "Todo"})

		todos := readFromFile()
		for i, todo := range todos {
			status := "[❎]"
			if strings.Contains(todo, "[✅]") {
				status = "[✅]"
			}
			todosTable.AppendRow([]interface{}{i, status, todo[5:]})
		}
		todosTable.Render()
		fmt.Println(color.GreenString("\n----Go select one----"))
		fmt.Println("1. Create todo")
		fmt.Println("2. Mark as done")
		fmt.Println("3. Update todo")
		fmt.Println("4. Delete todo")
		fmt.Println("5. Celebrate")
		fmt.Println("6. Exit")

		var choice int
		fmt.Printf("%v", color.GreenString("\nYour choice: "))
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Print("What do you want to do: ")
			var todo string
			todo, _ = br.ReadString('\n')
			WriteToFile(todo)
		case 2:
			var id int
			fmt.Print("ID of todo to mark as done: ")
			fmt.Scan(&id)
			markDone(id)
		case 3:
			var id int
			fmt.Print("ID of todo to update: ")
			fmt.Scan(&id)
			fmt.Print("Updated text: ")
			new, _ := br.ReadString('\n')
			updateTodo(id, new)
		case 4:
			var id int
			fmt.Print("ID of todo to delete: ")
			fmt.Scan(&id)
			deleteTodo(id)
		case 5:
			cmd := exec.Command("curl", "-s", "parrot.live")
			cmd.Stdout = os.Stdout
			go func() {
				timer := time.NewTimer(5 * time.Second)
				<-timer.C
				cmd.Process.Kill()
			}()
			cmd.Run()
		case 6:
			fmt.Println("Goodbye 👋")
			return
		default:
			fmt.Println("Invalid choice 🤬")
		}
	}
}
