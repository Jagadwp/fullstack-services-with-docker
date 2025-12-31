import json
import logging
import os
from flask import Flask, request, jsonify

# Logging configuration
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s"
)
logger = logging.getLogger(__name__)

# Configuration
DATA_DIR = os.getenv("DATA_DIR", "/data")
RECEIVED_DIR = os.path.join(DATA_DIR, "received")

# Ensure directory exists at startup
os.makedirs(RECEIVED_DIR, exist_ok=True)

app = Flask(__name__)

# Routes
@app.route("/process", methods=["POST"])
def process_user():
    """
    Receive user data from Go scheduler.
    Expected JSON:
    {
        "id": int,
        "name": str,
        "email": str
    }
    """
    if not request.is_json:
        logger.warning("Request content-type is not application/json")
        return jsonify({"error": "Invalid content type"}), 400

    data = request.get_json()

    # Basic validation
    if not isinstance(data, dict):
        logger.warning("Invalid JSON payload")
        return jsonify({"error": "Invalid JSON payload"}), 400

    user_id = data.get("id")
    name = data.get("name")
    email = data.get("email")

    if not user_id or not name or not email:
        logger.warning("Missing required fields in payload: %s", data)
        return jsonify({"error": "Missing required fields"}), 400

    file_name = f"user_{user_id}.json"
    file_path = os.path.join(RECEIVED_DIR, file_name)

    try:
        with open(file_path, "w") as f:
            json.dump(data, f, indent=2)
    except Exception as e:
        logger.error("Failed to write file %s: %s", file_path, e)
        return jsonify({"error": "Internal server error"}), 500

    logger.info("Processed user ID=%s Name=%s", user_id, name)

    return jsonify({"status": "ok"}), 200


# App entry point
if __name__ == "__main__":
    logger.info("Starting Python API")
    app.run(host="0.0.0.0", port=5000)
