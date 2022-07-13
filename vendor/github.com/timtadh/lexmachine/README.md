# `lexmachine` - Lexical Analysis Framework for Golang

By Tim Henderson

Copyright 2014-2017, All Rights Reserved. Made available for public use under
the terms of a BSD 3-Clause license.

[![GoDoc](https://godoc.org/github.com/timtadh/lexmachine?status.svg)](https://godoc.org/github.com/timtadh/lexmachine)
[![ReportCard](https://goreportcard.com/badge/github.com/timtadh/lexmachine)](https://goreportcard.com/report/github.com/timtadh/lexmachine)

## What?

`lexmachine` is a full lexical analysis framework for the Go programming
language. It supports a restricted but usable set of regular expressions
appropriate for writing lexers for complex programming languages. The framework
also supports sub lexers and non-regular lexing through an "escape hatch" which
allows the users to consume any number of further bytes after a match. So if you
want to support nested C-style comments or other paired structures you can do so
at the lexical analysis stage.

Subscribe to the [mailing
list](https://groups.google.com/forum/#!forum/lexmachine-users) to get
announcement of major changes, new versions, and important patches.

## Goal

`lexmachine` intends to be the best, fastest, and easiest to use lexical
analysis system for Go.

1. [Documentation Links](#documentation)
1. [Narrative Documentation](#narrative-documentation)
1. [Regular Expressions in `lexmachine`](#regular-expressions)
1. [History](#history)
1. [Complete Example](#complete-example)

## Documentation

-   [Tutorial](http://hackthology.com/writing-a-lexer-in-go-with-lexmachine.html)
-   [How It Works](http://hackthology.com/faster-tokenization-with-a-dfa-backend-for-lexmachine.html)
-   [Narrative Documentation](#narrative-documentation)
-   [![GoDoc](https://godoc.org/github.com/timtadh/lexmachine?status.svg)](https://godoc.org/github.com/timtadh/lexmachine)

### What is in Box

`lexmachine` includes the following components

1.  A parser for restricted set of regular expressions.
2.  A abstract syntax tree (AST) for regular expressions.
3.  A backpatching code generator which compiles the AST to (NFA) machine code.
4.  Both DFA (Deterministic Finite Automata) and a NFA (Non-deterministic Finite
    Automata) simulation based lexical analysis engines. Lexical analysis
    engines work in a slightly different way from a normal regular expression
    engine as they tokenize a stream rather than matching one string.
5.  Match objects which include start and end column and line numbers of the
    lexemes as well as their associate token name.
6.  A declarative "DSL" for specifying the lexers.
7.  An "escape hatch" which allows one to match non-regular tokens by consuming
    any number of further bytes after the match.

## Narrative Documentation

`lexmachine` splits strings into substrings and categorizes each substring. In
compiler design, the substrings are referred to as *lexemes* and the
categories are referred to as *token types* or just *tokens*. The categories are
defined by *patterns* which are specified using [regular
expressions](#regular-expressions). The process of splitting up a string is
sometimes called *tokenization*, *lexical analysis*, or *lexing*.

### Defining a Lexer

The set of patterns (regular expressions) used to *tokenize* (split up and
categorize) is called a *lexer*. Lexer's are first class objects in
`lexmachine`. They can be defined once and re-used over and over-again to
tokenize multiple strings. After the lexer has been defined it will be compiled
(either explicitly or implicitly) into either a Non-deterministic Finite
Automaton (NFA) or Deterministic Finite Automaton (DFA). The automaton is then
used (and re-used) to tokenize strings.

#### Creating a new Lexer

```go
lexer := lexmachine.NewLexer()
```

#### Adding a pattern

Let's pretend we want a lexer which only recognizes one category: strings which
match the word "wild" capitalized or not (eg. Wild, wild, WILD, ...). That
expression is denoted: `[Ww][Ii][Ll][Dd]`. Patterns are added using the `Add`
function:

```go
lexer.Add([]byte(`[Ww][Ii][Ll][Dd]`), func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	return 0, nil
})
```

Add takes two arguments: the pattern and a call back function called a *lexing
action*. The action allows you, the programmer, to transform the low level
`machines.Match` object (from `github.com/lexmachine/machines`) into a object
meaningful for your program. As an example, let's define a few token types, and
a token object. Then we will construct appropriate action functions.

```go
Tokens := []string{
	"WILD",
	"SPACE",
	"BANG",
}
TokenIds := make(map[string]int)
for i, tok := range Tokens {
	TokenIds[tok] = i
}
```

Now that we have defined a set of three tokens (WILD, SPACE, BANG), let's create
a token object:

```go
type Token struct {
	TokenType int
	Lexeme string
	Match *machines.Match
}
```

Now let's make a helper function which takes a `Match` and a token type and
creates a Token.

```go
func NewToken(tokenType string, m *machines.Match) *Token {
	return &Token{
		TokenType: TokenIds[tokenType], // defined above
		Lexeme: string(m.Bytes),
		Match: m,
	}
}
```

Now we write an action for the previous pattern

```go
lexer.Add([]byte(`[Ww][Ii][Ll][Dd]`), func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	return NewToken("WILD", m), nil
})
```

Writing the action functions can get tedious, a good idea is to create a helper
function which produces these action functions:

```go
func token(tokenType string) func(*lexmachine.Scanner, *machines.Match) (interface{}, error) {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return NewToken(tokenType, m), nil
	}
}
```

Then adding patterns for our 3 tokens is concise:

```go
lexer.Add([]byte(`[Ww][Ii][Ll][Dd]`), token("WILD"))
lexer.Add([]byte(` `), token("SPACE"))
lexer.Add([]byte(`!`), token("BANG"))
```

#### Built-in Token Type

Many programs use similar representations for tokens. `lexmachine` provides a
completely optional `Token` object you can use in lieu of writing your own.

```go
type Token struct {
    Type        int
    Value       interface{}
    Lexeme      []byte
    TC          int
    StartLine   int
    StartColumn int
    EndLine     int
    EndColumn   int
}
```

Here is an example for constructing a lexer Action which turns a machines.Match
struct into a token using the scanners Token helper function.

```go
func token(name string, tokenIds map[string]int) lex.Action {
    return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
        return s.Token(tokenIds[name], string(m.Bytes), m), nil
    }
}
```

#### Adding Multiple Patterns

When constructing a lexer for a complex computer language often tokens have
patterns which overlap -- multiple patterns could match the same strings. To
address this problem lexical analysis engines follow 2 rules when choosing which
pattern to match:

1. Pick the pattern which matches the longest prefix of unmatched text.
2. Break ties by picking the pattern which appears earlier in the user supplied
   list.

For example, let's pretend we are writing a lexer for Python. Python has a bunch
of keywords in it such as `class` and `def`. However, it also has identifiers
which match the pattern `[A-Za-z_][A-Za-z0-9_]*`. That pattern also matches the
keywords, if we were to define the lexer as:

```go
lexer.Add([]byte(`[A-Za-z_][A-Za-z0-9_]*`), token("ID"))
lexer.Add([]byte(`class`), token("CLASS"))
lexer.Add([]byte(`def`), token("DEF"))
```

Then, the keywords class and def would never be found because the ID token would
take precedence. The correct way to solve this problem is by putting the
keywords first:

```go
lexer.Add([]byte(`class`), token("CLASS"))
lexer.Add([]byte(`def`), token("DEF"))
lexer.Add([]byte(`[A-Za-z_][A-Za-z0-9_]*`), token("ID"))
```

#### Skipping Patterns

Sometimes it is advantageous to not emit tokens for certain patterns and to
instead skip them. Commonly this occurs for whitespace and comments. To skip a
pattern simply have the action `return nil, nil`:

```go
lexer.Add(
	[]byte("( |\t|\n)"),
	func(scan *Scanner, match *machines.Match) (interface{}, error) {
		// skip white space
		return nil, nil
	},
)
lexer.Add(
	[]byte("//[^\n]*\n"),
	func(scan *Scanner, match *machines.Match) (interface{}, error) {
		// skip white space
		return nil, nil
	},
)
```

#### Compiling the Lexer

`lexmachine` uses the theory of [finite state
machines](http://hackthology.com/faster-tokenization-with-a-dfa-backend-for-lexmachine.html)
to efficiently tokenize text. So what is a finite state machine? A finite state
machine is a mathematical construct which is made up of a set of states, with a
labeled starting state, and accepting states. There is a transition function
which moves from one state to another state based on an input character. In
general, in lexing there are two usual types of state machines used:
Non-deterministic and Deterministic.

Before a lexer (like the ones described above) and be used it must be compiled
into either a Non-deterministic Finite Automaton (NFA) or a [Deterministic
Finite Automaton
(DFA)](http://hackthology.com/faster-tokenization-with-a-dfa-backend-for-lexmachine.html).
The difference between the two (from a practical perspective) is *construction
time* and *match efficiency*.

Construction time is the amount of time it takes to turn a set of regular
expressions into a state machine (also called a finite state automaton). For an
NFA it is O(`r`) which `r` is the length of the regular expression. However, for
DFA it could be as bad as O(`2^r`) but in practical terms it is rarely worse
than O(`r^3`). The DFA's in `lexmachine` are also automatically *minimized* which
reduces the amount of memory they consume which takes O(`r*log(log(r))`) steps.

However, construction time is an upfront cost. If your program is tokenizing
multiple strings it is less important than match efficiency. Let's say a string
has length `n`. An NFA can tokenize such a string in O(`n*r`) steps while a DFA
can tokenize the string in O(`n`). For larger languages `r` becomes a
significant overhead.

By default, `lexmachine` uses a DFA. To explicitly invoke compilation call
`Compile`:

```go
err := lexer.Compile()
if err != nil {
	// handle err
}
```

To explicitly compile a DFA (in case of changes to the default behavior of
Compile):

```go
err := lexer.CompileDFA()
if err != nil {
	// handle err
}
```

To explicitly compile a NFA:

```go
err := lexer.CompileNFA()
if err != nil {
	// handle err
}
```

### Tokenizing a String

To tokenize (lex) a string construct a `Scanner` object using the lexer. This
will compile the lexer if it has not already been compiled.

```go
scanner, err := lexer.Scanner([]byte("some text to lex"))
if err != nil {
	// handle err
}
```

The scanner object is an iterator which yields the next token (or error) by
calling the `Next()` method:

```go
for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
	if ui, is := err.(*machines.UnconsumedInput); is {
		// skip the error via:
		// scanner.TC = ui.FailTC
		//
		return err
	} else if err != nil {
		return err
	}
	fmt.Println(tok)
}
```

Let's break down that first line:

```go
for tok, err, eos := scanner.Next();
```

The `Next()` method returns three things, the token (`tok`) if there is one, an
error (`err`) if there is one, and `eos` which is a boolean which indicates if
the End Of String (EOS) has been reached.

```go
; !eos;
```

Iteration proceeds until the EOS has been reached.

```go
; tok, err, eos = scanner.Next() {
```

The update block calls `Next()` again to get the next token. In each iteration
of the loop the first thing a client **must** do is check for an error.

```go
	if err != nil {
		return err
	}
```

This prevents an infinite loop on an unexpected character or other bad token. To
skip bad tokens check to see if the `err` is a `*machines.UnconsumedInput`
object and reset the scanners text counter (`scanner.TC`) to point to the end of
the failed token.

```go
	if ui, is := err.(*machines.UnconsumedInput); is {
		scanner.TC = ui.FailTC
		continue
	}
```

Finally, a client can make use of the token produced by the scanner (if there
was no error:

```go
	fmt.Println(tok)
```

### Dealing with Non-regular Tokens

`lexmachine` like most lexical analysis frameworks primarily deals with patterns
which are represented by regular expressions. However, sometimes a language
has a token which is "non-regular." A pattern is non-regular if there is no
regular expression (or finite automata) which can express the pattern. For
instance, if you wanted to define a pattern which matches only consecutive
balanced parentheses: `()`, `()()()`, `((()()))()()`, ... You would quickly find
there is no regular expression which can express this language. The reason is
simple: finite automata cannot "count" or keep track of how many opening
parentheses it has seen.

This problem arises in many programming languages when dealing with nested
"c-style" comments. Supporting the nesting means solving the "balanced
parenthesis" problem. Luckily, `lexmachine` provides an "escape-hatch" to deal
with these situations in the `Action` functions. All actions receive a pointer
to the `Scanner`. The scanner (as discussed above) has a public modifiable field
called `TC` which stands for text counter. Any action can *modify* the text
counter to point at the desired position it would like the scanner to resume
scanning from.

An example of using this feature for tokenizing nested "c-style" comments is
below:

```go
lexer.Add(
	[]byte("/\\*"),
	func(scan *Scanner, match *machines.Match) (interface{}, error) {
		for tc := scan.TC; tc < len(scan.Text); tc++ {
			if scan.Text[tc] == '\\' {
				// the next character is skipped
				tc++
			} else if scan.Text[tc] == '*' && tc+1 < len(scan.Text) {
				if scan.Text[tc+1] == '/' {
					// set the text counter to point to after the
					// end of the comment. This will cause the
					// scanner to resume after the comment instead
					// of picking up in the middle.
					scan.TC = tc + 2
					// don't return a token to skip the comment
					return nil, nil
				}
			}
		}
		return nil,
			fmt.Errorf("unclosed comment starting at %d, (%d, %d)",
				match.TC, match.StartLine, match.StartColumn)
	},
)
```

## Regular Expressions

Lexmachine (like most lexical analysis frameworks) uses [Regular
Expressions](https://en.wikipedia.org/wiki/Regular_expression) to specify the
*patterns* to match when splitting the string up into categorized *tokens.*
For a more advanced introduction to regular expressions engines see Russ Cox's
[articles](https://swtch.com/~rsc/regexp/). To learn more about how regular
expressions are used to *tokenize* string take a look at Alex Aiken's [video
lectures](https://youtu.be/SRhkfvqeA1M) on the subject. Finally, Aho *et al.*
give a through treatment of the subject in the [Dragon
Book](http://www.worldcat.org/oclc/951336275) Chapter 3.

A regular expression is a *pattern* which *matches* a set of strings. It is made
up of *characters* such as `a` or `b`, characters with special meanings (such as
`.` which matches any character), and operators. The regular expression `abc`
matches exactly one string `abc`.

### Character Expressions

In lexmachine most characters (eg. `a`, `b` or `#`) represent themselves. Some
have special meanings (as detailed below in operators). However, all characters
can be represented by prefixing the character with a `\`.

#### Any Character

`.` matches any character.

#### Special Characters

1. `\` use `\\` to match
2. newline use `\n` to match
3. carriage return use `\r` to match
4. tab use `\t` to match
5. `.` use `\.` to match
6. operators: {`|`, `+`, `*`, `?`, `(`, `)`, `[`, `]`, `^`} prefix with a `\` to
   match.

#### Character Classes

Sometimes it is advantages to match a variety of characters. For instance, if
you want to ignore capitalization for the work `Capitol` you could write the
expression `[Cc]apitol` which would match both `Capitol` or `capitol`. There are
two forms of character ranges:

1. `[abcd]` matches all the letters inside the `[]` (eg. that pattern matches
   the strings `a`, `b`, `c`, `d`).
2. `[a-d]` matches the range of characters between the character before the dash
   (`a`) and the character after the dash (`d`) (eg. that pattern matches
   the strings `a`, `b`, `c`, `d`).

These two forms may be combined:

For instance, `[a-zA-Z123]` matches the strings {`a`, `b`, ..., `z`, `A`, `B`,
... `Z`, `0`, `2`, `3`}

#### Inverted Character Classes

Sometimes it is easier to specify the characters you don't want to match than
the characters you do. For instance, you want to match any character but a lower
case one. This can be achieved using an inverted class: `[^a-z]`. An inverted
class is specified by putting a `^` just after the opening bracket.

#### Built-in Character Classes

1. `\d` = `[0-9]` (the digit class)
2. `\D` = `[^0-9]` (the not a digit class)
3. `\s` = `[ \t\n\r\f]` (the space class). where \f is a form feed (note: \f is
   not a special sequence in lexmachine, if you want to specify the form feed
   character (ascii 0x0c) use []byte{12}.
4. `\S` = `[^ \t\n\r\f]` (the not a space class)
5. `\w` = `[0-9a-zA-Z_]` (the letter class)
5. `\W` = `[^0-9a-zA-Z_]` (the not a letter class)

### Operators

1. The pipe operator `|` indicates alternative choices. For instance the
   expression `a|b` matches either the string `a` or the string `b` but not `ab`
   or `ba` or the empty string.

2. The parenthesis operator `()` groups a subexpression together. For instance
   the expression `a(b|c)d` matches `abd` or `acd` but not `abcd`.

3. The star operator `*` indicates the "starred" subexpression should match zero
   or more times. For instance, `a*` matches the empty string, `a`, `aa`, `aaa`
   and so on.

4. The plus operator `+` indicates the "plussed" subexpression should match one
   or more times. For instance, `a+` matches `a`, `aa`, `aaa` and so on.

5. The maybe operator `?` indicates the "questioned" subexpression should match
   zero or one times. For instance, `a?` matches the empty string and `a`.

### Grammar

The canonical grammar is found in the handwritten recursive descent
[parser](https://github.com/timtadh/lexmachine/blob/master/frontend/parser.go).
This section should be considered documentation not specification.

Note: e stands for the empty string

```
Regex -> Alternation

Alternation -> AtomicOps Alternation'

Alternation' -> `|` AtomicOps Alternation'
              | e

AtomicOps -> AtomicOp AtomicOps
           | e

AtomicOp -> Atomic
          | Atomic Ops

Ops -> Op Ops
     | e

Op -> `+`
    | `*`
    | `?`

Atomic -> Char
        | Group

Group -> `(` Alternation `)`

Char -> CHAR
      | CharClass

CharClass -> `[` Range `]`
           | `[` `^` Range `]`

Range -> CharClassItem Range'

Range' -> CharClassItem Range'
        | e

CharClassItem -> BYTE
              -> BYTE `-` BYTE

CHAR -> matches any character expect '|', '+', '*', '?', '(', ')', '[', ']', '^'
        unless escaped. Additionally '.' is returned as the wildcard character
        which matches any character. Built-in character classes are also handled
        here.

BYTE -> matches any byte
```

## History

This library was started when I was teaching EECS 337 *Compiler Design and
Implementation* at Case Western Reserve University in Fall of 2014. It wrote two
compilers one was "hidden" from the students as the language implemented was
their project language. The other was [tcel](https://github.com/timtadh/tcel)
which was written initially as an example of how to do type checking. That
compiler was later expanded to explain AST interpretation, intermediate code
generation, and x86 code generation.

## Complete Example

### Using the Lexer

```go
package main

import (
    "fmt"
    "log"
)

import (
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func main() {
    s, err := Lexer.Scanner([]byte(`digraph {
  rankdir=LR;
  a [label="a" shape=box];
  c [<label>=<<u>C</u>>];
  b [label="bb"];
  a -> c;
  c -> b;
  d -> c;
  b -> a;
  b -> e;
  e -> f;
}`))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Type    | Lexeme     | Position")
    fmt.Println("--------+------------+------------")
    for tok, err, eof := s.Next(); !eof; tok, err, eof = s.Next() {
        if ui, is := err.(*machines.UnconsumedInput); is{
            // to skip bad token do:
            // s.TC = ui.FailTC
            log.Fatal(err) // however, we will just fail the program
        } else if err != nil {
            log.Fatal(err)
        }
        token := tok.(*lexmachine.Token)
        fmt.Printf("%-7v | %-10v | %v:%v-%v:%v\n",
            Tokens[token.Type],
            string(token.Lexeme),
            token.StartLine,
            token.StartColumn,
            token.EndLine,
            token.EndColumn)
    }
}
```

### Lexer Definition

```go
package main

import (
	"fmt"
	"strings"
)

import (
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

var Literals []string       // The tokens representing literal strings
var Keywords []string       // The keyword tokens
var Tokens []string         // All of the tokens (including literals and keywords)
var TokenIds map[string]int // A map from the token names to their int ids
var Lexer *lexmachine.Lexer // The lexer object. Use this to construct a Scanner

// Called at package initialization. Creates the lexer and populates token lists.
func init() {
	initTokens()
	var err error
	Lexer, err = initLexer()
	if err != nil {
		panic(err)
	}
}

func initTokens() {
	Literals = []string{
		"[",
		"]",
		"{",
		"}",
		"=",
		",",
		";",
		":",
		"->",
		"--",
	}
	Keywords = []string{
		"NODE",
		"EDGE",
		"GRAPH",
		"DIGRAPH",
		"SUBGRAPH",
		"STRICT",
	}
	Tokens = []string{
		"COMMENT",
		"ID",
	}
	Tokens = append(Tokens, Keywords...)
	Tokens = append(Tokens, Literals...)
	TokenIds = make(map[string]int)
	for i, tok := range Tokens {
		TokenIds[tok] = i
	}
}

// Creates the lexer object and compiles the NFA.
func initLexer() (*lexmachine.Lexer, error) {
	lexer := lexmachine.NewLexer()

	for _, lit := range Literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		lexer.Add([]byte(r), token(lit))
	}
	for _, name := range Keywords {
		lexer.Add([]byte(strings.ToLower(name)), token(name))
	}

	lexer.Add([]byte(`//[^\n]*\n?`), token("COMMENT"))
	lexer.Add([]byte(`/\*([^*]|\r|\n|(\*+([^*/]|\r|\n)))*\*+/`), token("COMMENT"))
	lexer.Add([]byte(`([a-z]|[A-Z]|[0-9]|_)+`), token("ID"))
	lexer.Add([]byte(`[0-9]*\.[0-9]+`), token("ID"))
	lexer.Add([]byte(`"([^\\"]|(\\.))*"`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			x, _ := token("ID")(scan, match)
			t := x.(*lexmachine.Token)
			v := t.Value.(string)
			t.Value = v[1 : len(v)-1]
			return t, nil
		})
	lexer.Add([]byte("( |\t|\n|\r)+"), skip)
	lexer.Add([]byte(`\<`),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			str := make([]byte, 0, 10)
			str = append(str, match.Bytes...)
			brackets := 1
			match.EndLine = match.StartLine
			match.EndColumn = match.StartColumn
			for tc := scan.TC; tc < len(scan.Text); tc++ {
				str = append(str, scan.Text[tc])
				match.EndColumn += 1
				if scan.Text[tc] == '\n' {
					match.EndLine += 1
				}
				if scan.Text[tc] == '<' {
					brackets += 1
				} else if scan.Text[tc] == '>' {
					brackets -= 1
				}
				if brackets == 0 {
					match.TC = scan.TC
					scan.TC = tc + 1
					match.Bytes = str
					x, _ := token("ID")(scan, match)
					t := x.(*lexmachine.Token)
					v := t.Value.(string)
					t.Value = v[1 : len(v)-1]
					return t, nil
				}
			}
			return nil,
				fmt.Errorf("unclosed HTML literal starting at %d, (%d, %d)",
					match.TC, match.StartLine, match.StartColumn)
		},
	)

	err := lexer.Compile()
	if err != nil {
		return nil, err
	}
	return lexer, nil
}

// a lexmachine.Action function which skips the match.
func skip(*lexmachine.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

// a lexmachine.Action function with constructs a Token of the given token type by
// the token type's name.
func token(name string) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIds[name], string(m.Bytes), m), nil
	}
}
```
