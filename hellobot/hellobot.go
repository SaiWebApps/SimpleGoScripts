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
            // Process user input.
            reader := bufio.NewReader(os.Stdin)
            fmt.Print("Who am I saying 'Hello' to? ")
            text, _ := reader.ReadString('\n')

            // Replace trailing \n in Linux and \r\n in Windows.
            replacer := strings.NewReplacer("\r", "", "\n", "")
            text = replacer.Replace(text)
            
            // If an empty string was specified, then use "World" by default.
            if (len(text) == 0) {
                text = "World"
            }

            // Say "Hello" to the specified subject.
            fmt.Println("Hello " + text + "!")

        // Use the value specified with the "-to" flag.
        default:
            fmt.Println("Hello " + *namePtr + "!")
    }
}