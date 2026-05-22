# ZemenLink Regional & Compliance Packs

To meet the diverse regulatory and cultural needs of global enterprises, ZemenLink uses a **Pack-based Modular System**.

## 1. The Regional Pack Concept
A "Regional Pack" is a collection of modules, configurations, and assets tailored for a specific jurisdiction.

### Example: Ethiopia Compliance Pack (`zemenlink_et_pack`)
- **Data Sovereignty Module:** Ensures that for tenants flagged as "ET-Resident," all S3 and PostgreSQL traffic is routed to servers located physically within Ethiopia (or approved regional nodes).
- **Localization Module:** Overrides default UI strings with Amharic, Oromo, Tigrinya, etc.
- **Regulatory Reporting:** Includes specialized auditing templates that meet the National Bank of Ethiopia (NBE) requirements.

## 2. Manifest-Driven Features
Each module is controlled by a `manifest.json`. The ZemenLink Kernel loads these manifests based on the tenant's subscription and compliance level.

### Sample Manifest (`et_compliance_pack/manifest.json`)
```json
{
  "name": "et_compliance_pack",
  "version": "1.0.0",
  "enabled": true,
  "config": {
    "data_sovereignty": "strict_ethiopia",
    "local_languages": ["am", "om", "ti"],
    "regulatory_reporting": "nbe_standard_v2"
  }
}
```

## 3. Module Interceptors
Modules can register "Interceptors" into the ZemenLink Kernel's pipeline.
- *Example:* The `et_compliance_pack` might register an interceptor that logs every file upload specifically for a regulatory audit trail.
