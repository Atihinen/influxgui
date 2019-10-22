# InfluxGUI


Standalone GUI for influxDB for debugging and verify Influx Query


# Installation

Download correct binary for your OS from releases.

## macOS

You need to give permissions to the binary with `chmod u+x` before running the application.

# Usage

After download, douple click the binary to open the GUI.

Connect to database by inputting the hostname, port and credentials (if there are none set, leave empty). Then press "Connect".

Then run queries against correct database.


## Example queries

`SHOW MEASUREMENTS`

`SELECT * FROM <measurement> LIMIT 1;`



