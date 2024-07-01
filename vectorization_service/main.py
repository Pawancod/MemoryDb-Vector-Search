# app.py
from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer

app = Flask(__name__)
model = SentenceTransformer('all-MiniLM-L6-v2')

@app.route('/vectorize', methods=['POST'])
def vectorize():
    data = request.json
    text = data['text']
    embeddings = model.encode([text])
    return jsonify({'vector': embeddings[0].tolist()})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)
