import os

# SECRET_KEY = os.urandom(32)
# app.config['SECRET_KEY'] = SECRET_KEY

class Config(object):
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'you-will-never-guess'