#!/usr/bin/env bash
# ── Muse Backup Script ─────────────────────────────────────────────
# Backs up PostgreSQL and Redis data volumes to a configurable directory.
#
# Usage:
#   ./backup.sh                          # backups to ./backups/
#   ./backup.sh /path/to/backups         # custom directory
#   BACKUP_DIR=/custom ./backup.sh       # env var override
#
# Schedule with cron (daily at 3am):
#   0 3 * * * /opt/music/deployments/scripts/backup.sh /mnt/backups
#
# Restore:
#   pg_restore -U postgres -d musicapp -F custom latest_postgres.dump
#   cat latest_redis.rdb > /var/lib/redis/dump.rdb

set -euo pipefail

BACKUP_DIR="${1:-${BACKUP_DIR:-./backups}}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
mkdir -p "${BACKUP_DIR}"

echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting backup to ${BACKUP_DIR}"

# ── PostgreSQL ─────────────────────────────────────────────────────
# Uses pg_dump with custom format (compressed, parallel-restore capable).
# Requires PG* env vars or running container.
# Assumes docker compose context - adjust if running on bare metal.

PG_CONTAINER="${PG_CONTAINER:-musicapp_postgres}"
PG_USER="${PG_USER:-postgres}"
PG_DB="${PG_DB:-musicapp}"

if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^${PG_CONTAINER}$"; then
  echo "  Backing up PostgreSQL (container: ${PG_CONTAINER})..."
  docker exec "${PG_CONTAINER}" pg_dump -U "${PG_USER}" -d "${PG_DB}" -F custom \
    --compress=9 \
    --no-owner \
    --clean \
    --if-exists \
    > "${BACKUP_DIR}/${TIMESTAMP}_postgres.dump"
  echo "  PostgreSQL: ${TIMESTAMP}_postgres.dump ($(du -h "${BACKUP_DIR}/${TIMESTAMP}_postgres.dump" | cut -f1))"
else
  echo "  WARNING: PostgreSQL container '${PG_CONTAINER}' not running. Skipping."
fi

# ── Redis ──────────────────────────────────────────────────────────
REDIS_CONTAINER="${REDIS_CONTAINER:-musicapp_redis}"

if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^${REDIS_CONTAINER}$"; then
  echo "  Backing up Redis (container: ${REDIS_CONTAINER})..."
  docker exec "${REDIS_CONTAINER}" redis-cli SAVE
  docker cp "${REDIS_CONTAINER}:/data/dump.rdb" "${BACKUP_DIR}/${TIMESTAMP}_redis.rdb"
  echo "  Redis: ${TIMESTAMP}_redis.rdb ($(du -h "${BACKUP_DIR}/${TIMESTAMP}_redis.rdb" | cut -f1))"
else
  echo "  WARNING: Redis container '${REDIS_CONTAINER}' not running. Skipping."
fi

# ── Retention ──────────────────────────────────────────────────────
RETENTION_DAYS="${RETENTION_DAYS:-7}"
echo "  Cleaning backups older than ${RETENTION_DAYS} days..."
find "${BACKUP_DIR}" -name "*_postgres.dump" -type f -mtime "+${RETENTION_DAYS}" -delete
find "${BACKUP_DIR}" -name "*_redis.rdb" -type f -mtime "+${RETENTION_DAYS}" -delete

# ── Latest symlink ─────────────────────────────────────────────────
ln -sf "${TIMESTAMP}_postgres.dump" "${BACKUP_DIR}/latest_postgres.dump"
ln -sf "${TIMESTAMP}_redis.rdb" "${BACKUP_DIR}/latest_redis.rdb"

echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup complete"
echo "  Latest: ${BACKUP_DIR}/latest_postgres.dump, ${BACKUP_DIR}/latest_redis.rdb"
