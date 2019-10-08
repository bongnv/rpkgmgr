# rpkgmgr

`rpkgmgr` helps to index R packages from online sources. It parses PACKAGES in each repository then parse DESCRIPTION in each package to populate details into a database.

## How to use

In order to use, `docker-compose` and `docker` are required. Following linkg to install them:

- https://docs.docker.com/compose/install/
- https://www.docker.com/get-started

Before start, we need to migrate schema:

    $docker-compose up migrate
    
Then, we can start the service:

    $ docker-compose up go
    
It will start the job to index packages in `http://cran.r-project.org/src/contrib/` as well as to set a scheduled job at 12PM.
