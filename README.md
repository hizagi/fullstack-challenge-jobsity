# fullstack-challenge-jobsity
A tasklist app that allows the user to manage a list of tasks and the status of them.

The project is divided in two folders:

`/frontend` - contains an angular app called "tasklist", that runs the UI of the application

`/backend` - contains a golang app, is responsible for running an http server that connects to a mongo database

In order to run the application please run the following command:

- `make up`

To remove the containers, run this one:

- `make down`