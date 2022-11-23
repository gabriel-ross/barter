### TODO
* Get WithOffset and WithLimit working for list query builder

### PLAN
## Model
User:
id string
name string
email string

Account:
id string
user User
Dollars float64
Apples int

// rather than bank account it's a marketplace
market place 
object id: 
things u would take for it 
user: account: id

Transaction:
id string
timestamp datetime
unit id/string
quantity float64 // would need some way to validate currencies that can be subdivided vs currencies that can't. Maybe all   
    currencies CANNOT be subdivided and then there's some internal conversion done (like USD is all in cents and then 
    divided by 100)
sender Account
recipient Account


Currency:
id string
humanReadable string
unit string

## Spec
* Documentation
* Request validation
    - includes MIME format validation
        - probably only need to support JSON though
* Pagination
* User creation/JWT issuer
* Authorization (JWT authorization)
* Containerization (deployed via cloud run)
* proper use of the following status codes:
    200
    201
    204
    401
    403
    405
    406
* Postman collection & env
* Self links in responses 

## Features
* Only admins can make changes to the schedule and class resources
    - though maybe this should be an admin user rather than an admin flag 
* Reference maintainance
* List requests for user-related resources must only display resources associated with the provided JWT

### ARCHITECTURE
* JWT issuer
    - Issues JWT, registers user in user db if they don't exist
* scheduler
    - wrapped firestore client
    - chi router
    - validation middleware
    - auth middleware
    - rendering service?
        - separate rendering for individual and list responses
* User
* Class
* Timeslot

### QUESTIONS
* Do I need separate database, request, and response models?
    - What about patch that makes a database call and then unmarshals on top of the fetched resource? Do I need to write an in-house update-non-zero method?