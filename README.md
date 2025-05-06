# DDDStructure

DDDStructure is an example SaaS application for invoicing. It includes a fully built and structured backend API built in Golang, and a fully built frontend application built in Vue.

There are a few goals with this project:

1. Showcase a mix of DDD and SOA principles for how to structure a Golang backend application/API.
2. Provide a baseline codebase, for both the frontend and backend of the application, where others can build from it.

Original concepts for this idea go back to my [gotodo](https://github.com/beeker1121/gotodo) project. This new project is a more fleshed out example of those same concepts, with examples of how to handle problems like sharing services, deploying with Kubernetes, etc.

## Preview

**Login**

![Login page](https://raw.githubusercontent.com/beeker1121/dddstructure/master/preview/login.png)

**Invoices**

![Invoices page](https://raw.githubusercontent.com/beeker1121/dddstructure/master/preview/invoices-none.png)

**Create new invoice**

![Invoices page](https://raw.githubusercontent.com/beeker1121/dddstructure/master/preview/invoice-new.png)

**Pay invoice**

![Pay invoice](https://raw.githubusercontent.com/beeker1121/dddstructure/master/preview/invoice-pay.png)

# Structure

This project uses a mix of DDD (Domain Driven Design) and SOA (Service Oriented Architecture) principles to structure the application.

Mainly, we want to think of our project and the packages contained therein in terms of their responsibilities, or 'domains'.

For example, we main goal of this project is to provide an HTTP based API. The API should be viewed as its own 'domain' - it should really only care about routing, handling incoming HTTP requests, sending back responses, so on and so forth. For this reason, all API specific logic is within the `cmd/api` package. It's in the `cmd` folder since ultimately, our API will be a binary we run.

All business logic then is handled within the `services` package or 'domain'. The API does not care at all about our business logic, it just cares about which service method to call that implements the business logic. For example, if we have a `POST` endpoint for `/api/v1/signup` to create a new user, the API creates this route to handle the incoming requests and then calls the `services.User.Create` service.

Going a layer lower, we have the `storage` 'domain'. Storage has no other concern other than creating, reading, updating, and deleting entries from whatever database is being used. The services domain then calls the storage domain to get data, persist data, etc.

Again the main point is for each domain to manage its own responsibilities - API handles HTTP requests, services handles our business logic, and storage handles getting and persisting data.

## Storage

At the very lowest level, we have the storage interfaces. We create a new struct of type `Storage`, and this struct has fields for each of the various tables/entities we want to store for our application. For instance, `Invoice` and `User`.

These subfields are interfaces - the `Invoice` field is the `storage/invoice.Database` interface type. All this interface cares about is implementing the `Create` and `GetByID` methods. This is the same for the `User` field and interface type.

We then have a `storage/mysql` package. This package returns a new `storage.Storage` type, and implements the `Invoice` and `User` interface field types using MySQL as the backend database.

## Services

Services implement the business logic side of the project.

We have a main `service` package at the root level, which uses the `storage` interfaces to get and save data. These service level methods should be used in the API, custom cli programs, etc to handle all of the business side calls. For instance, when creating a merchant, part of the business logic is to create a new random ID if one isn't passed in - the `services.Merchant.Create()` method handles this business logic.

The question of how to handle services importing other services and the inherent cyclical dependency issue that may arise. We could either create other service packages who's only goal is to import other service packages, but this can become cluttered given larger applications, and there's a good chance code will be duplicated across 'top level' services.

One solution, which is implemented in this example, is to create interfaces for each service, and create another root level package that implements each interface but can be shared between individual services. Each service struct implements a separate interface which defines other services. When a new `services` is created, this also sets this service interface to be used by each individual service, ie each individual service can call `services.<service>.<method>` from its struct method.

With this 'dependency' idea implemented, we can call top level service methods within other top level service methods. While we have to worry about infinite recursion, ie function A calls function B which also may call function A again, this would be present in any flat package structure.

This basically gives us the ability to keep services separated into their own packages, whill still being able to essnetially cyclically import top level services so we're not duplicating already built business logic.

# Concerns

### 1. Infinite Recursion

It's possible with this structure that service method A can call dependency method B, while the service that implements dependency method B calls dependency method A, which is implement by service method A. When this happens, the program will be in an infinite recursion loop. It can be argued this goes into developer competency, as infinite recursion remains an issue even with flat packages (function A calls function B and vice versa).

### 2. Splitting Into Microservices

Compared to a true DDD implementation, it may be more difficult to decouple services from one another to split them into microservices. This structure does allow us though to easily split microservices into a service domain. An example of this can be found in the `microservice` branches.

# Sample cmd/invoice/main.go Output

```sh
>go run cmd/invoice/main.go
running...
[+] Creating new MySQL storage implementation...
[+] Creating new service...
[+] Registering dependencies...
[+] Creating merchant...
Created merchant and added to MySQL database...
Created user and added to MySQL database...
[+] Creating invoice...
Created invoice and added to MySQL database...
[+] New invoice: {ID:1 MerchantID:1 ProcessorType:achcom BillTo:Bill Smith PayTo:John Doe AmountDue:100 AmountPaid:0 Status:pending}
[+] Paying invoice...
Got invoice from MySQL database...
Processing transaction via ACHCom...
Created transaction and added to MySQL database...
[+] Paid invoice: {ID:1 MerchantID:1 ProcessorType:achcom BillTo:Bill Smith PayTo:John Doe AmountDue:0 AmountPaid:100 Status:paid}
[+] Processing a separate transaction...
Processing transaction via ACHCom...
Created transaction and added to MySQL database...
Got invoice from MySQL database...
[+] Separate transaction processed: {ID:2 MerchantID:1 Type:refund ProcessorType:achcom CardType:visa AmountCaptured:100 InvoiceID:1}
[+] Getting invoice...
Got invoice from MySQL database...
[+] Invoice after transaction refund: {ID:1 MerchantID:1 ProcessorType:achcom BillTo:Bill Smith PayTo:John Doe AmountDue:100 AmountPaid:0 Status:pending}
```

# Running for Development

We use Docker for development, which runs and initializes the MySQL database.

Export the following variables:

```
DB_HOST=localhost
DB_PORT=3306
DB_NAME=dddstructure
DB_USER=root
DB_PASS=jLiEo34@3!%k
```

Run `docker-compose up` in a terminal from the root application directory.

Once the container has finished loading, open a separate terminal, browse to `cmd/api` and run `go run main.go`.

The API will now be running on `http://localhost:8080`. You can now run the frontend application that will use this API.

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

# Updating MySQL Models with SQLBoiler

We use SQLBoiler to generate the Go structs (models) based on our MySQL database tables. This ORM also allows us to easily query MySQL.

Install SQLBoiler and the MySQL driver if you haven't already:

```
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
```

Then run SQLBoiler to update the models based on the MySQL database schema:

1. Create a new database migration file in the `./db` folder, so we know what to update in production.

2. Update the `./db/init.sql` file with the migration, so when we start Docker it has the most up to date database version.

3. Run `sqlboiler mysql` in a terminal in the root folder (where `sqlboiler.toml` is).

4. The output will be to the `./storage/mysql/models` folder.

# Deployments

A guide on how to deploy this application using Kubernetes can be found in the README at `cmd/deployments/kubemysql/README.md`. It uses the `minikube` application to test a local cluster, running the backend API, MySQL database, and frontend.

# TODO

- :heavy_check_mark: Figure out how to handle `time.Time`, ie time coming in from the API, how to convert that to `time.Time` for service level, and finally how to store via storage layer.
  - API request struct has it commented out currently.
  - :heavy_check_mark: For storage side, the top level storage struct should still probably just use `time.Time` and the database itself can convert to and from that.
- :heavy_check_mark: Add validate methods where needed.
  - :heavy_check_mark: On invoice create, validate at least one line item is passed in.
    - :heavy_check_mark: Validate at least one payment method is passed in.
- Use xid for all IDs instead of an unsigned int.
- :heavy_check_mark: Determine if we want to refactor services.
  - :heavy_check_mark:Right now, we have certain functions like `Invoice.UpdateByIDAndUserID` to ensure a user is updating an invoice they own. Should this just be coded into the main `Invoice.Update()` function? Then, have another function like `Invoice.UpdateRaw()` that will let you pass any field of an invoice to update to that service method, just use the ID on the params? Or should that not be a 'service' level method, and if that needs to happen, then use `storage` directly?
    - :heavy_check_mark: Changed `Invoice.UpdateByIDAndUserID` to just `Invoice.UpdateForUser` - this much more clearly explains that the method is to update an invoice for a user, and the user ID, role, permissions, etc can all be checked within this function. The standard `Invoice.Update` method can be used as a general update (that can call other services, like the transaction service if needed), and `Invoice.UpdateForUser` uses the general update method after its own business logic.
- :heavy_check_mark: Move all main error logging to the service layer instead of the API layer.
  - :heavy_check_mark: Keep logs at the API layer as well though.
- :heavy_check_mark: Maybe local service package functions, such as `CalculateAmounts` for invoices, should be an invoice service method instead. Reason being is it can then access the logger, once we set it up so it's part of the service struct. This way we can easily log errors wherever in the service level, exactly where we should be handling errors, and sending back custom general errors like `ErrInvoiceNotFound`. Or, we just handle all errors in the main service level methods only (check return from `CalculateAmounts` inside `Update` for instance, log error there instead of in `CalculateAmounts`).
  - :heavy_check_mark: Decided to just keep `CalculateAmounts` as a separate function, not a method on the service. The service methods that call it log any errors returned.

# Thanks

Full credit to the following people for their ideas and help on how to implement this structure.

**borosr** ([https://github.com/borosr](https://github.com/borosr))  
**hajnalandor** ([https://github.com/hajnalandor](https://github.com/hajnalandor))  
**PumpkinSeed** ([https://github.com/PumpkinSeed](https://github.com/PumpkinSeed))  
**brianvoe** ([https://github.com/brianvoe](https://github.com/brianvoe))  
**jschweihs** ([https://github.com/jschweihs](https://github.com/jschweihs))