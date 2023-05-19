import requests

function_url = 'http://ec2-15-188-193-232.eu-west-3.compute.amazonaws.com/invoke'


def invoke(id, args):
    if not isinstance(id, str) or not isinstance(args, dict):
        raise Exception("'id' must be string and 'args' must be a dict")

    requests.post({"id": id, "args": args})
