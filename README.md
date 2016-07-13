# do-all
A CloudFoundry CLI Plugin that runs a given command on all apps in a given space

## Install

Use the CF Community Plugin repo at https://plugins.cloudfoundry.org

`cf install-plugin do-all -r "CF-Community"`

## Run

Run `cf do-all COMMAND ARGS` and it will run the command for all apps in the current space. 
Where the app name would go in the arguments, use `{}`. For example, if we have this:

```
$ cf apps
Getting apps in org jghiloni / space development as admin...
OK

name                       requested state   instances   memory   disk   urls
app1                       started           1/1         1G       1G     app1.apps.labb.jghiloni.io
app2                       started           1/1         128M     1G     app2.apps.labb.jghiloni.io
app3                       started           3/3         512M     1G     app3.apps.labb.jghiloni.io
app4                       started           2/2         256M     128M   app4.apps.labb.jghiloni.io
```

Running `cf do-all set-env {} FOO bar` would set the environment variable `FOO` equal to `bar` on `app1`, `app2`, `app3`,
and `app4`.
