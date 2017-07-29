# ApiKeyBinding

| Field | Required | Description |
| ----- | -------- | ----------- |
| apiVersion<br />*string*   | `true`       |   APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values.   |
| kind<br />*string*   | `true`      |    Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase.         |
| metadata<br />*[ObjectMeta](https://kubernetes.io/docs/api-reference/v1.6/#objectmeta-v1-meta)*  | `true`    |     Standard object's metadata.        |
| spec<br />*[ApiKeyBindingSpec](#apikeybindingspec)*   | `true`     |      Defines an ApiKeyBinding   |

# ApiKeyBindingSpec

| Field | Required | Description |
| ----- | -------- | ----------- |
| proxy<br />*string*   | `true`  |  The name of the `ApiProxy` that this binding applies to. |
| keys<br />*[Key](#key) array*   | `true`    |   List of `ApiKey`s that belong to this binding.  |

# Key

| Field | Required | Description |
| ----- | -------- | ----------- |
| name<br />*string*   | `true`  |  Name of the `ApiKey` |
| quota<br />*integer*   | `false`    |  Number of requests that this `ApiKey` is granted.  |
| rate<br />*[Rate](#rate)*   | `false`    |  The rate limiting policy for this `ApiKey`  |
| defaultRule<br />*[Rule](#rule)*   | `false`    | The default rule this `ApiKey` has for fine grained access. Default is `false` |
| subpaths<br />*[Path](#path) array* | `false` | Defines find grained authorization based on subpath. If not defined, falls back to the `defaultRule` for any subpath |

# Rate

| Field | Required | Description |
| ----- | -------- | ----------- |
| amount<br />*integer*   | `true`       | Scalar value for the defined `unit`  |
| unit<br />*string*    | `true`       | Unit of rate limit. Valid values are `second`, `minute`, `hour`   |

# Rule

| Field | Required | Description |
| ----- | -------- | ----------- |
| global<br />*boolean*   | If undefined, *granular* must be defined.       |   If true, access to all HTTP methods is granted.  |
| granular<br />*[Granular](#granular)*   | If undefined, *global* must be defined.        |    Defines granular rules    |

# Granular

| Field | Required | Description |
| ----- | -------- | ----------- |
| verbs<br />*string array*   | `true`       |    List of http verbs that this `ApiKeyBinding` has access to.    |

# Path

| Field | Required | Description |
| ----- | -------- | ----------- |
| path<br />*string*   | `true`       | The subpath  |
| rule<br />*[Rule](#rule)*    | `true`       |  The rules defined for this subpath  |