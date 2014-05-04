## pg2librato

pg2librato periodically queries a Postgres database and emits the results
as metrics to Librato. In this way you can see business metrics from
Postgres side-by-side with your existing system and service metrics in
Librato.

### Usage

Build pg2librato:

```
$ go get
```

Write queries as .sql files in ./queries. The file name of the query
determines the resulting metric name in Librato. Each query must return
either:

* One column, in which case the first row's value gives the metric
  value.

* Two columns, in which case for each row the first column gives the
  metric source and the second column the metric value.
 
For example:

```
$ cat > queries/users.sql <<EOF
select count(*)
from users
EOF

$ cat > queries/multi-source-metric-name.sql <<EOF
select country, count(*)
from users
group by country
EOF
```

Set the database you want to query, the reporting interval, and your
Librato credentials:

```console
$ export DATABASE_URL=...     # postgres://user:pass@host:port/database
$ export QUERY_INTERVAL=...  # time in seconds between queries for each metric
$ export LIBRATO_AUTH=...     # token from your Librato dashboard
```

Run the reporter:

```
./pg2librato
```

Then open up your Librato dashboard to see the results!

### Error Reporting

If you'd like, you can use Rollbar to report errors in your pg2librato
reporter. Configure Rollbar with:

```
$ export ROLLBAR_TOKEN=...
```

If Rollbar config is not provided, you can still see errors in the logs.

### Deploying to Heroku

Since pg2librato uses only one dyno, it's free to run on Heroku. Use the
provided setup-heroku script to set up your reporter app:

```
$ export DATABASE_URL=...     # perhaps `heroku config:get DATABASE_URL -a business-app`
$ export QUERY_INTERVAL=...
$ export APP_NAME=...
$ export ROLLBARP_TOKEN=...   # optional
$ ./setup-heroku
```
