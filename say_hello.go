package main

import (
    "bufio"
    "fmt"
    "flag"
    "os"
    "strings"
)

func main() {
    // Process command-line arguments.
    namePtr := flag.String("to", "World", "Say Hello to this entity.")
    interactive := flag.Bool("i", false, "Say Hello to typed-in value.")
    flag.Parse()
    
    // Print "Hello " + specified value.
    switch (*interactive) {
        case true:
            reader := bufio.NewReader(os.Stdin)
            fmt.Print("Who am I saying 'Hello' to? ")
            text, _ := reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
            if (len(text) == 0) {
                text = "World"
            }
            fmt.Println("Hello " + text + "!")

        default:
            fmt.Println("Hello " + *namePtr + "!")
    }
}