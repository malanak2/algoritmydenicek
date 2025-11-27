package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Entry struct {
	date     time.Time
	entry    string
	next     *Entry
	previous *Entry
}

func (e *Entry) Next() *Entry {
	return e.next
}
func (e *Entry) Previous() *Entry {
	return e.previous
}

type List struct {
	head *Entry
	tail *Entry
}

func (l *List) Len() int {
	count := 0
	if l.head == nil {
		return count
	}
	current := l.head
	for {
		count++
		if current.next == nil {
			break
		}
		current = current.next
	}
	return count
}
func (l *List) debugLog() {
	fmt.Println("--- DEBUG LOG START ---")
	fmt.Printf("Head is %v, tail is %v.\n", l.head, l.tail)
	if l.head == nil {
		fmt.Println("List is empty.")
		fmt.Println("--- DEBUG LOG END ---")
		return
	}

	current := l.head
	count := 0
	for current != nil {
		fmt.Printf("Index: %d, Date: %s, Entry: \"%s...\"\n", count, current.date.Format(time.DateOnly), strings.Split(current.entry, "\n")[0])
		count++
		current = current.next
	}

	fmt.Printf("Total entries: %d\n", count)
	fmt.Println("--- DEBUG LOG END ---")
}

func (l *List) addAfter(element *Entry, newElement *Entry) {
	if element.next == nil {
		element.next = newElement
		newElement.previous = element
		l.tail = newElement
		return
	}
	nextOrigin := element.next
	element.next = newElement
	newElement.previous = element
	nextOrigin.previous = newElement
	newElement.next = nextOrigin
}
func (l *List) addBefore(element *Entry, newElement *Entry) {
	if element.previous == nil {
		element.previous = newElement
		newElement.next = element
		l.head = newElement
		return
	}
	prevOrigin := element.previous
	element.previous = newElement
	newElement.next = element
	prevOrigin.next = newElement
	newElement.previous = prevOrigin
}
func (l *List) addFirst(newElement *Entry) {
	if l.head == nil {
		l.head = newElement
		l.tail = newElement
		return
	}
	l.head.previous = newElement
	newElement.next = l.head
	l.head = newElement
}
func (l *List) addLast(newElement Entry) {
	if l.tail == nil {
		l.head = &newElement
		l.tail = &newElement
		return
	}
	l.tail.next = &newElement
	newElement.previous = l.tail
	l.tail = &newElement
}

func (l *List) removeFirst() {
	if l.head == l.tail {
		l.head = nil
		l.tail = nil
		return
	}
	l.head = l.head.next
}
func (l *List) removeLast() {
	if l.head == l.tail {
		l.head = nil
		l.tail = nil
		return
	}
	l.tail = l.tail.previous
}

func (l *List) remove(element *Entry) {
	if l.head == l.tail {
		l.head = nil
		l.tail = nil
		return
	}
	if element.next == nil {
		l.tail = element.previous
		element.previous.next = nil
		return
	}
	if element.previous == nil {
		l.head = element.next

		element.next.previous = nil
		return
	}
	next := element.next
	element.previous.next = next
	next.previous = element.previous
}

func (l *List) First() *Entry {
	return l.head
}
func (l *List) Last() *Entry {
	return l.tail
}

var (
	l                *List
	currentOperation string
	msg              string
	index            int
)

func (l *List) atIndex(indx int) *Entry {
	if l.Len() <= indx {
		return nil
	}
	countr := 0
	listElement := l.First()
	for {
		if countr == indx {
			return listElement
		}
		countr++
		listElement = listElement.Next()
	}
}
func main() {
	l = &List{}
	msg = ""
	newEntry := Entry{}
	currentOperation = "command"
	index = 0
	for {
		fmt.Print("\033[H\033[2J")                             // clear term
		fmt.Printf("%s, index: %d\n", currentOperation, index) // debug
		if msg != "" {
			fmt.Print(msg + "\n")
			msg = ""
		}
		//l.debugLog()
		fmt.Printf(`
---------------------------------
Deník se ovládá následujícími příkazy:
- predchozi: Přesunutí na předchozí záznam
- dalsi: Přesunutí na další záznam
- zacatek: přenese mě na první záznam
- konec: přenese mě na poslední záznam
- novy: Vytvoření nového záznamu
- uloz: Uložení vytvořeného záznamu
- smaz: Odstranění záznamu
- zavri: Zavření deníku
---------------------------------

Počet záznamů: %d
`, l.Len())
		switch currentOperation {
		case "command":
			if l.Len() != 0 {
				entr := l.atIndex(index)
				if entr == nil {
					msg = "Chyba při načítání záznamu."
					currentOperation = "command"
					continue
				}
				fmt.Printf(`
Datum: %s

%s
---------------------------------
`, entr.date.String(), entr.entry)
			}

			fmt.Print("Zadej příkaz: ")
		case "newDate":
			fmt.Printf("Datum (formát %s): ", time.DateOnly)

		case "newEntryText":
			fmt.Printf(`
Datum: %s

%s
---------------------------------
`, newEntry.date.String(), newEntry.entry)
		}

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)

		switch currentOperation {
		case "newDate":
			date, err := time.Parse(time.DateOnly, line)
			if err != nil {
				msg = err.Error()
				currentOperation = "command"
				continue
			}
			newEntry.date = date
			currentOperation = "newEntryText"
		case "newEntryText":
			if line == "uloz" {
				l.addLast(newEntry)
				newEntry = Entry{}
				index = l.Len() - 1
				if index < 0 {
					index = 0
				}
				currentOperation = "command"
				continue
			}
			newEntry.entry += line + "\n"
			continue
		case "command":
			switch line {
			case "novy":
				currentOperation = "newDate"
			case "dalsi":
				if index >= l.Len()-1 {
					msg = "Jste na konci listu."
					continue
				}
				index++
			case "predchozi":
				if index <= 0 {
					msg = "Jste na začátku listu."
					continue
				}
				index--
				if index < 0 {
					index = 0
				}
			case "zacatek":
				index = 0
			case "konec":
				index = l.Len() - 1
				if index < 0 {
					index = 0
				}
			case "smaz":
				if l.Len() == 0 {
					msg = "List je prázdný."
				}
				entr := l.atIndex(index)
				l.remove(entr)
				if index >= l.Len() {
					index--
				}
				if index < 0 {
					index = 0
				}
			case "zavri":
				return
			default:
				msg = "This command is not valid, please try again."
				fmt.Printf("%s", line)
			}
		}
	}
}
