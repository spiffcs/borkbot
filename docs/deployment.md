# DEPLOYMENT BORKBOT GCP

## APP

### - KubeDeployment -

#### Create Container Cluster of 2 of the smallest google machine:

```bash
gcloud container clusters create borkbot-prod --machine-type=g1-small --disk-size=50 --num-nodes=2
```

#### Export the project_id for use in internal commands:

```bash
export PROJECT_ID="$(gcloud config get-value project -q)"
```

#### Build the production image of the borkbot container:

```bash
docker build -t gcr.io/${PROJECT_ID}/borkd:v1 .
```

#### Push the image to your private repo

```bash
gcloud docker -- push gcr.io/${PROJECT_ID}/borkd:v1
```

#### Create kubectl secret for our new image to use for slack app token called "borkbot-slack"

```bash
kubectl create secret generic borkbot-slack --from-literal verification_token=<TOKEN_HERE_PLS>
```

#### Run the deployment file to create a new deployment resource

```bash
kubectl create -f ./production-app/borkbot-app-deployment.yaml
```

#### Check to make sure pods are running no errors pls

```bash
kubectl get pods
```

#### If pods are misbehaving run this to see what the container is printing

```bash
kubectl describe pods <pod_name>
```

#### Scale up the deployment if everything is fine

```bash
kubectl scale deployment borkbot-prod --replicas=3
```

#### Check to make sure all 3 replicas are healthy

```bash
kubectl get pods
```

#### SSH forward to a pod and check the healthcheck endpoint

```bash
kubectl port-forward <POD> 1443:8080
```

#### Expose deployment as a service

```bash
kubectl create -f ./production-app/borkbot-service-prod.yaml
kubectl describe service
```

## TRANSPORT

### - Nginx INGRESS CONTROLLER - YAY HTTP MVP -

#### Install helm if neede

```bash
brew search helm; brew install helm
```

#### init helm on our new kubectl cluster with RBAC
##### (Role based access controls - this gives us default name space access)

```bash
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account tiller --upgrade
```

#### Check for tiller system

```bash
kubectl get deployments -n kube-system
```
You should see the tiller service running here after configuring the above. More information here: https://cloud.google.com/community/tutorials/nginx-ingress-gke

#### Deploy the nginx ingress controller via helm

```bash
helm install --name nginx-ingress stable/nginx-ingress --set rbac.create=true
```

#### Check for the service
```bash
kubectl get service
```
You should see the nginx-ingress setup as a service with a controller

#### Apply ingress to interface with controller

```bash
kubectl apply -f ./production-ingress-http/bork-ingress-resource-http.yaml
```

#### Configure google dns through console more information here: https://cloud.google.com/dns/quickstart

1) Create domain that aligns with one you own under cloud dns

        EX: I own example.com I could create a sub domain like api.production.example.com

2) Update domain registers records under your OWN (not google but where you got your domain from) to point to zone

3) Add a new A record in the created zone to point to our deployment nginx controllers load balancer

4) Test dns

## SECURE - SSL CERTIFICATE ISSUANCE

### Using Cert Manager with Let's Encrypt

#### Install cert-manager

```bash
helm install --name cert-manager --namespace kube-system stable/cert-manager
```

#### Get certbot locally to get signed up

```bash
brew install certbot
```

#### Register with certbot to get an account and key

```bash
cerbot register
```

#### Backup created folder

Certbot registery generates a folder. Don't lose this and back it up

#### Create issuer secret from certbot account

```bash
kubectl create secret generic letsencrypt-staging --from-file=tls.key=private_key.json
```

This private_key.json is in the folder of the account generate in the previous register step

#### Configure issuer/clusterissuer

```bash
kubectl create -f production/bork-ingress-staging.yaml
```

#### Configure cert manager service account

```bash
  gcloud iam service-accounts create kube-cert-manager --display-name "kube-cert-manager"
  gcloud iam service accounts keys create ~/.config/kube-cert-manager.json \
      --iam-account kube-cert-manager@default-project-200117.iam.gserviceaccount.com
  gcloud iam service-accounts keys create ~/.config/kube-cert-manager.json \
      --iam-account kube-cert-manager@default-project-200117.iam.gserviceaccount.com
  gcloud projects add-iam-policy-binding default-project-200117 \
      --member serviceAccount:kube-cert-manager@default-project-200117.iam.gserviceaccount.com \
      --role roles/dns.admin
```

#### Configure kubernetes secret with account credentials for DNS challenge

```bash
kubectl create secret generic kube-cert-manager-google --from-file=${HOME}/.config/kube-cert-manager.json
```

#### Create Staging Certificate

```bash
kubectl create -f production-ingress-http/bork-ingress-certificate.yaml
```

#### Check status of certificate

```bash
kubectl describe certificates
```

You should see events in the output of this command that confirm the certificate has been generated

After confirming the certificate was created create the production issuer and create the production certificate

Let's Encrypt uses the staging to test the process works since the production api is rate limited heavily. The production api result is the only "VALID" ca issued cert.