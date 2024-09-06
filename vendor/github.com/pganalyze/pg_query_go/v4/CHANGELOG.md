# Changelog

## Unreleased

* ...


## 4.2.1     2023-05-25

* Update to libpg_query 15-4.2.1 ([#85](https://github.com/pganalyze/pg_query_go/pull/85))
  - Deparser: Handle INTERVAL correctly when used in SET statements [#184](https://github.com/pganalyze/libpg_query/pull/184)
  - Deparser: Ensure index names are quoted as identifiers [#182](https://github.com/pganalyze/libpg_query/pull/182)
* Suppress -Wdeprecated-non-prototype warnings ([#83](https://github.com/pganalyze/pg_query_go/pull/83) and [#84](https://github.com/pganalyze/pg_query_go/pull/84))
* Return full structured error info instead of just error message ([#76](https://github.com/pganalyze/pg_query_go/pull/76))


## 4.2.0     2023-02-08

* Update to libpg_query 15-4.2.0
  - Update to PostgreSQL 15.1


## 2.2.0     2022-11-02

* Update to libpg_query 13-2.2.0 ([#64](https://github.com/pganalyze/pg_query_go/pull/64))
  - Fingerprinting version 3.1
    - Fixes issue with "SELECT DISTINCT" having the same fingerprint as "SELECT"
      (fingerprints for "SELECT DISTINCT" will change with this revision)
    - Group additional DDL statements together that otherwise generate a lot of
      unique fingerprints (ListenStmt, UnlistenStmt, NotifyStmt, CreateFunctionStmt,
      FunctionParameter and DoStmt)
  - Update to Postgres 13.8 patch release
  - Backport Xcode 14.1 build fix from upcoming 13.9 release
  - Normalize additional DDL statements


## 2.1.2      2022-06-28

* Update libpg_query to 13-2.1.2 ([#58](https://github.com/pganalyze/pg_query_go/pull/58))
  - Add support for analyzing PL/pgSQL code inside DO blocks [#142](https://github.com/pganalyze/libpg_query/pull/142)
  - Fix memory leak in pg_query_fingerprint error handling [#141](https://github.com/pganalyze/libpg_query/pull/141)
  - PL/pgSQL parser:
    - Add support for Assert [#135](https://github.com/pganalyze/libpg_query/pull/135)
    - Add support for SET, COMMIT, ROLLBACK and CALL [#130](https://github.com/pganalyze/libpg_query/pull/130)
  - Add support for parsing more operators that include a `?` character (special cased to support old pg_stat_statements query texts)
    - ltree extension [#136](https://github.com/pganalyze/libpg_query/pull/136)
    - promscale extension [#133](https://github.com/pganalyze/libpg_query/pull/133)
  - Deparser improvements
    - Prefix errors with "deparse", and remove some asserts [#131](https://github.com/pganalyze/libpg_query/pull/131)
    - Fix potential segfault when passing invalid protobuf (RawStmt without Stmt) [#128](https://github.com/pganalyze/libpg_query/pull/128)


## 2.1.0      2021-10-12

* Update libpg_query to 13-2.1.0 ([#53](https://github.com/pganalyze/pg_query_go/pull/53))
  - Normalize: add funcname error object
  - Normalize: Match GROUP BY against target list and re-use param refs
  - PL/pgSQL: Setup namespace items for parameters, support RECORD types
    - This significantly improves parsing for PL/pgSQL functions, to the
      extent that most functions should now parse successfully


## 2.0.5      2021-07-16

* Update libpg_query to 13-2.0.7 ([#49](https://github.com/pganalyze/pg_query_go/pull/49))
  - Normalize: Don't modify constants in TypeName typmods/arrayBounds fields
  - Don't fail builds on systems that have strchrnul support (FreeBSD)


## 2.0.4      2021-06-29

* Update libpg_query to 13-2.0.6 ([#47](https://github.com/pganalyze/pg_query_go/pull/47))
  - Normalize: Don't touch "ORDER BY 1" expressions, keep original text


## 2.0.3      2021-06-28

* Update libpg_query to 13-2.0.5 ([#45](https://github.com/pganalyze/pg_query_go/pull/45))
  - Update to Postgres 13.3 patch release
  - Normalize: Don't touch "GROUP BY 1" type statements, keep original text
  - Fingerprint: Cache list item hashes to fingerprint complex queries faster
  - Deparser: Emit the RangeVar catalogname if present
  - Fix crash in pg_scan function when encountering backslash escapes
* Add pg_query.Scan to access Postgres scanner ([#43](https://github.com/pganalyze/pg_query_go/pull/43))


## 2.0.2      2021-04-02

* Update to libpg_query 13-2.0.3
  - Normalize: Fix handling of two subsequent DefElem elements
* Parser CFLAGS: Build with -std=gnu99 to ensure CentOS compatibility


## 2.0.1      2021-03-30

* Update to libpg_query 13-2.0.2
  - Fix ARM builds: Avoid dependency on cpuid.h header
  - Simplify deparser of TableLikeClause
  - Fix asprintf warnings by ensuring _GNU_SOURCE is set early enough
* Add FingerprintToUInt64 method for callers that prefer to handle uint64s
  - This prevents a caller from having to do a hex string to uint64 conversion
    by simply returning the uint64 version of the fingerprint instead, which
    is already always provided by libpg_query.
* Add HashXXH3_64 helper method for generating XXH3 64-bit hash values
  - This can be useful when trying to fit other values into a data structure
    sized for the fingerprint, such as when encountering a parsing error
    during fingerprinting, and wanting to encode that fact uniquely into
    the fingerprint value.


## 2.0.0      2021-03-18

* Update libpg_query to 13-2.0.0
* Switch to use Protobuf generated nodes instead of custom logic
  - WARNING: This is breaking API change!
  - This is a breaking change in the API, but necessary in order to
    significantly improve the performance of parsing a query into Go structs,
    as well as allowing future bi-directional passing of parse trees between
    Go and C, such as for a future addition of a deparser.
* Rename pg_query.FastFingerprint to pg_query.Fingerprint
  - WARNING: This is breaking API change!
  - We no longer have direct support for running the fingerprint in the Go
    library, as its unnecessarily complex to support, and we can instead
    rely on the C fingerprinting method.
* Add support for deparsing parse trees back into a SQL statement
  - Call the new `pg_query.Deparse` method to get back the SQL text for
    a particular parse tree structure
  - This relies on the new deparser in libpg_query to transform a given parse
    tree back into a SQL statement. Note that this is currently considered
    experimental, and should not be used on unsanitized input, due to the risk
    of crashes in the C code with unexpected conditions.
* Update import path to `github.com/pganalyze/pg_query_go`


## 1.0.2      2021-02-18

* Update libpg_query to 10-1.0.5
  - This resolves memory leak problems, adds PPC architecture support,
    and refreshes the Postgres minor version to 10.16.


## 1.0.1      2020-11-07

* Update libpg_query to 10-1.0.3
  - This fixes ARM builds, and refreshes the Postgres patch version. This
    commit does not change the Postgres major version yet (still at PG 10).


## 1.0.0      2019-01-11

* Initial release with a tagged version
  - Note that 1.X releases will reflect Postgres 10 parser based releases from
    now on. When the Postgres 11 parser is released, that will be a new major
    version.
