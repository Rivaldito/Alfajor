# 🍬 Alfajor Logger

[](https://www.google.com/search?q=https://goreportcard.com/report/github.com/Rivaldito/alfajor)
[](https://opensource.org/licenses/MIT)

A sweet, powerful, and structured logging package for Go. ✨ Built with an **abstraction-first (Hexagonal)** approach, it uses `zap` for high-performance file/console logging and seamlessly integrates with SQL databases for persistent auditing.

### 🌟 Key Features

  * **🔌 Pluggable Architecture**: Built on a `Logger` interface. Swap implementation details without touching your business logic.
  * **🛢️ SQL Database Support**: Native support for logging directly to **PostgreSQL**, **MySQL**, and **MariaDB**.
  * **✨ Hybrid Logging**: Automatically combines Console/File logging with Database logging when a DB connection is provided. Best of both worlds\!
  * **⚡ Async & Non-Blocking**: SQL insertions happen in background goroutines (Fire-and-Forget), ensuring your application performance isn't affected by database latency.
  * **🗂️ Structured Context**: Add key-value context (`map[string]interface{}`) to your logs. In databases, this is stored as JSON text for easy querying.
  * **🖥️ & 📄 Dual Output**: Log to console (colored) and file (JSON) simultaneously.
  * **🔄 Automatic Rotation**: Built-in log rotation based on size, age, and backups using `lumberjack`.

-----

### 💾 Installation

```bash
go get github.com/Rivaldito/alfajor@v1.0.0
```

-----

### 🗄️ Database Setup

To use the SQL logging feature, you must create the destination table in your database.

**PostgreSQL (Recommended Schema)**

```sql
CREATE TABLE alfajor_logs (
    id SERIAL PRIMARY KEY,
    level VARCHAR(10),
    message TEXT,
    error TEXT,
    context JSONB, -- Use JSON for MySQL/MariaDB
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

> **Note for MySQL/MariaDB users:** Change `JSONB` to `JSON` or `TEXT` depending on your version.

-----

### 🚀 Usage Examples

#### 1\. The "All-in-One" Hybrid Mode (Recommended)

This is the most powerful feature. By simply passing your DB connection, **Alfajor** detects it and automatically logs to **Console + File + Database** simultaneously.

```go
package main

import (
    "database/sql"
    "github.com/Rivaldito/alfajor"
    _ "github.com/lib/pq" // Don't forget your driver!
)

func main() {
    // 1. Setup Database Connection
    db, _ := sql.Open("postgres", "user:pass@localhost/dbname?sslmode=disable")

    // 2. Configure Alfajor
    cfg := alfajor.NewDefaultConfig()
    
    // 3. Initialize Hybrid Logger
    // By passing WithSQLDB, Alfajor activates the "MultiLogger" automatically.
    // We specify DialectPostgres to use '$1' placeholders.
    dulce, err := alfajor.New(cfg, "zap", 
        alfajor.WithSQLDB(db, "alfajor_logs", alfajor.DialectPostgres),
    )
    if err != nil {
        panic(err)
    }
    defer dulce.Sync()

    // 4. Log once, write everywhere!
    // This appears in your terminal, your log file, AND your 'alfajor_logs' table.
    dulce.Info("New user registered", map[string]interface{}{
        "user_id": 42,
        "plan":    "premium",
    })
}
```

#### 2\. MySQL / MariaDB Logging

If you are using MySQL, the placeholder syntax (`?`) is handled automatically (default) or explicitly via `DialectMySQL`.

```go
import (
    "database/sql"
    "github.com/Rivaldito/alfajor"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, _ := sql.Open("mysql", "user:password@/dbname")
    cfg := alfajor.NewDefaultConfig()

    // Explicitly requesting "sql" logger only (No console/file)
    dulce, _ := alfajor.New(cfg, "sql", 
        alfajor.WithSQLDB(db, "alfajor_logs", alfajor.DialectMySQL),
    )

    dulce.Warn("High memory usage detected", map[string]interface{}{
        "usage_percent": 85,
    })
}
```

#### 3\. Basic File & Console (Classic Mode)

No database needed? No problem. It works as a lightweight wrapper around Zap.

```go
func main() {
    cfg := alfajor.NewDefaultConfig()
    cfg.EnableFile = true
    cfg.Rotation.Filename = "logs/server.log"

    // No DB options passed -> Only creates the Zap adapter
    dulce, _ := alfajor.New(cfg, "zap") 

    dulce.Info("Server started on port 8080")
}
```

-----

### ⚙️ Configuration Options

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `Level` | `string` | `"info"` | Log level (`debug`, `info`, `warn`, `error`, `fatal`) |
| `EnableConsole` | `bool` | `true` | Print logs to stdout |
| `EnableFile` | `bool` | `false` | Save logs to a local file |
| `Rotation` | `struct` | `{...}` | Settings for `lumberjack` (MaxSize, MaxAge, etc.) |
| `WithSQLDB` | `func` | `nil` | Inject DB connection, table name, and SQL Dialect |

#### Supported SQL Dialects

  * `alfajor.DialectMySQL` (Default): Uses `?` syntax. Compatible with **MySQL**, **MariaDB**, and **SQLite**.
  * `alfajor.DialectPostgres`: Uses `$1, $2` syntax. Compatible with **PostgreSQL** and **CockroachDB**.