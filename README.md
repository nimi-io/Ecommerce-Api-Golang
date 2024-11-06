
# Ecommerce API Documentation

Welcome to the **E-commerce API** documentation! This API allows users to interact with an online store by registering, logging in, managing products (admins only), and placing orders. The key features include user authentication, product management, and order processing. 

To help with understanding and testing the endpoints, a Postman API collection is included: **"Ecommerce_API_Postman_Collection.json"**.

## Table of Contents
1. [Introduction](#introduction)
2. [Authentication Endpoints](#authentication-endpoints)
   - Registering a User
   - Logging In
3. [User Endpoints](#user-endpoints)
   - Get User Information
4. [Admin Product Management Endpoints](#admin-product-management-endpoints)
   - Create a Product (Admin Only)
   - Get a Single Product by ID
   - List All Products
   - Update a Product (Admin Only)
   - Delete a Product (Admin Only)
5. [Order Management Endpoints](#order-management-endpoints)
   - Place an Order
   - List My Orders
   - Cancel an Order
   - Update Order Status (Admin Only)
6. [Miscellaneous Endpoints](#miscellaneous-endpoints)
   - Test Endpoints

## Introduction

The **E-commerce API** is designed to simulate the functionality of an online store, providing essential features like user authentication, product management, and order handling. This project includes different levels of access: regular users and admin users.

- Regular users can:
  - Register and log in.
  - View their profile.
  - Place and manage orders.
- Admin users can:
  - Manage products (create, view, update, and delete products).
  - Update order status.

### Pre-saved Admin User

To test admin features, use the following pre-saved admin credentials:

```json
{
  "email": "admin@admin.com",
  "password": "admin"
}
```

For more detailed usage and testing, refer to the included Postman collection.

## Authentication Endpoints
### Registering a User
**Endpoint**: `POST /api/v1/auth/register`

**Description**: Register a new user with a unique username, email, and password.

**Request Payload**:
```json
{
  "username": "johndoe",
  "email": "johndoe@example.com",
  "password": "password123"
}
```
**Responses**:
- `200 OK`: Successfully registered user.
- `400 Bad Request`: Missing or invalid fields.

### Logging In
**Endpoint**: `POST /api/v1/auth/login`

**Description**: Log in with an existing user's email and password to receive a JWT token.

**Request Payload**:
```json
{
  "email": "johndoe@example.com",
  "password": "password123"
}
```
**Responses**:
- `200 OK`: Successfully logged in, returns JWT token.
- `401 Unauthorized`: Invalid email or password.

## User Endpoints
### Get User Information
**Endpoint**: `GET /api/v1/user`

**Description**: Retrieve the details of the authenticated user.

**Responses**:
- `200 OK`: Returns user details.
- `401 Unauthorized`: User is not authenticated.

## Admin Product Management Endpoints
Admin users can manage products by performing CRUD operations.

### Create a Product (Admin Only)
**Endpoint**: `POST /api/v1/admin/products/create`

**Description**: Create a new product.

**Request Payload**:
```json
{
  "name": "Sample Product",
  "description": "Product description",
  "price": 49.99,
  "stock": 100
}
```
**Responses**:
- `201 Created`: Product successfully created.
- `403 Forbidden`: Only admins can create products.

### Get a Single Product by ID
**Endpoint**: `GET /api/v1/admin/products/:id`

**Description**: Retrieve details of a product using its ID.

**Responses**:
- `200 OK`: Returns product details.
- `404 Not Found`: Product with specified ID not found.

### List All Products
**Endpoint**: `GET /api/v1/admin/products`

**Description**: Retrieve a list of all products available.

**Responses**:
- `200 OK`: Returns a list of products.
- `403 Forbidden`: Only admins can list products.

### Update a Product (Admin Only)
**Endpoint**: `PATCH /api/v1/admin/products/:id`

**Description**: Update an existing product's details.

**Request Payload**:
```json
{
  "name": "Updated Product Name",
  "price": 59.99
}
```
**Responses**:
- `200 OK`: Product successfully updated.
- `404 Not Found`: Product with specified ID not found.
- `403 Forbidden`: Only admins can update products.

### Delete a Product (Admin Only)
**Endpoint**: `DELETE /api/v1/admin/products/:id`

**Description**: Delete a product using its ID.

**Responses**:
- `200 OK`: Product successfully deleted.
- `404 Not Found`: Product with specified ID not found.
- `403 Forbidden`: Only admins can delete products.

## Order Management Endpoints
### Place an Order
**Endpoint**: `POST /api/v1/order/`

**Description**: Place an order for multiple products.

**Request Payload**:
```json
{
  "products": [
    {
      "productId": 1,
      "quantity": 2
    },
    {
      "productId": 3,
      "quantity": 1
    }
  ]
}
```
**Responses**:
- `201 Created`: Order successfully placed.
- `400 Bad Request`: Invalid request data.

### List My Orders
**Endpoint**: `GET /api/v1/order/list`

**Description**: Get a list of all orders placed by the authenticated user.

**Responses**:
- `200 OK`: Returns a list of the user's orders.
- `401 Unauthorized`: User is not authenticated.

### Cancel an Order
**Endpoint**: `POST /api/v1/order/cancel/:id`

**Description**: Cancel an order by its ID.

**Responses**:
- `200 OK`: Order successfully canceled.
- `404 Not Found`: Order with specified ID not found.

### Update Order Status (Admin Only)
**Endpoint**: `PATCH /api/v1/order/status/:id`

**Description**: Update the status of an order (e.g., marked as shipped or delivered).

**Responses**:
- `200 OK`: Order status successfully updated.
- `404 Not Found`: Order with specified ID not found.
- `403 Forbidden`: Only admins can update the order status.

## Miscellaneous Endpoints
### Test Endpoints
- `GET /api-1`: Testing endpoint 1.
- `GET /api-2`: Testing endpoint 2.

These endpoints are primarily intended for basic server functionality testing.

## Conclusion
This E-commerce API provides a simple structure for handling e-commerce operations, with robust user and admin functionalities. Make sure to test the API using the provided Postman collection to experience all features firsthand.

If you have any questions or run into issues, please don't hesitate to reach out!

Happy coding! 🎉
