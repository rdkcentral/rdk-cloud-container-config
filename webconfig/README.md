# WebConfig - Flask Configuration Web Application

A Flask-based web application for managing device configuration through a web UI. This application handles device subdocument configurations, validates MAC addresses, normalizes JSON data, and generates msgpack files for configuration updates.

## Overview

The webconfig application provides:
- Web-based UI for submitting device configurations
- REST API endpoints for programmatic configuration updates
- JSON data normalization and validation
- msgpack file generation via Go integration
- MAC address validation and sanitization

## Directory Structure

```
webconfig/
└── srv-ref-app/
    ├── app.py                      # Flask application main file
    ├── create_msgpack_main.go      # Go program for msgpack file generation
    ├── template/
    │   └── index.html              # Web UI form
    └── 404.html                    # Error page template (optional)
```

## Dependencies

### Python
- Flask - Web framework
- json - JSON processing (built-in)
- re - Regular expressions (built-in)
- subprocess - External process execution (built-in)

### System
- Go runtime - Required for running `create_msgpack_main.go`

## Installation

### Prerequisites
- Python 3.6+
- Flask (`pip install flask`)
- Go 1.11+

### Setup

1. Navigate to the srv-ref-app directory:
   ```bash
   cd ~/rdk-docker/rdk-cloud-container-config/webconfig/srv-ref-app
   ```

2. Install Python dependencies:
   ```bash
   pip install flask
   ```

3. Ensure Go is installed:
   ```bash
   go version
   ```

## Running the Application

Start the Flask application with:

```bash
python3 -m flask run -h <host_ip> -p <port>
```

### Examples

- Run on localhost port 5000:
  ```bash
  python3 -m flask run -h 127.0.0.1 -p 5000
  ```

- Run on all interfaces port 8080:
  ```bash
  python3 -m flask run -h 0.0.0.0 -p 8080
  ```

The application will be accessible at `http://<host_ip>:<port>/app1/`

## API Documentation

### Routes

#### 1. Web UI Form (`/app1/`)
- **Method**: GET
- **Description**: Serves the web configuration form
- **Response**: HTML form for submitting device configuration

#### 2. Form Submission (`/app1/send`)
- **Method**: POST
- **Description**: Processes form data from the web UI
- **Form Parameters**:
  - `subdoc_name` (string, required): Subdocument name
  - `subdoc_data` (JSON string, required): Configuration data as JSON
  - `param_name` (string, required): TR181 parameter name
  - `mac_address` (string, required): Device MAC address (format: `XX:XX:XX:XX:XX:XX`)
- **Response**: 
  - Success: "Request Sumbitted Successfully"
  - Error: "Invalid MAC Address" or error message from msgpack creation

#### 3. REST API Endpoint (`/api/v1/device/<mac>/document/<doc_name>`)
- **Method**: POST
- **Description**: Programmatic configuration update endpoint
- **URL Parameters**:
  - `mac`: Device MAC address (without colons)
  - `doc_name`: Document name
- **Query Parameters**:
  - `param_name`: TR181 parameter name
- **Request Body**: JSON configuration data
- **Response**: JSON with success/error message

### Example API Usage

```bash
curl -s -i "http://webconfig.rdkcentral.com:9008/api/v1/device/AABBCCDDEEFF/document/privatessid?param_name=Device.WiFi.Private" \
  -H 'Content-type: application/json' \
  -X POST \
  --data '{
    "private_ssid_2g": {
        "SSID": "private_2_rdkm",
        "Enable": true,
        "SSIDAdvertisementEnabled": true
    },
    "private_security_2g": {
        "EncryptionMethod": "AES",
        "ModeEnabled": "WPA2-Personal",
        "Passphrase": "*****"
    },
    "private_ssid_5g": {
        "SSID": "private_5_rdkm",
        "Enable": true,
        "SSIDAdvertisementEnabled": true
    },
    "private_security_5g": {
        "EncryptionMethod": "AES",
        "ModeEnabled": "WPA2-Personal",
        "Passphrase": "*****"
    }
}'
```

**Note:** Replace the passphrase values with proper credentials as per WiFi component requirements.

## Features

### MAC Address Validation
- Formats: Accepts both colon-separated (`AA:BB:CC:DD:EE:FF`) and no-separator formats
- Sanitizes to lowercase without colons
- Returns error for invalid MAC addresses

### JSON Data Normalization
- **String Boolean Conversion**: Converts "true"/"false" strings to boolean values
- **Float to Integer Conversion**: Converts floats that are whole numbers to integers
- **Recursive Processing**: Handles nested JSON objects and arrays
- Ensures data compatibility with Python program requirements

### Error Handling
- 404 error page for undefined routes
- JSON error responses with HTTP status codes
- Go program error output capturing and reporting
- Invalid MAC address detection

## Key Functions

### `sanitize_mac_address(mac_address)`
- Removes colons from MAC address
- Converts to lowercase
- Validates hexadecimal format
- Returns -1 for invalid input

### `normalize_json_data(data)`
- Recursively normalizes JSON structure
- Converts string booleans to boolean types
- Converts whole-number floats to integers
- Handles nested structures

### `receive_data()` / `handle_request(mac, doc_name)`
- Processes configuration submissions
- Creates temporary `subdoc_data.json` file
- Invokes Go program for msgpack generation
- Returns status/error messages

## Output Files

Generated during request processing:
- **`subdoc_data.json`**: Temporary JSON file with normalized configuration data (created in current working directory)
- Msgpack files: Generated by `create_msgpack_main.go`

## Configuration Files

The application generates a `subdoc_data.json` file when processing requests. This file contains the normalized configuration data in JSON format before msgpack conversion.

## Troubleshooting

### Port Already in Use
```bash
python3 -m flask run -h 0.0.0.0 -p 8080
```
Use a different port if 5000 is already occupied.

### Go Program Not Found
Ensure Go is installed and in PATH:
```bash
which go
```

### Invalid MAC Address Error
- Use format: `XX:XX:XX:XX:XX:XX` (with colons) or `XXXXXXXXXX` (without)
- Ensure all characters are valid hexadecimal (0-9, a-f, A-F)

### JSON Parsing Error
- Ensure `subdoc_data` is valid JSON
- Check that all required fields are properly formatted

## Files

- **[app.py](srv-ref-app/app.py)** - Main Flask application with route handlers
- **[create_msgpack_main.go](srv-ref-app/create_msgpack_main.go)** - Go program for msgpack file generation
- **[template/index.html](srv-ref-app/template/index.html)** - Web UI form

## Security Considerations

- MAC address validation prevents injection attacks
- JSON normalization ensures type safety
- Subprocess execution uses parameterized arguments to prevent command injection
- 404 blocking on root endpoint prevents information disclosure

## Development Notes

- Debug mode is enabled in the Flask app by default (`debug=True`)
- All console output is printed for debugging purposes
- Temporary `subdoc_data.json` files are created in the current working directory
- Both form-based and REST API submissions follow the same processing pipeline

## License

See project repository for license information.
