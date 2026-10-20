import pytest
import requests
import time

@pytest.fixture(scope="session")
def url_prefix():
    return "http://127.0.0.1:8080/api"

@pytest.fixture
def created_director(url_prefix):
    payload = {"Name": "Test Director", "DateOfBirth": int(time.time()), "Biography": "Test bio"}
    response = requests.post(f"{url_prefix}/directors", json=payload)
    assert response.status_code == 201
    director_id = response.json()["DirectorID"]
    
    yield director_id
    
    requests.delete(f"{url_prefix}/directors/{director_id}")

@pytest.fixture
def created_media(url_prefix, created_director):
    payload = {
        "Title": "Test Movie", 
        "DirectorID": created_director, 
        "ReleasedAt": int(time.time()), 
        "Synopsis": "Test synopsis"
    }
    response = requests.post(f"{url_prefix}/media", json=payload)
    assert response.status_code == 201
    media_id = response.json()["MediaID"]
    
    yield media_id
    
    requests.delete(f"{url_prefix}/media/{media_id}")

@pytest.fixture
def created_comment(url_prefix, created_media):
    payload = {"UserID": 1, "MediaID": created_media, "Content": "Test comment", "Rating": 5}
    response = requests.post(f"{url_prefix}/comments", json=payload)
    assert response.status_code == 201
    comment_id = response.json()["CommentID"]
    
    yield comment_id
    
    requests.delete(f"{url_prefix}/comments/{comment_id}")
