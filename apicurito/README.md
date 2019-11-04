# Apicurito Operator

An operator for installing [Apicurito](https://github.com/Apicurio/apicurito), a small/minimal version of Apicurio used for standalone editing of API designs.

The aicurito operator can:
 - Install apicurito
 - Reconcile replica count
 - Upgrade apicurito

## Installing the operator

Before running the operator, the CRD must be registered with the Kubernetes apiserver:
```
$ kubectl create -f deploy/crds/apicur_v1alpha1_apicurito_crd.yaml
```

Setup RBAC and deploy the apicurito-operator:
```
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
$ kubectl create -f deploy/operator.yaml
```

Verify that the apicurito-operator is up and running:
```
$ kubectl get deployment
NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
apicurito-operator       1         1         1            1           1m
```

## Start an apicurito deployment
Edit the example Apicurito CR at deploy/crds/apicur_v1alpha1_apicurito_cr.yaml:
```
$ cat deploy/crds/apicur_v1alpha1_apicurito_cr.yaml
apiVersion: apicur.io/v1alpha1
kind: Apicurito
metadata:
  name: apicurito-service
spec:
  size: 3
  image: apicurio/apicurito-ui:latest

$ kubectl apply -f deploy/crds/apicur_v1alpha1_apicurito_cr.yaml
```
Ensure that the apicurito-operator creates the deployment for the Apicurito CR:
```
$ kubectl get deployment
NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
apicurito-operator       1         1         1            1           2m
apicurito-service        3         3         3            3           1m
```

# Upgrade apicurito
In order to upgrade apicurito, you need to install the desired version of the operator. Once the newer version is installed, an upgrade of the operand will kick in.

## Building the operator

In the apicurito directory issue the following command: 

```bash
make
```

## Upload to a container registry

e.g.

```bash
docker push quay.io/apicurito-operator/:<version>
```

