#!/bin/bash
# Cron runs in its own isolated environment (usually using only /etc/environment )
# So when the container starts up, we will do a dump of the runtime environment into a .env file that we
# will then source into the crontab file (/etc/cron.d/scrutiny.sh)

printenv | sed 's/^\(.*\)$/export \1/g' > /env.sh

# adding ability to customize the cron schedule.
COLLECTOR_CRON_SCHEDULE=${COLLECTOR_CRON_SCHEDULE:-"0 0 * * *"}
sed -i 's|{COLLECTOR_CRON_SCHEDULE}|'"${COLLECTOR_CRON_SCHEDULE}"'|g' /etc/cron.d/scrutiny

# now that we have the env start cron in the foreground
echo "starting cron"
su -c "cron -l 8 -f" root
