# Loverly : A Dating Apps

## Prerequsites

- Go v1.22
- PostgreSQL
- Redis
- Docker & Docker Compose (for dev)

## How to Prepare the Environment

Before you can start `loverly` up, you need to prepare the environment. First one you should set up is the configuration file. To create the configuration file, you need to create `.env` file by running this commands:

```shell
cp .env.example .env
```

Next, you need to start the database & redis for `loverly`. You can use the docker-compose prepared inside the `etc` directory. To start the docker-compose, run:

```shell
cd etc
docker-compose up
```

or use `-d` flag to run docker-compose in the background:

```shell
docker-compose up -d
```

Once you have verified that the database & redis are started, update the `.env` configuration file with the correct database and Redis credentials. Then, run the schema and seed migration by:

```shell
make migrate.up
```

And don't forget to run `go mod tidy` after cloning this repo to install all the necessary dependencies.

## How to Run the Application

Start the application by running:
  
```shell
make run
```

Run unit test :
```shell
make test
```

And that's it! `loverly` is already running on your machine. To access it try these:

- `GET:     http://localhost:3003/health` -> for check service running or not
- `POST:    http://localhost:3003/v1/register` -> for registering new users
- `POST:    http://localhost:3003/v1/login` -> for login using your credentials. use `handsome@gmail.com`, password `password` for demo.

- `GET:     http://localhost:3003/v1/discovery` -> for get list profile for dating
- `POST:    http://localhost:3003/v1/swipe` -> for like (right) or pass (left)
- `GET:     http://localhost:3003/v1/match` -> for list of profile match with you

- `GET:     http://localhost:3003/v1/profile` -> for get detail profile
- `GET:     http://localhost:3003/v1/subscription` -> for get detail subscription plan you have
- `POST:    http://localhost:3003/v1/subscription` -> for subscribe a package plan

Or, import the collection JSON (`loverly.json`) into Postman for easy endpoint testing.