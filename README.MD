<!-- Written by Lawrence Colman, (C) 2020 -->

# Fabrikam Auth Challenge
## Background
This challenge is focused on a fictional company "Fabrikam", who want to build a platform to allow HR staff such as recruiters to interact with job candidates in order to perform screening inverviews with them. The goal of the project is to enable Fabrikam customers (typically recruiters) to have video conferencing calls with other Fabrikam customers (typically job candidates) to screen call them on the fly using their existing identities stored in the Fabrikam identity provider without requiring from them to use and share external party identities like Zoom, or Skype.  In this challenge only the identity provision, application access and user management parts of the project are considered.

## Requirements
To meet the challenge for this project, three distinct parts must be either implemented or provided by the solution.  These consist of:
1. Fabrikam identity provider
2. Screening web application
3. Screening application backend/API

The use of pre-existing tools, libraries, platforms and frameworks is highly encouraged wherever possible.

## Exclusions
Excluded from the scope of the challenge are parts not strictly related to the identity access and management for the project, such as:
 * Any of the actual video capture, transfer or rendering that would be required in order to actually facilitate calls,
 * The storage or processing of candidate personal information or career profile information, such as real name, street address, work history, etc,
 * Dealing with and details of actual jobs, such as position descriptions, tasks and duties, etc,
 * Actually facilitating any type of interaction between users once they are authenticated and authorized (i.e. if a user can login and access some API resource, the challenge is met, there is no need for this to be visible to any other user),
 * Implementing the persistance and/or data storage and retrieval logic, such as databases, etc.


# Solution
## Architecture
In terms of application architecture, I opted to go with a single page application (SPA) approach for the frontend, making use of React-router for in-app navigation.  Beyond this the backend API was written as a [Golang](https://golang.org/) microservice making use of the [Gorilla toolkit](https://www.gorillatoolkit.org/pkg/mux).
These are expanded upon in more detail in the following sections:

## Identity Provider (IdP)
As was encouraged by the challenge case study, [Auth0](https://auth0.com) was chosen for the identity provider for the solution.
This means that users, their permissions, and the responsibility for the secure storage of their credentials is all managed via that Auth0 platform.
More details about this aspect of the solution are available within the [/docs/Auth0.MD](./docs/Auth0.MD) page.

## Web Application (Frontend)
Despite my own lack of first-hand experience or familiarity, [React.JS](https://reactjs.org/) was chosen as the framework for the rapid development of a minimal web app frontend.  This was enabled by making use of the popular [create-react-app](https://github.com/facebook/create-react-app/) starter kit from Facebook.  This was then easily augmented with code taken from [the sample code provided by Auth0 for React](https://github.com/auth0-samples/auth0-react-samples/tree/master/Sample-01#license)), which made life a lot eaier thanks to their awesome documentation.

## Application Backend (API)
Given the discrete requirements of the challenge in terms of backend fnctionality, I opted for a RESTful API microservices solution.  As such, the Go programming language seemed an obvious choice for the development, and I had knowledge of many of the componenets and paterns required.
A more full discusion of the API implementation and design decisions is available in [docs/API.MD](./docs/API.MD).


# Notes
## Project Structure
The repository is organised as a polyglot monorepo, meaning all of the (poly) code in many different languages (JavaScript for the `webapp` and Go for the `screening-api`) resides in one single (mono) repository.
The directory layout within the frontend `webapp` is almost exactly that produced by the `create-react-app` starter kit (with some extras for auth).
While the layout of the Golang `screening-api` follows conventions of modern (i.e. >v1.11) modules-based Go code, with separation of the module's public (`pkg/**`) and private (`internal/**`) functions and structures.  More detail on this structuring is availabile within the [docs/API.MD documentation file](./docs/API.MD).

## How to run it
### The easy way
(from repo root:)
`docker-compose up`

The React frontend will then be available at [http://loalhost:3000/], and the API can be interrogated at [http://localhost:8080/v1/screening/screening] (via cURL, Postman, etc).
### Developer mode
Or, to start everything up locally (i.e. not dockerized), the usual commands for each sub-repo:
**API**
1. `cd screening-api`
2. `go run main.go` (requires a working Go environment, I was using Go v1.15)
**WWW**
1. `cd webapp`
2. `npm install && npm start` (requires working NodeJS/NPM environmnt)

## Copyright
Except as noted in the acknowledgements below, all souce code within this solution is authored and owned by me.
Where other sources have been used in predominant part (such as create-react-app an the auth0-react example code),
those parts along with my modifications are released under the original license of the other sources used (which is
the [MIT license](https://opensource.org/licenses/MIT) in the case of both [create-react-app](https://github.com/facebook/create-react-app/blob/master/LICENSE),
and [the Auth0 React sample code](https://github.com/auth0-samples/auth0-react-samples/tree/master/Sample-01#license), and the [The 3-Clause BSD License](https://opensource.org/licenses/BSD-3-Clause) for the [Gorilla Toolkit](https://github.com/gorilla/mux/blob/master/LICENSE)).

## Acknowledgements
### Splash/Hero image
The front homepage "Hero" image ([this one](https://c.pxhere.com/photos/7a/bb/bridge_light_suspension_bridge_east_river_relexion-87244.jpg!d) of the Brooklyn Bridge,) is from [here](https://pxhere.com/en/photo/87244), and has been released for unlimited free personal and commercial use with no attribution required under the [CC0 (public domain) license](https://creativecommons.org/share-your-work/public-domain/cc0/).

### Challenge Case Study
The challenge case study itself (found within [docs/case-study.pdf](./docs/case-study.pdf)) is not my own work, and is the property of its original authors.
It is included within this respository in order to give a more complete context for the rest of the source code and solution.

## The Dodgy Hack Bits
(You knew they were coming).

### Frontend
Due predominantly to time constraints, (and also not being a React developer _at all_,)
the frontend is very basic and rushed.  From time-to-tome you may see weird React-related warnings in the developer
console of the browser - this is expected due to me not knowing React well enough to fix them in the available time.
The screening page:
 * Has almost _zero_ styling whatsoever,
 * Has an awful layout, particularly for the "Book New Screening" form,
 * Has no input validation on the form _at all_ (you can submit everything blank if you so desire),
 * Seems to have a bug in the `DatePicker` (probably how I have used it, not given it the rights props or whatever) where the component does not populate when clicked - as a workaround, arbitrary numbers can be typed in and it works fine.

### Backend
Unfortunately (and, unlike many other IdP providers,) Auth0 generates JWT tokens that contain _multiple_ audience (`aud`) claims, in the form of a JSON array.  While this is [completely within the specification](https://tools.ietf.org/html/rfc7519#section-4.1.3), it is not the expectation that seems to exist within the [jwt-go](https://github.com/dgrijalva/jwt-go) library (which itself seems to be a _de facto_ standard in Golang for dealing with JWTs).  Given this shortcoming, I was unable to make use of the `VerifyAudience` function of the library, and have instead had to implement my own array looping logic to parse and validate the audience claim (see the `handleMultipleTokenAudiences` function within the `screening-api/internal/auth/validate.go` file).  However, the standard library logic has been used for all other standard claim validation (such as token expiry, issuer, etc).
Apart from this hiccup I am quite happy with the backend code overall, given the available time and level of detail given in the challenge case study.
