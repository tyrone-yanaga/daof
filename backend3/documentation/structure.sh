issue - database not needed

project-root/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── checkout.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── cors.go
│   ├── models/
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── order.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── product.go
│   │   │   └── cart.go
│   │   └── interfaces.go
│   └── services/
│       ├── odoo/
│       │   └── client.go
│       ├── adyen/
│       │   └── client.go
│       └── product.go
├── pkg/
│   └── utils/
│       ├── logger.go
│       └── validator.go
├── migrations/
│   ├── 000001_create_products_table.up.sql
│   └── 000001_create_products_table.down.sql
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── .env.example
├── go.mod
└── Makefile