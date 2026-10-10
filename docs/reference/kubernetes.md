# Kubernetes helpers

These packages are small wrappers used by operators. They do not talk to the ioFog Controller.

## `pkg/k8s`

Import `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/k8s`.

```go
func New(kubeConfig string) (*Client, error)
func NewInCluster() (*Client, error)
```

`Client` embeds `*kubernetes.Clientset`. `New` loads a kubeconfig file (`clientcmd.BuildConfigFromFlags`). `NewInCluster` uses the pod service account. The package blank-imports the GCP auth plugin.

```go
func (cl *Client) WaitForLoadBalancer(namespace, name string, timeoutSeconds int64) (string, error)
func (cl *Client) WaitForPod(namespace, name string, timeoutSeconds int64) error
```

`WaitForLoadBalancer` watches Services in the namespace until the named Service has an ingress IP or hostname. When both are set, the hostname is returned. Watch `Error` or `Deleted` returns an error. An empty address after the watch ends returns `IP and Hostname values were empty`.

`WaitForPod` watches Pods and waits until a Pod whose name, minus the last two hyphen segments, equals `name` is `Running` and every condition is `True`. Pods with fewer than three hyphen segments are ignored. The watch stops once that Pod is ready. A watch error or delete returns an error.

Both watches pass `timeoutSeconds` as the Kubernetes list/watch timeout.

## `pkg/k8s/operator`

Import `github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/k8s/operator`.

`Reconciliation` is the value a reconcile step returns before it becomes a controller-runtime result.

| Field | Meaning |
| --- | --- |
| `End` | Stop this reconcile without requeue. |
| `Requeue` | Requeue after `Delay`. |
| `Delay` | Passed to `RequeueAfter`. |
| `Err` | Returned to the runtime. |

`IsFinal` is true when `Err != nil` or `Requeue` or `End`.

`Result()` precedence: error (`RequeueWithError`), then `End` (`DoNotRequeue`), otherwise `RequeueWithDelay(Delay)`. A zero `Reconciliation` from `Continue()` therefore requeues after a zero delay (`RequeueAfter: 0`, `Requeue: true`).

| Constructor | Result |
| --- | --- |
| `ReconcileWithError(err)` | `Err` set. |
| `ReconcileWithRequeue(delay)` | `Requeue` and `Delay` set. |
| `Continue()` | zero value. |
| `Reconcile()` | `End: true`. |

Direct `ctrl.Result` helpers:

| Function | Result |
| --- | --- |
| `DoNotRequeue()` | empty result, nil error |
| `RequeueWithError(err)` | empty result, the error |
| `Requeue()` | `Requeue: true` |
| `RequeueWithDelay(delay)` | `Requeue: true`, `RequeueAfter: delay` |

## Related

- [Operations](../operations.md)
- [NATS bootstrap](nats.md) (`GetNatsBootstrap` is the Controller call operators use before these Kubernetes helpers create NATS)
