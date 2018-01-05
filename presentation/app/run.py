#By Sam Kreter
from flask import render_template
from flask import Flask
import requests
import json
import os

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0',port=8080)
