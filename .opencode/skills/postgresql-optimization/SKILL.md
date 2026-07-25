---
name: postgresql-optimization
description: 'PostgreSQL-specific development and optimization covering JSONB, arrays, custom types, full-text search, window functions, extensions, and performance tuning.'
---

# PostgreSQL Optimization

Expert PostgreSQL guidance for optimization, advanced features, and performance tuning.

## Advanced Features
- **JSONB**: Containment queries, GIN indexes, path queries (#>>)
- **Arrays**: ANY, && (overlap), array_agg, unnest
- **Window functions**: PARTITION BY, ROWS BETWEEN, LAG/LEAD
- **Full-text search**: tsvector, ts_rank, GIN indexes
- **Custom types**: ENUM, DOMAIN, composite types (address_type)
- **Range types**: tstzrange, EXCLUDE USING gist
- **Geometric**: POINT, GiST indexes

## Performance Tuning
- `EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)` for analysis
- Composite indexes for multi-column queries
- Partial indexes for filtered queries
- Expression indexes for computed values
- Covering indexes with INCLUDE
- pg_stat_statements for slow query identification
- Cursor-based pagination (not OFFSET)

## Extensions
- uuid-ossp, pgcrypto, unaccent, pg_trgm, btree_gin

## Monitoring
- Database/table/index size queries
- Index usage statistics (unused indexes)
- Connection monitoring
- VACUUM and ANALYZE maintenance
