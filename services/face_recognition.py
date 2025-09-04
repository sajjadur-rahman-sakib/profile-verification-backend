from flask import Blueprint, request, jsonify
import face_recognition

verify_faces_route = Blueprint('verify_faces_route', __name__)

@verify_faces_route.route('/verify', methods=['POST'])
def verify_faces():
    try:
        image1 = request.files['image1']
        image2 = request.files['image2']

        img1 = face_recognition.load_image_file(image1)
        img2 = face_recognition.load_image_file(image2)

        encoding1 = face_recognition.face_encodings(img1)
        encoding2 = face_recognition.face_encodings(img2)

        if not encoding1 or not encoding2:
            return jsonify({"error": "No faces detected in one or both images"}), 400

        results = face_recognition.compare_faces([encoding1[0]], encoding2[0], tolerance=0.6)
        return jsonify({"is_match": bool(results[0])})

    except Exception as e:
        return jsonify({"error": str(e)}), 500
