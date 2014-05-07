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

Write queries as .sql files in ./queries. Each query must return result
sets with exactly three columns, in the following order:

* A string giving the name of the metric.

* A string giving the source of the metric, or null if the source is
  unassigned.

* A numeric giving the value for the corresponding {name, source} pair.
 
For example:

```
$ cat > queries/test-without-source.sql <<EOF
select
  'users.total', null, count(*)
from
  users
EOF

$ cat > queries/test-with-source.sql <<EOF
select
  'users.per-country', country, count(*)
from
  users
group by
  country
EOF
```

Set the database you want to query, the reporting interval, and your
Librato credentials:

```console
$ export DATABASE_URL=...    # postgres://user:pass@host:port/database
$ export QUERY_INTERVAL=...  # time in seconds between queries for each metric
$ export LIBRATO_AUTH=...    # email:token
```

Run the reporter:

```
$ pg2librato
```

Then open up your Librato dashboard to see the results!

### Error Reporting

If you'd like, you can use Rollbar to report errors in your pg2librato
reporter. Configure Rollbar with:

```
$ export ROLLBAR_ACCESS_TOKEN=...  # provided by Rollbar service
$ export ROLLBAR_ENVIRONMENT=      # e.g. development or production
```

If Rollbar config is not provided, you can still see errors in the logs.

### Deploying to Heroku

Since pg2librato uses only one dyno, it's free to run on Heroku. Deploy
with something like the following:

```
$ heroku create
$ heroku config:set DATABASE_URL=...         # perhaps `heroku config:get DATABASE_URL -a business-app`
$ heroku config:set QUERY_INTERVAL=...
$ heroku addons:add rollbar                  # optional
$ heroku config:set ROLLBAR_ENVIRONMENT=...
$ git push heroku master
$ heroku scale pg2librato=1
```
