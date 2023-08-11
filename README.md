# Bootcamp Product Service
This service is have feature to 
1. Get All Products
2. Create Products (only admin)
3. Add to Cart 
4. Get User Cart
5. Checkout
6. Get All Order align with role with pagination

## Setup and Installation

1. clone this repository 
2. create new database to store bootcamp.sql
3. import database to mysql 
```
mysql -u username -p database_name < path/to/03-finalproject.sql
```
4. copy .env.example file and rename to .env 
5. fill the env with your credentials, database credentials and jwt secret specially
```
APP.AUTH_SERVICE_BASE_URL=http://localhost:port
APP.AUTH_SERVICE_VALIDATE_URL=/v1/users/validate
```
6. run go generate command in root project to setup project
```
go generate ./...
```

## Run and Test
To run this program, run this command in root terminal 
1. First run Auth Service
2. run this command
```
go run . 
```