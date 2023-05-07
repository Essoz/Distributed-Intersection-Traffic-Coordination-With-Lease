from flask import Flask
import flask
port = 11001

def get_cars_dict():
    return {
        'car1': 0,
        'car2': 1,
        'car3': 2,
    }

# run http server on port 11001
app = Flask(__name__)
@app.route('/perception/getSurrounding', methods=['GET'])
def getSurrounding():
    return flask.jsonify({'data': get_cars_dict()})

# run the server
app.run(port=port, host='localhost', debug=True, use_reloader=False)