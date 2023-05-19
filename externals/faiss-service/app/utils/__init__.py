import logging
import time
from flask import current_app, jsonify

from app.config import get_config

config = get_config()

class AbtractError:

    Error = 'Err0r'


def catch_error_for_all_methods(Cls):
    """This is a decoration to try whenever get errors while communicating with database
    
    Arguments:
        Cls {Class} -- Catch error from all methods in Class
    
    Returns:
        [Function]
    """
    def catch_error(func, init_func):

        def func_wrapper(*args, **kwargs):
            if config.TESTING:
                return
            result = AbtractError.Error
            while result == AbtractError.Error:
                try:
                    result = func(*args, **kwargs)
                except Exception as e:
                    try:
                        init_func()
                    except:
                        pass
                    logging.error('Caught something: ')
                    logging.error(e)
                    time.sleep(5)
            return result

        return func_wrapper

    class ClassWrapper(object):
        def __init__(self, *args, **kwargs):
            self._intance = Cls(*args, **kwargs)

        def __getattribute__(self, s):
            """
            This is called whenever any attribute of a NewCls object is accessed. This function first tries to
            get the attribute off NewCls. If it fails then it tries to fetch the attribute from self._intance (an
            instance of the decorated class). If it manages to fetch the attribute from self._intance, and
            the attribute is an instance method then `time_this` is applied.
            """
            try:
                x = super(ClassWrapper, self).__getattribute__(s)
            except AttributeError:
                pass
            else:
                return x

            x = self._intance.__getattribute__(s)
            init_func = self._intance.__init__

            if type(x) == type(self.__init__):  # it is an instance method
                # this is equivalent of just decorating the method with time_this
                return catch_error(x, init_func)
            else:
                return x

    return ClassWrapper


def make_response(body, status_code=200):
    """
    Convenient wrapper for the `make_response` function to return a json response
    """
    response = current_app.make_response(jsonify(body))
    response.content_type = 'application/json'
    response.status_code = status_code
    return response
