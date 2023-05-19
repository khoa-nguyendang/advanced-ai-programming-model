import sys

from flask import Flask, jsonify, request


class AbstractServer():
    """
    Base class for Flask server, this class is for inheritance only
    """
    def __init__(self):
        self.app = Flask(__name__)

        self.request = request
        self.add_endpoint()
        self.init()

    def init(self):
        # subclass do building pipeline, ... in here
        pass

    def add_endpoint(self):
        # subclass add api endpoint in this function
        pass

    def run(self):
        self.app.run(debug=True)
