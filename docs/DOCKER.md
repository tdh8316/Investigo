# operative-framework
[![Go Report Card](https://goreportcard.com/badge/github.com/tdh8316/Investigo)](https://goreportcard.com/report/github.com/tdh8316/Investigo) [![GoDoc](https://godoc.org/github.com/tdh8316/Investigo?status.svg)](http://godoc.org/github.com/tdh8316/Investigo) [![GitHub release](https://img.shields.io/github/release/tdh8316/Investigo.svg)](https://github.com/tdh8316/Investigo/releases/latest) [![LICENSE](https://img.shields.io/github/license/tdh8316/Investigo.svg)](https://github.com/tdh8316/Investigo/blob/master/LICENSE)

## Installing

### Running as a Docker Container

#### Pre-requisite
You can run operative-framework in a Docker container to avoid installing Golang locally. To install Docker check out the [official Docker documentation](https://docs.docker.com/engine/getstarted/step_one/#step-1-get-docker).

#### Pull pre-built image
```
docker pull tdh8316/Investigo
```

#### Start a container
Once you have docker installed you can run operative-framework:

    $ docker run -ti --rm tdh8316/Investigo

If you are running this command often you will probably want to define an alias:

    $ alias investigo="docker run -ti --rm tdh8316/Investigo"

To build the Docker image from sources:

    $ git clone https://github.com/tdh8316/Investigo.git
    $ cd Investigo
    $ docker build -t investigo . (or make docker)

### Running as a Docker-Compose container
```
docker-compose run investigo
```
