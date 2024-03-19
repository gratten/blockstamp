from app import app
from flask import render_template, flash, redirect, url_for
from app.forms import DateForm
from app import block



@app.route('/', methods=['GET', 'POST'])
@app.route('/home', methods=['GET', 'POST'])
def home():

    form = DateForm()
    if form.validate_on_submit():
        year = form.year.data
        month = form.month.data
        day = form.day.data
        hour = form.hour.data
        minute = form.minute.data
        second = form.second.data

        block_height = year


        return redirect(url_for('block_height', block_height=block_height))

    return render_template('home.html', form=form)

@app.route('/block_height/<int:block_height>')
def block_height(block_height):
    return render_template('block_height.html', block_height=block_height)