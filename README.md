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

---

### 2. Choose framework

```bash
go-init new myapp --framework=fiber
```

Available frameworks:

- fiber
- gin
- net-http (default)

---

### 3. Add authentication (JWT)

```bash
go-init new myapp --auth
```

---

### 4. Add database

```bash
go-init new myapp --db=postgres
```

---

### 5. Add Docker support

```bash
go-init new myapp --docker
```

---

### 6. Full setup

```bash
go-init new myapp --framework=fiber --auth --db=postgres --docker
```

---

## Generated Structure

```bash
myapp/
 ├── main.go
 ├── go.mod
 ├── .env
 ├── Dockerfile
 ├── config/
 ├── internal/
 └── pkg/
```

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

---

## Example

```bash
go-init new api-service --framework=fiber --auth --db=postgres --docker
cd api-service
go run main.go
```
