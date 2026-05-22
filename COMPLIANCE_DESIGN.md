# Compliance Dashboard & Corporate Escrow Design

## 1. Corporate Key Escrow Mechanism
ZemenLink ensures that while messages are End-to-End Encrypted (E2EE) for privacy, the organization retains the legal right to audit communications.

### Encryption Protocol
- **Sender (S)** sends a message to **Recipient (R)**.
- **S** generates a random Symmetric Key (**K**).
- Message content is encrypted with **K**.
- **K** is then encrypted multiple times:
    1. `Enc(Public_Key_R, K)` -> For the recipient.
    2. `Enc(Public_Key_Escrow, K)` -> For corporate escrow.
- Both encrypted versions of **K** and the encrypted message are sent to the server.

### Key Management
- `Public_Key_Escrow` is a tenant-specific public key.
- The corresponding `Private_Key_Escrow` is stored in a secure Hardware Security Module (HSM) or a highly secure Key Management Service (KMS) like AWS KMS or HashiCorp Vault.
- Access to `Private_Key_Escrow` is strictly gated by software logic.

## 2. Compliance Dashboard
The Compliance Dashboard is a restricted UI available only to users with the `compliance_officer` or `admin` role.

### Features:
- **Audit Logs:** Real-time view of all administrative and compliance actions.
- **Data Export:** Export chat history for specific users/timeframes in standardized formats (JSON, CSV, PDF) for legal discovery.
- **Retention Policies:** Set per-tenant rules for how long messages are kept (e.g., "Delete after 7 years" or "Never delete").
- **Escrow Access Management:** The interface to request and approve access to encrypted content.

## 3. Escrow Access Flow (The Event-Driven Gate)
The "Two-Person Gate" is implemented as an event-driven workflow engine, ensuring that escrow access is never a static permission but a transient, multi-sig state.

### Workflow Stages:
1. **Event: `EscrowKeyRequested`**
    - A Compliance Officer (CO) initiates a request.
    - System triggers the "Escrow Workflow Module."
2. **State: `PendingApproval`**
    - Notifications are dispatched to the `AdminGroup` (e.g., Legal Counsel, CTO) via ZemenLink's own system channels.
    - The request is locked in a pending state with a fixed expiration (e.g., 24 hours).
3. **Transition: `On(ApprovalCount >= 2)`**
    - As approvers sign off, the workflow state is updated.
    - Once the threshold is met, the system triggers the `ReleaseKey` event.
4. **Grant:** The system temporarily grants the CO access to a decryption service.
5. **Decryption:** The decryption service uses the `Private_Key_Escrow` (from KMS) to decrypt the symmetric keys for the requested messages.
6. **Logging:** Every step—request, approval, and actual decryption/viewing—is recorded in the immutable `audit_logs` table.

## 4. Institutional Knowledge Retention
- When a user is offboarded via the OIDC/LDAP integration, their `users.is_active` flag is set to `false`.
- Their messages remain in the tenant database, preserved as part of the corporate record.
- Group chats and threaded conversations remain intact for the remaining participants.
