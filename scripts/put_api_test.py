import pytest
import requests
import time

def test_put_director(url_prefix, created_director):
    url = f"{url_prefix}/directors/{created_director}"
    payload = {
        "Name": "Updated Director Name",
        "DateOfBirth": int(time.time()),
        "DateOfDeath": None,
        "Biography": "Updated biography"
    }
    
    response = requests.put(url, json=payload)
    assert response.status_code == 200
    
    data = response.json()
    assert data["Name"] == "Updated Director Name"
    assert data["Biography"] == "Updated biography"

def test_put_media(url_prefix, created_media, created_director):
    url = f"{url_prefix}/media/{created_media}"
    payload = {
        "Title": "Updated Movie Title",
        "DirectorID": created_director,
        "ReleasedAt": int(time.time()),
        "Synopsis": "Updated synopsis"
    }
    
    response = requests.put(url, json=payload)
    assert response.status_code == 200
    
    data = response.json()
    assert data["Title"] == "Updated Movie Title"
    assert data["Synopsis"] == "Updated synopsis"

def test_put_comment(url_prefix, created_comment, created_media):
    url = f"{url_prefix}/comments/{created_comment}"
    get = requests.get(url)

    payload = {
        "UserID": 1,
        "MediaID": created_media,
        "Content": "Updated comment content",
        "Rating": 5,
        "CreatedAt": get.json()["CreatedAt"]

    }
    
    response = requests.put(url, json=payload)
    assert response.status_code == 200
    
    data = response.json()
    assert data["Content"] == "Updated comment content"
    assert data["Rating"] == 5

