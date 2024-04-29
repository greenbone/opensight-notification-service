# only intended to be used for the local setup

# values for which we don't have a default
export DB_USERNAME=postgres
export DB_NAME=notification_service
export DB_SSL_MODE=disable
if [[ -z "$DB_PASSWORD" ]]; then # this is a secret, set this env var by other means
    echo "warning: database password not set, but is required" > /dev/stderr; return 1
else
    export DB_PASSWORD 
fi

# log level for convenience
export LOG_LEVEL=debug
