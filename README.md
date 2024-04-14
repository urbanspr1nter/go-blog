# go-blog

![Example](./docs/images/example.png)

## What is this?
**go-blog** is a standalone executable that hosts a customizable blog!
Configure some settings, and run the binary, and you will have a blog which people can visit!

## Requirements
1. You will need at least `go` **1.20**. This is the version I was able to test it.
2. Access to a MySQL database (connection string)

## Build and Run
To build the program, you can just run the included shell script, `build.sh`.

```shell
./build.sh
```

After that, the `server` executable will be emitted. This is the binary to run the service.

Ensure that the `configuration.json` has the correct settings, and just run the server like this:

```shell
./server
```

The blog is **live**! 🎉🥳

## Structure
The blog is designed to highly customizable. If you want to change the static content to fit your tastes, here are some good places to start:

* All static files are within the `public` directory. The common use case is to place static assets such as JavaScript, and CSS here.
* All template files are in the `views` directory. Modify these to customize your pages.

## Schema
You'd minimally need a `post` table with the following MySQL schema:

```mysql
CREATE TABLE post (
    id PRIMARY KEY AUTO_INCREMENT, 
    slug TEXT,
    title TEXT,
    content TEXT,
    created DATETIME,
    published BOOLEAN
);
```

## configuration.json

It's relatively straightforward to configure the `configuration.json` file. The 2 main things to make the app work are:

* `database_url_env` - The environment variable to reference so that the connection string to the database. So for example, common setting value would be `DATABASE_URL`, which would reference the `env` variable called `DATABASE_URL` within the system to find the DB connection string.
* `http_port` = The port to listen to.

## Building the Connection String
The connection string is in the format:
```text
username:password@tcp(hostname:port)/database_name
```
So you'd have an environment variable like this:
```shell
export DATABASE_URL="username:password@tcp(hostname:port)/database_name"
```

And in `configuration.json`, you can set it like this:

```json
{
  "database_url_env": "DATABASE_URL"
}
```

## Other
This project will be updated from time to time. I use it for real things. 
You can contact me at [urbanspr1nter@gmail.com](mailto:urbanspr1nter@gmail.com) too!