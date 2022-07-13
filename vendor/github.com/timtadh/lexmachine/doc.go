// Package lexmachine is a full lexical analysis framework for the Go
// programming language. It supports a restricted but usable set of regular
// expressions appropriate for writing lexers for complex programming
// languages. The framework also supports sub-lexers and non-regular lexing
// through an "escape hatch" which allows the users to consume any number of
// further bytes after a match. So if you want to support nested C-style
// comments or other paired structures you can do so at the lexical analysis
// stage.
//
// For a tutorial see
// http://hackthology.com/writing-a-lexer-in-go-with-lexmachine.html
//
// Example of defining a lexer
//
//     // CreateLexer defines a lexer for the graphviz dot language.
//     func CreateLexer() (*lexmachine.Lexer, error) {
//         lexer := lexmachine.NewLexer()
//
//         for _, lit := range Literals {
//             r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
//             lexer.Add([]byte(r), token(lit))
//         }
//         for _, name := range Keywords {
//             lexer.Add([]byte(strings.ToLower(name)), token(name))
//         }
//
//         lexer.Add([]byte(`//[^\n]*\n?`), token("COMMENT"))
//         lexer.Add([]byte(`/\*([^*]|\r|\n|(\*+([^*/]|\r|\n)))*\*+/`), token("COMMENT"))
//         lexer.Add([]byte(`([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*`), token("ID"))
//         lexer.Add([]byte(`"([^\\"]|(\\.))*"`), token("ID"))
//         lexer.Add([]byte("( |\t|\n|\r)+"), skip)
//         lexer.Add([]byte(`\<`),
//             func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
//                 str := make([]byte, 0, 10)
//                 str = append(str, match.Bytes...)
//                 brackets := 1
//                 match.EndLine = match.StartLine
//                 match.EndColumn = match.StartColumn
//                 for tc := scan.TC; tc < len(scan.Text); tc++ {
//                     str = append(str, scan.Text[tc])
//                     match.EndColumn += 1
//                     if scan.Text[tc] == '\n' {
//                         match.EndLine += 1
//                     }
//                     if scan.Text[tc] == '<' {
//                         brackets += 1
//                     } else if scan.Text[tc] == '>' {
//                         brackets -= 1
//                     }
//                     if brackets == 0 {
//                         match.TC = scan.TC
//                         scan.TC = tc + 1
//                         match.Bytes = str
//                         return token("ID")(scan, match)
//                     }
//                 }
//                 return nil,
//                     fmt.Errorf("unclosed HTML literal starting at %d, (%d, %d)",
//                         match.TC, match.StartLine, match.StartColumn)
//             },
//         )
//
//         err := lexer.Compile()
//         if err != nil {
//             return nil, err
//         }
//         return lexer, nil
//     }
//
//     func token(name string) lex.Action {
//         return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
//             return s.Token(TokenIds[name], string(m.Bytes), m), nil
//         }
//     }
//
// Example of using a lexer
//
//     func ExampleLex() error {
//         lexer, err := CreateLexer()
//         if err != nil {
//             return err
//         }
//         scanner, err := lexer.Scanner([]byte(`digraph {
//           rankdir=LR;
//           a [label="a" shape=box];
//           c [<label>=<<u>C</u>>];
//           b [label="bb"];
//           a -> c;
//           c -> b;
//           d -> c;
//           b -> a;
//           b -> e;
//           e -> f;
//         }`))
//         if err != nil {
//             return err
//         }
//         fmt.Println("Type    | Lexeme     | Position")
//         fmt.Println("--------+------------+------------")
//         for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
//             if err != nil {
//                 return err
//             }
//             token := tok.(*lexmachine.Token)
//             fmt.Printf("%-7v | %-10v | %v:%v-%v:%v\n",
//                 dot.Tokens[token.Type],
//                 string(token.Lexeme),
//                 token.StartLine,
//                 token.StartColumn,
//                 token.EndLine,
//                 token.EndColumn)
//         }
//         return nil
//     }
//
package lexmachine
