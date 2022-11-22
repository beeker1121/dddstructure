# Structure

This explains the structure.

## Storage

At the very lowest level, we have the storage interfaces. We create a new struct of type `Storage`, and this struct has fields for each of the various tables/entities we want to store for our application. For instance, `Merchant` and `User`.

These subfields are interfaces - the `Merchant` field is the `storage/merchant.Database` interface type. All this interface cares about is implementing the `Create` and `GetByID` methods. This is the same for the `User` field and interface type.

We then have a `storage/mysql` package. This package returns a new `storage.Storage` type, and implements the `Merchant` and `User` interface field types using MySQL as the backend database.

## Services

This is split up into two parts - the "core" services, and the "business logic" (root level) services.

The core services handle the most rudimentary level of business logic - ie we need to create a merchant, and if an ID isn't passed in, we create one on the fly (storage doesn't care about this handling of an ID, since storage is not technically part of business logic). The core service for merchant creation, however, does not involve the core service of user creation.

Merchant creation at the true business logic level though is handled within the main merchant creation service. This service uses *both* the core service of merchant creation and user creation (when creating a new merchant, we also need to create a new user for that merchant).

This allows us to have business logic level services use core level services, and get around possible cyclical dependencies.

## Billing Example

Accounting is a database table, and hence has a storage level package for it. However, when handling billing... this will need to take in merchant, user, and accounting data - this should not be handled by the accounting top level service.

Instead we can branch this idea out into its own packages - billing. So at the storage level, we will have a new billing package. Even though there is no billing table, this billing storage package will handle doing the custom JOIN query we need to grab merchant, user, and account data and placing it into one struct.

The top level billing service will then use the other core services to handle billing (ie on payment we will want to call `service.Billing.AddAmountPaid`, and this service in part will call `core.Accounting.Update`).

# Flow

Storage implements interfaces to make it easy to use whatever backend database, or mix of them, we want to implement.

A MySQL implementation will simply then import the `storage` package interfaces and make it's `Database` struct type which implements all of the interface methods.

The core service (`service/core`) then takes in the `Storage` sturct, which gives it access to backend persistence.

The top level, business logic service (`service`) take in the core service struct. It can then mix and match whatever functions it needs to perform business logic.

Finally, an API command-line program then can just create a new storage backend, create a new core service, which is passed into service, and merchant creation becomes a single one line command of `service.Merchant.Create(...)`.

Each package along the way only cares about what it needs to - if the API wants to accept a different request and response JSON object for instance when creating a merchant, it can create these types and simple map to and from the main services.

# Sample cmd/billing Output

```sh
➜  dddstructure git:(billing) go run cmd/billing/main.go
running...
[+] Creating new MySQL storage implementation...
[+] Creating new core service...
[+] Creating new service...
[+] Creating a merchant via the service...
Created merchant and added to MySQL database...
Created user and added to MySQL database...
[+] Created merchant: {ID:1 Name:John Doe Email:johndoe@gmail.com}
[+] Getting merchant via the service...
Got merchant from MySQL database...
[+] Got merchant: {ID:1 Name:John Doe Email:johndoe@gmail.com}
[+] Getting user for merchant...
Got user from MySQL database...
[+] Got user: {ID:1 Name:John Doe Email:johndoe@gmail.com}
[+] Adding amount owed via billing service...
Created accounting entry and added to MySQL database...
[+] Getting billing information for all merchants...
Got billing entry from MySQL database...
[+] Got billing information for merchant: {ID:1 MerchantID:1 MerchantName:John Doe AmountDue:100}
[+] Adding 100 as amount paid...
Got accounting entry from MySQL database...
Updated accounting entry '1'...
[+] Getting accounting entry for merchant '1'...
Got accounting entry from MySQL database...
[+] Got accounting entry: {ID:1 MerchantID:1 UserID:1 AmountDue:0}
```