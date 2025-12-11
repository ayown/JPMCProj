#!/bin/bash

# Database backup script
# Usage: ./scripts/backup.sh

set -e

echo "💾 Starting database backup..."

# Configuration
BACKUP_DIR="database/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_$DATE.sql.gz"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Check if PostgreSQL is running
if ! docker-compose ps | grep -q "postgres.*Up"; then
    echo "❌ PostgreSQL is not running"
    exit 1
fi

# Perform backup
echo "📦 Creating backup: $BACKUP_FILE"
docker-compose exec -T postgres pg_dump -U frauddetection frauddetection_db | gzip > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Backup completed successfully"
    echo "📁 Backup file: $BACKUP_FILE"
    
    # Get backup size
    SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo "📊 Backup size: $SIZE"
    
    # Keep only last 7 backups
    echo "🧹 Cleaning old backups (keeping last 7)..."
    cd "$BACKUP_DIR"
    ls -t backup_*.sql.gz | tail -n +8 | xargs -r rm
    
    echo "✅ Backup process completed"
else
    echo "❌ Backup failed"
    exit 1
fi

# Optional: Upload to cloud storage
# TODO: Implement cloud backup
# aws s3 cp "$BACKUP_FILE" s3://your-bucket/backups/
# gsutil cp "$BACKUP_FILE" gs://your-bucket/backups/

