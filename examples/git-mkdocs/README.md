# git-mkdocs <!-- omit in toc -->

- [Use Case](#use-case)
- [Prerequisites](#prerequisites)
- [Step 1: Prepare your git repository](#step-1-prepare-your-git-repository)
  - [Step 1.1: Setup a SSH key](#step-11-setup-a-ssh-key)
  - [Step 1.2: Setup a webhook](#step-12-setup-a-webhook)
- [Step 2: Prepare your Kubernetes manifests](#step-2-prepare-your-kubernetes-manifests)
  - [Step 2.1 : Update the secret.yaml manifest](#step-21--update-the-secretyaml-manifest)
  - [Step 2.2 : Update the configmap.yaml manifest](#step-22--update-the-configmapyaml-manifest)
  - [Step 2.3 : Update the ingress.yaml manifest](#step-23--update-the-ingressyaml-manifest)
  - [Step 2.4 : Update the deployment.yaml manifest](#step-24--update-the-deploymentyaml-manifest)
- [Step 3: Deploy](#step-3-deploy)

## Use Case

This example deploys a static documentation website generated with [mkdocs](https://www.mkdocs.org/). The documentation's source is stored in a GitHub repository. The website regenerates itself when changes are pushed to the `main` branch.

## Prerequisites

- A repository containing documentation structured for mkdocs.
- A Docker an image of the `git-volume-reloader` available in a registry. Currently, no public image is available.
- A domain name that resolves to your cluster's public IP address.
- A valid certificate for your domain name stored in a Kubernetes secret. Alternatively, you can use [`cert-manager`](https://cert-manager.io/docs/) to generate the secret for you.

## Step 1: Prepare your git repository

### Step 1.1: Setup a SSH key

Generate a SSH key pair without a passphrase.

```zsh
ssh-keygen -t RSA -b 2048
```

In you documentation repository, go to **Settings > Deploy keys** and click on **Add deploy key** to add the public SSH key you just generated. Keep your private key somewhere, you will need it later.

### Step 1.2: Setup a webhook

Generate a secret key.

```zsh
openssl rand -base64 32
```

In you documentation repository, go to **Settings > Webhooks** and click on **Add webhook**. For the webhook parameters:

- The **Payload URL** is the host of the documentation with `/webhook` appended. For example, if you want your documentation to be available on `docs.padok.cloud`, the payload URL will be `https://docs.padok.cloud/webhook`.
- The **Content Type** must be set to `application/json`.
- The **Secret** is the secret key you just generated above.
- Leave the other parameters with their default values.

## Step 2: Prepare your Kubernetes manifests

Copy the examples manifests avaialble in `./manifests`.

### Step 2.1 : Update the secret.yaml manifest

- Put the webhook's secret key you generated in step 1.2 under the `GITHUB_SECRET` key of the Kubernetes secret.
- Put the SSH private key you generated in step 1.1 under the `SSH_PRIVATE_KEY` key of the Kubernetes secret.

### Step 2.2 : Update the configmap.yaml manifest

In the `docs` ConfigMap, update all the values. The keys are self-explanatory.

### Step 2.3 : Update the ingress.yaml manifest

- Change all the `host` values to the hostname you want your documentation to be available at.
- Update the `tls` field of both Ingress resources to match your TLS certificate.

### Step 2.4 : Update the deployment.yaml manifest

Change the `image` of the `git-volume-reloader` container to the one you have in your registry. If needed, add an `imagePullSecrets` field in your Deployment.

## Step 3: Deploy

You can now deploy your documentation. Go to the folder you put your manifests in and run:

```zsh
kubectl apply -f .
```

Your documentation should be available at the URL you configured. You can also try to push a change to the documentation to the branch you configured. The documentation should update itself.
