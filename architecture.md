### Model
User:
ID string
Name string
Email string
PhoneNumber string
Token string

Account:
ID string
User User
Funds map[string]float64
TransactionVolume int
CreationDate datetime

Transaction:
ID string
Timestamp datetime
Quantities map[string]float64
Sender Account
Recipient Account

Currency:
ID string
Name string

// ???
Listing
ID string
Acc Account
Object CurrencyQuantity
Price []CurrencyQuantity // list of things you'd sell

### Spec
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

### Features
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
* why wasn't db working? could it be different regions? or bad service acc?