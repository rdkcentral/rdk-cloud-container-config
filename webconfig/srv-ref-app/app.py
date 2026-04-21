from flask import Flask, Blueprint, render_template, request, Response, jsonify
import json
import re
import subprocess

app = Flask(__name__)
app1 = Blueprint('app1', __name__, url_prefix='/app1')


@app1.route('/')
def index():
    return render_template('index.html')

def sanitize_mac_address(mac_address):
    mac_address = mac_address.replace(":", "").lower()
    # Validate MAC address format (optional)
    if not re.match(r'^([0-9a-fA-F]{2}){5}[0-9a-fA-F]{2}$', mac_address):
        return -1
    return mac_address

def normalize_json_data(data):
    """
    Normalize JSON data to ensure correct types for Python program.
    Handles both list-based and map-based JSON structures.
    """
    if isinstance(data, dict):  # If data is a map
        for key, value in data.items():
            if isinstance(value, str):
                # Convert "true"/"false" strings to boolean
                if value.lower() == "true":
                    data[key] = True
                elif value.lower() == "false":
                    data[key] = False
            elif isinstance(value, float):
                # Convert float to int if applicable
                if value.is_integer():
                    data[key] = int(value)
            elif isinstance(value, (dict, list)):
                # Recursively normalize nested maps or lists
                data[key] = normalize_json_data(value)
    elif isinstance(data, list):  # If data is a list
        for i in range(len(data)):
            if isinstance(data[i], (dict, list)):
                # Recursively normalize nested maps or lists
                data[i] = normalize_json_data(data[i])
    return data

@app1.route('/send', methods=['POST'])
def receive_data():
    if request.method == 'POST':
        print("Received data from Webconfig UI")
        subdoc_name = request.form.get('subdoc_name')
        print("Received subdoc_name data from UI:" ,subdoc_name)
        subdoc_data = request.form.get('subdoc_data')
        print("Received subdoc_data from UI:" ,subdoc_data)
        param_name = request.form.get('param_name')
        print("Received param name from UI:" ,param_name)
        mac_address = request.form.get('mac_address')
        MAC = sanitize_mac_address(mac_address)
        if MAC == -1:
            return "Invalid MAC Address"
        else:
            print("Received MAC address from UI:" ,MAC)


        dict_data = json.loads(subdoc_data)
        normalized_data = normalize_json_data(dict_data)
        with open('subdoc_data.json', 'w') as f:
            json.dump(normalized_data, f, indent=4)


        print("Creating msgpack.........")
        cmd = [
             "go", "run", "create_msgpack_main.go", "subdoc_data.json", subdoc_name, param_name, MAC
        ]
        result = subprocess.run(cmd, capture_output=True, text=True)

        if result.returncode != 0:
            return jsonify({"error": result.stderr}), 500

        if result.returncode == 0:
            return "Request Submitted Successfully"
        else:
            return "Request Failed"

@app1.route('/api/v1/device/<mac>/document/<doc_name>', methods=['POST'])
def handle_request(mac, doc_name):
  try:
    print("Got POST request.........")
    # Get JSON data from the request
    json_data = request.get_json()

    MAC = sanitize_mac_address(mac)
    if MAC == -1:
        return "Invalid MAC Address"
    else:
        print("Received MAC address from UI:" ,MAC)

    param_name = request.args.get('param_name')

    normalized_data = normalize_json_data(json_data)
    with open('subdoc_data.json', 'w') as f:
        json.dump(normalized_data, f, indent=4)


    # Pass the MAC, document name, and JSON data as arguments to sample.py
    print("Creating msgpack.........")
    cmd = [
       "go", "run", "create_msgpack_main.go", "subdoc_data.json", doc_name, param_name, MAC
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)

    if result.returncode != 0:
        return jsonify({"error": result.stderr}), 500

    return jsonify({"message": "Request successful"}), 200

  except Exception as e:
     return jsonify({"error": str(e)}), 500


app.register_blueprint(app1)

@app.route('/')
def root_blocked():
    abort(404)

# Optional: a custom 404 page
@app.errorhandler(404)
def not_found(e):
    return render_template('404.html'), 404


if __name__ == '__main__':
    app.run(debug=False)

