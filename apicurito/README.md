# Apicurito Operator

An operator for installing [Apicurito](https://github.com/Apicurio/apicurito), a small/minimal version of Apicurio used for standalone editing of API designs.

The aicurito operator can:
 - Install apicurito
 - Reconcile replica count
 - Upgrade apicurito

## Installing the operator

Before running the operator, the CRD must be registered with the Kubernetes apiserver and RBAC permissions configured:
```
# Execute as a cluster-admin equivalent account
$ make -C config setup
```

Deploy the apicurito-operator:
```
$ make -C config operator

#
# Alternative:
# To modify the image and tag being used for the operator
#
$ IMAGE=<image url> TAG=<tag version> make -C config operator
```

Verify that the apicurito-operator is up and running:
```
$ kubectl get deployment
NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
apicurito-operator       1         1         1            1           1m
```

## Start an apicurito deployment
Edit the example Apicurito CR at config/samples/apicur_v1alpha1_apicurito_cr.yaml:
```
$ cat config/samples/apicur_v1alpha1_apicurito_cr.yaml
apiVersion: apicur.io/v1alpha1
kind: Apicurito
metadata:
  name: apicurito-service
spec:
  size: 3

$ make -C config app
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


## CSV Bundle Generation

```bash
make -C config bundle

# OR
# w/ sha lookup/replacement against registry.redhat.io
DIGESTS=true REDHATIO_TOKEN="<username>:<password>"  make csv
