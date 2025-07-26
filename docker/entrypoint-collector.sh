#!/bin/bash

# Cron runs in its own isolated environment (usually using only /etc/environment )
# So when the container starts up, we will do a dump of the runtime environment into a .env file that we
# will then source into the crontab file (/etc/cron.d/scrutiny)
(set -o posix; export -p) > /env.sh

# adding ability to customize the cron schedules.
COLLECTOR_CRON_SCHEDULE=${COLLECTOR_CRON_SCHEDULE:-"0 0 * * *"}
COLLECTOR_ZFS_CRON_SCHEDULE=${COLLECTOR_ZFS_CRON_SCHEDULE:-"*/5 * * * *"}
COLLECTOR_RUN_STARTUP=${COLLECTOR_RUN_STARTUP:-"false"}
COLLECTOR_RUN_STARTUP_SLEEP=${COLLECTOR_RUN_STARTUP_SLEEP:-"1"}

# if the cron schedules have been overridden via env variable (eg docker-compose) we should make sure to strip quotes
[[ "${COLLECTOR_CRON_SCHEDULE}" == \"*\" || "${COLLECTOR_CRON_SCHEDULE}" == \'*\' ]] && COLLECTOR_CRON_SCHEDULE="${COLLECTOR_CRON_SCHEDULE:1:-1}"
[[ "${COLLECTOR_ZFS_CRON_SCHEDULE}" == \"*\" || "${COLLECTOR_ZFS_CRON_SCHEDULE}" == \'*\' ]] && COLLECTOR_ZFS_CRON_SCHEDULE="${COLLECTOR_ZFS_CRON_SCHEDULE:1:-1}"

# replace placeholders with correct values
sed -i 's|{COLLECTOR_CRON_SCHEDULE}|'"${COLLECTOR_CRON_SCHEDULE}"'|g' /etc/cron.d/scrutiny
sed -i 's|{COLLECTOR_ZFS_CRON_SCHEDULE}|'"${COLLECTOR_ZFS_CRON_SCHEDULE}"'|g' /etc/cron.d/scrutiny

if [[ "${COLLECTOR_RUN_STARTUP}" == "true" ]]; then
    sleep ${COLLECTOR_RUN_STARTUP_SLEEP}
    echo "starting scrutiny collector (run-once mode. subsequent calls will be triggered via cron service)"
    /opt/scrutiny/bin/scrutiny-collector-metrics run
fi


# now that we have the env start cron in the foreground
echo "starting cron"
exec su -c "cron -f -L 15" root
