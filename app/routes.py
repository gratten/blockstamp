from app import app
from flask import render_template, flash, redirect, url_for
# from app.forms import ReceiveForm
# from app.forms import PayForm
# from flask_qrcode import QRcode
from app import block

# import requests

# qrcode = QRcode(app)

@app.route('/')
@app.route('/home')
def home():
    print('hello console', flush=True)
    best_block = block.get_best_block()
    # sys.stdout.flush()
    # r = requests.get('https://gratten.duckdns.org/api/v1/wallet', \
        # headers={'accept': 'application/json', 'X-API-KEY': '2a5c2c5417284bd994413e85358350d2'}, \
        # verify=False).json()
    # wallet = r['name']
    # balance = round(r['balance']/1000)
    return render_template('home.html', best_block=best_block)