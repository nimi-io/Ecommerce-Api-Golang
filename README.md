# Ecommerce API Documentation

This API allows users to register, log in, view their profiles, manage products (admin-only), and place and manage orders.

PreSaved - Admin User
```json
{
  "email": "admin@admin.com",
  "password": "admin"
}
```

Please Note that a postman api documentation had been added - "Ecommerce_API_Postman_Collection.json" 

## Endpoints

### Authentication Endpoints

#### Registering a User
- **Endpoint:** `POST /api/v1/auth/register`
- **Description:** Register a new user by sending a POST request with the following JSON payload:

  ```json
  {
    "username": "johndoe",
    "email": "johndoe@example.com",
    "password": "password123"
  }
  ```

- **Responses:**
  - `200 OK`: User successfully registered.
  - `400 Bad Request`: Invalid or missing fields.

#### Logging In
- **Endpoint:** `POST /api/v1/auth/login`
- **Description:** Log in by sending a POST request with the following JSON payload:

  ```json
  {
    "email": "johndoe@example.com",
    "password": "password123"
  }
  ```

- **Responses:**
  - `200 OK`: User successfully logged in, returns a JWT token.
  - `401 Unauthorized`: Invalid email or password.

### User Endpoints

#### Get User Information
- **Endpoint:** `GET /api/v1/user`
- **Description:** Retrieve information of the authenticated user.

- **Responses:**
  - `200 OK`: Returns user details.
  - `401 Unauthorized`: User not authenticated.

### Admin Product Management Endpoints

#### Create a Product (Admin Only)
- **Endpoint:** `POST /api/v1/admin/products/create`
- **Description:** Create a new product by sending a POST request with the following JSON payload:

  ```json
  {
    "name": "Sample Product",
    "description": "Product description",
    "price": 49.99,
    "stock": 100
  }
  ```

- **Responses:**
  - `201 Created`: Product successfully created.
  - `403 Forbidden`: Unauthorized access for non-admin users.

#### Get a Single Product by ID
- **Endpoint:** `GET /api/v1/admin/products/:id`
- **Description:** Retrieve details of a specific product by its ID.

- **Responses:**
  - `200 OK`: Returns the product details.
  - `404 Not Found`: Product not found.

#### List All Products
- **Endpoint:** `GET /api/v1/admin/products`
- **Description:** Retrieve a list of all products.

- **Responses:**
  - `200 OK`: Returns a list of products.
  - `403 Forbidden`: Unauthorized access for non-admin users.

#### Update a Product (Admin Only)
- **Endpoint:** `PATCH /api/v1/admin/products/:id`
- **Description:** Update a product's details by sending a PATCH request with the fields to update. Example payload:

  ```json
  {
    "name": "Updated Product Name",
    "price": 59.99
  }
  ```

- **Responses:**
  - `200 OK`: Product successfully updated.
  - `404 Not Found`: Product not found.
  - `403 Forbidden`: Unauthorized access for non-admin users.

#### Delete a Product (Admin Only)
- **Endpoint:** `DELETE /api/v1/admin/products/:id`
- **Description:** Delete a product by its ID.

- **Responses:**
  - `200 OK`: Product successfully deleted.
  - `404 Not Found`: Product not found.
  - `403 Forbidden`: Unauthorized access for non-admin users.

### Order Management Endpoints

#### Place an Order
- **Endpoint:** `POST /api/v1/order/`
- **Description:** Place an order by sending a POST request with the following JSON payload:

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
    ],
  }
  ```

- **Responses:**
  - `201 Created`: Order successfully placed.
  - `400 Bad Request`: Invalid data provided.

#### List My Orders
- **Endpoint:** `GET /api/v1/order/list`
- **Description:** Retrieve a list of orders for the authenticated user.

- **Responses:**
  - `200 OK`: Returns a list of user's orders.
  - `401 Unauthorized`: User not authenticated.

#### Cancel an Order
- **Endpoint:** `POST /api/v1/order/cancel/:id`
- **Description:** Cancel an existing order by its ID.

- **Responses:**
  - `200 OK`: Order successfully canceled.
  - `404 Not Found`: Order not found.

#### Update Order Status (Admin Only)
- **Endpoint:** `PATCH /api/v1/order/status/:id`
- **Description:** Update the status of an order (e.g., mark as shipped or delivered). Admin-only access.

- **Responses:**
  - `200 OK`: Order status updated.
  - `404 Not Found`: Order not found.
  - `403 Forbidden`: Unauthorized access for non-admin users.

### Miscellaneous Endpoints

#### Test Endpoints
- **GET /api-1** - Testing endpoint 1.
- **GET /api-2** - Testing endpoint 2.



