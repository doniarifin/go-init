# go-init

A CLI tool to quickly generate Go backend boilerplate projects.

---

## Installation

```bash
go install github.com/doniarifin/go-init@latest
```

Make sure `$GOPATH/bin` is in your PATH.

---

## Usage

### 1. Create a new project

```bash
go-init new myapp
```

### 2. Interactive mode

Run without flags (or with `-i`) to get a guided setup:

```bash
go-init new
go-init new myapp -i
```

### 3. Choose framework

```bash
go-init new myapp --framework=fiber
```

Available frameworks:

- fiber
- gin
- echo
- chi
- mux
- net-http (default)

### 4. Add authentication (JWT)

```bash
go-init new myapp --auth
```

### 5. Add database

```bash
go-init new myapp --db=postgres
```

Available databases:

- postgres
- mysql
- sqlite
- mongo
- redis
- none (default)

### 6. Generate CRUD boilerplate

```bash
# entity with default fields (name, email)
go-init new myapp --crud=user

# entity with custom fields, repeatable for multiple entities
go-init new myapp --crud=user:name,email --crud=product:title,price
```

Generates model, repository, service, handler, and routes for each entity.

### 7. Choose project structure

```bash
go-init new myapp --structure=clean
```

Available layouts:

- standard (default) - `internal/{model,repository,service,handler}`
- clean - `internal/{entities,repositories,usecases,delivery/http}`
- hexagonal - `internal/{core/domain,core/service,adapters/out,adapters/in/http}`

### 8. Add Docker support

```bash
go-init new myapp --docker
```

### 9. Add Makefile, .gitignore, or migrations

```bash
go-init new myapp --makefile --gitignore --migrations
```

### 10. Full setup

```bash
go-init new myapp --framework=fiber --auth --db=postgres --docker --crud=user:name,email --structure=clean --makefile --gitignore --migrations
```

---

## Other commands

```bash
go-init list      # show all available options
go-init version   # print version
```

---

## Generated Structure

```bash
myapp/
 ├── main.go
 ├── go.mod
 ├── .env
 ├── Dockerfile          (with --docker)
 ├── Makefile            (with --makefile)
 ├── .gitignore          (with --gitignore)
 ├── config/
 │   └── database.go     (with --db)
 ├── migrations/         (with --migrations)
 ├── internal/
 │   ├── model/          (CRUD entities)
 │   ├── repository/     (in-memory CRUD, swap for your DB)
 │   ├── service/        (business logic)
 │   ├── handler/        (HTTP handlers)
 │   └── route/
 └── pkg/
```

The exact layout changes with `--structure`.

---

## Run the Project

```bash
cd myapp
go run main.go
```

---

## Tips

- Edit `.env` to configure database and JWT settings
- Use Docker for easier setup
- The project structure is ready for production use
- CRUD repository starts in-memory (mutex + map); replace it with your DB driver

---

## Example

```bash
go-init new api-service --framework=fiber --auth --db=postgres --docker --crud=user
cd api-service
go run main.go
```
