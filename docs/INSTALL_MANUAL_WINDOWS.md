# Manual Windows Install

This guide is specifically for people who are on a Windows machine using [WSL](https://learn.microsoft.com/en-us/windows/wsl/about) with Docker.

Scrutiny is made up of three components: an influxdb Database, a collector and a webapp/api. Docker will be used for
the influxdb and webapp/API, the collector component will be facilitated by [Windows Task Scheduler](https://learn.microsoft.com/en-us/windows/win32/taskschd/task-scheduler-start-page).

> **NOTE:** If you are **NOT** using WSL with docker, then the easiest way to get started with [Scrutiny is the omnibus Docker image](https://github.com/AnalogJ/scrutiny#docker).

## InfluxDB and Webapp/API (Docker)

1. Copy the [example.hubspoke.docker-compose.yml](https://github.com/AnalogJ/scrutiny/blob/master/docker/example.hubspoke.docker-compose.yml) 
file and delete the collector section near the bottom of the file.
2. Run `docker-compose up -d` to verify that the DB and webapp are working correctly and once its completed, your webapp
should be up and running but the dashboard will be empty (default location is `localhost:8080`)

## Collector (Windows Task Scheduler)

1. Download the latest `scrutiny-collector-metrics-windows-amd64.exe` from the [releases page](https://github.com/AnalogJ/scrutiny/releases) (under assets)
2. On your windows host, open [Windows Task Scheduler](https://www.wikihow.com/Open-Task-Scheduler-in-Windows-10) as **Administrator**
   1. In the **Start Menu** (Windows key), type `Task Scheduler` and then right click `Run as Administrator` to open
3. On the status bar (under the `action` tab), click `Create Task...`
4. A new window should open with the `General` Tab open, enter relevant information into the `Name` and `Description` fields
   1. Under **Security Options** check:
      1. **Run whether user is logged on or not**
      2. **Run with highest privileges**
5. Next, click the `Triggers` tab and then click `New...` (bottom left-hand side of the window)
6. Here you can set how often you want this task to run, example settings are the following:
   1. **Settings:**
      1. `Daily`, start at `TODAYS_DATE` `12:00:00 AM`, Recur every `1` days,
   2. **Advanced Settings:**
      1. Repeat Task every: `1 hour` for a duration of `Indefinitely`
      2. Stop task if it runs longer than: `30 minutes`
   3. Click Ok when satisfied with your schedule
   > **NOTE:** The above settings will trigger the task **every day at midnight** and then **run every hour after that** (modify as needed)
7. Next, click the `Actions` tab and then click `New...` (bottom left-hand side of the window)
   1. **Action Settings:**
      1. In the **Program/Script** field, put: `scrutiny-collector-metrics-windows-amd64.exe`
      2. In the **Add arguments (optional)** field, put: `run --api-endpoint "http://localhost:8080" --config collector.yaml`
         > **NOTE:** 
         >  * Make sure that you put the correct port number (as specified in the docker-compose file) for the webapp (default is `8080`)
         > * The `--config` param is optional and is not needed if you just want to use the default collector config, see [example.collector.yaml](https://github.com/AnalogJ/scrutiny/blob/master/example.collector.yaml) for more info on the collector config.
      3. In the **Start in (optional)** field, put: FOLDER_PATH_TO_YOUR `scrutiny-collector-metrics-windows-amd64.exe` file
          > **NOTE:** Must be exact and do not include `scrutiny-collector-metrics-windows-amd64.exe` in the path
      4. Click Ok when finished
8. Next, click the `Conditions` tab and make sure that everything is unchecked (unless you want to specify otherwise)
9. Next, click the `Settings` tab and check everything except for the last checkbox
   1. **Examples for the following settings:**
      1. If the task fails, restart every: `5 minutes`
      2. Attempt restart up to: `3` times
      3. Stop the task if it runs longer than `1 hour`
10. Next, once satisfied with everything, click Ok
11. Then, find your newly created task (by its name) in the scheduler task list and then manually run it (right click it and then click `Run`)
12. Finally, refresh your dashboard after a minute or two and your drive information should have populated the webapp dashboard.




