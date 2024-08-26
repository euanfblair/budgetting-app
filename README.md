# Project Setup and Installation Guide

## Prerequisites

Ensure you have the following tools and dependencies installed:

- **Air**: For live reloading during development.
    - Installation: Follow the [Air Installation Guide](https://github.com/air-verse/air).

- **Tailwind CLI Standalone**: For Tailwind CSS without additional build tools.
    - Installation: Follow the [Tailwind CLI Installation Guide](https://tailwindcss.com/blog/standalone-cli).

- **PostgreSQL**: For the database setup.

## Database Setup

1. **Create PostgreSQL Database and Tables**

   Use the following SQL commands to set up the required tables in your PostgreSQL database:

   ```sql
   CREATE TABLE sessions (
       token VARCHAR(43) NOT NULL PRIMARY KEY,
       data BYTEA,
       expiry TIMESTAMP WITHOUT TIME ZONE
   );

   CREATE TABLE transactions (
       transactionid SERIAL PRIMARY KEY,
       name VARCHAR(100),
       type BOOLEAN,
       amount INTEGER,
       transaction_date DATE,
       user_id INTEGER,
       category VARCHAR
   );

   CREATE TABLE users (
       user_id SERIAL PRIMARY KEY,
       first_name VARCHAR(255),
       surname VARCHAR(255),
       email VARCHAR(255),
       password_hash VARCHAR(255),
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
   );
   ```

## Installation Steps

1. **Install Air**

   Follow the instructions on the [Air GitHub repository](https://github.com/air-verse/air) to install Air.

2. **Install Tailwind CLI**

   Visit the [Tailwind CLI blog post](https://tailwindcss.com/blog/standalone-cli) for detailed installation instructions on the Tailwind CLI.

## Configuration

Ensure your application is properly configured to connect to the PostgreSQL database and utilize Tailwind CSS as needed. Refer to the documentation of your respective tools for configuration details.

## Usage

In 2 separate Terminal tabs run 
 - `./tailwindcss -i ui/static/styles/base.css -o ui/static/styles/output.css --watch`
- `air`

