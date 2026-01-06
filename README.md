# Passport
Passport is an extremely simple self hosted password manager, that consists of a backend/server, a TUI, a CLI(WIP) and an android application(TBD).

# Getting started
If it is the first time you are using `passport`, you will need to sign a SSL certificate, for the HTTPS server to run. The
easiest way to do this is by running this command while in the `backend` directory.
```bash
$ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes
```
This certificate will however expire after one year so, if you are planing to use passport for a longer period of time, be prepared to renew it using the same command in one year (or modify the value 
passed to the `-days` flag to however long you want it to be valid).

 ## Setting up the server
 Setting up the server is a very simple process. The first time you use `passport` it requires only a few steps. If it's not your first time, it is as simple as just compiling, and running the server.
 ### Compilation
 `passport` comes with a server Makefile for an easier setup process. To compile the server, navigate to the `backend` directory and simply run:
 ```bash
 $ make build
 ```
This will compile the server and place the binary in `target/bin/passport-server`. It will also create a log file (`target/log/file.log`) and `touch` two `.pem` files for the RSA keys. 
NOTE! The log file will be overwritten if the code is recomplied. The `.pem` files will however not. 

For all upcomming steps on how to set up the server, you need to be located in the `backend` directory for them too work as they should.

If it is not your first time using `passport`, you can now skip ahead to `running`

### Setting up RSA
Once the code is compiled, it is now time to prepare the RSA keys to be ready for encryption of passwords. This is done, once again while in the `backend` directory, by running this command:
```bash
$ ./target/bin/passport-server keygen
```
The paths for the keys are relative, which makes your location in your file system very important for `passport` to function.

### Migration of database
It is now time to set up the database correctly. As of now, the migration adds some dummy data for debugging purposes. If you want to preload it wiht other data, 
or none at all, edit the [migrate.go](./backend/db/migrate.go) file and add your own data. If you don't want any data in the database to begin with, you can remove `line 35` (or comment it out)
and then recompile the code, as done in `Compilation` (this does not require you to redo step two `Setting up RSA`). 

To run the migration, simply run this command:
```bash
$ ./target/bin/passport-server migrate
```
It is once again important that you are located in the `backend` directory

### Running
Running the server requires a `superuser` in order to open the port `:443` (HTTPS). To start the server, simply run:
```bash
$ sudo ./target/bin/passport-server
```
and enter your sudo password. If everything was done correctly, you should now have the `passport` server up and runnning, ready for use.
