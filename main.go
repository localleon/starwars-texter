package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/micmonay/keybd_event"
)

var kb keybd_event.KeyBonding

func main() {
	var err error
	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatal("Cant write to keyboard")
	}

	fmt.Println("Select text input field to write. Wait 5 sec to start writting...")
	time.Sleep(5 * time.Second)

	file, err := os.Open("./revenge-of-the-sith")
	//file, err := os.Open("./test")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writeTextFromFile(file)
}

func writeTextFromFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	// Delay before we start writing

	for scanner.Scan() {
		message := scanner.Text()
		if message != "" {
			err := KeyboardWrite(message)
			// Press Enter
			kb.SetKeys(keybd_event.VK_ENTER)
			kb.Launching()
			if err != nil {
				fmt.Printf("Could not write as keyboard input. Error: %s", err.Error())
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}

type keySet struct {
	code  int
	shift bool
}

//KeyboardWrite emulate keyboard input from string
func KeyboardWrite(textInput string) error {

	//Should we skip next character in string
	//Used if we found some escape sequence
	skip := false
	for i, c := range textInput {
		if !skip {
			if c != '\\' {
				kb.SetKeys(names[string(c)].code)
				kb.HasSHIFT(names[string(c)].shift)
			} else {
				//Found backslash escape character
				//Check next character
				switch textInput[i+1] {
				case 'n':
					//Found newline character sequence
					kb.SetKeys(names["ENTER"].code)
					skip = true
				case '\\':
					//Found backslash character sequence
					kb.SetKeys(names["\\"].code)
					kb.HasSHIFT(names["\\"].shift)
					skip = true
				case 'b':
					//Found backspace character sequence
					kb.SetKeys(names["BACKSPACE"].code)
					skip = true
				case 't':
					//Found tab character sequence
					kb.SetKeys(names["TAB"].code)
					skip = true
				case '"':
					//Found double quote character sequence
					kb.SetKeys(names["\""].code)
					kb.HasSHIFT(names["\""].shift)
					skip = true
				case '`':
					//Found single quote character sequence
					kb.SetKeys(keybd_event.VK_SP3)
					kb.HasSHIFT(true)
					skip = true
				default:
					//Nothing special, jsut backslash output
					kb.SetKeys(names["\\"].code)
					kb.HasSHIFT(names["\\"].shift)
				}

			}
			err := kb.Launching()
			if err != nil {
				return err
			}
		} else {
			skip = false
		}

	}
	return nil

}
