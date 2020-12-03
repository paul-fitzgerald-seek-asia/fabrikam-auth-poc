# Auth0 Configuration
## Identity Platform
### Auth flow
The webapp frontend makes use of the `@auth0/auth0-react` npm library, which supports single page application's in implementing the **Authorization Code Flow with PKCE** auth flow.
Essentially, the React frontend, after login, requests an auth code from Auth0 and then sends this as a Bearer token to the API resource owner.  The token contains the authorized scopes for the authenticated user, along with  unique `subject` clim representing their identity.

## Users & Passwords
### Users
As the Fabrikam challenge deals with two types of users, two demo users have been setup:
 * `recruiter@fabrikam.example`
 * `candidate@fabrikam.example`

However, due to time constraints both have identical permissions and functionality within the implemented application and API.  However, inspection of ID tokens for each will reveal they are members of disctinct `recruiter` and `candidate` groups respectively.

### Default password
Each of these users has been setup with the password of `P4ssw0rd`, in order to meet complexity requirements.

## Permissions
The two users listed above have been preconfigured to each have the following:

    "permissions": [ "create:screening", "list:screening" ]

which enables them each to view and submit new screening entities to the API.

## Signup & new users
Although new user signup is still enabled, and can be performed by the Auth0 Lock widget included in the webapp, new users will not have any permissions granted by default.
As such, requests from newly signed up users will result in no action being performed by the API, with a log output similar to the following being output by the server:
> INFO rejecting as forbidden due to authenticated user not having scope create:screening
