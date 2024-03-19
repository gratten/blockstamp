from flask_wtf import FlaskForm
from wtforms import SubmitField, SelectField, IntegerField 
from wtforms.validators import DataRequired

# units = ['sats', 'bitcoin']

class DateForm(FlaskForm):
    year = IntegerField('Year', validators=[DataRequired()])
    month = IntegerField('Month', validators=[DataRequired()])
    day = IntegerField('Day', validators=[DataRequired()])
    hour = IntegerField('Hour', validators=[DataRequired()])
    minute = IntegerField('Minute', validators=[DataRequired()])
    second = IntegerField('Second', validators=[DataRequired()])
    submit = SubmitField('Submit')
