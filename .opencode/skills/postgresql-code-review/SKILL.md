---
name: postgresql-code-review
description: 'PostgreSQL-specific code review focusing on best practices, anti-patterns, JSONB, arrays, custom types, schema design, function optimization, and RLS.'
---

# PostgreSQL Code Review

Expert PostgreSQL code review focusing on PostgreSQL-specific best practices, anti-patterns, and quality standards.

## Review Areas
- **JSONB**: GIN indexes, containment operators (@>), validation constraints
- **Array operations**: GIN indexes, @> operator for contains
- **Schema design**: BIGSERIAL, CITEXT, TIMESTAMPTZ, JSONB
- **Custom types**: ENUM, DOMAIN for constrained values
- **Functions/Triggers**: WHEN (OLD.* IS DISTINCT FROM NEW.*)

## Anti-Patterns
- No GIN/GiST indexes for JSONB/arrays
- VARCHAR instead of TEXT/CITEXT
- Missing CHECK constraints
- Unstructured JSONB without validation

## Security
- Row Level Security (RLS) for multi-tenant data
- Granular permissions (not GRANT ALL)
- Built-in encryption (pgcrypto)

## Checklist
- [ ] Appropriate data types (CITEXT, JSONB, arrays)
- [ ] ENUM types for constrained values
- [ ] CHECK constraints for validation
- [ ] TIMESTAMPTZ not TIMESTAMP
- [ ] GIN indexes for JSONB/arrays
- [ ] RLS implemented where needed
