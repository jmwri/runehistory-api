from runehistory_api import make_app

app = make_app(__name__)

if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)
