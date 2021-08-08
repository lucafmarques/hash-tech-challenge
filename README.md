# Hash Test

Repository for Hash's Go Backend Developer technical test

More info on the technical test can be found here: https://github.com/hashlab/hiring/tree/master/challenges/pt-br/new-backend-challenge

# Overview

This project contains the `checkout` service responsible for creating a checkout order out of a cart of products. For this, we must consult the `discount` gRPC service, responsible for providing discount percentage for products.


## Architecture

```mermaid
flowchart LR
    subgraph checkout
        direction LR
        server{{server}} <-.-> repository[/repository\]
    end
    
    server <---> discount{discount}
```

## Sequence Diagram

```mermaid
sequenceDiagram
    participant consumer
    participant checkout
    participant discount
    participant repository
    consumer->>checkout: request checkout with cart of products
    checkout->repository: fetch products data
    loop Each product
        checkout->>discount: request discount
        alt is up
            discount->>checkout: provide discount
        else is down
            Note right of checkout: Discount service is down<br>Don't apply any discounts
        end
    end

    checkout->>consumer: respond with checkout order

    opt Blackfriday
        Note left of checkout: Include blackfriday gift  
    end
```