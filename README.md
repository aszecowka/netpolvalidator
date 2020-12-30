# Network Policy Validator

**Netpolvalidator** allows you to validate Kubernetes Network Policies. By default, in Kubernetes, every component can
communicate with everything else, inside or outside of the cluster. Network Policy allows us to explicitly specify,
which communication is allowed. This increases **security**
of your cluster.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: ingress-to-component-a
spec:
  podSelector:
    matchLabels:
      app: component-a
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: orders
          podSelector:
            matchLabels:
              app: component-b
```

Unfortunately, when defining network policy, you can make a lot of mistakes:

- incorrect pod selector
- incorrect ingress/egress namespace or pod selector
- block access to DNS services that disable all communication

**Netpolvalidator** helps you detect such a kind of issues on the early stage.

## Usage

`make run`

## Development

- To build, tests and check quality of code, execute: `make all`