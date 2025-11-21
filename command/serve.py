from flask import Flask
import sys
import os

sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from services.face_recognition import verify_faces_route


def serve():
    app = Flask(__name__)

    env_path = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "app.env")

    configuration = {}

    if os.path.exists(env_path):
        with open(env_path, "r") as f:
            for line in f:
                line = line.strip()
                if "=" in line and not line.startswith("#"):
                    key, value = line.split("=", 1)
                    configuration[key.strip()] = value.strip()
    else:
        raise FileNotFoundError("app.env file not found")

    python_port = int(configuration.get("PYTHON_PORT"))

    app.register_blueprint(verify_faces_route)

    app.run(host='0.0.0.0', port=python_port)
