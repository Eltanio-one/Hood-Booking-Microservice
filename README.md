# Hood Booking Microservice
This project was created to allow easier booking of the laminar flow hoods for members of my research institute.
Currently, we operate using an excel sheet which can easily be edited or lost/corrupted, and so it felt necessary to develop a more efficient and less easily corrupted solution.
This project also stood as a first attempt at programming using Golang, which I have been learning through reading but had not put into practice with a project.
Within this README are explanations of each file within the project, and also a roadmap that will set out future goals for this project which I will be continually working on.
Currently I use Postman to send HTTP requests to this microservice as I have yet to design a front-end for increased usability, but as this project is a proof of concept I felt more inclined to ensure the back-end was fully functional before starting work on the front.
I have currently designed this booking system to operate on a full-day booking system, but hope to implement specific timings soon.

This microservice allows registration and login of users. Session cookies are used post-login to validate that users have an account and are only trying to send requests related to their own profile.
Currently, I have used maps to store data, but plan to increase security by connecting a database to the microservice, which will be mentioned in the roadmap.
Once logged in, the main use of the microservice is to see which hoods are booked and at what time, and also to make bookings which can then be seen by other users.
The rest of the README will cover how edge cases are handled within each file for requests, and the intricacies of these requests.

The Hood Booking Microservice was created by Dan Haver, who can be contacted via [email](mailto:haverd08@gmail.com).

## Technologies
Golang

## Database Setup
(ALERT: STILL IN DEVELOPMENT, TUTORIAL WILL BE UPDATED WHEN COMPLETED)

This project is connected to a PostgreSQL database that is hosted in PgAdmin4. A sql file containing the CREATE TABLE commands is present in the code for this project, and below will show you how to execute these commands using **psql** to instantiate your tables ready for population.

This tutorial assumes that you have created a server and relevant database within PgAdmin 4, in which you will have chosen a host (e.g. localhost), a port (e.g. 5432), a username, and a name for the database within your server. These will be used to initialise a connection via your command line with the database.

- Open up your operating systems command line (terminal, powershell et)

- Initialise a connection with your database.
  
  ``` psql -h your_host -p your_port -U your_username -d your_database ```

- Enter the password that you created for the database when prompted.

- Once a connection is initialised, execute the SQL commands using \i and the path to your sql file.

  ``` \i path/to/sqlfile.sql ```

- For good measure, check pgAdmin to ensure your tables have been created under the schema tab of your database!

# Features

## Registration
### Handler Package
- Handles POST HTTP requests for user registration.
- Validates data input from the user:
    - Checks for missing data.
    - Checks if the user already has an account.
- Uses bcrypt to hash passwords before storage for increased security.
- Adds users to map for storage.

## login
### Handler Package
- Handles POST HTTP requests for user login.
- Validates data input from the user:
    - Checks that the username provided exists within the user data map.
    - Once a username is matched, ensures the hash of the password provided matches the stored hash.
- Generates a session token that is stored in a map and linked with a value of the users ID, to verify in later requests that the user is only attempting to send requests involving their profile.
- Stores the session token as a cookie that will be sent in any further requests in the http.Request.

## Bookings
### Handler Package
### GET requests
- Session cookies are verified and the list of current bookings for all users/hoods/times is returned.
### POST requests
- Session cookies are verified, and the session token map is consulted to ensure that the user is only trying to create a booking for themselves.
- Validates data input from the user:
    - Checks for missing data.
    - Ensures both the hood and user profiles exist.
    - Validates that both the user and the hood are not already booked for the date provided.
- Adds the booking to the booking list, which can then be queried by all users to inform whether they need to book a different hood or shift work to a different day if all hoods booked.

## Updating User Profile
- Users can send PUT requests via the user handler package to update their details.
- Verification of data follows similar processes as above, where missing data is checked and the user can only edit their own profile data.

## Data Packages
- There are three data packages, 'user', 'hood' and 'booking'.
- Each of these contains helper functions and the data structures temporarily being used to house data until a database connection is implemented.

### Session Package
- This package contains helper functions involved in session token/cookie creation and authentication.
- Having created a separate session package, I plan to organise my handler/data packages in a similar way to aid with the organisation of my files and understanding.

# Roadmap
- [x] Implement user registration and login capabilities.
    - [x] Include session cookie creation and validation post-login.
    - [x] Validate data entry upon POST requests.
- [ ] Enable booking of a hood at a specific time:
    - [x] Full-day booking.
    - [ ] Specific time-slot booking.
- [ ] Allow editing of bookings and deletion of bookings.
- [ ] Reorganise packages to be centered around struct types.
- [ ] Add unit tests for the microservice (currently only testing manually using Postman).
- [ ] Create a database to host all data, and connect to the microservice.
- [ ] Create a front-end for increased user experience and functionality.
