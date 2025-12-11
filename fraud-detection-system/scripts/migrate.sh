#!/bin/bash

# Database migration script
# Usage: ./scripts/migrate.sh [up|down|create]

set -e

COMMAND=${1:-up}
MIGRATION_NAME=${2:-}

echo "📊 Database Migration Tool"
echo ""

case $COMMAND in
    up)
        echo "⬆️  Running migrations..."
        docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db << EOF
$(cat backend/internal/database/migrations/001_create_users.sql)
$(cat backend/internal/database/migrations/002_create_messages.sql)
$(cat backend/internal/database/migrations/003_create_verifications.sql)
$(cat backend/internal/database/migrations/004_create_reports.sql)
$(cat backend/internal/database/migrations/005_create_rbi_circulars.sql)
$(cat backend/internal/database/migrations/006_create_sender_registry.sql)
EOF
        echo "✅ Migrations completed"
        ;;
    
    down)
        echo "⬇️  Rolling back last migration..."
        echo "⚠️  This will drop tables!"
        read -p "Are you sure? (yes/no): " confirm
        if [ "$confirm" != "yes" ]; then
            echo "Rollback cancelled."
            exit 0
        fi
        
        # TODO: Implement proper rollback
        echo "❌ Rollback not implemented"
        exit 1
        ;;
    
    create)
        if [ -z "$MIGRATION_NAME" ]; then
            echo "❌ Migration name required"
            echo "Usage: ./scripts/migrate.sh create migration_name"
            exit 1
        fi
        
        # Get next migration number
        LAST_NUM=$(ls backend/internal/database/migrations/*.sql 2>/dev/null | tail -1 | grep -o '[0-9]\+' | head -1)
        NEXT_NUM=$(printf "%03d" $((10#$LAST_NUM + 1)))
        
        MIGRATION_FILE="backend/internal/database/migrations/${NEXT_NUM}_${MIGRATION_NAME}.sql"
        
        cat > "$MIGRATION_FILE" << EOF
-- Migration: $MIGRATION_NAME
-- Created: $(date)

-- TODO: Add your migration SQL here

-- Example:
-- CREATE TABLE example (
--     id UUID PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
EOF
        
        echo "✅ Created migration: $MIGRATION_FILE"
        echo "📝 Edit the file and add your SQL"
        ;;
    
    status)
        echo "📋 Migration Status:"
        echo ""
        echo "Applied migrations:"
        docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db -c "\dt"
        ;;
    
    *)
        echo "❌ Unknown command: $COMMAND"
        echo ""
        echo "Usage: ./scripts/migrate.sh [command]"
        echo ""
        echo "Commands:"
        echo "  up              - Run all migrations"
        echo "  down            - Rollback last migration"
        echo "  create <name>   - Create new migration"
        echo "  status          - Show migration status"
        exit 1
        ;;
esac

