# Deploy App and MySQL in Kubernetes

This will explain how to deploy the API application with both the application and the MySQL database hosted inside of a Kubernetes cluster. We will use `minikube` to test our Kubernetes configuration and that everything runs correctly, and eventually use GKE (Google Kubernetes Engine) to host a near production ready deployment of the application.

## Minikube

This requires that you have Docker, minikube, and kubectl installed. You can follow the guide in the `../deployment/README.md` file to see how to install these applications on a Debian VM.

### Build the Docker image

1. Clone the `dddstructure` repository to your VM and make it your working directory:

`git clone git@github.com:beeker1121/dddstructure.git`
`cd dddstructure`

OR

Upload to the VM using FTP.

2. Build the Docker image:

`docker build -t dddstructure:v1.0.0 -f cmd/deployments/kubemysql/Dockerfile .`
`docker tag dddstructure:v1.0.0 dddstructure:latest`

This is a little strange, since in the `Dockerfile`, we specify commands such as,

`COPY ../../../ ./`

You would think we should browse to the `cmd/deployments/kubemysql` directory and run the Docker build command there. However, Docker does not have parent navigation with the `COPY` and other commands. The way we get around this is running the `docker build` command from the parent directory of the application, and specifying the Dockerfile we want to use via the `-f` command.

The last command takes the image we just built, which has a `v1.0.0` version tag at the end, and tags it as the `latest` version as well.

### Start minikube

1. Start `minikube` using Docker as the driver:

`minikube start --driver=docker --mount-string="[path-to-dddstructure-directory]:/app/dddstructure" --mount`

For instance, if the path to your local `dddstructure` directory is `/home/my-user/dddstructure`, you would use:

`minikube start --driver=docker --mount-string="/home/my-user/dddstructure:/app/dddstructure" --mount`

We need to mount this directory so our MySQL Deployment that we start down the line will be able to pull in the `init.sql` file. Since we use the `hostPath` option for the `mysql-initdb` persistent volume, it ends up looking within the VM that `minikube` creates, and not our local machine.

We can verify this by SSH'ing into `minikube` and verifying the `/app` directory exists:

`minikube ssh`
`cd /`
`ls`

And we should see the `app` directory listed.

2. View `minikube` status:

`minikube status`

3. View the `kubectl` default configuration, which should now be set to work with `minikube`:

`kubectl config view`

This shows the configuration for the `minikube` cluster, as well as the context for it. Notice `minikube` uses a local certificate authority, as well as a local client certificate and private key for authentication with a user named `minikube`.

### Create kubectl Secrets

In order to connect to the MySQL database, our Go application will need the credentials for the given MySQL user. Instead of storing these details, like the username and password, in plaintext, we can use `kubectl` to store these variables as Secrets. This will then allow our application to use these secrets as environment variables.

While this is fine for this example application, in a production application, you would want to take more precaution on how you store and use Secrets. It's also possible to have Kubernetes user a third-party service for storing and using sensitive information within the cluster.

The secrets configuration file is `dddstructure-secrets.yaml`.

### Run on minikube

Next we will create the Secrets, Deployments, and Services we need to run the application and MySQL on minikube. The configurations for these Kubernetes components are located in the `.yaml` files.

1. Browse to the `cmd/deployments/kubemysql` folder in your terminal.

`cd cmd/deployments/kubemysql`

#### Debugging minikube

We can SSH into `minikube` using the following command:

`minikube ssh`

This will give us a shell into the `minikube` VM. This can be useful for things like checking of the directories are mounted.

#### ConfigMap

First, we'll apply the ConfigMap to the cluster. This is meant for non-sensitive variables such as the database host and port. Our API and MySQL deployments will be able to pull the variables set in the ConfigMap via environment variables.

Apply the ConfigMap:

`kubectl apply -f dddstructure-configmap.yaml`

Which should show:

```sh
configmap/dddstructure-configmap created
```

You should now be able to view the ConfigMap via `kubectl`:

`kubectl get configmaps`

And it should display:

```sh
NAME                     DATA   AGE
dddstructure-configmap   2      14s
kube-root-ca.crt         1      2d
```

#### Secret

Now we'll apply the Secret component to the cluster. This allows both our API and MySQL Deployments to pull certain settings via environment variables, like ConfigMap, but is meant for use with sensitive information such as the database password.

Apply the Secret:

`kubectl apply -f dddstructure-secrets.yaml`

Which should show:

```sh
secret/dddstructure-secrets created
```

You should now be able to view the Secret via `kubectl`:

`kubectl get secrets`

And it should display:

```sh
NAME                   TYPE     DATA   AGE
dddstructure-secrets   Opaque   3      6s
```

#### MySQL deployment

Next we will deploy MySQL.

**NOTE**: Skip this step if you're using Debian, or you just mounted using the `minikube start` command.

1. Mount the `dddstructure/db` directory into `minikube`:

`minikube mount [path-to-db-directory]:/app/dddstructure/db`

For instance, if you had the `dddstructure` application located in `/home/my-user/dddstructure`, you would want to run the command as:

`minikube mount /home/my-user/dddstructure/db:/app/dddstructure/db`

This is needed since when we use `hostPath` within the persisten volume manifest, it's actually looking at the VM that `minikube` creates rather than our local machine/VM.

1. Apply the MySQL initdb persistent volume and persistent volume claim:

`kubectl apply -f mysql-initdb-pv.yaml`

Which should show:

```sh
persistentvolume/mysql-initdb-pv-volume created
persistentvolumeclaim/mysql-initdb-pv-claim created
```

This allows us to mount the `init.sql` file to the `/docker-entrypoint-initdb.d` directory via the persistent volume, which allows the MySQL container to initialize the database.

Now we can see information about the persistent volume:

`kubectl get pv`

Should display:

```sh
NAME                     CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                           STORAGECLASS   VOLUMEATTRIBUTESCLASS   REASON   AGE
mysql-initdb-pv-volume   1Mi        ROX            Retain           Bound    default/mysql-initdb-pv-claim   manual         <unset>                          11s
```

As well as the persistent volume claim:

`kubectl get pvc`

Should display:

```sh
NAME                    STATUS   VOLUME                   CAPACITY   ACCESS MODES   STORAGECLASS   VOLUMEATTRIBUTESCLASS   AGE
mysql-initdb-pv-claim   Bound    mysql-initdb-pv-volume   1Mi        ROX            manual         <unset>                 27s
```

2. Deploy MySQL:

`kubectl apply -f dddstructure-mysql-deployment.yaml`

Which should show:

```sh
deployment.apps/dddstructure-mysql-deployment created
service/dddstructure-mysql-service created
```

View deployments:

`kubectl get deployments`

Which should show:

```sh
NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
dddstructure-mysql-deployment   1/1     1            1           24s
```

Desribe the MySQL deployment:

`kubectl describe deployment dddstructure-mysql-deployment`

Get the pods:

`kubectl get pods -l app=mysql`

Which should show a result like this:

```sh
NAME                                             READY   STATUS    RESTARTS   AGE
dddstructure-mysql-deployment-848b8cdfb5-w9kkk   1/1     Running   0          30s
```

MySQL is now running within the `minikube` cluster, on a single node, on a single Pod that is running a single container (the MySQL image from Docker).

If we run the get Pods command with `-o wide` option, we can see the internal IP of the Pod:

`kubectl get pods -o wide`

Which should display:

```sh
NAME                                             READY   STATUS    RESTARTS   AGE   IP           NODE       NOMINATED NODE   READINESS GATES
dddstructure-mysql-deployment-848b8cdfb5-w9kkk   1/1     Running   0          58s   10.244.0.8   minikube   <none>           <none>
```

Where we can see the internal IP of the MySQL Pod is `10.244.0.8`. Since Pod IPs are ephemeral, this means that there is no guarantee the IP will remain the same if a Pod is recreated or restarted. In order for us to tie the MySQL Pod(s) to a static, internal IP address, we will want to use a Service component - this is what the second half of the `dddstructure-mysql-deployment.yaml` does, is set up a Service component for the MySQL deployment.

3. View the Kubernetes services currently available:

`kubectl get services`

Which should show:

```sh
NAME                         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
dddstructure-mysql-service   ClusterIP   10.109.186.77   <none>        3306/TCP   92s
kubernetes                   ClusterIP   10.96.0.1       <none>        443/TCP    2d
```

Note we can now just use the `10.109.186.77` cluster IP address in this example to route requests to from our API application, which will then be forwarded to the correct Pod(s) for the MySQL deployment. Now if a Pod is restarted or goes down and is recreated, it will still be tied to this static internal IP address.

We can also in fact use the `dddstructure-mysql-service` Service name as the hostname in our application, which is why we have the `db-host` ConfigMap variables set to the Service name. This is because both our API and MySQL applications will be running within the same namespace within Kubernetes, and Kubernetes sets this hostname us as the Service name. If we tried to access another Service under a different namespace, we would have to include that namespace within the hostname.

Also, if we describe the service, we can see that it is linked to our MySQL Pod via the `Endpoints` info:

`kubectl describe service dddstructure-mysql-service`

Which should show:

```sh
Name:                     dddstructure-mysql-service
Namespace:                default
Labels:                   <none>
Annotations:              <none>
Selector:                 app=mysql
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.109.186.77
IPs:                      10.109.186.77
Port:                     <unset>  3306/TCP
TargetPort:               3306/TCP
Endpoints:                10.244.0.8:3306
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
```

And you can see the `Endpoints` points to our MySQL Pod IP of `10.244.0.8` on port `3306`. You can also see the `Selector` references our `app` label for our MySQL Deployment and Pod as well.

##### Debugging MySQL Deployment

At the basic level, here's how we can create a new Pod that will run MySQL client and allow us to interact with our MySQL database in the cluster:

`kubectl run -it --rm --image=mysql:5.7 --restart=Never mysql-client -- mysql -u root -h dddstructure-mysql-service -p`

Once we enter the password, we can use MySQL client as normal.

When we exit the shell (ctrl + C/X), the Pod will automatically be deleted.

Another weird occurence, if we try to use:

```
volumeMounts:
        - name: mysql-initdb
          mountPath: /docker-entrypoint-initdb.d
```

Ie `mountPath: /docker-entrypoint-initdb.d` for the `mountPath`, we will get this error:

```sh
2025-05-05 00:14:49+00:00 [Note] [Entrypoint]: /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/init.sql
ERROR: Can't initialize batch_readline - may be the input source is a directory or a block device.
```

This is fixed by changing the `mountPath` to `/docker-entrypoint-initdb.d/`:

```
volumeMounts:
        - name: mysql-initdb
          mountPath: /docker-entrypoint-initdb.d/
```

#### API Docker image and Deployment

Finally, we can deploy our API application.

##### Build Docker image on minikube

Our API application is going to be an image contained within Docker. In order for our `minikube` cluster to use a local image, we have to make sure that image exists within the Docker engine being used by `minikube`. `minikube` in facts runs in a separate, self-contained VM on your machine, ans also uses its own Docker daemon. This means taht if we build an image locally on our own Docker daemon, it unfortunately won't be available on the `minikube` Docker daemon.

1. Set our local Docker commands to use the `minikube` Docker daemon:

`eval $(minikube docker-env)`

This command will set up our environment variables to point to the `minikube` Docker daemon.

2. Verify we are using the `minikube` Docker daemon:

`docker images`

Which should show something like:

```sh
REPOSITORY                                TAG        IMAGE ID       CREATED         SIZE
mysql                                     8.0        00a697b8380c   2 weeks ago     772MB
registry.k8s.io/kube-apiserver            v1.32.0    c2e17b8d0f4a   4 months ago    97MB
registry.k8s.io/kube-controller-manager   v1.32.0    8cab3d2a8bd0   4 months ago    89.7MB
registry.k8s.io/kube-scheduler            v1.32.0    a389e107f4ff   4 months ago    69.6MB
registry.k8s.io/kube-proxy                v1.32.0    040f9f8aac8c   4 months ago    94MB
registry.k8s.io/etcd                      3.5.16-0   a9e7e6b294ba   7 months ago    150MB
registry.k8s.io/coredns/coredns           v1.11.3    c69fa2e9cbf5   9 months ago    61.8MB
registry.k8s.io/pause                     3.10       873ed7510279   11 months ago   736kB
gcr.io/k8s-minikube/storage-provisioner   v5         6e38f40d628d   4 years ago     31.5MB
```

3. Browse to our application's root folder, `dddstructure`, since we need to run the Dockerfile from here for the relative paths.

4. Follow the [Build the Docker image](#build-the-docker-image) process starting from step 2.

5. Check that image is now available to the `minikube` Docker daemon:

`docker images`

Which should now show:

```sh
dddstructure                              latest     449d56985aae   37 seconds ago   109MB
dddstructure                              v1.0.0     449d56985aae   37 seconds ago   109MB
mysql                                     8.0        00a697b8380c   2 weeks ago      772MB
registry.k8s.io/kube-apiserver            v1.32.0    c2e17b8d0f4a   4 months ago     97MB
registry.k8s.io/kube-controller-manager   v1.32.0    8cab3d2a8bd0   4 months ago     89.7MB
registry.k8s.io/kube-scheduler            v1.32.0    a389e107f4ff   4 months ago     69.6MB
registry.k8s.io/kube-proxy                v1.32.0    040f9f8aac8c   4 months ago     94MB
registry.k8s.io/etcd                      3.5.16-0   a9e7e6b294ba   7 months ago     150MB
registry.k8s.io/coredns/coredns           v1.11.3    c69fa2e9cbf5   9 months ago     61.8MB
registry.k8s.io/pause                     3.10       873ed7510279   11 months ago    736kB
gcr.io/k8s-minikube/storage-provisioner   v5         6e38f40d628d   4 years ago      31.5MB
```

Notice we have the `dddstructure` image now present.

##### API Deployment

1. Apply the API Deployment:

`kubectl apply -f dddstructure-api-deployment.yaml`

Which should show:

```sh
deployment.apps/dddstructure-api-deployment created
```

View deployments:

`kubectl get deployments`

Which should show:

```sh
NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
dddstructure-api-deployment     1/1     1            1           5s
dddstructure-mysql-deployment   1/1     1            1           133m
```

Desribe the API deployment:

`kubectl describe deployment dddstructure-api-deployment`

Get the pods:

`kubectl get pods -l app=dddstructure`

Which should show a result like this:

```sh
NAME                                           READY   STATUS    RESTARTS   AGE
dddstructure-api-deployment-5869cbf996-pn54d   1/1     Running   0          36s
```

2. Let's test a `curl` command:

Get the IP address of the Pod:

`kubectl get pods -o wide`

Which will display:

```sh
NAME                                             READY   STATUS    RESTARTS   AGE     IP            NODE       NOMINATED NODE   READINESS GATES
dddstructure-api-deployment-5869cbf996-78lrw     1/1     Running   0          116s    10.244.0.23   minikube   <none>           <none>
dddstructure-mysql-deployment-848b8cdfb5-frr64   1/1     Running   0          4m37s   10.244.0.22   minikube   <none>           <none>
```

We can see the internal IP of our API Deployment Pod in the cluster is `10.244.0.23` - this is the IP we'll use to issue a `curl` command.

Spin up a new Pod within the cluster that has `curl` installed, and when we exit, the shell will be removed:

`kubectl run -it --rm --image=curlimages/curl curly -- /bin/sh`

Get the IP address of the Pod

Run the signup example request:

```sh
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{
    "email": "test@test.com",
    "password": "TestPassword123"
}' \
http://10.244.0.23:8080/api/v1/signup
```

Notice we're using the `10.244.0.23` IP, and the port the container is listening on of `8080`.

If we set up everything correctly, we should get back a successful response with a JWT:

```sh
~ $ curl -X POST \
>     -H 'Content-Type: application/json' \
>     -d '{
>     "email": "test@test.com",
>     "password": "TestPassword123"
> }' \
> http://10.244.0.23:8080/api/v1/signup
{
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDcwMjEwNTgsImlhdCI6MTc0NjQxNjI1OH0.xlV4eZDc5NDRjttw267aCXsyzZiLtNwzIg0PHNOQDt4"
}
```

###### API External Service

Now that our API application is running within our `minikube` cluster and it's able to communicate our MySQL deployment internally, we want to make our API application available externally to the outside internet. In order to do this, we'll want to create a new Service for our API Deployment.

To break it down a bit, a "service" is defined as the combination of a group of Pod(s), and a policy to access them. A service needs three things - a name (myapp-service), a way to identify the Pod(s) in its group (typically a label like app=myapp), and a way to access those Pod(s) (port 3333 via TCP).

Once a service has been established, Kubernetes assigns it a ClusterIP, which is an internal IP address accessible only within the Kubernetes cluster. Now, other containers in the cluster can access the service through its ClusterIP (or hostname, that resolves to it).

The root Service type in Kubernetes is `ClusterIP`, and this will be the type of the Service if no type is explicitly defined.

There's a few Service types we can use, for example:

**ClusterIP**: This is the default Service type. It creates an internal IP that can be reached from inside the cluster - there is no external access.

**NodePort**: Exposes a Service on a port on each worker node in the cluster, and each node redirects traffic to the given Pod(s) of the Service. Cluster networking handles routing the traffic to the target Service's Pod(s). This is the most basic way to expose a service externally, and does not support anything fancy such as SSL/TLS, load balancing, etc. You would connect to this using the node's IP address and the specified port.

**LoadBalancer**: Exposes a service using a cloud-native Load Balancer. It uses a single IP address of an external load balancer and routes traffic to the Service(s) Pod(s).

**Ingress**: *Not a service type, but a rule to chart external access to services*

Ingress works in conjunction with Services and exposes HTTP and HTTPS routes from outside the cluster to Service(s) within the cluster. Traffic routing is controlled by rules defined on the Ingress resource/component. Ingress may also be configured to give Service(s) externally-reachable URLs, load balance traffic, handle SSL/TLS, etc.

1. Let's create a new `LoadBalancer` Service for our API application:

`kubectl apply -f dddstructure-api-service.yaml`

Get the current services:

`kubectl get services`

Which should show:

```sh
NAME                         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
dddstructure-api-service     LoadBalancer   10.106.197.53   <pending>     8080:30000/TCP   4s
dddstructure-mysql-service   ClusterIP      10.106.44.32    <none>        3306/TCP         20h
kubernetes                   ClusterIP      10.96.0.1       <none>        443/TCP          20h
```

We can now see our API Service is running, with the internal IP on port `8080` and the external IP on port `30000`.

2. Get the external IP.

If we were using a cloud cluster, we could use this command:

`kubectl describe services dddstructure-api-service`

Which would give us the IP listed next to `LoadBalancer Ingress`.

Since we're using `minikube` for this example, we can use the following command:

`minikube service dddstructure-api-service --url`

Which should show:

```sh
http://192.168.49.2:30000
```

We can now access our API application via the IP address and port via our terminal, or other tools like Postman, etc.

###### Debugging API Deployment

We shouldn't run into any issues, but knowing these commands can help debug why a Deployment/Pod isn't starting.

First, view all Pods to get the Pod name:

`kubectl get pods`

We should see Pod name.

We can then get logs specifically for that Pod. This will output anything send to `stdout` and `stderr`:

`kubectl logs [pod-name]`

Get logs by label:

`kubectl logs -l app=dddstructure -f`

Which also tails the logs for all of the Pods/instances with the given label `app=dddstructure`.

The logs though might not show us what we need, or anything at all. We can also use the `describe` command with `kubectl` to get more information about a Pod:

`kubectl describe pod [pod-name]`

This may contain more information on general errors when trying to start the Pod and container.

Also, we could try getting a bash terminal for the Pod:

`kubectl exec --stdin --tty [pod-name] -- /bin/bash`

### Deleting deployment and associated components.

Delete the deployment, service, persistent volume, and persistent volume claim:

```sh
kubectl delete service dddstructure-api-deployment
kubectl delete service dddstructure-api-service
kubectl delete deployment dddstructure-mysql-deployment
kubectl delete service dddstructure-mysql-service
kubectl delete pvc mysql-initdb-pv-claim
kubectl delete pv mysql-initdb-pv-volume
kubectl delete secret dddstructure-secrets
kubectl delete configmap dddstructure-configmap
```

Then if we want, we can also tell the `minikube` cluster to stop:

`minikube stop`

And finally, delete the `minikube` machine:

`minikube delete`