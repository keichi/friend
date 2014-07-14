# friend-client command line specs

`friend signup`

Create account on server.

`friend login`

Login to server.

`friend <user>`

Reads stdin and send to `user`.

`friend <user> <filename>`

Send `filename` to `user`.

`friend`

Wait for incoming file and write out to stdout.

`friend -o <filename>`

Wait for incoming file and save to `filename`.

`friend -d`

Launch as daemon and automatically download files in background.
