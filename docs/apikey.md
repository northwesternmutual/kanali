The `ApiKey` spec represents an rsa encrypted API key.

*NOTE:* while you can generate this spec yourself, it is recommended that you generate an `ApiKey` spec using the [`kanalictl`](https://github.com/northwesternmutual/kanalictl) tool. See example below for usage details.

Here are some useful commands to generate an rsa key pair:

```sh
# generate private key - using 4096 bit long modulus
$ openssl genrsa -out private.pem 4096
# calculate public key
$ openssl rsa -in private.pem -pubout -out public.pem
# use the kanalictl tool to generate ApiKey spec
$ kanalictl generate apikey -k public.pem -o apikey.yml -n my-test-api-key
Here is your api key (you will only see this once): ksAR0xqSKjh9UGSBvhP2IxDDC9Ckou0S
Corresponding Kubernetes config written to apikey.yml
$ cat apikey.yml
apiVersion: kanali.io/v1
kind: ApiKey
metadata:
  creationTimestamp: null
  name: my-test-api-key
  namespace: default
spec:
  data: 2778ac7127f97212dbc27cb9a8b7fb5a51f49bcefc8a23eb26fd9e3b8673faa9dc3e98597dcc62d1dcc01a5c054b28268c30b206c5e5296e058fded2458905382b15ba30eef7596ae46248b0958442e03e38ec1097a96a9fe6420fb671a06ed7782deadd0bb35f9ef1debb5693a34d20108647364834939a12f8a9959864c52f3df4d0cebbae60a27facf0b75bbae5e91077c2e013179810a7cdca77bca6c8d1a48acc3e6b3af72119f4886cc9c483063b5e42f660095d4e3f69c35a6511c9ecbe59c5893eb176c208c6d00c0eda2315416e856fd264ab886ee18527f6cc0c5311953a79ad2c1695780d322bb5d6cacac61d808bfbca531614084d7caac6a11f310127adb319ba53fcd91835d0bcf318f85242563ec555e2d3c0cefbff31585ec6a631f893a5dd57725002b4e9ac5d68ac4ba849d9f3314968ea63d1f8520060cf800fcded379a1353b6f018f431e5206018b9a5c81d52c13069a7621ca6b02de302ad830279f9963c957ef73a6e170f17883eefac405bae03796fcbb02e07b7ff1b691bd320a8a72a35203898664206ac386f730787160f94739459d11ab3b0019648414c6a4b9bcc7121a17a42aa8bd2e3e7a64234f9e78503833dd208c8a4a948b51491a0a4fec15f17a213c0ae4a5d87d002b8047f9aa235c9f32b052301e499d64d7650a1cb3a201f7342028d6b5f50e0f4ab7d3d3b3c4bbb410aa81e04
```

# ApiKey

| Field | Required | Description |
| ----- | -------- | ----------- |
| apiVersion<br />*string*   | `true`       |   APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values.   |
| kind<br />*string*   | `true`      |    Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase.         |
| metadata<br />*[ObjectMeta](https://kubernetes.io/docs/api-reference/v1.6/#objectmeta-v1-meta)*  | `true`    |     Standard object's metadata.        |
| spec<br />*[ApiKeySpec](#apikeyspec)*   | `true`     |      Defines an ApiKey   |

# ApiKeySpec

| Field | Required | Description |
| ----- | -------- | ----------- |
| data<br />*string*   | `true`   |  encrypted api key  |