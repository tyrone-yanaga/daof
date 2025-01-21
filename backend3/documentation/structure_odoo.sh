structureyour-project/
├── cmd/
│   └── api/
│       └── main.go               # Application entry point
│
├── internal/
│   ├── config/
│   │   ├── config.go            # Configuration structures and loading
│   │   └── env.go               # Environment variable handling
│   │
│   ├── handlers/
│   │   ├── product.go           # Product HTTP handlers
│   │   ├── cart.go              # Cart HTTP handlers
│   │   ├── checkout.go          # Checkout HTTP handlers
│   │   └── order.go             # Order HTTP handlers
│   │
│   ├── middleware/
│   │   ├── auth.go              # Authentication middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── logger.go            # Logging middleware
│   │   └── recovery.go          # Error recovery middleware
│   │
│   ├── models/
│   │   ├── product.go           # Product model structs
│   │   ├── cart.go              # Cart model structs
│   │   ├── order.go             # Order model structs
│   │   └── user.go              # User model structs
│   │
│   ├── services/
│   │   ├── product.go           # Product business logic
│   │   ├── cart.go              # Cart business logic
│   │   ├── checkout.go          # Checkout business logic
│   │   └── order.go             # Order business logic
│   │
│   └── utils/
│       ├── logger.go            # Logging utilities
│       ├── hash.go              # Hashing utilities
│       └── validator.go         # Validation utilities
│
├── pkg/
│   └── odoo/
│       ├── client.go            # Odoo API client
│       └── types.go             # Odoo data types
│
├── docker/
│   ├── Dockerfile              # API service Dockerfile
│   └── docker-compose.yml      # Docker compose configuration
│
├── odoo/
│   ├── config/
│   │   └── odoo.conf           # Odoo configuration
│   └── addons/
│       └── .gitkeep            # Custom Odoo modules go here
│
├── scripts/
│   └── migrate.sh              # Database migration script
│
├── .env.example                # Example environment variables
├── .gitignore                  # Git ignore file
├── go.mod                      # Go module file
├── go.sum                      # Go dependencies checksum
├── Makefile                    # Build and development commands
└── README.md                   # Project documentation