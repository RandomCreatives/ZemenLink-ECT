# ZemenLink: Enterprise Communication Platform

ZemenLink is an enterprise communication platform that bridges the gap between consumer-grade usability and enterprise-grade compliance. It provides a "Telegram-like" messaging experience while implementing the security, auditability, and data sovereignty required by modern enterprises.

## 🚀 Vision
The goal of ZemenLink is to eliminate "Shadow IT" by providing an intuitive, fast, and familiar UI/UX that employees love, backed by a robust architecture that CISOs trust.

## 🏗️ Architecture

### Core Pillars
- **Absolute Tenant Isolation:** Uses a "Database-per-Tenant" strategy. Each enterprise client operates in a siloed data environment with separate PostgreSQL databases and S3 buckets.
- **Federated Identity:** Integrates directly with enterprise Identity Providers (OIDC/SAML/Active Directory). No standalone login system.
- **E2EE with Corporate Escrow:** Implements the Signal Protocol for message transport, with a unique "Corporate Key Escrow" mechanism. Access to encrypted data requires a 2-person approval gate.
- **Offline-First Resilience:** The application remains hyper-responsive during network drops, with locally encrypted storage (SQLCipher).

### Tech Stack
- **Backend:** Golang (High concurrency messaging)
- **Database:** PostgreSQL (JSONB for flexible schemas)
- **Cache/Real-time:** Redis (Presence/Pub-Sub) & WebSockets
- **Frontend:** React Native (Mobile) and Electron/Web (Desktop)

## 📂 Project Structure
```text
.
├── cmd/server/             # Backend entry point
├── internal/
│   ├── tenant/             # Tenant management and DB routing
│   └── db/                 # Database utilities
├── ARCHITECTURAL_DECISIONS.md # Detailed architectural reasoning
├── DIAGRAMS.md            # Mermaid.js system flow diagrams
├── SCHEMA.sql             # Database schemas (Global & Tenant)
├── COMPLIANCE_DESIGN.md    # E2EE and Escrow mechanism details
└── go.mod                 # Go module definition
```

## 🛠️ Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL
- Redis

### Setup
1. **Initialize the Global Database:**
   Execute the first part of `SCHEMA.sql` to create the `zemenlink_global` database and the `tenants` table.

2. **Environment Variables:**
   Set the following environment variables:
   ```bash
   export GLOBAL_DB_URL="host=localhost user=postgres password=your_password dbname=zemenlink_global sslmode=disable"
   ```

3. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the Backend:**
   ```bash
   go run cmd/server/main.go
   ```

## 🛡️ Compliance & Security
For details on the 2-person approval flow for corporate escrow and data retention policies, refer to [COMPLIANCE_DESIGN.md](./COMPLIANCE_DESIGN.md).

## 🗺️ Roadmap
- **Phase 1:** Core Messenger (1:1 & Group chats, threading, file sharing).
- **Phase 2:** Enterprise Hooks ("Quiet Hours", Server Status, Institutional Knowledge).
- **Phase 3:** Integration Layer (Webhooks for Jira, GitHub, etc.).

## 📄 Documentation
- [Architectural Decisions](./ARCHITECTURAL_DECISIONS.md)
- [System Diagrams](./DIAGRAMS.md)
- [Database Schema](./SCHEMA.sql)
- [Compliance Design](./COMPLIANCE_DESIGN.md)
