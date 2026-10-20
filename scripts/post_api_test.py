import pytest
import requests
import time

def test_create_director(url_prefix):
    url = url_prefix + "/directors"
    new_director = {
        "Name": "Quentin Tarantino",
        "DateOfBirth": int(time.time()),
        "Biography": "Director of Pulp Fiction"
    }
    response = requests.post(url, json=new_director)

    assert response.status_code == 201
    data = response.json()  
    assert "DirectorID" in data
    assert data["Name"] == "Quentin Tarantino"

def test_create_media(url_prefix, created_director):
    url = url_prefix + "/media"
    new_media = {
        "Title": "Pulp Fiction", 
        "DirectorID": created_director,
        "ReleasedAt": int(time.time()),
        "Synopsis": "Good Movie"
    }
    response = requests.post(url, json=new_media)

    assert response.status_code == 201
    data = response.json()  
    assert "MediaID" in data
    assert data["Title"] == "Pulp Fiction"
    assert data["DirectorID"] == created_director
    assert data["Synopsis"] == "Good Movie"

def test_create_comment(url_prefix, created_media):
    url = url_prefix + "/comments"
    new_comment = {
        "UserID": 1, 
        "MediaID": created_media,
        "Content": "It was alright. Should have more action.", 
        "Rating": 3
    }
    response = requests.post(url, json=new_comment)

    assert response.status_code == 201
    data = response.json()  
    assert "CommentID" in data
    assert "CreatedAt" in data
    assert data["UserID"] == 1
    assert data["MediaID"] == created_media
    assert data["Content"] == "It was alright. Should have more action."
    assert data["Rating"] == 3

