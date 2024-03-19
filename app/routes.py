from app import app
from flask import render_template, flash, redirect, url_for
from app.forms import DateForm
from app import block
import datetime


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

        target_time = int(datetime.datetime(year, month, day, hour, minute, second).timestamp())


        # block_height = given_datetime
        max_block_height = block.get_best_block()
        block_height = block.get_block_height(max_block_height, target_time)


        return redirect(url_for('block_height', block_height=block_height))

    return render_template('home.html', form=form)

@app.route('/block_height/<int:block_height>')
def block_height(block_height):
    return render_template('block_height.html', block_height=block_height)