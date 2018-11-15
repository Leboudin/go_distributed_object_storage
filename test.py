import requests
import datetime
import os
import requests

addr = '0.0.0.0:8030'
if os.getenv('SERVER_ENV_ADDR') is not None:
    addr = os.getenv('SERVER_ENV_ADDR')


def test_put_object():
    with open('./test.txt', 'w') as fp:
        ts = datetime.datetime.now().strftime('%Y-%m-%d-%H-%M-%S')
        content = 'This is a test file, created at: {}'.format(ts)
        fp.write(content)

    with open('./test.txt', 'rb') as fp:
        # test put object
        object_name = 'obj-{}'.format(ts)
        api = 'http://{}/objects/{}'.format(addr, object_name)
        try:
            resp = requests.put(api, data=fp)
        except Exception as e:
            print('- erro: {}'.format(e))
            return

        if resp.status_code != 200:
            print('- non 200 response: {}, {}'.format(
                resp.status_code, 
                resp.content.decode()
            ))
        else:
            print('- seems PUT succeed')

    # test get object
    try:
        resp = requests.get(api)
    except Exception as e:
        print('- errro: {}'.format(e))
        return

    if resp.status_code != 200:
            print('- non 200 response: {}, {}'.format(
                resp.status_code, 
                resp.content.decode()
            ))
    else:
        got_content = resp.content.decode()
        if got_content == content:
            print('- seems GET succeed')
        else:
            print('- seems GET failed')




if __name__ == '__main__':
    test_put_object()

