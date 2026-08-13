Team A                     Team B                    Team C
AWS                        Azure                     On-prem
Payments API               Customer API              Legacy API
     │                          │                         │
     └──────────────────────────┼─────────────────────────┘
                                │
                         Gloo Gateway
                                │
                    Authentication / Security
                     ┌──────────┼──────────┐
                     │          │          │
                    OIDC       JWT       Basic Auth


## Sample Example CR expected from feature team

apiVersion: platform.mac.com/v1
kind: GatewayAPI
metadata:
  name: payments-api
spec:
  teamName: payments
  apiName: payments-api

  backend:
    url: https://payments.aws.example.com

  authentication:
    type: oidc


## CR Prep
Developer
   │
   │ kubectl apply
   ▼
GatewayAPI CR
   │
   ▼
Kubernetes API Server
   │
   │ validates against
   ▼
GatewayAPI CRD
   │
   ├── Invalid → REJECT
   │
   └── Valid
         │
         ▼
       etcd