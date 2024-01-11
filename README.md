## Container resources policy
This policy is designed to enforce constraints on the resource requirements  of
Kubernetes containers. It follows a two-step verification process: initially
checking whether the container has defined resource requirements, and
subsequently ensuring that these resources fall within the permissible range
set by the maximum resource requirements configured in the policy settings.

## Configuration
Users can configure the policy using the following parameters in YAML format:

```yaml
# optional - the maximum permissible CPU limit
maxCPULimit: "2"
# optional - the maximum permissible memory limit
maxMemoryLimit: "4G"
# optional - the maximum permissible CPU requested
maxCpuRequest: "2"
# optional - the maximum permissible memory requested
maxMemoryRequest: "4G"
# optional - define the container image where the resource validation should be skipped
ignoreImages: ["myimage:latest", "registry.k8s.io/pause"]
```

Users can skip the some of the configuration. But an empty configuration is not
allowed. Thus, at least one of the configurations,  `maxCPULimit`,
`maxMemoryLimit`, `maxCpuRequest` or `maxMemoryRequest` should be defined. All
CPU and memory configuration should be expressed using the [quantity
definitions](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/quantity/)
of Kubernetes.

Full example of policy definition:

```yaml
apiVersion: policies.kubewarden.io/v1
kind: ClusterAdmissionPolicy
metadata:
  name: container-resources-policy
spec:
  policyServer: container-resources-policy
  module: registry://ghcr.io/kubewarden/policies/container-resources:latest
  rules:
  - apiGroups: [""]
    apiVersions: ["v1"]
    resources: ["pods"]
    operations:
    - CREATE
    - UPDATE
  mutating: true
  settings:
    maxCPULimit: "2"
    maxMemoryLimit: "4G"
    maxCpuRequest: "1"
    maxMemoryRequest: "1G"
    ignoreImages: ["myimage:latest", "registry.k8s.io/pause"]
```

## Constraints
The policy assumes that users want to defined both CPU and memory resource
limits for their containers. If the specified CPU or memory limit exceeds the
maximum permissible limit configured in the policy, the policy will intervene
to mutate the container's settings. It adjusts the limit to comply with the
configured maximum values. If the require requirements is missing in the
container definition, the policy will also mutate adding the missing
configuration

## Possible Scenarios
This policy operates as follows:

| Scenario | Configuration | Outcome |
| ---- | ---- | ---- |
| Resource limits within range | `maxCPULimit: "2"`<br>`maxMemoryLimit: "4G"` | Container settings accepted without modification if both CPU and memory limits are defined within permissible ranges. |
| Resource limits exceeds permissible max | `maxCPULimit: "2"`<br>`maxMemoryLimit: "4G"` | Policy intervenes if either CPU or memory limit exceeds predefined maximums. Mutates container settings to comply. |
| Undefined resource limits | `maxCPULimit: "2"`<br>`maxMemoryLimit: "4G"` | Policy intervenes if either CPU or memory limits are undefined . |
| Requests resources within range | `maxCpuRequest: "2"`<br>`maxMemoryRequest: "4G"` | Container settings accepted without modification if both CPU and memory requested are defined within permissible ranges. |
| Requested resources exceeds permissible max | `maxCpuRequest: "2"`<br>`maxMemoryRequest: "4G"` | Policy intervenes if either CPU or memory requested exceeds predefined maximums. Mutates container settings to comply. |
| Undefined requested resources | `maxCpuRequest: "2"`<br>`maxMemoryRequest: "4G"` | Policy intervenes if either CPU or memory requested are undefined . |

