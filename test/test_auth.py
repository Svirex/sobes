import requests
from server import addr

def test_invalid_content_type():
    r = requests.post(f"{addr}/api/auth", data={'key': 'value'})
    assert r.status_code == 400

def test_empty_body():
    r = requests.post(f"{addr}/api/auth", headers={'Content-Type': "application/json"})
    assert r.status_code == 400

def test_invalid_data():
    r = requests.post(f"{addr}/api/auth", data="asdfasdfsfd", headers={'Content-Type': "application/json"})
    assert r.status_code == 400

def test_invalid_user_password():
    r = requests.post(f"{addr}/api/auth", json={'login': 'admin', 'password': 'admin'}, headers={'Content-Type': "application/json"})
    assert r.status_code == 400
    body = r.json()
    assert body['error'] == 'invalid login or password'

def test_good():
    r = requests.post(f"{addr}/api/auth", json={'login': 'test', 'password': 'test'}, headers={'Content-Type': "application/json"})
    assert r.status_code == 200
    body = r.json()
    assert 'token' in body

def test_invalid_field_names():
    r = requests.post(f"{addr}/api/auth", json={'login': 'test', 'pass': 'test'}, headers={'Content-Type': "application/json"})
    assert r.status_code == 400
    body = r.json()
    assert body['error'] == 'invalid login or password'