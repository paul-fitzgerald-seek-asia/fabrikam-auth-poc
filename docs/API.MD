# Fabrikam screening-api (Golang)
## RESTful API
### Entities
The **screening** API has been implemented with minimal scope for this challenge, and as such only deals with one type of entity: `Screening`.
A screening is essentally a schedule event, consisting of a starting time, a duration, and a title.  Once stored (ephemerally in memory only) within the service, a screening also contains a unique ID.
An example JSON encoding of a screening structure is as follows:

    {
        "startTime": "2020-`1-30T15:36:24.000Z",
        "duration": 60,
        "title": "Interview with Lawrence"
    }


### API endpoints
As a modern microservice pattern, the implemented API is versioned, with all paths begining with `/v1`.  Furthermore, the API's service name is then appended to the path so as to enable easier reverse-proxying without the risk of path conflicts.  As such, the base path for all endpoints is `/v1/screening`, with entity specific paths being children of this base.

#### GET Screening
The `GET` `/screening` endpoint _lists_ the existing screenings for a user presenting the `list:screening` scope.  The body of a successful HTTP 200 `OK` response will be in the format of a JSON encoding of the `screening` entity, as shown above, with an ID field added containing a UUID of the format shown below.

#### POST Screening
The `POST` `/screening` endpoint _creates_ a new screening for a user presenting the `create:screening` scope.  The body of this request must be a JSON encoding of the `screening` entity, as shown above.
Successful requests to this endpoint respond with a HTTP 201 `Created` response, with a response body in the form:

    {
        "screeningID": "123e4567-e89b-12d3-a456-426614174000"
    }


### Authentication & Authorization
The API is implemented with a custom middleware for the validation of HTTP Bearer token JWT's, with the verification of these tokens being handled by the [jwt-go](https://github.com/dgrijalva/jwt-go) library.
This middleware, along with a basic in-memory JWKS cache is contained within the `auth` package of the API codebase.
More information regarding the particulars of the authentication and authorizaion used is given within the [docs/Auth0.MD](./Auth0.MD) file.


## Golang code
### Package discipline
As the API has been implemented using modern Go (such as modules, Go v1.11 onwards), the packages have been structured
to make use of [Go's (as of v1.14) `internal` packages](https://golang.google.cn/doc/go1.4#internalpackages).  Meaning other code bases cannot import the internal workings of this repository, apart from the one exported (_public_) package "`definitions`".
The `definitions` package contains the reusable `struct` for un/marhalling `Screeing` entities from this API.  All other codebases wishing to interact with this API should do so making use of the RESTful endpoints outlined above, rather than sharing any internal code or logic.

### Package structure
The packages have been structured to follow [Go's convention of being "_short, concise, evocative_"](https://golang.org/doc/effective_go.html#package-names) and each refer to a single **responsibility**, rather than being outlined in terms of function or level of abstraction.
For example, rather than having `handlers`, `business` (or "service"), and `persistance` packages, each with various responsibilities spread throughout, all three layers of abstraction are present in the single `screening` package.  Which is responsible for _handling_ requests to the `/screening` endpoint, _servicing_ these requessts and managing the `storage` of screeing entity data (only using an in-memory data structure, for now).

> screening-api
> ├── internal
> │   ├── auth
> │   ├── config
> │   ├── screening
> │   └── server
> └── pkg
>     └── definitions
