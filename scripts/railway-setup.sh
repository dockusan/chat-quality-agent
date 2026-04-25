#!/usr/bin/env bash
set -euo pipefail

APP_SERVICE="${APP_SERVICE:-chat-agent}"
DB_SERVICE="${DB_SERVICE:-mysql}"
ENVIRONMENT="${RAILWAY_ENVIRONMENT:-production}"
VOLUME_PATH="${VOLUME_PATH:-/data}"
DOMAIN_PORT="${DOMAIN_PORT:-8080}"
DEPLOY="${DEPLOY:-true}"
GENERATE_DOMAIN="${GENERATE_DOMAIN:-true}"
ADD_VOLUME="${ADD_VOLUME:-true}"

require_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "Error: missing required command: $1" >&2
        exit 1
    fi
}

random_hex() {
    openssl rand -hex "$1"
}

ensure_linked_project() {
    if ! railway status >/dev/null 2>&1; then
        echo "No Railway project is linked to this directory."
        echo "Run one of these first:"
        echo "  railway login"
        echo "  railway link"
        echo ""
        echo "Or create/link a project from Railway dashboard, then run this script again."
        exit 1
    fi
}

ensure_service() {
    local service="$1"
    local add_args=("${@:2}")

    if railway service link "$service" >/dev/null 2>&1; then
        echo "Service exists: $service"
        return
    fi

    echo "Creating service: $service"
    railway add "${add_args[@]}"
}

set_app_variables() {
    local jwt_secret
    local encryption_key

    jwt_secret="$(random_hex 32)"
    encryption_key="$(random_hex 16)"

    echo "Setting app variables on service: $APP_SERVICE"
    railway variable set \
        --service "$APP_SERVICE" \
        --environment "$ENVIRONMENT" \
        --skip-deploys \
        "APP_ENV=production" \
        "JWT_SECRET=$jwt_secret" \
        "ENCRYPTION_KEY=$encryption_key" \
        "TZ=Asia/Ho_Chi_Minh" \
        "MYSQLHOST=\${{${DB_SERVICE}.MYSQLHOST}}" \
        "MYSQLPORT=\${{${DB_SERVICE}.MYSQLPORT}}" \
        "MYSQLUSER=\${{${DB_SERVICE}.MYSQLUSER}}" \
        "MYSQLPASSWORD=\${{${DB_SERVICE}.MYSQLPASSWORD}}" \
        "MYSQLDATABASE=\${{${DB_SERVICE}.MYSQLDATABASE}}" \
        "DATABASE_URL=\${{${DB_SERVICE}.MYSQL_URL}}"
}

add_volume() {
    if [ "$ADD_VOLUME" != "true" ]; then
        return
    fi

    echo "Adding volume to service $APP_SERVICE at $VOLUME_PATH"
    railway service link "$APP_SERVICE" >/dev/null
    railway volume add --mount-path "$VOLUME_PATH" || true
}

generate_domain() {
    if [ "$GENERATE_DOMAIN" != "true" ]; then
        return
    fi

    echo "Generating Railway domain for service: $APP_SERVICE"
    railway domain --service "$APP_SERVICE" --port "$DOMAIN_PORT" || true
}

deploy_app() {
    if [ "$DEPLOY" != "true" ]; then
        return
    fi

    echo "Deploying app service: $APP_SERVICE"
    railway up --service "$APP_SERVICE" --environment "$ENVIRONMENT" --detach
}

main() {
    require_cmd railway
    require_cmd openssl

    ensure_linked_project

    echo "Using Railway environment: $ENVIRONMENT"
    echo "App service: $APP_SERVICE"
    echo "MySQL service: $DB_SERVICE"

    ensure_service "$DB_SERVICE" --database mysql --service "$DB_SERVICE"
    ensure_service "$APP_SERVICE" --service "$APP_SERVICE"

    set_app_variables
    add_volume
    generate_domain
    deploy_app

    echo ""
    echo "Railway setup complete."
    echo "Check logs with:"
    echo "  railway logs --service $APP_SERVICE"
    echo ""
    echo "If deployment still reports DB_PASSWORD is required, verify the MySQL service name is exactly:"
    echo "  $DB_SERVICE"
}

main "$@"
