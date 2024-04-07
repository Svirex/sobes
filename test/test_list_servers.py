import requests
from server import addr

first_list = sorted(["Adele Ingram",
	"Imran Bryan",
	"Cecilia Odom",
	"Jaydon Gould",
	"Elodie Hendrix"])

second_list = sorted(["Sumaiya Cruz",
	"Ernest Stafford",
	"Lorraine House",
	"Gregory O'Doherty",
	"Mikolaj Dale"])

def test_not_exists_header():
    r = requests.get(f"{addr}/api/servers")
    assert r.status_code == 401

def test_invalid_authorization_header():
    r = requests.get(f"{addr}/api/servers", headers={"Authorization": "sdfsdfsd"})
    assert r.status_code == 401

def test_invalid_authorization_scheme():
    r = requests.get(f"{addr}/api/servers", headers={"Authorization": "Hello fsdfds"})
    assert r.status_code == 401

def test_invalid_authorization_token():
    r = requests.get(f"{addr}/api/servers", headers={"Authorization": "Bearer fsdfds"})
    assert r.status_code == 401

def test_check_servers_list():
    r = requests.post(f"{addr}/api/auth", json={'login': 'test', 'password': 'test'}, headers={'Content-Type': "application/json"})
    body = r.json()
    r = requests.get(f"{addr}/api/servers", headers={"Authorization": f"Bearer {body['token']}"})
    assert r.status_code == 200

    body_first = sorted(r.json())
    r = requests.get(f"{addr}/api/servers", headers={"Authorization": f"Bearer {body['token']}"})
    assert r.status_code == 200
    body_second = sorted(r.json())

    assert (body_first == first_list and body_second == second_list) or (body_first == second_list and body_second == first_list)
