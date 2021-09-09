// Package reform is a better ORM for Go, based on non-empty interfaces and code generation.
//
// See README (https://github.com/go-reform/reform/blob/main/README.md) for quickstart information.
//
//
// Context
//
// Querier object, embedded into DB and TX types, contains context which is used by all its methods.
// It defaults to context.Background() and can be changed with WithContext method:
//
//  // for a single call
//  projects, err := DB.WithContext(ctx).SelectAllFrom(ProjectTable, "")
//
//  // for several calls
//  q := DB.WithContext(ctx)
//  projects, err := q.SelectAllFrom(ProjectTable, "")
//  persons, err := q.SelectAllFrom(PersonTable, "")
//
// Methods Exec, Query, and QueryRow use the same context.
// Methods ExecContext, QueryContext, and QueryRowContext are just compatibility wrappers for
// Querier.WithContext(ctx).Exec/Query/QuyeryRow to satisfy various standard interfaces.
//
// DB object methods Begin and InTransaction start transaction with the same context.
// Methods BeginTx and InTransactionContext start transaction with a given context without changing
// DB's context:
//
//  var projects, persons []Struct
//  err := DB.InTransactionContext(ctx, nil, func(tx *reform.TX) error {
//      var e error
//
//      // uses ctx
//      if projects, e = tx.SelectAllFrom(ProjectTable, ""); e != nil {
//          return e
//      }
//
//      // uses ctx too
//      if persons, e = tx.SelectAllFrom(PersonTable, ""); e != nil {
//          return e
//      }
//
//      return nil
//  }
//
// Note that several different contexts can be used:
//
//  DB.InTransactionContext(ctx1, nil, func(tx *reform.TX) error {
//      _, _ = tx.SelectAllFrom(PersonTable, "")                    // uses ctx1
//      _, _ = tx.WithContext(ctx2).SelectAllFrom(PersonTable, "")  // uses ctx2
//      ...
//  })
//
// In theory, ctx1 and ctx2 can be entirely unrelated. Although that construct is occasionally useful,
// the behavior on context cancelation is entirely driver-defined; some drivers may just close the whole
// connection, effectively canceling unrelated ctx2 on ctx1 cancelation. For that reason mixing several
// contexts is not recommended.
//
//
// Tagging
//
// reform allows one to add tags (comments) to generated queries with WithTag Querier method.
// They can be used to track queries from RDBMS logs and tools back to application code. For example, this code:
//  id := "baron"
//  project, err := DB.WithTag("GetProject:%v", id).FindByPrimaryKeyFrom(ProjectTable, id)
// will generate the following query:
//  SELECT /* GetProject:baron */ "projects"."name", "projects"."id", "projects"."start", "projects"."end" FROM "projects" WHERE "projects"."id" = ? LIMIT 1
// Please keep in mind that dynamic tags can affect RDBMS query cache. Consult your RDBMS documentation for details.
// Some known links:
//  MySQL / Percona Server: https://www.percona.com/doc/percona-server/5.7/performance/query_cache_enhance.html#ignoring-comments
//  Microsoft SQL Server: https://msdn.microsoft.com/en-us/library/cc293623.aspx
//
//
// Short example
//
// This example shows some reform features.
// It uses https://github.com/AlekSi/pointer to get pointers to values of build-in types.
package reform // import "gopkg.in/reform.v1"

// Version defines reform version.
const Version = "v1.5.0"
