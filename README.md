# Structure

This explains the structure.

## Storage

At the very lowest level, we have the storage interfaces. We create a new struct of type `Storage`, and this struct has fields for each of the various tables/entities we want to store for our application. For instance, `Merchant` and `User`.

These subfields are interfaces - the `Merchant` field is the `storage/merchant.Database` interface type. All this interface cares about is implementing the `Create` and `GetByID` methods. This is the same for the `User` field and interface type.

We then have a `storage/mysql` package. This package returns a new `storage.Storage` type, and implements the `Merchant` and `User` interface field types using MySQL as the backend database.

## Services

Services implement the business logic side of the project.

We have a main `service` package at the root level, which uses the `storage` interfaces to get and save data. These service level methods should be used in the API, custom cli programs, etc to handle all of the business side calls. For instance, when creating a merchant, part of the business logic is to create a new random ID if one isn't passed in.

The question of how to handle services importing other services and the inherent cyclical dependency issue that may arise. We could either create other service packages who's only goal is to import other service packages, but this can become cluttered given larger applications, and there's a good chance code will be duplicated across 'top level' services.

One solution, which is implemented in this example, is to create interfaces for each service, and create another root level package that implements each interface but can be shared between individual services. Each service struct implements a separate interface which defines other services. When a new `services` is created, this also sets this service interface to be used by each individual service, ie each individual service can call `s.service.<service>.<method>` from its struct method.

With this 'dependency' idea implemented, we can call top level service methods within other top level service methods. While we have to worry about infinite recursion, ie function A calls function B which also may call function A again, this would be present in any flat package structure.

This basically gives us the ability to keep services separated into their own packages, whill still being able to essnetially cyclically import top level services.

## Billing Example

**This will be implemented soon**

Accounting is a database table, and hence has a storage level package for it. However, when handling billing... this will need to take in merchant, user, and accounting data - this should not be handled by the accounting top level service.

Instead we can branch this idea out into its own packages - billing. So at the storage level, we will have a new billing package. Even though there is no billing table, this billing storage package will handle doing the custom JOIN query we need to grab merchant, user, and account data and placing it into one struct.

The top level billing service will then use the other core services to handle billing (ie on payment we will want to call `service.Billing.AddAmountPaid`, and this service in part will call `core.Accounting.Update`).

# Flow

Storage implements interfaces to make it easy to use whatever backend database, or mix of them, we want to implement.

A MySQL implementation will simply then import the `storage` package interfaces and make it's `Database` struct type which implements all of the interface methods.

The top level, business logic service (`service`) take in the storage struct. It can then mix and match whatever top level service (via the `dep` package in this example) and/or storage functions it needs to perform business logic.

Finally, an API command-line program then can just create a new storage backend, which is passed into service, and merchant creation becomes a single one line command of `service.Merchant.Create(...)`.

Each package along the way only cares about what it needs to - if the API wants to accept a different request and response JSON object for instance when creating a merchant, it can create these types and simple map to and from the main services.

# Concerns

### 1. Infinite Recursion

It's possible with this structure that service method A can call dependency method B, while the service that implements dependency method B calls dependency method A, which is implement by service method A. When this happens, the program will be in an infinite recursion loop. It can be argued this goes into developer competency, as infinite recursion remains an issue even with flat packages (function A calls function B and vice versa).

### 2. Splitting Into Microservices

Compared to a true DDD implementation, it may be more difficult to decouple services from one another to split them into microservices.

# Sample cmd/invoice/main.go Output

```sh
>go run cmd/invoice/main.go
running...
[+] Creating new MySQL storage implementation...
[+] Creating new service...
Created merchant and added to MySQL database...
Created user and added to MySQL database...
Created invoice and added to MySQL database...
[+] New invoice: {ID:1 MerchantID:1 BillTo:Bill Smith PayTo:John Doe AmountDue:100 AmountPaid:0 Status:pending}
Got invoice from MySQL database...
Created transaction and added to MySQL database...
[+] Paid invoice: {ID:1 MerchantID:1 BillTo:Bill Smith PayTo:John Doe AmountDue:0 AmountPaid:100 Status:paid}
Created transaction and added to MySQL database...
Got invoice from MySQL database...
[+] New transaction processed: {ID:2 MerchantID:1 Type:refund CardType:visa AmountCaptured:100 InvoiceID:1}
Got invoice from MySQL database...
[+] Invoice after transaction refund: {ID:1 MerchantID:1 BillTo:Bill Smith PayTo:John Doe AmountDue:100 AmountPaid:0 Status:pending}
```

# Running the API

Clone the project and then browse to the `cmd/api` folder.

Run `go run main.go`

The API will be running on `http://localhost:8080`

## Create a New User

Run this cURL request:

```sh
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{
    "email": "test@test.com",
    "password": "TestPassword123"
}' \
http://localhost:8080/api/v1/signup
```

## Create a New Invoice

```sh
curl -X POST \
    -H 'Authorization: Bearer <TOKEN>' \
    -H 'Content-Type: application/json' \
    -d '{
    "bill_to": "John Doe",
    "pay_to": "John Smith",
    "amount_due": 100
}' \
http://localhost:8080/api/v1/invoice
```

## Get Invoices

```sh
curl -X GET \
    -H 'Authorization: Bearer <TOKEN>' \
http://localhost:8080/api/v1/invoice
```

# Thanks

Full credit to the following people for their ideas and help on how to implement this structure.

**borosr** ([https://github.com/borosr](https://github.com/borosr))  
**hajnalandor** ([https://github.com/hajnalandor](https://github.com/hajnalandor))  
**PumpkinSeed** ([https://github.com/PumpkinSeed](https://github.com/PumpkinSeed))  
**brianvoe** ([https://github.com/brianvoe](https://github.com/brianvoe))  
**jschweihs** ([https://github.com/jschweihs](https://github.com/jschweihs))