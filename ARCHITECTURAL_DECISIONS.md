# ZemenLink High-Level System Architecture

## 1. Overview
ZemenLink is an enterprise-grade messaging platform designed to provide a "Telegram-like" user experience with the security and compliance required by large organizations. The system follows a modular monolith or microservices architecture depending on scale, with a strong focus on tenant isolation.

## 2. Tenant Silo Strategy: Database-per-Tenant
To ensure absolute data isolation and comply with strict enterprise requirements, ZemenLink employs a **Database-per-Tenant** strategy.

### Implementation Details:
- **Logical Isolation:** Each tenant (Company A, Company B) has its own dedicated PostgreSQL database.
- **Connection Management:** The API Gateway/Backend uses a "Tenant Router" to determine the correct database connection string based on the user's authenticated context (Tenant ID).
- **Schema Management:** All tenant databases share the same schema. Migrations are managed centrally and applied to all tenant databases in a rolling fashion to ensure consistency.
- **Scalability:**
    - Small/Medium Tenants: Multiple logical databases can reside on a single RDS/PostgreSQL instance.
    - Large Tenants: High-usage tenants can be migrated to a dedicated RDS instance without affecting others, solving the "Noisy Neighbor" problem.
- **Maintenance & Maintainability:**
    - **Connection Pooling:** We use a centralized connection pooler (like PgBouncer) or a managed pool within the Go backend that dynamically opens/closes connections to tenant databases based on demand.
    - **Uniform Schema:** We enforce a strictly uniform schema across all tenants. No custom schema changes are allowed for individual tenants to prevent "snowflake" databases.
    - **Automated Migrations:** Tools like `golang-migrate` or `atlas` are integrated into the CI/CD pipeline. When a schema change is merged, the pipeline iterates through the list of active tenant databases and applies the migration.
    - **Health Monitoring:** Centralized monitoring (Prometheus/Grafana) tracks performance metrics per tenant database to proactively identify scaling needs.

## 3. System Flow (Mermaid)

```mermaid
graph TD
    Client[ZemenLink Client Mobile/Desktop]
    Gateway[API Gateway / Load Balancer]
    AuthService[Auth Service OIDC/SAML]
    MessagingService[Messaging Service Go]
    GlobalDB[(Global Management DB)]
    Redis[(Redis Presence/Cache)]

    subgraph "Tenant A Silo"
        DB_A[(PostgreSQL DB A)]
        S3_A[S3 Bucket A]
    end

    subgraph "Tenant B Silo"
        DB_B[(PostgreSQL DB B)]
        S3_B[S3 Bucket B]
    end

    Client --> Gateway
    Gateway --> AuthService
    AuthService -- Validate JWT/Tenant --> MessagingService
    MessagingService -- Lookup Connection --> GlobalDB
    MessagingService -- Pub/Sub --> Redis
    MessagingService -- Query/Write --> DB_A
    MessagingService -- Query/Write --> DB_B
    MessagingService -- Upload/Download --> S3_A
    MessagingService -- Upload/Download --> S3_B
```

## 4. Core Components

### Identity & Access Management (IAM)
- **Federated Identity:** ZemenLink acts as a Service Provider (SP). Integration with OIDC/SAML/Active Directory is the primary onboarding mechanism.
- **Zero Standalone Login:** Authentication is delegated to the enterprise Identity Provider (IdP) (e.g., Okta, Azure AD).

### Messaging Service (Backend)
- **Language:** Golang (for high concurrency).
- **Real-time:** WebSockets for message delivery and presence updates, backed by Redis for pub/sub and transient state.
- **Storage:** Metadata and messages are stored in the tenant's dedicated PostgreSQL database.
- **Attachments:** Files are stored in tenant-specific S3 buckets (or prefixed folders within a shared bucket, depending on compliance tier).

### E2EE with Corporate Escrow
- **Transport:** Signal Protocol for E2EE.
- **Escrow:** Messages are multi-encrypted. In addition to the recipient's public key, they are encrypted with a Corporate Escrow Public Key.
- **Access:** Decryption requires an "Escrow Access Request" through the Compliance Dashboard, mandating a 2-person approval gate.

### Frontend (Client)
- **Stack:** React Native (Mobile) and Electron (Desktop).
- **Offline-First:** Local storage using SQLite with **SQLCipher** for full-disk encryption.
- **Key Management:** Local encryption keys are derived from the user's session, never stored in plain text.

## 4. Maintenance & Compliance
- **Global Management DB:** A centralized (highly restricted) database stores tenant metadata, connection strings, and global configurations.
- **Audit Logs:** Every action in the system, especially compliance-related ones, is logged to an immutable audit trail.
