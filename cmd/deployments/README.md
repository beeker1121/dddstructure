# Deployments

This folder holds the various deployment types,

1. Manual - run app via supervisor or some other service, host MySQL on same or separate VM.
2. Kube w/out database - run app via Kubernetes (GKE), host MySQL on separate VM.
3. Kube w/database - run both the app and MySQL via Kubernetes.

# Testing with Minikube

To test the Kubernetes deployments locally, you can use minikube either locally or on another VM.

## Install Docker

Minikube can run on macOS, Windows, and Linux, but there are some different requirements on each platform.

On macOS it runs through the Hypervisor.framework, on Windows it runs through Hyper-V, and Linux it runs either natively (without a virtual machine), or through Docker, or a KVM. A hypervisor is software or a hardware compononent that allows multiple VMs to run on a single physical server, allocating resources like CPU, memory, storage, to the VMs.

**NOTE**: This guide will assume we're using a GCE VM running a Debian image.

Since we're using Debian, we can follow this guide from Docker to install it via the apt repository: https://docs.docker.com/engine/install/debian/

General rundown is as follows:

1. Set up Docker's `apt` repository:

```sh
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
```

2. Install Docker packages:

`sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin`

It's possible we need to add the local user to the `docker` group, to run Docker commands without `sudo`.

3. Add current user to `docker` group:

`sudo usermod -aG docker $USER`

4. Start new shell in current window with group applied for user:

`newgrp docker`

5. Verify installation:

`sudo docker run hello-world`

OR

`docker run --rm busybox date`

This download a test image and creates a container. When the container runs, it prints a confirmation message and exits.

## Install Minikube

Minikube allows us to test our Kubernetes configuration locally.

1. Install minikube:

```sh
curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64
```

2. Verify installation:

`minikube version`

## Install Kubectl

Kubectl is the command-line interface tool to interact with a Kubernetes cluster, whether it's through minikube or GKE or any cluster.

Upon starting minikube (`minikube start`), all `kubectl` commands will go to minikube. You can switch `kubectl` clusters using the `kubectl config get-contexts` and `kubectl config use-context CONTEXT_NAME` commands.

1. Download the kubectl binary:

`curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"`

2. Download sha256 for checksum to validate:

`curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"`

3. Validate checksum:

`echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check`

Which should output:

`kubectl: OK`

4. Install kubectl:

`sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl`

5. Verify installation:

`kubectl version --client`

6. Remove local files:

```sh
rm kubectl
rm kubectl.sha256
```

You can view the default configuration with:

`kubectl config view`

Which should be empty given we just downloaded it and have not started minikube, or created any contexts to a cluster.

Once we run `minikube start --driver docker`, `minikube` will modify the default configuration to include a context we can use for it. This allows us to switch between `minikube` and a remote cluster. Review documentation for `$HOME/.kube/config` and contexts for more information. Also look into `KUBECONFIG` environment variable.

Can view the config names via `grep 'name:' ~/.kube/config` (should show `minikube` after starting it), and for example to switch to it `minikube` context you would do `kubectl config use-context minikube`. Now all `kubectl` commands will run on the `minikube` cluster.

View for more information: https://cloud.google.com/kubernetes-engine/multi-cloud/docs/aws/how-to/configure-cluster-access-for-kubectl and https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/ and https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl

## Gcloud CLI

In this project example, we're using GCP (Google Cloud Platform).

To interact with GCP products and services, we two options - the web based GCP UI, and the `gcloud` CLI.

### Web based UI

Using this option, GCP will create a new VM on our project with `gcloud` installed.

### Local CLI

Install `gcloud` CLI using this guide: https://cloud.google.com/sdk/docs/install#linux

For basic info on the `gcloud` CLI including authentication, can view this Medium article: https://medium.com/@danilo.drobac/become-a-gcp-master-get-comfortable-with-the-gcloud-cli-ae99a0b629e4

After installation, can run `gcloud init` to set some basic configuration such as project ID, and then will open a browser session to authenticate.

To switch the account `gcloud` points to, run `gcloud auth login` to open a new browser window.

If you'd rather authorize without a web browser but still interact with the command line, use the `--no-browser` flag. To authorize without a web browser and non-interactively, create a service account with the appropriate scopes using the Google Cloud Console and use gcloud auth activate-service-account with the corresponding JSON key file.

To view authenticated accounts, run `gcloud auth list`.

To change the currently active `gcloud` account, run `gcloud config set account [account]`

Another option is to create custom `gcloud` configurations that we can switch between. For instance, can run `gcloud config configurations create testconfig`, then `gcloud init` since it's a new config. Then we can switch the active account using `gcloud config configurations activate [config name]`. For more information, view this post: https://stackoverflow.com/a/52728878

## Set up Google Artifact Repository

In order to easily pull our application container image into Kubernetes (either using minikube or GKE) on Google Cloud, we can use their artifact repository service. This service is basically a package and image manager within GCP. Docker Hub or any other image repository service can of course be used in replacement of this at any time.

### Local images

Local images can also be used, although this is a much more complex topic.

NOTE: Post on how to download OCI (Open Container Initiative) images locally: https://www.reddit.com/r/docker/comments/1f845ug/guide_to_manual_downloading_docker_images_wget/

To use a local image, look into using the `docker image save` command (https://docs.docker.com/reference/cli/docker/image/save/)

`docker image save busybox > busybox.tar`

Once uploaded to the VM, we can load it in Docker using:

`docker image load -i busybox.tar`

The image should now be listed in available images:

`docker images`

If you're using `minikube` and Docker locally and not on a remote VM, then these images can be used directly after building locally without saving/loading.

If we're using GKE to run this, then using local images is more work than what it's worth - just enable and use the Google Artifacts Registry or Docker Hub or some other image repository service.

### Build images locally on VM

If running `minikube` on a remote VM, the above can be sort of a pain to get a Docker image on the VM.

Another option is to just transfer the source of the application with the Dockerfile and build the image on the VM itself. This way, so long as Docker is installed on the VM, we don't have to build locally and then transfer an image.