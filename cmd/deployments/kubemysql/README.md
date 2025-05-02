# Deploy App and MySQL in Kubernetes

This will explain how to deploy the application with both the application and the MySQL database hosted inside a Kubernetes cluster. We will use `minikube` to test our Kubernetes configuration and that everything runs correctly, and eventually use GKE (Google Kubernetes Engine) to host a near production ready deployment of the application.

## Minikube

This requires that you have Docker, minikube, and kubectl installed. You can follow the guide in the `../deployment/README.md` file to see how to install these applications on a Debian VM.

### Build the Docker image

1. Clone the `dddstructure` repository to your VM and make it your working directory:

`git clone git@github.com:beeker1121/dddstructure.git`
`cd dddstructure`

2. Build the Docker image:

`docker build -t dddstructure:v1.0.0 -f cmd/deployments/kubemysql/Dockerfile .`
`docker tag dddstructure:v1.0.0 dddstructure:latest`

This is a little strange, since in the `Dockerfile`, we specify commands such as,

`COPY ../../../ ./`

You would think we should browse to the `cmd/deployments/kubemysql` directory and run the Docker build command there. However, Docker does not have parent navigation with the `COPY` and other commands. The way we get around this is running the `docker build` command from the parent directory of the application, and specifying the Dockerfile we want to use via the `-f` command.

The last command takes the image we just built, which has a `v1.0.0` version tag at the end, and tags it as the `latest` version as well.

### Start minikube

Start `minikube` using Docker as the driver:

`minikube start --driver docker`

### Create kubectl Secrets

In order to connect to the MySQL database, our Go application will need the credentials for the given MySQL user. Instead of storing these details, like the username and password, in plaintext, we can use `kubectl` to store these variables as Secrets. This will then allow our application to use these secrets as environment variables.

While this is fine for this example application, in a production application, you would want to take more precaution on how you store and use Secrets. It's also possible to have Kubernetes user a third-party service for storing and using sensitive information within the cluster.

### Run on minikube

Next we will create the Deployment and StatefulSet we need to run the application and MySQL on minikube. The configuration for these Kubernetes components are located in the `dddstructure-depl.yaml` and `dddstructure-mysql.yaml` files.