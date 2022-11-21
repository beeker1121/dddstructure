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

# Flow

Storage implements interfaces to make it easy to use whatever backend database, or mix of them, we want to implement.

A MySQL implementation will simply then import the `storage` package interfaces and make it's `Database` struct type which implements all of the interface methods.

The core service (`service/core`) then takes in the `Storage` sturct, which gives it access to backend persistence.

The top level, business logic service (`service`) take in the core service struct. It can then mix and match whatever functions it needs to perform business logic.

Finally, an API command-line program then can just create a new storage backend, create a new core service, which is passed into service, and merchant creation becomes a single one line command of `service.Merchant.Create(...)`.

Each package along the way only cares about what it needs to - if the API wants to accept a different request and response JSON object for instance when creating a merchant, it can create these types and simple map to and from the main services.