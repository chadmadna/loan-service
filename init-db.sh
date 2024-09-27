#!/bin/bash

# Change ownership of data dir
chown -R postgres:postgres /var/lib/postgresql/data

# Initialize Postgres
if [ ! -d "/var/lib/postgresql/data" ] || [ ! -s "/var/lib/postgresql/data/PG_VERSION" ]; then
  su-exec postgres initdb
fi

# Start Postgres and wait until ready
PGDATA="/var/lib/postgresql/data"
PIDFILE="$PGDATA/postmaster.pid"
if [ -f "$PIDFILE" ]; then
  if ps -p "$(head -n 1 "$PIDFILE")" > /dev/null 2>&1; then
    echo "postgres is already running with PID $(head -n 1 "$PIDFILE")."
  else
    echo "postgres is not running, but pid file exists. removing pid file.."
    rm -f "$PIDFILE"
    echo "starting postgres server..."
    su-exec postgres pg_ctl start
  fi
else
  echo "starting postgres server..."
  su-exec postgres pg_ctl start
fi

until su-exec postgres psql -U "$POSTGRES_USER" -d postgres -c '\q'; do
  >&2 echo "postgres is unavailable - sleeping"
  sleep 1
done

# Create database
echo "creating database $POSTGRES_DB.."
PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -p 5432 -U $POSTGRES_USER -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$POSTGRES_DB'" | grep -q 1 || (PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -p 5432 -U $POSTGRES_USER -d postgres -c "CREATE DATABASE $POSTGRES_DB" && echo "database $POSTGRES_DB created")
