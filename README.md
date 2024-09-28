# Loan Service
Practice for creating a loan service. Create loans from available loan products, register for loans as borrowers, invest in loans as an investor, approve and disburse loans as staff, and assist in loan approval and disbursement as field validator.

## Prerequisites
- Go 1.19
- Docker 20
- Docker Compose 3
- Postgres 13

## Quickstart
- Download Go deps and other dependencies `make dep` and `make setup`
- Fill in `.env` by copying from defaults `env.sample`
- Register additional users (borrowers and investors) `make register-user`
- Build Go app with `make build` before running via Docker
- Run service as multi-container with docker-compose.yml `make run`
- Seed database with initial data (users, roles, loan products, loans) `make seed-db`
- App runs on `localhost:8080` by default, with Postgres in `localhost:5432`

## Definition
The term loan refers to a type of credit vehicle in which a sum of money is lent to another party in exchange for future repayment of the value or principal amount. In many cases, the lender also adds interest or finance charges to the principal value, which the borrower must repay in addition to the principal balance.

There are several important terms that determine the size of a loan and how quickly the borrower can pay it back:
- **Principal**: This is the original amount of money that is being borrowed.
- **Loan Term**: The amount of time that the borrower has to repay the loan.
- **Interest Rate**: The rate at which the amount of money owed increases, usually expressed in terms of an annual percentage rate (APR, or 'per annum').
- **Loan Payments**: The amount of money that must be paid every month or week in order to satisfy the terms of the loan.

Source: [Investopedia](https://www.investopedia.com/terms/l/loan.asp)

## Functional requirements (and assumptions)
### Users
Under assumption, there are 4 main types of users:
- Borrower, a person that requests a loan.
    - A borrower can request loans, view all loans requested by them (including loan details, remaining principal, installment details, payments, etc.), request disbursement, and pay loan installments.
- Investor, a person that contributes to the loan principal to be disbursed.
    - An investor can view all approved loan applications with their details (borrower profile, principal amount, rate, ROI, etc.), invest in loans, view loans they invested on including payments and status, and
- Field validator, a person that liaises with the borrower and approves, collects loan agreement letter, and hands the disbursed money to the borrower.
- Staff, a person that administrates and coordinates the process of loan approval, investment, and disbursement.
- Additionally, there is a superuser for the application with all read-write permissions.

All users can log in, edit their profile, and log out.
- For simplicity purposes, let's assume that registration is out of scope, and the following users are already registered in the system:
    - 1 superuser (user: `admin@loanservice.io`, pass: `@admin`)
    - 1 staff (user: `staff@loanservice.io`, pass: `@staff`)
    - 1 field validator (user: `field.validator@loanservice.io`, pass: `@field.validator`)
- Additional users with custom emails can be registered by executing `make register-user`. See the [Quickstart](#quickstart) section for more details.

### Loans
- A loan can be in the following states: `[proposed, approved, invested, disbursed]`. The state change must move forward in that order.
- When a loan is created it will have `proposed` as the initial state.
- A loan can be approved by a staff, which will change the state into `approved`.
    - A loan approval must contain several information:
        - An image proof that a field validator has visited the borrower
        - The employee ID of field validator
        - Date of approval
    - Once a loan is approved, it cannot go back to the proposed state.
    - Once approved, a loan is ready to be offered to investors.
- A loan is considered invested when the total invested amount is equal to the loan principal amount. Once that amount is reached, the state will change to `invested`.
    - A loan can have multiple investors, each with their own amount.
    - Total of invested amount must not exceed the loan principal amount.
    - Once the loan is invested, all investors will receive an email containing a link to the agreement letter in PDF.
- A loan is disbursed when the loan is given to the borrower, which changes the state to `disbursed`.

A loan object should consist of the following information:
- Borrower ID number
- Principal amount
- Interest rate, which defines tbe total interest that borrower will pay.
    - Let's assume that the interest is a simple interest with the rate expressed in APR (per annum), and an additional attribute 'loan term' will be needed to calculate total interest.
- ROI (return of investment), will define total profit received by investors.
- A link to the generated loan agreement letter.
- *Not listed above, but let's assume we will have additional attributes needed*:
    - Loan term: borrowers can chooses between the predetermined loan periods: 1 month, 3 months, 6 months, and 1 year.
    - Total interest: Will define the total interest that will be accrued from the loan principal given the interest rate and loan term.

## Non-functional requirements (and assumptions):
Assume this is an MVP project, and traffic is expected to be low as the application will be subjected to internal user acceptance testing (UAT) before being further developed with production-ready infrastructure.
- Therefore, scaling isn't a concern. Expect minimal requests per second (rps), at around 1 rps.
- Therefore, this application is NOT a distributed system, and will be a monolith architecture for MVP purposes.
    - The solution will be deployed as a multi-container application (using Docker Compose) to a cloud VM, with only a single instance per application server.
- The system will expose a REST-ful backend API only.
- Monitoring will solely consist of logging from the STDOUT of each Docker container.
- Tests will include unit tests for core functions in the loan module, and end-to-end tests for user actions for the core loan functionality.
- Latency should be minimal, and we would like to prioritize consistency over availability.
    - We will use Echo webserver for Golang (personal pick), and PostgreSQL for the database in order to enforce ACID properties.
- Reliability and fault tolerance will be enforced to the best of the system's single-instance ability.
    - We will use the `restart` flag on the `docker-compose.yml` for the application service, so that it restarts if it panics.
    - We will use database transactions to ensure atomicity and ensure consistency, so that a bad write will not propagate throughout the system.
- Security will be implemented with a simple role-based access control, as well as rate limiting and JWT authentication with short-lived tokens (5 minutes).

## Out of scope
- Loan payments
- Loan installment calculation
- Payment and disbursement transactions (actual bank transactions via bank transfer or payment gateways)
- User registration

## Entity design
The data entities will include the following objects:
- Users
- Roles (borrower, investor, staff, field validator, superuser)
- Products (different loan products have different principals, interest rates, and terms)
- Loans (a single loan of a specific product, requested by borrower and paid for by investors)
    - Approvals and disbursements will be embedded into the main loan entity since it has 1-to-1 relationship.
- Investments (money pitched in by investors to loans)











