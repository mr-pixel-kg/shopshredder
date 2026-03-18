# mpXsandbox

This application allows you to create demo shops in a docker environment.


Start application with docker compose:
```
docker-compose up --build
```


## Configuration Guide

This document provides an overview of the configuration settings for the application.  
The configuration is loaded from a `config.yml` file and can be overridden using environment variables.

### Configuration File (`config.yml`)

The application configuration is structured as follows:

```yaml
server:
  port: 8080
  allowed_origins:
    - "http://localhost:5173"
    - "https://www.shopshredder.de"
    - "http://localhost:8000"
    - "http://localhost:80"
  app_url: "http://localhost"

sandbox:
  url_prefix: "sandbox-"
  url_suffix: ".shopshredder.de"
  default_lifetime: 60

auth:
  username: "admin"
  password: "password"

guard:
  max_total_sandboxes: 32
  max_sandboxes_per_ip: 5
  max_sandbox_lifetime: 60

#database:
#  host: "localhost"
#  port: 5432
#  user: "postgres"
#  password: "password"
#  name: "appdb"
```


### Configuration Overview

| Section   | 	Key                  | 	Type    | 	Default            | 	Description                                            |
|-----------|-----------------------|----------|---------------------|---------------------------------------------------------|
| server    | 	port                 | 	int	    | 8080                | 	The port on which the server runs.                     |
|           | allowed_origins       | 	array   | 	[]                 | 	List of allowed CORS origins.                          |
|           | app_url               | 	string  | 	"http://localhost" | 	The base URL of the application.                       |
| sandbox   | 	url_prefix           | 	string  | 	"sandbox-"         | 	Prefix for sandbox URLs.                               |
|           | url_suffix            | 	string  | 	".shopshredder.de" | 	Suffix for sandbox URLs.                               |
|           | default_lifetime      | 	int	    | 60                  | 	Default sandbox lifetime in minutes.                   |
| auth      | 	username             | 	string  | 	"admin"            | 	Username for basic authentication.                     |
|           | password              | 	string	 | "password"          | 	Password for basic authentication.                     |
| guard     | 	max_total_sandboxes  | 	int	    | 32                  | 	Maximum number of sandboxes allowed in total.          |
|           | max_sandboxes_per_ip  | 	int	    | 5                   | 	Maximum number of concurrent sandboxes per IP address. |
|           | max_sandbox_lifetime	 | int	     | 60                  | 	Maximum sandbox lifetime in minutes.                   |
| database* | 	host                 | 	string  | 	"localhost"        | 	Database host address.                                 |
|           | port                  | 	int     | 	5432               | 	Database port.                                         |
|           | user                  | 	string  | 	"postgres"         | 	Database username.                                     |
|           | password              | 	string  | 	"password"         | 	Database password.                                     |
|           | name                  | 	string  | 	"appdb"            | 	Database name.                                         |
| sentry**  | dsn                   | 	string  | 	                   | 	Sentry DSN.                                            |

*Database configuration is currently disabled. Uncomment it in config.yml to enable.
If no database is specified a SQLite database is used automatically.
**Sentry configuration is currently disabled. Uncomment it in config.yml to enable.

### Overriding Configuration with Environment Variables
The application supports environment variables to override configuration values.
Here is a list of environment variables and their corresponding settings:

| Environment Variable        | Corresponding Config Key   |
|-----------------------------|----------------------------|
| SERVER_PORT	                | server.port                |
| SANDBOX_URL_PREFIX          | sandbox.url_prefix         |
| SANDBOX_URL_SUFFIX          | sandbox.url_suffix         |
| SANDBOX_DEFAULT_LIFETIME    | sandbox.default_lifetime   |
| AUTH_USERNAME               | auth.username              |
| AUTH_PASSWORD               | auth.password              |
| GUARD_MAX_TOTAL_SANDBOXES   | guard.max_total_sandboxes  |
| GUARD_MAX_SANDBOXES_PER_IP	 | guard.max_sandboxes_per_ip |
| GUARD_MAX_SANDBOX_LIFETIME	 | guard.max_sandbox_lifetime |
| SENTRY_DSN                  | sentry.dsn                 |


## Development

This project uses automated CI/CD pipeline and pushes production ready images to Github Packages.
The configuration and build options can be found at ./github/workflows/release.yml

### Backend

#### Change Docker Host
If you have multiple users on your system, you maybe have to specify your docker host and run the application as root:
```
DOCKER_HOST=unix:///Users/<username>/.docker/run/docker.sock
```


#### Swagger Docs
Command to compile swagger documentation page under http://localhost:8080/swagger/index.html
```
swag init
```


### Frontend
Start development server
```
yarn dev
```

The `VITE_BACKEND_URL` can be set in the .env files. For production use, the pipeline will set the correct value by using docker build arguments.

### Custom Shopware Image Build

In order to provide a Shopware sandbox image, you can not use the default dockware image because it has missing configuration options behind proxies.
To avoid getting Mixed-Content and HSTS problems in your browser, you have to extend the dockware image and build a custom sandbox image.
Edit the Dockerfile in docker/images/.../Dockerfile and change the Shopware version in the FROM statement.
In addition you can add your custom configuration files to the image by copying them into the image.


Then, you can build your new image.
```bash
docker build . -t mr-pixel/sw-sandbox:6.7.0.0-rc1
```

After that, go to the administration page in the application and add this new image that you just created.
Now we have a plain shopware image that can be used for demo and testing purposes.

#### Custom Configuration

However, in order to create demonstration images with custom configurations, you have to first start a new sandbox and choose your base image.
Then you can open your sandbox in the browser and install plugins, custom themes and configure the entire store as you like.

After that you can create a new snapshot of the running sandbox container. Enter your image name and tag.
Congrats, you have created a custom Shopware demo sandbox image. Don't forget to stop the running sandbox container.


