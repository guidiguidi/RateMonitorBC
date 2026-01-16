# RateMonitorBC

RateMonitorBC is a service that finds the best exchange rate for a given currency pair, amount, and a set of exchanger marks using the BestChange API.

## Features

-   Finds the best exchange rate from BestChange.
-   Filters rates by the amount and required marks.
-   Provides both an HTTP API and a CLI for searching.
-   Handles errors from the BestChange API.

## Getting Started

### Prerequisites

-   Go 1.22 or higher

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/guidiguidi/RateMonitorBC.git
    cd RateMonitorBC
    ```

2.  Install dependencies:
    ```sh
    go mod tidy
    ```

### Running the API Server

To start the API server, run the following command:

```sh
go run cmd/api/main.go
```

The server will start on port `8080` by default.

## API Endpoints

### Find Best Rate

-   **Method:** `POST`
-   **Endpoint:** `/api/v1/best-rate`
-   **Description:** Finds the best exchange rate for a given currency pair and amount.

#### Request Body

```json
{
  "from_id": 1,
  "to_id": 10,
  "amount": 1000,
  "marks": ["manual", "reg"]
}
```

-   `from_id` (int, required): The ID of the currency you are giving.
-   `to_id` (int, required): The ID of the currency you are receiving.
-   `amount` (float, required): The amount you want to exchange.
-   `marks` (array[string], optional): A list of required exchanger marks.

#### Success Response (200 OK)

```json
{
  "from_id": 1,
  "to_id": 10,
  "amount": 1000,
  "marks": ["manual", "reg"],
  "best_rate": {
    "exchanger_id": "123",
    "rate": "0.00002345",
    "rankrate": "0.00002360",
    "inmin": "100",
    "inmax": "5000",
    "reserve": "123456.78",
    "marks": ["manual", "reg"],
    "from_amount": "1000.00000000",
    "to_amount": "0.02345000"
  },
  "source": "bestchange"
}
```

#### Error Responses

-   `400 Bad Request`: Invalid JSON, missing required fields, or amount <= 0.
-   `422 Unprocessable Entity`: No suitable rates found after filtering.
-   `502 Bad Gateway`: Error communicating with the BestChange API.

## CLI Usage

The CLI provides a way to find the best exchange rate from the command line.

### Command

```sh
go run cmd/cli/main.go --from <from_id> --to <to_id> --amount <amount> [--marks <marks>]
```

### Example

```sh
go run cmd/cli/main.go --from 1 --to 10 --amount 1000 --marks manual,reg
```

### Output

```
Best Rate Found:
  Exchanger ID: 123
  Rate: 0.00002345
  To Amount: 0.02345000
  Marks: manual, reg
```
