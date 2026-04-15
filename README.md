**This is not the official Atlas repo, this is a fork!!** For the official repo, [see here](https://github.com/ariga/atlas).

# Sam's Atlas Fork

This is a fork of [Atlas](https://github.com/ariga/atlas), an open source schema migration tool written in Go.
Unfortunately some of the features are hidden behind a paywall. This fork aims to implement the ones I need
in public under an open source license (Apache 2).

None of the code written here is written with access to the non-public, paid source code for the paid features.
It's solely implemented based on public documentation. No contributor who has access to the code may derive
their implementations from Atlas' private, BSL-licensed code. All code must be original.

Thanks goes to the original Atlas authors. Just like the original repo, this is licensed under Apache2.

## Currently implemnted features

None so far, but more to be added here.

## Rest of the readme

Atlas is a language-agnostic tool for managing and migrating database schemas using modern DevOps principles.
It offers two workflows:

- **Declarative**: Similar to Terraform, Atlas compares the current state of the database to the desired state, as
  defined in an [HCL], [SQL], or [ORM] schema. Based on this comparison, it generates and executes a migration plan to
  transition the database to its desired state.

- **Versioned**: Unlike other tools, Atlas automatically plans schema migrations for you. Users can describe their desired
  database schema in [HCL], [SQL], or their chosen [ORM], and by utilizing Atlas, they can plan, lint, and apply the
  necessary migrations to the database.


## Supported Databases

[PostgreSQL](https://atlasgo.io/guides/postgres) ·
[MySQL](https://atlasgo.io/guides/mysql) ·
[SQLite](https://atlasgo.io/guides/sqlite) ·

## Installation

Currently no supported installation. Working on it!

## Key Features

- **[Declarative schema migrations](https://atlasgo.io/declarative/apply)**: The `atlas schema` command offers various options for [inspecting](https://atlasgo.io/inspect), diffing, comparing, [planning](https://atlasgo.io/declarative/plan) and applying migrations using standard Terraform-like workflows.
- **[Versioned migrations](https://atlasgo.io/versioned/intro)**: The `atlas migrate` command provides a state-of-the-art experience for [planning](https://atlasgo.io/versioned/diff), [linting](https://atlasgo.io/lint/analyzers), and [applying](https://atlasgo.io/versioned/apply) migrations.
- **[Schema as Code](https://atlasgo.io/atlas-schema)**: Define your desired database schema using [SQL], [HCL], or your chosen [ORM]. Atlas supports [16 ORM loaders](https://atlasgo.io/orms) across 6 languages.
- **[50+ safety analyzers](https://atlasgo.io/lint/analyzers)**: Database-aware migration linting that detects destructive changes, data-dependent modifications, table locks, backward-incompatible changes, and more.
- **[Multi-tenancy](https://atlasgo.io/guides/multi-tenancy)**: Built-in support for multi-tenant database migrations.

## Getting Started

Get started with Atlas by following the [Getting Started](https://atlasgo.io/getting-started/) docs.

Inspect an existing database schema:
```shell
atlas schema inspect -u "postgres://localhost:5432/mydb"
```

Apply your desired schema to the database:
```shell
atlas schema apply \
  --url "postgres://localhost:5432/mydb" \
  --to file://schema.hcl \
  --dev-url "docker://postgres/16/dev"
```

📖 [Getting Started docs](https://atlasgo.io/getting-started/)

## Migration Linting

Atlas ships with 50+ built-in [analyzers](https://atlasgo.io/lint/analyzers) that review your migration files
and catch issues before they reach production. Analyzers detect [destructive changes](https://atlasgo.io/lint/analyzers#destructive-changes)
like dropped tables or columns, [data-dependent modifications](https://atlasgo.io/lint/analyzers#data-dependent-changes)
such as adding non-nullable columns without defaults, and database-specific risks like table locks
and table rewrites that can cause downtime on busy tables. You can also define
your own [custom policy rules](https://atlasgo.io/lint/rules).

```bash
atlas migrate lint --dev-url "docker://postgres/16/dev"
```

## Schema Testing

[Test](https://atlasgo.io/testing/schema) database logic (functions, views, triggers, procedures) and
[data migrations](https://atlasgo.io/testing/migrate) with `.test.hcl` files:

```hcl
test "schema" "postal" {
  # Valid postal codes pass
  exec {
    sql = "SELECT '12345'::us_postal_code"
  }
  # Invalid postal codes fail
  catch {
    sql = "SELECT 'hello'::us_postal_code"
  }
}

test "schema" "seed" {
  for_each = [
    {input: "hello", expected: "HELLO"},
    {input: "world", expected: "WORLD"},
  ]
  exec {
    sql    = "SELECT upper('${each.value.input}')"
    output = each.value.expected
  }
}
```

```bash
atlas schema test --dev-url "docker://postgres/16/dev"
```

📖 [Testing docs](https://atlasgo.io/testing/schema)

## ORM Support

Define your schema in any of the 16 supported ORMs. Atlas reads your models and generates migrations:

| Language | ORMs |
|----------|------|
| Go | [GORM](https://atlasgo.io/guides/orms/gorm), [Ent](https://atlasgo.io/guides/orms/ent), [Bun](https://atlasgo.io/guides/orms/bun), [Beego](https://atlasgo.io/guides/orms/beego), [sqlc](https://atlasgo.io/guides/frameworks/sqlc-versioned) |
| TypeScript | [Prisma](https://atlasgo.io/guides/orms/prisma), [Drizzle](https://atlasgo.io/guides/orms/drizzle), [TypeORM](https://atlasgo.io/guides/orms/typeorm), [Sequelize](https://atlasgo.io/guides/orms/sequelize) |
| Python | [Django](https://atlasgo.io/guides/orms/django), [SQLAlchemy](https://atlasgo.io/guides/orms/sqlalchemy) |
| Java | [Hibernate](https://atlasgo.io/guides/orms/hibernate) |
| .NET | [EF Core](https://atlasgo.io/guides/orms/efcore) |
| PHP | [Doctrine](https://atlasgo.io/guides/orms/doctrine) |

📖 [ORM integration docs](https://atlasgo.io/orms)

## Integrations

Lint, test, and apply migrations automatically in your CI/CD pipeline or infrastructure-as-code workflow:

| Integration | Docs |
|-------------|------|
| GitHub Actions | [Versioned guide](https://atlasgo.io/guides/ci-platforms/github-versioned) · [Declarative guide](https://atlasgo.io/guides/ci-platforms/github-declarative) |
| GitLab CI | [Versioned guide](https://atlasgo.io/guides/ci-platforms/gitlab-versioned) · [Declarative guide](https://atlasgo.io/guides/ci-platforms/gitlab-declarative) |
| CircleCI | [Versioned guide](https://atlasgo.io/guides/ci-platforms/circleci-versioned) · [Declarative guide](https://atlasgo.io/guides/ci-platforms/circleci-declarative) |
| Bitbucket Pipes | [Versioned guide](https://atlasgo.io/guides/ci-platforms/bitbucket-versioned) · [Declarative guide](https://atlasgo.io/guides/ci-platforms/bitbucket-declarative) |
| Azure DevOps | [GitHub repos](https://atlasgo.io/guides/ci-platforms/azure-devops-github) · [Azure repos](https://atlasgo.io/guides/ci-platforms/azure-devops-repos) |
| Terraform Provider | [atlasgo.io/integrations/terraform-provider](https://atlasgo.io/integrations/terraform-provider) |
| Kubernetes Operator | [atlasgo.io/integrations/kubernetes](https://atlasgo.io/integrations/kubernetes) |
| ArgoCD | [atlasgo.io/guides/deploying/k8s-argo](https://atlasgo.io/guides/deploying/k8s-argo) |
| Flux | [atlasgo.io/guides/deploying/k8s-flux](https://atlasgo.io/guides/deploying/k8s-flux) |
| Crossplane | [atlasgo.io/guides/deploying/crossplane](https://atlasgo.io/guides/deploying/crossplane) |
| Go SDK | [pkg.go.dev/ariga.io/atlas-go-sdk/atlasexec](https://pkg.go.dev/ariga.io/atlas-go-sdk/atlasexec) |

## CLI Usage

### `schema inspect`

_**Easily inspect your database schema by providing a database URL and convert it to HCL, JSON, SQL, ERD, or other formats.**_

Inspect a specific MySQL schema and get its representation in Atlas DDL syntax:
```shell
atlas schema inspect -u "mysql://root:pass@localhost:3306/example" > schema.hcl
```

<details><summary>Result</summary>

```hcl
table "users" {
  schema = schema.example
  column "id" {
    null = false
    type = int
  }
  ...
}
```
</details>

Inspect the entire MySQL database and get its JSON representation:
```shell
atlas schema inspect \
  --url "mysql://root:pass@localhost:3306/" \
  --format '{{ json . }}' | jq
```

<details><summary>Result</summary>

```json
{
  "schemas": [
    {
      "name": "example",
      "tables": [
        {
          "name": "users",
          "columns": [
            ...
          ]
        }
      ]
    }
  ]
}
```
</details>

Inspect a specific PostgreSQL schema and get its ERD representation in Mermaid syntax:
```shell
atlas schema inspect \
  --url "postgres://root:pass@:5432/test?search_path=public&sslmode=disable" \
  --format '{{ mermaid . }}'
```

```mermaid
erDiagram
    users {
      int id PK
      varchar name
    }
    blog_posts {
      int id PK
      varchar title
      text body
      int author_id FK
    }
    blog_posts }o--o| users : author_fk
```

Use the [split format](https://atlasgo.io/inspect/database-to-code) for one-file-per-object output:

```bash
atlas schema inspect -u '<url>' --format '{{ sql . | split | write }}'
```

```
├── schemas
│   └── public
│       ├── public.sql
│       ├── tables
│       │   ├── profiles.sql
│       │   └── users.sql
│       ├── functions
│       └── types
└── main.sql
```

📖 [Schema inspection docs](https://atlasgo.io/inspect)

### `schema diff`

_**Compare two schema states and get a migration plan to transform one into the other. A state can be specified using a
database URL, HCL, SQL, or ORM schema, or a migration directory.**_

```shell
atlas schema diff \
  --from "postgres://postgres:pass@:5432/test?search_path=public&sslmode=disable" \
  --to file://schema.hcl \
  --dev-url "docker://postgres/15/test"
```

📖 [Declarative workflow docs](https://atlasgo.io/declarative/apply)

### `schema apply`

_**Generate a migration plan and apply it to the database to bring it to the desired state. The desired state can be
specified using a database URL, HCL, SQL, or ORM schema, or a migration directory.**_

```shell
atlas schema apply \
  --url mysql://root:pass@:3306/db1 \
  --to file://schema.hcl \
  --dev-url docker://mysql/8/db1
```

<details><summary>Result</summary>

```shell
-- Planned Changes:
-- Modify "users" table
ALTER TABLE `db1`.`users` DROP COLUMN `d`, ADD COLUMN `c` int NOT NULL;
Use the arrow keys to navigate: ↓ ↑ → ←
? Are you sure?:
  ▸ Apply
    Abort
```
</details>

📖 [Declarative workflow docs](https://atlasgo.io/declarative/apply)

### `migrate diff`

_**Write a new migration file to the migration directory that brings it to the desired state. The desired state can be
specified using a database URL, HCL, SQL, or ORM schema, or a migration directory.**_

```shell
atlas migrate diff add_blog_posts \
  --dir file://migrations \
  --to file://schema.hcl \
  --dev-url docker://mysql/8/test
```

📖 [Versioned workflow docs](https://atlasgo.io/versioned/diff)

### `migrate apply`

_**Apply all or part of pending migration files in the migration directory on the database.**_

```shell
atlas migrate apply \
  --url mysql://root:pass@:3306/db1 \
  --dir file://migrations
```

📖 [Versioned workflow docs](https://atlasgo.io/versioned/apply)
