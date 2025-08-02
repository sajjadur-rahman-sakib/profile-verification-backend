from flask import Flask
import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from services.face_recognition import verify_faces_route

app = Flask(__name__)
app.register_blueprint(verify_faces_route)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
