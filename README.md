[![Testing and building the go code.](https://github.com/Pik-9/fixadated/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Pik-9/fixadated/actions/workflows/go.yml)
# fixadated
`fixadated` is the backend for the [Caerus](https://github.com/Pik-9/caerus) (Not publicly
available yet) online tool to organize meetups and find dates.
This is meant to be kept very simple - rather a proof of concept.

`fixadated` can embed an entire static web frontend and serve it under `/` thus eliminating the
need for a separate hosting plan. Again: Suitable for a small user base; if you need to scale up
you will probably want a different solution.

## Building
You need to have Go (>= 1.19) installed to build.

First you need to put your static frontend in `res/webapp`. It is imparative that you put at least
one file here. **The build will fail otherwise**.

Then:
```bash
go build
```

## Operating
Once you start the daemon it will listen on `localhost:8080` for incoming http requests.
TLS/SSL is not implemented since this daemon is meant to be run behind a proxy that handles
this kind of stuff.

On `/` a file server will run that serves the files from `res/webapp`.
On `/api` a very barebones REST API can be found that handles the important stuff.


## REST API
The following endpoints are implemented so far:

Method |            Route             |    Expected JSON Request   | Answer
-------|------------------------------|----------------------------|-----------------------
GET    | /api/event/:eventid          | _None_                     | {name, description, dates, participants}
POST   | /api/event                   | {name, description, dates} | {name, description, dates, id, editId, participants}
PATCH  | /api/event/:editeventid      | {name, description}        | {name, description, dates, id, editId, participants
POST   | /api/event/:eventid/register | {name, declarations}       | {name, id, editId, declarations}
PATCH  | /api/participation/:partid   | {name, declarations}       | {name, id, editId, declarations}
